/*
channel:
goroutine和goroutine之间的通信的通道

{进程 [调度器] [线程 (goroutine)<-channel->(goroutine)]}

channel
buffered channel
close and range

理论基础：Communication Sequential Process (CSP) model 通信顺序进程

虽然我们在 Go 语言中也能使用共享内存加互斥锁进行通信，
但是 Go 语言提供了一种不同的并发模型，即通信顺序进程（Communicating sequential processes，CSP）。
Goroutine 和 Channel 分别对应 CSP 中的实体和传递信息的媒介，Goroutine 之间会通过 Channel 传递数据。

apply channel
Don't communicate by sharing memory; share memory by communicating.
不要通过共享内存来通信，通过通信来共享内存.

直接发送 #
如果目标 Channel 没有被关闭并且已经有处于读等待的 Goroutine，那么 runtime.chansend 会从 recvq(接收队列)中取出最先陷入等待的 recv Goroutine 并直接向它发送数据：

	if sg := c.recvq.dequeue(); sg != nil {
		send(c, sg, ep, func() { unlock(&c.lock) }, 3)//需要注意的是，发送数据的过程只是将接收方的 Goroutine 放到了处理器的 runnext 中，程序没有立刻执行该 Goroutine。
		return true
	}

发送数据时会调用 runtime.send，该函数的执行可以分成两个部分：

1.调用 runtime.sendDirect 将发送的数据直接拷贝到 x = <-c 表达式中变量 x 所在的内存地址上；等待 recv（接收方） 来拷贝
2.调用 runtime.goready 将等待接收数据的 Goroutine 标记成可运行状态 Grunnable.
  并把该 Goroutine 放到发送方所在的处理器的 runnext 上等待执行，该处理器在下一次调度时会立刻唤醒数据的接收方；
func send(c *hchan, sg *sudog, ep unsafe.Pointer, unlockf func(), skip int) {
	if sg.elem != nil {
		sendDirect(c.elemtype, sg, ep)
		sg.elem = nil
	}
	gp := sg.g
	unlockf()
	gp.param = unsafe.Pointer(sg)
	goready(gp, skip+1)
}
我们在这里可以简单梳理和总结一下使用 ch <- i 表达式向 Channel 发送数据时遇到的几种情况：

1.如果当前 Channel 的 recvq（接收队列） 上存在已经被阻塞的 Goroutine，那么会直接将数据发送给当前 Goroutine 并将其设置成下一个运行的 Goroutine；
2.如果 Channel 存在缓冲区并且其中还有空闲的容量，我们会直接将数据存储到缓冲区 sendx 所在的位置上；
3.如果不满足上面的两种情况，会创建一个 runtime.sudog 结构并将其加入 Channel 的 sendq（注意是发送队列） 队列中，当前 Goroutine 也会陷入阻塞等待其他的协程从 Channel 接收数据；

发送数据的过程中包含几个会触发 Goroutine 调度的时机：

1.发送数据时发现 Channel 上存在等待接收数据的 Goroutine，立刻设置处理器的 runnext 属性，但是并不会立刻触发调度；
2.发送数据时并没有找到接收方并且缓冲区已经满了，这时会将自己加入 Channel 的 sendq（发送队列） 队列并调用 runtime.goparkunlock 触发 Goroutine 的调度让出处理器的使用权；让消费 goroutine 去消费



直接接收 #
当 Channel 的 sendq 队列中包含处于等待状态的 Goroutine 时，那么runtime.chanrecv 会从 sendq（接收队列）中取出最先陷入等待的 send Goroutine 并直接向它接收数据，
处理的逻辑和发送时相差无几，只是发送数据时调用的是 runtime.send 函数，而接收数据时使用 runtime.recv：

	if sg := c.sendq.dequeue(); sg != nil {
		recv(c, sg, ep, func() { unlock(&c.lock) }, 3)
		return true, true
	}

func recv(c *hchan, sg *sudog, ep unsafe.Pointer, unlockf func(), skip int) {
	if c.dataqsiz == 0 {
		if ep != nil {
			recvDirect(c.elemtype, sg, ep)
		}
	} else {
		qp := chanbuf(c, c.recvx)
		if ep != nil {
			typedmemmove(c.elemtype, ep, qp)
		}
		typedmemmove(c.elemtype, qp, sg.elem)
		c.recvx++
		c.sendx = c.recvx // c.sendx = (c.sendx+1) % c.dataqsiz
	}
	gp := sg.g
	gp.param = unsafe.Pointer(sg)
	goready(gp, skip+1)
}
该函数会根据缓冲区的大小分别处理不同的情况：

如果 Channel 不存在缓冲区；
	1.调用 runtime.recvDirect 将 Channel 发送队列中 Goroutine 存储的 elem 数据拷贝到目标内存地址中；
如果 Channel 存在缓冲区；
	1.将队列中的数据拷贝到接收方的内存地址；
	2.将发送队列头的数据拷贝到缓冲区中，释放一个阻塞的发送方；
无论发生哪种情况，运行时都会调用 runtime.goready 将当前处理器的 runnext 设置成发送数据的 Goroutine，在调度器下一次调度时将阻塞的发送方唤醒。

1.如果 Channel 为空，那么会直接调用 runtime.gopark 挂起当前 Goroutine；
2.如果 Channel 已经关闭并且缓冲区没有任何数据，runtime.chanrecv 会直接返回；
3.如果 Channel 的 sendq 队列中存在挂起的 Goroutine，会将 recvx 索引所在的数据拷贝到接收变量所在的内存空间上并将 sendq 队列中 Goroutine 的数据拷贝到缓冲区；
4.如果 Channel 的缓冲区中包含数据，那么直接读取 recvx 索引对应的数据；
5.在默认情况下会挂起当前的 Goroutine，将 runtime.sudog 结构加入 recvq 队列并陷入休眠等待调度器的唤醒；

6.4.6 关闭管道 #
编译器会将用于关闭管道的 close 关键字转换成 OCLOSE 节点以及 runtime.closechan 函数。

当 Channel 是一个空指针或者已经被关闭时，Go 语言运行时都会直接崩溃并抛出异常：

func closechan(c *hchan) {
	if c == nil {
		panic(plainError("close of nil channel"))
	}

	lock(&c.lock)
	if c.closed != 0 {
		unlock(&c.lock)
		panic(plainError("close of closed channel"))
	}
处理完了这些异常的情况之后就可以开始执行关闭 Channel 的逻辑了，下面这段代码的主要工作就是将 recvq 和 sendq 两个队列中的数据加入到 Goroutine 列表 gList 中，
与此同时该函数会清除所有 runtime.sudog 上未被处理的元素：

	c.closed = 1

	var glist gList
	for {
		sg := c.recvq.dequeue()
		if sg == nil {
			break
		}
		if sg.elem != nil {
			typedmemclr(c.elemtype, sg.elem)
			sg.elem = nil
		}
		gp := sg.g
		gp.param = nil
		glist.push(gp)
	}

	for {
		sg := c.sendq.dequeue()
		...
	}
	for !glist.empty() {
		gp := glist.pop()
		gp.schedlink = 0
		goready(gp, 3)
	}
}
Go
该函数在最后会为所有被阻塞的 Goroutine 调用 runtime.goready 触发调度。

*/
package main

import (
	"fmt"
	"time"
)

func worker(id int, c chan int) {
	// continuous receiving
	//for element := range c {
	for {
		// receive an element from this channel
		// element为具体的元素，ok表示是否还有值
		element, ok := <-c
		// break if ok is false
		if !ok {
			break
		}
		fmt.Printf("Worker %d received %v\n", id, element)
	}
}

// channel as the return value
// <- 修饰符放到chan后面表示这个channel是用来送数据的，收到这个channel的人只能给它发数据
// <- 修饰符放到chan后面表示这个channel只能来收数据，不能用来发
// 返回一个channel 常见的写法 创建一个channel 就把它返回；
// 这个channel做什么事情，我们交给一个goroutine协程去做
// 建了一个channel。开了一个goroutine。立马就发回。
func createWorker(id int) chan<- int {
	// create channel
	c := make(chan int)

	// channel is interaction between goroutine and goroutine
	// must be received by another goroutine, so create an goroutine
	go worker(id, c)

	return c
}

func chanDemo() {
	//define an array of type chan int as an registrar
	var channels [10]chan<- int

	for i := 0; i < 10; i++ {
		// add an element to the registrar
		// element as a goroutine
		channels[i] = createWorker(i)
	}
	for i := 0; i < 10; i++ {
		// send a data to the channel
		// character type can be operated,string type cannot
		// example: 'a'+1 got 'b' ;"a"+"b" got "ab"
		channels[i] <- 'a' + i

		channels[i] <- 'A' + i
		// 编译错误这个channel只能发数据不能收数据
		//n:=<-channels[i]
	}

	// 阻塞1毫秒 不然上面的goroutine还没来的及打印第二个数据，主goroutine就会退出
	time.Sleep(time.Millisecond)
}

func bufferedChannel() {
	// create an channel ，并规定这个channel缓冲区中可以存3个元素,会提高性能
	bc := make(chan int, 3)
	go worker(0, bc)
	bc <- 'a'
	bc <- 'b'
	bc <- 'c'
	// 如果不用另外的goroutine去接数据，只有三个数据被允许发送到缓冲区，那么下面是发送不了数据的。
	bc <- 'd'
	// 等一下go worker 打印
	time.Sleep(time.Millisecond)
}

// 如果发送的数据有明显的结尾，channel是可以close的
// channel close 永远是发送方去close
// 发送方Close，通知接收方方，我没有数据需要发了
func channelClose() {
	// 生产者往channel里发一个消息，必须得有消费者去收。如果没人收就会阻塞
	// 创建channel时，加第二个参数，生产者往 channel 里发 1024 个消息，再让消费者去收，这样就大大提高了效率
	cc := make(chan int, 3)
	go worker(0, cc)
	cc <- 'a'
	cc <- 'b'
	cc <- 'c'
	// 如果不用另外的goroutine去接数据，只有三个数据被允许发送到缓冲区，那么下面是发送不了数据的。
	cc <- 'd'
	// close channel 通知接收方我发完了
	close(cc)

	// 等一下go worker 打印
	time.Sleep(time.Millisecond)
}
func main() {
	fmt.Println("channel as first-class citizen")
	chanDemo()

	fmt.Println("buffered channel")
	bufferedChannel()

	fmt.Println("channel close and range")
	channelClose()
}
