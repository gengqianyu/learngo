package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	//设计原理
	//在 Goroutine 构成的树形结构中对信号进行同步以减少计算资源的浪费是 context.Context 的最大作用。
	//Go 服务的每一个请求都是通过单独的 Goroutine 处理的，HTTP/RPC 请求的 处理器 会启动新的 Goroutine 访问数据库和其他服务。
	//我们可能会创建多个 Goroutine 来处理一次请求，而 context.Context 的作用是在不同 Goroutine 之间同步请求特定数据、取消信号以及处理请求的截止日期。

	//每一个 context.Context 都会从最顶层的 Goroutine 一层一层传递到最下层。context.Context 可以在上层 Goroutine 执行出现错误时，将信号及时同步给下层。
	//当最上层的 Goroutine 因为某些原因执行失败时，下层的 Goroutine 由于没有接收到这个信号所以会继续工作；但是当我们正确地使用 context.Context 时，就可以在下层及时停掉无用的工作以减少额外资源的消耗：

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second) //创建了一个过期时间为 1s 的上下文
	defer cancel()
	//因为过期时间大于处理时间，所以我们有足够的时间处理该请求，因此输出如下
	//process request with 500ms
	//main context deadline exceeded
	//go handle(ctx, 500*time.Millisecond)

	//如果我们将处理请求时间增加至 1500ms，整个程序都会因为上下文的过期而被中止，因此输出如下
	//main context deadline exceeded
	//handle context deadline exceeded
	go handle(ctx, 1500*time.Millisecond)

	select { //main goroutine 会阻塞，直到 1 秒后 ctx 超过最后期限
	case <-ctx.Done():
		fmt.Println("main", ctx.Err())
	}

	//main goroutine 和 handle goroutine  多个 Goroutine 同时订阅 ctx.Done() 管道中的消息，

	//context 包中最常用的方法还是 context.Background、context.TODO，这两个方法都会返回预先初始化好的私有变量 background 和 todo，它们会在同一个 Go 程序中被复用
	//两个私有变量都是通过 new(emptyCtx) 语句初始化的，它们是指向私有结构体 context.emptyCtx 的指针，这是最简单、最常用的上下文类型
	//.context.Background 是上下文的默认值，所有其他的上下文都应该从它衍生出来；
	//.context.TODO 应该仅在不确定应该使用哪种上下文时使用；
	//在多数情况下，如果当前函数没有上下文作为入参，我们都会使用 context.Background 作为起始的上下文向下传递。

	//利用 context.Background() 默认 context ，衍生出一个新的子上下文 ctx，并返回用于取消该上下文的函数 cancel。
	//一旦我们执行返回的取消函数 cancel（），当前上下文以及它的子上下文都会被取消，所有的 Goroutine 都会同步收到这一取消信号。
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	//Go 语言中的 context.Context 的主要作用还是在多个 Goroutine 组成的树中同步取消信号以减少对资源的消耗和占用，虽然它也有传值的功能，但是这个功能我们还是很少用到。

	//在真正使用传值的功能时我们也应该非常谨慎，使用 context.Context 进行传递参数请求的所有参数一种非常差的设计，比较常见的使用场景是传递请求对应用户的认证令牌以及用于进行分布式追踪的请求 ID。
}

func handle(ctx context.Context, duration time.Duration) {
	select {
	case <-ctx.Done():
		fmt.Println("handle", ctx.Err())
	case <-time.After(duration):
		fmt.Println("process request with", duration)
	}
}
