package recursivemutex

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/petermattis/goid"
)

type RecursiveMutex struct {
	sync.Mutex
	owner     int64 // 当前持有锁的g
	recursion int32 // g 重入次数
}

func (m *RecursiveMutex) Lock() {
	gid := goid.Get()
	if atomic.LoadInt64(&m.owner) == gid {
		m.recursion++
		return
	}

	m.Mutex.Lock()
	// 获得锁的g第一次调用，记录它的gid
	atomic.StoreInt64(&m.owner, gid)
	m.recursion = 1
}

func (m *RecursiveMutex) Unlock() {
	gid := goid.Get()
	if atomic.LoadInt64(&m.owner) != gid {
		panic(fmt.Sprintf("wrong the owner(%d): %d!", m.owner, gid))
	}

	m.recursion--
	if m.recursion != 0 {
		return
	}

	// g 最后一次调用，释放锁
	atomic.StoreInt64(&m.owner, -1)
	m.Mutex.Unlock()
}
