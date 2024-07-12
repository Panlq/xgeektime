package xchan

import (
	"fmt"
	"time"
)

type ChMutex struct {
	ch chan struct{}
}

func NewChMutex() *ChMutex {
	ch := make(chan struct{}, 1)
	ch <- struct{}{}
	return &ChMutex{ch}
}

func (c *ChMutex) Lock() {
	<-c.ch
}

func (c *ChMutex) Unlock() {
	select {
	case c.ch <- struct{}{}:
	default:
		panic("unlock of unlocked mutex")
	}
}

func (c *ChMutex) TryLock() bool {
	select {
	case <-c.ch:
		return true
	default:
		return false
	}
}

func (c *ChMutex) LockTimeout(timeout time.Duration) bool {
	timer := time.NewTicker(timeout)
	select {
	case <-c.ch:
		timer.Stop()
		return false
	case <-timer.C:
	}
	return true
}

func (c *ChMutex) IsLocked() bool {
	return len(c.ch) == 0
}

func main() {
	m := NewChMutex()
	ok := m.TryLock()
	fmt.Printf("locked v %v\n", ok)
	m.Unlock()
}
