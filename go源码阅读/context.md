## context是什么？  
context：上下文。  
golang 中的context是官方给出的协程间上下文包。协程之间可通过context的kv来传递消息，通过Done来进行协程间通信，通过Err来获取取消原因
## context解决了什么问题？  
context最主要解决的是协程并发问题。当一个上游业务依赖多个下游业务返回的数据，下游业务没有相互依赖，因此这里会并发获取下游业务数据，当某个下游业务响应变慢，会导致上游业务大量的协程堆积，这会造成内存占用变多，gc压力变大，从而影响服务性能，而最严重的情况下，协程无限堆积，最终导致雪崩，服务不可用。
## context如何解决这些问题的？
context通过在不同协程间共享变量（方式是通信）来达到超时控制，优雅退出所有goroutine，并释放资源。  
context包提供上下文机制，在不同goroutine间传递deadline，取消信号，或其它请求相关的信息。
## 源码解释
```
// 1. Context
type Context interface {
	Deadline() (deadline time.Time, ok bool) // 如果设置了deadline，则返回deadline，否则，返回false
	Done() <-chan struct{} // channel信号，如果工作完成（完成，超时，异常取消）会发出close信号。信号由WithCancel，WithDeadline，WithTimeout方法来触发
	Err() error // Done的错误原因，如果还未Done，返回nil，如果通过Cancel Done，返回canceled，如果通过Deadline Done，返回DeadlineExceeded
	Value(key interface{}) interface{}
}
// 2. emptyCtx 空的上下文,不支持设置k、v，不支持返回错误，不支持设置超时时间，不支持Done，即不会被WithCancel，withDeadline，withTimeout触发Done掉。  
// 两种场景使用：1.明确不需要超时控制或者结束控制的场景，使用context.Background() 2. 不知道使用场景，用context.TODO()方法先占位，比如重构的时候，未来再传来自上游的ctx
type emptyCtx int
func (*emptyCtx) Deadline() (deadline time.Time, ok bool) {
	return
}
func (*emptyCtx) Done() <-chan struct{} {
	return nil
}
func (*emptyCtx) Err() error {
	return nil
}
func (*emptyCtx) Value(key interface{}) interface{} {
	return nil
}
3. cancelCtx：可取消的上下文, 注意cancelCtx没有实现Deadline方法，因此并没有实现Context 接口
type canceler interface {
	cancel(removeFromParent bool, err error)
	Done() <-chan struct{}
}
type cancelCtx struct {
	Context

	mu       sync.Mutex            // protects following fields, mu保证了并发安全，即修改children，err，done都是并发安全的
	done     chan struct{}         // created lazily, closed by first cancel call
	children map[canceler]struct{} // set to nil by the first cancel call
	err      error                 // set to non-nil by the first cancel call
}
func (c *cancelCtx) Value(key interface{}) interface{} {
	if key == &cancelCtxKey { // 判断内存地址是否一致
		return c
	}
	return c.Context.Value(key)
}
func (c *cancelCtx) Done() <-chan struct{} {
	c.mu.Lock()
	if c.done == nil { // created lazily
		c.done = make(chan struct{})
	}
	d := c.done
	c.mu.Unlock()
	return d
}
func (c *cancelCtx) Err() error {
	c.mu.Lock()
	err := c.err
	c.mu.Unlock()
	return err
}
// 移除所有的child，如果当前ctx删除事件是来自parent发起的，那么从parent中删除自身ctx
func (c *cancelCtx) cancel(removeFromParent bool, err error) {
	if err == nil {
		panic("context: internal error: missing cancel error")
	}
	c.mu.Lock()
	if c.err != nil {
		c.mu.Unlock()
		return // already canceled
	}
	c.err = err
	if c.done == nil {
		c.done = closedchan
	} else {
		close(c.done)
	}
	// 将上下文的所有后代协程都取消了
	for child := range c.children {
		// NOTE: acquiring the child's lock while holding parent's lock.
		child.cancel(false, err)
	}
	c.children = nil
	c.mu.Unlock()

	if removeFromParent {
		removeChild(c.Context, c)
	}
}
4. timerCtx：计时器上下文, timerCtx继承了cancelCtx并实现Deadline，即timerCtx实现了Context接口, 当使用可取消的Context，实际使用的是timerCtx
type timerCtx struct {
	cancelCtx
	timer *time.Timer // Under cancelCtx.mu.

	deadline time.Time
}
func (c *timerCtx) Deadline() (deadline time.Time, ok bool) {
	return c.deadline, true
}
func (c *timerCtx) cancel(removeFromParent bool, err error) {
	c.cancelCtx.cancel(false, err)
	if removeFromParent {
		// Remove this timerCtx from its parent cancelCtx's children.
		removeChild(c.cancelCtx.Context, c)
	}
	c.mu.Lock()
	if c.timer != nil {
		c.timer.Stop()
		c.timer = nil
	}
	c.mu.Unlock()
}
// 将当前节点从parent中移除, 即接触parent和child关系，释放内存
func removeChild(parent Context, child canceler) {
	p, ok := parentCancelCtx(parent)
	if !ok {
		return
	}
	p.mu.Lock()
	if p.children != nil {
		delete(p.children, child)
	}
	p.mu.Unlock()
}
// 返回有截止时间的context，d：期望截止时间
func WithDeadline(parent Context, d time.Time) (Context, CancelFunc) {
	if parent == nil {
		panic("cannot create context from nil parent")
	}
	if cur, ok := parent.Deadline(); ok && cur.Before(d) {
		// The current deadline is already sooner than the new one.
		return WithCancel(parent)
	}
	// 通过基于父协程创建新的有截止时间的子协程来设定子协程的截止时间
	c := &timerCtx{
		cancelCtx: newCancelCtx(parent),
		deadline:  d,
	}
	propagateCancel(parent, c)
	dur := time.Until(d)
	if dur <= 0 {
		c.cancel(true, DeadlineExceeded) // deadline has already passed
		return c, func() { c.cancel(false, Canceled) }
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.err == nil {
		// 设定结束时间，执行的任务
		c.timer = time.AfterFunc(dur, func() {
			c.cancel(true, DeadlineExceeded)
		})
	}
	return c, func() { c.cancel(true, Canceled) }
}
// 主要作用是关联parent和child，并在parent done的时候，cancel掉child
func propagateCancel(parent Context, child canceler) {
	done := parent.Done()
	if done == nil {
		return // parent is never canceled，父节点以上的路径没有可取消的
	}

	select {
	case <-done:
		// parent is already canceled
		child.cancel(false, parent.Err())
		return
	default:
	}
	// 获取最近可取消的父节点
	if p, ok := parentCancelCtx(parent); ok {
		p.mu.Lock()
		if p.err != nil { // 报错取消child
			// parent has already been canceled
			child.cancel(false, p.err)
		} else { // 没有报错，表示parent没有被取消，关联parent和child
			if p.children == nil {
				p.children = make(map[canceler]struct{})
			}
			p.children[child] = struct{}{}
		}
		p.mu.Unlock()
	} else {
		atomic.AddInt32(&goroutines, +1)
		go func() { // 发起协程，监听parent的结束状态
			select { 
			case <-parent.Done():
				child.cancel(false, parent.Err())
			case <-child.Done():
			}
		}()
	}
}
// 获取最近类型是cancelCtx的parent context，用来执行取消操作
func parentCancelCtx(parent Context) (*cancelCtx, bool) {
	done := parent.Done()
	if done == closedchan || done == nil {
		return nil, false
	}
	p, ok := parent.Value(&cancelCtxKey).(*cancelCtx)
	if !ok {
		return nil, false
	}
	p.mu.Lock()
	ok = p.done == done
	p.mu.Unlock()
	if !ok {
		return nil, false
	}
	return p, true
}
5. valueCtx： 值上下文可设置一组kv数据
type valueCtx struct {
	Context // 继承具体的context（Background，TODO，WithDeadline, WithTimeout, WithCancel）
	key, val interface{}
}
```

propagate：传播


参考：  
https://36kr.com/p/1721518997505  
https://zhuanlan.zhihu.com/p/68792989
