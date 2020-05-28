package main

import (
	"fmt"
	"runtime"
	"sync"
)

var counter int = 0 //共享内存

func Count(lock *sync.Mutex) {
	lock.Lock()         // 上锁
	defer lock.Unlock() // 解锁
	// 临界区：只允许一个进程进入，进入另一个进程意味着什么? 不能被调度。
	// 被调度: 另一个进程只有被调度才能执行，才可能进入临界区，如何阻止调度? cli(); 临界区  sti(); 剩余区
	// 什么时候不好使? 多CPU(多核)…
	// 临界区保护的硬件原子指令法 sync.Mutex.lock
	counter++ // 修改临界区

	fmt.Println("counter =", counter)
}

func main() {
	// 互斥锁 应该是临界区保护的硬件原子指令法，属于一个硬件实现估计。
	lock := &sync.Mutex{}
	//启动10个goroutine 去改变counter的值
	for i := 0; i < 10; i++ {
		// 在goroutine 中上锁 执行改变共享变量counter的值
		go Count(lock)
	}

	loop := 0

	for {
		// 利用互斥锁实现阻塞等待
		lock.Lock() // 上锁

		c := counter // 临界区 操作

		lock.Unlock() // 解锁
		// go语言的runtime系统主动出让时间片，进行cpu调度。 用loop记录循环次数。
		// 关闭主动调度。循环次数明显增加。可以测试。
		// 在某些情况下打开主动调度，显然更节省资源
		runtime.Gosched()
		loop++

		if c >= 10 { //满足条件退出 当全部goroutine全部执行完毕 counter等于10时
			break
		}
	}
	fmt.Println(loop)
}
