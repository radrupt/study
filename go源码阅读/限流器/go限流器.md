> 限流器是微服务中必不可少的一环，可以起到保护下游服务，防止服务过载等作用。time/rate 是基于 Token Bucket(令牌桶) 算法实现的限流。以下是源码解析
# 总览
time/rate 和常见令牌桶最大的区别是采用了lazyload的方式，即不会起一个携程定时王令牌池里加令牌，而是在每次使用的时候计算是否令牌数足够，不足够需要等待多久。
因此使用了两个数据结构Limiter和Reservation，分别表示限流器以及预定执行某个时间所需要的相关信息（如预期执行时间，是否可执行）
1. 数据结构
```
type Limiter struct {
	mu     sync.Mutex // 锁
	limit  Limit // 速率限制，即多久多少个令牌
	burst  int // 桶大小
	tokens float64 // 当前剩余tokens数量
	// last is the last time the limiter's tokens field was updated
	last time.Time // 最后一次取令牌时间
	// lastEvent is the latest time of a rate-limited event (past or future)
	lastEvent time.Time // 指的是最近一次消费的 timeToAct 值
}
```
2. Every方法,Limit方法
```
type Limit float64
// 计算得到每秒会产生几个令牌
func Every(interval time.Duration) Limit {
	if interval <= 0 {
		return Inf
	}
	return 1 / Limit(interval.Seconds())
}
// 获得每秒产生的令牌数
func (lim *Limiter) Limit() Limit {
	lim.mu.Lock()
	defer lim.mu.Unlock()
	return lim.limit
}
```
3. 获得桶大小
```
func (lim *Limiter) Burst() int {
	lim.mu.Lock()
	defer lim.mu.Unlock()
	return lim.burst
}
```
4. 创建限流器，简单工厂方法
```
// 参数r：每秒几个令牌，b：桶容量，即令牌数量达到桶容量后不再增加，桶提供了一种灵活应对突发流量的场景，如果桶容量为0，则降级为类似漏桶的方式，每秒的并发数量固定。真正类似漏桶需使用AllowN方案，AllowN方案会计算当前的tokens数是否满足想要获取的N个令牌
func NewLimiter(r Limit, b int) *Limiter {
	return &Limiter{
		limit: r,
		burst: b,
	}
}
```
5. Reservation（预定）保存了被限流器限制到一段时间后发生的事件
```
// A Reservation holds information about events that are permitted by a Limiter to happen after a delay.
// A Reservation may be canceled, which may enable the Limiter to permit additional events.
type Reservation struct {
	ok        bool // 是否可执行
	lim       *Limiter // 限流器
	tokens    int // 预定事件使用的令牌数
	timeToAct time.Time // 预计预定事件可执行的时间 
	// This is the Limit at reservation time, it can change later.
	limit Limit // 并发
}
```
6. 判断限流器是否允许该reservation预定N个数量令牌的事件执行
```
func (r *Reservation) OK() bool {
	return r.ok
}
```
7. 判断该限流器执行还需要多久时间
```
func (r *Reservation) Delay() time.Duration {
	return r.DelayFrom(time.Now())
}
func (r *Reservation) DelayFrom(now time.Time) time.Duration {
	if !r.ok { // 表示限流器不能保证在最大等待时间下，可以执行，即等待时间超过了最大时间
		return InfDuration
	}
	delay := r.timeToAct.Sub(now)
	if delay < 0 {
		return 0
	}
	return delay
}
```
8. 取消预定, 显然，如果取消预定，那就意味着有类似回滚的逻辑
```
// Cancel is shorthand for CancelAt(time.Now()).
func (r *Reservation) Cancel() {
	r.CancelAt(time.Now())
}

// CancelAt indicates that the reservation holder will not perform the reserved action
// and reverses the effects of this Reservation on the rate limit as much as possible,
// considering that other reservations may have already been made.
func (r *Reservation) CancelAt(now time.Time) {
	if !r.ok {
		return
	}

	r.lim.mu.Lock()
	defer r.lim.mu.Unlock()

	if r.lim.limit == Inf || r.tokens == 0 || r.timeToAct.Before(now) {
		return
	}

	// 计算需要恢复的token，即被该revervation预定了的tokens，需要恢复，算法是，上一次被执行了的事件时间+reservation预计可以执行的时间 = 多产生的tokens 设为a，当前limit的tokens设为b则，目标tokens c为=a+b, 需要恢复的即为b = c - a
	// The duration between lim.lastEvent and r.timeToAct tells us how many tokens were reserved
	// after r was obtained. These tokens should not be restored.
	// 即预定了的tokens=限流其中原有的tokens+（预计可执行时间+限流器最后发放令牌时间差释放的令牌）
	// 那么需要恢复的就是:限流其中原有的tokens = 预定了的tokens - 预计可执行时间+限流器最后发放令牌时间差释放的令牌）
	restoreTokens := float64(r.tokens) - r.limit.tokensFromDuration(r.lim.lastEvent.Sub(r.timeToAct))
	if restoreTokens <= 0 {
		return
	}
	// 此刻需求了预定，那么将取消作为一个事件来看，需要以当前时间重置限流器的token数量
	// advance time to now
	now, _, tokens := r.lim.advance(now)
	// calculate new number of tokens
	tokens += restoreTokens
	if burst := float64(r.lim.burst); tokens > burst {
		tokens = burst
	}
	// update state
	r.lim.last = now
	r.lim.tokens = tokens
	if r.timeToAct == r.lim.lastEvent {
		prevEvent := r.timeToAct.Add(r.limit.durationFromTokens(float64(-r.tokens)))
		if !prevEvent.Before(now) {
			r.lim.lastEvent = prevEvent
		}
	}
}
// 由于小数*小数精度问题，需要分别计算整数部分和小数部分，再相加
func (limit Limit) tokensFromDuration(d time.Duration) float64 {
	// Split the integer and fractional parts ourself to minimize rounding errors.
	// See golang.org/issues/34861.
	sec := float64(d/time.Second) * float64(limit)
	nsec := float64(d%time.Second) * float64(limit)
	return sec + nsec/1e9
}
```
9. 开启预定
```
func (lim *Limiter) Reserve() *Reservation {
	return lim.ReserveN(time.Now(), 1)
}
func (lim *Limiter) ReserveN(now time.Time, n int) *Reservation {
	r := lim.reserveN(now, n, InfDuration)
	return &r
}
// 返回reservation，标明获取n个tokens需要多久，以及是否可执行，比如n > bruse ok = false
func (lim *Limiter) reserveN(now time.Time, n int, maxFutureReserve time.Duration) Reservation {
	lim.mu.Lock()

	if lim.limit == Inf {
		lim.mu.Unlock()
		return Reservation{
			ok:        true,
			lim:       lim,
			tokens:    n,
			timeToAct: now,
		}
	}

	now, last, tokens := lim.advance(now) // 获取当前时间下桶里的令牌数

	// Calculate the remaining number of tokens resulting from the request.
	tokens -= float64(n) // 判断还需要多少令牌才能执行该N个事件

	// Calculate the wait duration
	var waitDuration time.Duration
	if tokens < 0 {
		waitDuration = lim.limit.durationFromTokens(-tokens)
	}

	// Decide result
	ok := n <= lim.burst && waitDuration <= maxFutureReserve

	// Prepare reservation
	r := Reservation{
		ok:    ok,
		lim:   lim,
		limit: lim.limit,
	}
	if ok {
		r.tokens = n
		r.timeToAct = now.Add(waitDuration)
	}

	// Update state
	if ok {
		lim.last = now
		lim.tokens = tokens
		lim.lastEvent = r.timeToAct
	} else {
		lim.last = last
	}

	lim.mu.Unlock()
	return r
}
func (lim *Limiter) advance(now time.Time) (newNow time.Time, newLast time.Time, newTokens float64) {
	last := lim.last
	if now.Before(last) {
		last = now
	}

	// Avoid making delta overflow below when last is very old.
	maxElapsed := lim.limit.durationFromTokens(float64(lim.burst) - lim.tokens) // 过多久会满
	elapsed := now.Sub(last)
	if elapsed > maxElapsed { // 如果当前时间距离上次事件执行时间，已经将桶放满了，则使用放满的时间，举个例子差5个放满，每秒1个，距离3秒，则返回3s，秒2个，则返回2.5秒
		elapsed = maxElapsed
	}

	// Calculate the new number of tokens, due to time that passed.
	delta := lim.limit.tokensFromDuration(elapsed) // 从时间换算新获取多少个tokens
	tokens := lim.tokens + delta
	if burst := float64(lim.burst); tokens > burst { // 最大只能获取桶容量个令牌
		tokens = burst
	}

	return now, last, tokens
}

```
10. 等待方法
```
// Wait is shorthand for WaitN(ctx, 1).
func (lim *Limiter) Wait(ctx context.Context) (err error) {
	return lim.WaitN(ctx, 1)
}

// WaitN blocks until lim permits n events to happen.
// It returns an error if n exceeds the Limiter's burst size, the Context is
// canceled, or the expected wait time exceeds the Context's Deadline.
// The burst limit is ignored if the rate limit is Inf.
func (lim *Limiter) WaitN(ctx context.Context, n int) (err error) {
	lim.mu.Lock()
	burst := lim.burst
	limit := lim.limit
	lim.mu.Unlock()

	if n > burst && limit != Inf {
		return fmt.Errorf("rate: Wait(n=%d) exceeds limiter's burst %d", n, burst)
	}
	// Check if ctx is already cancelled
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	// Determine wait limit
	now := time.Now()
	waitLimit := InfDuration
	if deadline, ok := ctx.Deadline(); ok {
		waitLimit = deadline.Sub(now)
	}
	// Reserve
	r := lim.reserveN(now, n, waitLimit)
	if !r.ok {
		return fmt.Errorf("rate: Wait(n=%d) would exceed context deadline", n)
	}
	// Wait if necessary
	delay := r.DelayFrom(now)
	if delay == 0 {
		return nil
	}
	t := time.NewTimer(delay) // 通过计时器+context来控制等待
	defer t.Stop()
	select {
	case <-t.C:
		// We can proceed.
		return nil
	case <-ctx.Done():
		// Context was canceled before we could proceed.  Cancel the
		// reservation, which may permit other events to proceed sooner.
		r.Cancel()
		return ctx.Err()
	}
}
```
11. 动态更新brust