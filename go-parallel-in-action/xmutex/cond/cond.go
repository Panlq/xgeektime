package cond

type Cond struct {
	noCopy noCopy

	// 当观察或者修改等待条件的时候需要加锁
	L Locker

	// 等待队列
	notify  notifyList
	checker copyChecker
}

func NewCond(l Locker) *Cond {
	return &Cond{L: l}
}

func (c *Cond) Wait() {
	c.checker.check()
	// 增加到等待队列中
	t := runtime_notifyListAdd(&c.notify)
	// Wait 把调用者加入到等待队列时会释放锁，在被唤醒之后还会请求锁。在阻塞休眠期间，调用者是不持有锁的，这样能让其他 goroutine 有机会检查或者更新等待变量。所以调用 Wait 前要先上锁。
	c.L.Unlock()
	// 阻塞休眠直到被唤醒
	runtime_notifyListWait(&c.notify, t)
	c.L.Lock()
}

func (c *Cond) Signal() {
	c.checker.check()
	runtime_notifyListNotifyOne(&c.notify)
}

func (c *Cond) Broadcast() {
	c.checker.check()
	runtime_notifyListNotifyAll(&c.notify)
}
