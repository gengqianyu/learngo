/*
select
channel 调度
使用select来进行调度
使用select
使用三种定时器
在select 中使用nil channel
*/
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func generator() chan int {
	out := make(chan int)
	go func() {
		i := 0
		for {
			// 随机一个范围到1500的整数，转成持续时间 单位是毫秒
			time.Sleep(time.Duration(rand.Intn(1500)) * time.Millisecond)
			out <- i
			i++
		}
	}()
	return out
}

func worker(id int, c chan int) {
	// continuous receiving
	for element := range c {
		// sleep 5秒，去模拟，c1,c2消息发送端，发的快，消息处理端worker处理打印慢的情况
		//time.Sleep(5 * time.Second)
		time.Sleep(time.Second)
		fmt.Printf("Worker %d received %v\n", id, element)
	}
}

func createWorker(id int) chan<- int {
	// create channel
	c := make(chan int)

	// channel is interaction between goroutine and goroutine
	// must be received by another goroutine, so create an goroutine
	// 接收并处理从channel（c）收到的数据
	go worker(id, c)

	return c
}

// 第一个版本，会有问题
func versionOne() {
	// nil channel 是可以在select中正常运行，但是它是不可能select到结果的。所以select永远执行default
	//var c1, c2 chan int // c1 and c2 zeroValue =nil

	// 定义两个 pipeline 用于两个 goroutine 通信
	c1, c2 := generator(), generator()
	// 创建另一个 pipeline (worker)去接收，通过 select 从 pipeline (c1,c2)接收的消息 n，
	// 在通过 select 转发收到的消息 n
	worker := createWorker(0)

	// receive data from c1 and c2
	//n1 <- c1
	//n2 <- c2
	//初始化接收者消息 n
	n := 0
	// 用它标识是否接收到了从c1和c2发来的消息。
	hasValue := false
	for {
		// 定义一个只能给它发送消息的worker
		var activeWorker chan<- int
		// 如果worker 有值了，就赋给activeWorker
		// 利用了activeWorker 初始值为nil 也能在select中正常执行。
		// 如果一收到消息，就让activeWorker去代理一下worker
		if hasValue {
			activeWorker = worker
		}
		//我想要从c1和c2内同时收，谁来的快，我要谁。
		// 这个怎么办呢，这个就是我们的select
		// 用select和default做了一个非阻塞的接收数据
		// channel 的接收/发送都是block的，这样接收不到数据，它不会在哪里等着，而是去执行default
		// 重点：select 通过不同的case 既可以收消息，又可以发消息.(好强大)
		select {
		/*
			这个版本虽然worker可以处理并打印，来自c1，c2发来的消息，但是有个问题。
			之所以能正常打印，那是因为 从c1，c2接收消息比较慢，c1,c2的发送端的goroutine通过channel向c1，c2发消息的时候block了，TimeSleep了
			worker处理这些消息比较快。所以没问题。
			那么如果c1,c2发送端的goroutine通过channel向c1，c2发消息非常快；那么worker正在等待I/O打印消息n。
			c1,c2接收端又收到一个消息，并把它赋给n。那么此时变量n，在未打印之前就被重写成下一个消息了。
			从而worker在打印的时候，会处理不了所有消息（n来回被新消息重写）
		*/
		case n = <-c1: // 从c1接收消息，然后发给w
			//在case里面，执行这个语句是阻塞的不是太好，
			// 解决方案是，在写一个case 让它去把收到的n发给w，这样就不阻塞了
			//w <- n
			hasValue = true
		case n = <-c2: // 从c2接收消息，然后发给w
			//w <- n
			hasValue = true
		case activeWorker <- n: // 将c1和c2收到的消息，发给activeWorker if hasValue true activeWorker is worker
			hasValue = false
			//default:
			//	fmt.Println("no value received")
		}
	}
}

// 改进版
func versionTwo() {
	// 定义两个channel
	c1, c2 := generator(), generator()
	// 创建另一个channel (worker)去接收，通过select从channel (c1,c2)接收的消息n，
	// 在通过select转发过来的消息n
	worker := createWorker(0)
	// 初始化接收者消息队列
	var messages []int
	// 定义一个计时器，time.After返回一个chan time；
	// 下面的意思是10秒钟以后，向channel（exitTime）发送一个时间消息
	exitTime := time.After(10 * time.Second)
	// 定时查看messages队列长度,每秒一次 time.Tick返回也是一个 chan time
	// 下面意思是说，会定时的向channel tick 发送一个时间消息。tick的类型为 <-chan Time 是一个channel
	tick := time.Tick(time.Second)
	for {
		// 定义一个只能给它发送消息的worker
		var activeWorker chan<- int
		var message int
		if len(messages) > 0 {
			activeWorker = worker
			message = messages[0]
		}
		// 注意：空的 select 语句会直接阻塞当前 Goroutine，导致 Goroutine 进入无法被唤醒的永久休眠状态。
		// 在通常情况下，如果 select 中没有 case 准备好，select 语句会阻塞当前 Goroutine 并等待多个 case 中的一个 channel 达到可以收发的状态。
		// 但是如果 select 控制结构中包含 default 语句，那么这个 select 语句在执行时会遇到以下两种情况：
		//	.当存在可以收发的 Channel 时，直接处理该 Channel 对应的 case；
		//	.当不存在可以收发的 Channel 时，执行 default 中的语句；
		// select 在遇到多个 <-ch 同时满足可读或者可写条件时会随机选择一个 case 执行其中的代码。因此 select 执行顺序是无法预测的。
		// select 语句一次只能执行一个 case 分支，因此想要执行所有 case 分支，必须在外层加 for 循环。
		select {

		case n := <-c1: // 从c1接收消息,加入消息队列
			messages = append(messages, n)
		case n := <-c2: // 从c2接收消息，加入消息队列
			messages = append(messages, n)
		case activeWorker <- message: // 将c1和c2收到的消息，发给activeWorker
			messages = messages[1:]
		case <-time.After(800 * time.Millisecond): // 等待消息超过8毫秒就，打印一个超时
			fmt.Println("time out")
		case <-tick: // 每次select到channel（tick）的发来的消息就打印一下，message length
			fmt.Println("the messages len is ", len(messages))
		case <-exitTime: // 等select到channel（exitTime）的消息退出时间,就结束程序
			fmt.Println("bye")
			return
			//default:
			//	fmt.Println("no value received")
		}
	}
}

func main() {
	//versionOne()
	versionTwo()
}
