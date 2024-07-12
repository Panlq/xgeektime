package main

import (
	"sync"
	"sync/atomic"
)

type Once struct {
	m    sync.Mutex
	done uint32
}

func (x *Once) Do(f func()) {
	if atomic.LoadUint32(&x.done) == 0 {
		x.doSlow(f)
	}
}

func (x *Once) doSlow(f func()) {
	x.m.Lock()
	defer x.m.Unlock()

	if x.done == 0 {
		defer atomic.StoreUint32(&x.done, 1)
		f()
	}
}

func main() {
}
