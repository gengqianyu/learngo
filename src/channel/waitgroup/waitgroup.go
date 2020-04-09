/*
使用waitGroup 来等待goroutine结束
对比 done 这个更加简洁/抽象

channel:
goroutine和goroutine之间的通信的通道

{进程 [调度器] [线程 (goroutine)<-channel->(goroutine)]}

理论基础：Communication Sequential Process (CSP) model 通信顺序过程

apply channel
Don't communicate by sharing memory; share memory by communicating.
不要通过共享内存来通信，通过通信来共享内存.

*/
package main

import (
	"fmt"
	"sync"
)

func doWork(id int, worker worker) {
	// continuous receiving
	for element := range worker.in {
		fmt.Printf("Worker %d received %c->%v\n", id, element, element)
		// 通知任务结束
		worker.done()
	}
}

type worker struct {
	in   chan int
	done func() // 函数式编程
}

// channel as the return value
func createWorker(id int, wg *sync.WaitGroup) worker {
	// create channel
	w := worker{
		in: make(chan int),
		done: func() {
			wg.Done()
		},
	}

	// channel is interaction between goroutine and goroutine
	// must be received by another goroutine, so create an goroutine
	go doWork(id, w)

	return w
}

func chanDemo() {
	//define an array of type worker as an registrar
	var workers [10]worker
	// create an sync wait group
	//waitGroup 只能定义一个 所以worker中的wg要通过传参，传过去
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		// add an element to the registrar
		// element as a goroutine work struct
		workers[i] = createWorker(i, &wg)
	}

	// 向waitGroup中添加20个任务
	wg.Add(20)

	for i, worker := range workers {
		// send a task to the channel
		// character type can be operated,string type cannot
		// example: 'a'+1 got 'b' ;"a"+"b" got "ab"
		worker.in <- 'a' + i
	}
	for i, worker := range workers {
		worker.in <- 'A' + i
	}

	// 等待所有任务结束 waiting for the end of all tasks
	wg.Wait()
}

func main() {
	fmt.Println("channel as first-class citizen")
	chanDemo()
}
