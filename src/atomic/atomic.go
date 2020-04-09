package main

import (
	"fmt"
	"sync"
	"time"
)

type atomicInt struct {
	value int
	lock  sync.Mutex //
}

func (a *atomicInt) increment() {
	fmt.Println("safe increment")
	//  使用锁保护代码区,使用匿名函数 defer 控制在函数体里面
	func() {
		// 使用锁来保护
		a.lock.Lock()
		defer a.lock.Unlock()
		a.value++
	}()
}

func (a *atomicInt) get() int {
	a.lock.Lock()
	defer a.lock.Unlock()
	return int(a.value)
}

//  使用waitGroup 就不用枷锁了
func (a *atomicInt) syncAtomic() {
	var wg sync.WaitGroup

	// a increment
	a.increment()

	wg.Add(1)
	// create a goroutine ,let a increment
	go func() {
		a.increment()
		wg.Done()
	}()
	wg.Wait()
}

func (a *atomicInt) lockAtomic() {
	a.increment()

	go func() {
		a.increment()
	}()
}

func main() {
	var a atomicInt
	//a.syncAtomic()
	a.lockAtomic()
	// block一下要不 goroutine 没来的及increment 主程序就退出了 替代方案使用sync.WaitGroup
	time.Sleep(time.Millisecond)
	fmt.Println(a.get())
}
