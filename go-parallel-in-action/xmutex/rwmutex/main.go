package rwmutex

import (
	"sync"
	"sync/atomic"
)

type RWMutex struct {
	w           sync.Mutex // held if there are pending writers
	writerSem   uint32     // 写 阻塞信号
	readerSem   uint32     // 读 阻塞信号
	readerCount int32      // 正在读的调用者数量/ 当为负数时 表示有write持有锁
	readerWait  int32      // writer持有锁之前正等待解锁的数量
}

const rwmutexMaxReaders = 1 << 30

func (rw *RWMutex) RLock() {
	if atomic.AddInt32(&rw.readerCount, 1) < 0 {
		// 写端 持有锁， 读端阻塞
		runtime_SemacquireMutex(&rw.readerSem, false, 0)
	}
}

func (rw *RWMutex) RUnlock() {
	if r := atomic.AddInt32(&rw.readerCount, -1); r < 0 {
		rw.rUnlockSlow(r)
	}
}

func (rw *RWMutex) rUnlockSlow(r int32) {
	if r+1 == 0 || r+1 == -rwmutexMaxReaders {
		fatal("sync: RUnlock of unlocked RWMutex")
	}

	if atomic.AddInt32(&rw.readerWait, -1) == 0 {
		// 无读者等待，唤醒写端等待者
		runtime_Semrelease(&rw.writerSem, false, 1)
	}
}

func (rw *RWMutex) Lock() {
	// 1. 先尝试获取互斥锁
	rw.w.Lock()
	// 2. 看是否有其他正持有锁的读者，有则阻塞
	r := atomic.AddInt32(&rw.readerCount, -rwmutexMaxReaders) + rwmutexMaxReaders
	if r != 0 && atomic.AddInt32(&rw.readerWait, r) != 0 {
		// rc - rwmutexMaxReaders + rwmutexMaxReaders > 0说明还有等待者, 写端阻塞
		runtime_SemacquireMutex(&rw.writerSem, false, 0)
	}
}

func (rw *rwMutex) UnLock() {
	r := atomic.AddInt32(&rw.readerCount, rwmutexMaxReaders)
	if r >= rwmutexMaxReaders {
		fatal("sync: Unlock of unlocked RWMutex")
	}

	// 如果有等待的读者，先唤醒
	for i := 0; i < int(r); i++ {
		runtime_Semrelease(&rw.readerSem, false, 0)
	}

	// 释放互斥锁
	rw.w.Unlock()
}
