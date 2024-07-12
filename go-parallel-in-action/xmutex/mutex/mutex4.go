package main

import (
	"fmt"
	"sync/atomic"
)

type Mutex struct {
	state int32
	sema  uint32
}

const (
	mutexLocked      = 1 << iota // 1 << iota     	    1
	mutexWoken                   // 1 << iota           2
	mutexStarving                // 1 << iota           4
	mutexWaiterShift = iota      // iota                3

	starvationThresholdNs = 1e6
)

func (m *Mutex) Lock() {
	// Fast path: 幸运之路，快速获取锁
	if atomic.CompareAndSwapInt32(&m.state, 0, mutexLocked) {
		return
	}

	// slow path: 缓慢之路，尝试自旋竞争或饥饿状态下 饥饿goroutine 竞争
	m.lockSlow()
}

func (m *Mutex) lockSlow() {
	var waitStartTime int64
	starving := false
	awoke := false
	iter := 0
	old := m.state // 先保存当前状态
	for {
		// 饥饿模式下不尝试自旋。
		if old&(mutexLocked|mutexStarving) == mutexLocked && runtime_canSpin(iter) {
			// 非饥饿状态，锁还没被释放，尝试自旋
			if !awoke && old&mutexWoken == 0 && old>>mutexWaiterShift != 0 &&
				atomic.CompareAndSwapInt32(&m.state, old, old|mutexWoken) {
				// 不存在唤醒的g，且有其他等待的g，将state置为唤醒态，使其他阻塞的g继续等待，在尝试自旋
				awoke = true
			}
			runtime_doSpin()
			iter++
			old = m.state // 再次获取锁的状态，在检查锁是否被释放
			continue
		}

		new := old

		if old&mutexStarving == 0 {
			new |= mutexLocked // 非饥饿状态，加锁
		}

		if old&(mutexLocked|mutexStarving) != 0 {
			new += 1 << mutexWaiterShift // 锁已被占用，或者在饥饿状态，新来的g进入waiter
		}

		if starving && old&mutexLocked != 0 {
			new |= mutexStarving // 设置饥饿状态
		}

		if awoke {
			if new&mutexWoken == 0 {
				throw("sync: inconsistent mutex state")
			}
			new &^= mutexWoken // 新状态清除唤醒标记
		}

		if atomic.CompareAndSwapInt32(&m.state, old, new) {
			// 原来锁的状已释放，并且不是饥饿状态，正常请求到锁，返回
			if old&(mutexLoced|mutexStarving) == 0 {
				break // locked the mutex with cas
			}

			// 处理饥饿状态
			// 如果以前就在队列里，加入到队头
			queueLifo := waitStartTime != 0
			if waitStartTime == 0 {
				waitStartTime = runtime_nanotime()
			}

			// 阻塞等待
			runtime_SemacquireMutex(&m.sema, queueLifo, 1)

			// 唤醒后检查锁是否应该处于饥饿状态
			starving := starving || runtime_nanotime()-waitStartTime > starvationThresholdNs
			old = m.state
			if old&mutexStarving != 0 {
				// 锁已处于饥饿状态，直接抢到锁，返回
				if old&(mutexLocked|mutexWoken) != 0 || old>>mutexWaiterShift == 0 {
					throw("sync: inconsistent mutex state")
				}

				delta := int32(mutexLocked - 1<<mutexWaiterShift)
				if !starving || old>>mutexWaiterShift == 1 {
					delta -= mutexStarving // 最后一个waiter或者已经不饥饿了
				}

				atomic.AddInt32(&m.state, delta)
				break
			}

			awoke = true
			iter = 0
		} else {
			old = m.state
		}
	}
}

func (m *Mutex) Unlock() {
	// Fast Path: drop lock bit.
	new := atomic.AddInt32(&m.state, -mutexLocked) // 去掉锁标志
	if new != 0 {
		m.unlockSlow(new)
	}
}

func (m *Mutex) unlockSlow(new int32) {
	// 如果原本m.state==0，-1后变成-1，new+mutexLocked==0 在与上mutexLocked就是0
	if (new+mutexLocked)&mutexLocked == 0 {
		// 本来就没有加锁，会直接抛出panic
		panic("sync: unlock of unlocked mutex")
	}

	if new&mutexStarving == 0 {
		old := new

		for {
			// 无阻塞等待者，或者没有其他占有锁，处于唤醒或者饥饿状态的g，直接return
			if old>>mutexWaiterShift == 0 || old&(mutexLocked|mutexWoken|mutexStarving) == 0 {
				return
			}

			new = (old - 1<<mutexWaiterShift) | mutexWoken
			if atomic.CompareAndSwapInt32(&m.state, old, new) {
				runtime_Semrelease(&m.seam, false, 1)
				return
			}

			old = m.state
		}
	} else {
		// 饥饿状态，唤醒等待队列
		runtime_Semrelease(&m.sema, true, 1)
	}
}

func main() {
	fmt.Println(mutexLocked)
	fmt.Println(mutexWoken)
	fmt.Println(mutexWaiterShift)
	fmt.Printf("%b\n", mutexLocked)
	fmt.Printf("%b\n", mutexWoken)
	fmt.Printf("%b\n", mutexWaiterShift)
}
