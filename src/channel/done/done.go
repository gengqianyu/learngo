/*
使用channel 来等待goroutine结束

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
)

func doWork(id int, c chan int, done chan bool) {
	// continuous receiving
	for element := range c {
		fmt.Printf("Worker %d received %c->%v\n", id, element, element)

		// 不用通过共享内存来通信，通过通信来共享内存
		// 通知外面事情做完
		// 为了防止，循环逻辑端，无法收到 done，我们开一个goroutine，让它并行去处理这个done
		//done <- true
		go func() {
			done <- true
		}()
	}
}

type worker struct {
	in   chan int
	done chan bool
}

// channel as the return value
func createWorker(id int) worker {
	// create channel
	w := worker{
		in:   make(chan int),
		done: make(chan bool),
	}

	// channel is interaction between goroutine and goroutine
	// must be received by another goroutine, so create an goroutine
	go doWork(id, w.in, w.done)

	return w
}

func chanDemo() {
	//define an array of type worker as an registrar
	var workers [10]worker

	for i := 0; i < 10; i++ {
		// add an element to the registrar
		// element as a goroutine
		workers[i] = createWorker(i)
	}

	//一口气发20个任务
	for i, worker := range workers {
		// send a data to the channel
		// character type can be operated,string type cannot
		// example: 'a'+1 got 'b' ;"a"+"b" got "ab"
		worker.in <- 'a' + i

		// 加上这个就可以不用下面sleep了
		//<-worker.done
	}
	for i, worker := range workers {
		worker.in <- 'A' + i
		// worker done 放在这里有个问题，打印结果是顺序打印的。它有什么缺点呢
		// 仿佛并行的创建10个worker 变的毫无意义。
		// 这里变成了，发一个任务，让它打。打完了以后，我才打下一个。
		// 我们不希望让他发一个任务，再让它等它结束。希望啊一口气把这20个任务发出去
		// 因此这里接收done，不太合理，应该把它提出去放到最后
		//<-worker.done
	}

	//等到发送任务结束我在退出来,在这里一次性收20个work done
	// wait for all of them
	// 这样写也有一个错误，第一组小写的打印出来，第二组没出来，发生了deadlock错误（循环等待）
	// ！注意：我们所有向channel发送，都是阻塞（block）的。
	// 一端向channel发送一个任务 ,channel的另一端必须有人收。
	// doWork端中发了一个done 那另一端(循环任务端)必须去收这个done。而收的逻辑我们写在了最后面。
	//因此呢上面第一个 循环逻辑端，向channel发了一个任务后，doWork端处理完毕，却没有去接收 doWork端 发回来的done。
	// 这时第二个循环逻辑端，又向这个channel发送任务；而这时这个channel的 doWork端 正在等待，第一个循环端 接收它发送的done。
	// 所以第二次的任务发送后，就会阻塞在channel 的doWork端
	// 解决办法是，在doWork，发送done的时候，在去开一个goroutine 去并行的发送这个done；
	// 这样doWork的主逻辑就不被卡住了。

	// worker发一个任务，doWork完事发出来的worker done，
	// 但是这个done是最后收的
	for _, worker := range workers {
		// 上面每个worker发了，两遍任务。所以这里是两个
		<-worker.done
		<-worker.done
	}
	// 阻塞1毫秒 不然上面的goroutine还没来的及打印第二个数据，主goroutine就会退出
	// time.Sleep(time.Millisecond)
}

func main() {
	fmt.Println("channel as first-class citizen")
	chanDemo()
}
