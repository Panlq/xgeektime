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
	mutexWaiterShift = iota      // iota                2
)

func (m *Mutex) Lock() {
	// Fast Path: 幸运case， 能够直接获取到锁
	if atomic.CompareAndSwapInt32(&m.state, 0, mutexLocked) {
		return
	}

	awoke := false
	iter := 0
	for {
		old := m.state           // 先保存当前状态
		new := old | mutexLocked // 新状态加锁

		if old&mutexLocked != 0 {
			// 锁还没被释放
			if runtime_canSpin(iter) {
				// 还可以自旋
				if !awoke && old&mutexWoken == 0 && old>>mutexWaiterShift {
					atomic.CompareAndSwapInt32(&m.state, old, old|mutexWoken)
					awoke = true
				}

				runtime_doSpin()
				iter++
				// 继续自旋，尝试获取锁
				continue
			}
			// 等待着数量加1，即第三位+1
			new = old + 1<<mutexWaiterShift
		}

		if awoke {
			// 唤醒状态
			if new&mutexWoken == 0 {
				panic("sync: inconsistent mutex state")
			}
			// goroutine 是被唤醒的，新状态清楚唤醒标志
			new &^= mutexWoken
		}

		// 在执行加锁的情况下，可能已经被新的g抢占了锁，所以会加锁失败，for循环后重新进入休眠状态
		if atomic.CompareAndSwapInt32(&m.state, old, new) {
			// 设置新状态
			if old&mutexLocked == 0 {
				// 锁原状态未加锁, 新goroutine已加锁，直接退出自旋
				break
			}

			// 休眠信号量-> 中断
			runtime_Semacquire(&m.sema)

			// 中断被唤醒
			awoke = true
			iter = 0
		}
	}
}

func (m *Mutex) Unlock() {
	// Fast Path: drop lock bit.
	new := atomic.AddInt32(&m.state, -mutexLocked) // 去掉锁标志
	// 如果原本m.state==0，-1后变成-1，new+mutexLocked==0 在与上mutexLocked就是0
	if (new+mutexLocked)&mutexLocked == 0 {
		// 本来就没有加锁，会直接抛出panic
		panic("sync: unlock of unlocked mutex")
	}

	// 以下逻辑 处理等待获取锁的goroutine --> waiter，通过信号量方式唤醒其中一个
	old := new
	for {
		// 如果没有其他的waiter，说明对这个锁的竞争的g只有一个，直接return
		if old>>mutexWaiterShift == 0 || old&(mutexLocked|mutexWoken) != 0 {
			// old&(mutexLocked|mutexWoken)
			// 如果这个时候有唤醒的g，或者又被别人加锁了，那解锁的同志无需操劳，直接return，其他g自己干得2挺好
			return
		}

		// 有等待者且没有唤醒的waiter，则唤醒一个等待的g
		new = (old - 1<<mutexWaiterShift) | mutexWoken // 将waiter数量-1，并将mutexWoken标置位置1

		// 在执行唤醒waiter的时候，可能已经被新来的g上锁，重新进入for循环后，在70行直接return
		if atomic.CompareAndSwapInt32(&m.state, old, new) {
			// runtime.Semrelease(&m.sema)
			return
		}

		old = m.state
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
