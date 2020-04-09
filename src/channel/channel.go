/*
channel:
goroutine和goroutine之间的通信的通道

{进程 [调度器] [线程 (goroutine)<-channel->(goroutine)]}

channel
buffered channel
close and range

理论基础：Communication Sequential Process (CSP) model 通信顺序过程

apply channel
Don't communicate by sharing memory; share memory by communicating.
不要通过共享内存来通信，通过通信来共享内存.

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
	//// 生产者往channel里发一个消息，必须得有消费者去收。如果没人收就会阻塞
	//	// 创建channel时，加第二个参数，生产者往channel里发1024个消息，再让消费者是收，这样就大大提高了效率
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
