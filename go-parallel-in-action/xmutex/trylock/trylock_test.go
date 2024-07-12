package trylock

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"
	"unsafe"
)

const (
	mutexLocked      = 1 << iota // 1 << iota     	    1
	mutexWoken                   // 1 << iota           2
	mutexStarving                // 1 << iota           4
	mutexWaiterShift = iota      // iota                4
)

type Mutex struct {
	sync.Mutex
}

func (m *Mutex) TryLock() bool {
	if atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(&m.Mutex)), 0, mutexLocked) {
		return true
	}

	// 如果处于唤醒，加锁或者饥饿状态，这次请求就不参与竞争
	old := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	if old&(mutexLocked|mutexStarving|mutexWoken) != 0 {
		return false
	}

	// 尝试在竞争状态下请求锁
	new := old | mutexLocked

	return atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(&m.Mutex)), old, new)
}

func TestTryMutext(t *testing.T) {
	var mu Mutex
	go func() {
		mu.Lock()
		time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
		mu.Unlock()
	}()

	time.Sleep(time.Second)

	ok := mu.TryLock()
	if ok {
		t.Log("got the lock")
		mu.Unlock()
		return
	}

	t.Log("can't get the lock")
}
