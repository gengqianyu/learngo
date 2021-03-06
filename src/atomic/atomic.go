package main

import (
	"fmt"
	"sync"
)

// 原子 int
type atomicInt struct {
	value int
	lock  sync.Mutex //
}

// 自增
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

// 获取
func (a *atomicInt) get() int {
	a.lock.Lock()
	defer a.lock.Unlock()
	return int(a.value)
}

//  使用 waitGroup 就不用加锁了
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
	//利用 sync.WaitGroup 控制并发
	a.syncAtomic()

	//a.lockAtomic()
	// block 一下要不 goroutine 没来的及 increment 主程序就退出了 替代方案使用 sync.WaitGroup
	//time.Sleep(time.Millisecond)

	fmt.Println(a.get())
}
