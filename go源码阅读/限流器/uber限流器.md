> uber限流器使用漏桶策略来限制请求qps
大体逻辑是：
1. 创建
```
rl := ratelimit.New(100)
// 通过指定的rate创建限流器
func New(rate int, opts ...Option) Limiter {
	return newAtomicBased(rate, opts...)
}
// 工厂函数创建限流器
func newAtomicBased(rate int, opts ...Option) *atomicLimiter {
	// TODO consider moving config building to the implementation
	// independent code.
	config := buildConfig(opts)
	perRequest := config.per / time.Duration(rate) // 每两次请求之间的最小时间间隔
	l := &atomicLimiter{
		perRequest: perRequest,
		// 最大松弛量，最简单的漏桶策略是严格处理每两次请求之间的间隔，不错弹性处理
		// 但是实际互联网场景中，每次请求的响应时间总是不固定的，有些请求间隔长，有些请求间隔短，
		// 比如设定的是每10ms一次请求，a请求花费了15ms，b请求花费了10毫秒，c请求花费了5毫秒，那么如果严格处理，将需要花费35ms
		// 但是其实可以通过松弛量来优化到30ms，保证rate是正确的。即把之前请求比较短的不超过perRequest时间匀给后面的请求
		maxSlack:   -1 * time.Duration(config.slack) * perRequest, 
		clock:      config.clock,
	}

	initialState := state{
		last:     time.Time{},
		sleepFor: 0,
	}
	atomic.StorePointer(&l.state, unsafe.Pointer(&initialState))
	return l
}
// 创建配置
func buildConfig(opts []Option) config {
	c := config{
		clock: clock.New(),
		slack: 10, // 最大松弛量的单位，用于控制最多一个perRequest最多可以并发执行多少个请求
		per:   time.Second,
	}

	for _, opt := range opts {
		opt.apply(&c)
	}
	return c
}
```
2. 判断当前请求是否可以执行
```
func (t *atomicLimiter) Take() time.Time {
	var (
		newState state
		taken    bool
		interval time.Duration
	)
	for !taken {
		now := t.clock.Now()

		previousStatePointer := atomic.LoadPointer(&t.state)
		oldState := (*state)(previousStatePointer)

		newState = state{
			last:     now,
			sleepFor: oldState.sleepFor,
		}

		// If this is our first request, then we allow it.
		if oldState.last.IsZero() {
			taken = atomic.CompareAndSwapPointer(&t.state, previousStatePointer, unsafe.Pointer(&newState))
			continue
		}

		// sleepFor calculates how much time we should sleep based on
		// the perRequest budget and how long the last request took.
		// Since the request may take longer than the budget, this number
		// can get negative, and is summed across requests.
		// maxSlack=-100ms, 最大松弛量
		// 新的请求时间-上一次请求时间设为b
		// 允许的请求间隔设为a
		// 若a=10ms，b=5ms，a-b=5ms
		// 若a=10ms, b=500ms, a-b=-490ms
		// sleepFor： 表示相较请求时间间隔（perRequest）还需要睡眠多久，<0 表示不需要睡眠
		// sleepFor 被作为松弛量使用，用来处理请求时间不固定的场景
		// 即如果sleepFor=-90ms, a-b=5ms,则sleepFor = -85ms, 直接执行不用等待
		// 如果sleepFor=-3ms,a-b=5ms,则sleepFor=2ms,新的请求仅需要再等待2ms，即两次请求之间间隔7ms
		newState.sleepFor += t.perRequest - now.Sub(oldState.last)
		// We shouldn't allow sleepFor to get too negative, since it would mean that
		// a service that slowed down a lot for a short period of time would get
		// a much higher RPS following that.
		if newState.sleepFor < t.maxSlack {
			newState.sleepFor = t.maxSlack
		}
		if newState.sleepFor > 0 {
			newState.last = newState.last.Add(newState.sleepFor)
			interval, newState.sleepFor = newState.sleepFor, 0
		}
		taken = atomic.CompareAndSwapPointer(&t.state, previousStatePointer, unsafe.Pointer(&newState))
	}
	t.clock.Sleep(interval)
	return newState.last
}
```