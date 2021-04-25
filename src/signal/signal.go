package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	//（1）何为优雅（graceful）？
	//Linux Server端的应用程序经常会长时间运行，在运行过程中，可能申请了很多系统资源，也可能保存了很多状态。
	//在这些场景下，我们希望进程在退出前，可以释放资源或将当前状态 dump 到磁盘上或打印一些重要的日志，即希望进程优雅退出。
	//（2）从对优雅退出的理解不难看出：优雅退出可以通过捕获 SIGTERM 来实现。
	//A、注册SIGTERM信号的处理函数并在处理函数中做一些进程退出的准备，信号处理函数的注册sigaction()来实现。
	//B、在主进程的main()中，通过类似于while(!fQuit)的逻辑来检测那个flag变量，一旦fQuit在signal handler function中被置为true，
	//则主进程退出while()循环，接下来就是一些释放资源或dump进程当前状态或记录日志的动作，完成这些后，主进程退出。

	//创建监听退出chan
	c := make(chan os.Signal)
	//监听指定信号 ctrl+c kill
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {

		for s := range c {

			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				fmt.Println("Program Exit...", s)
				GraceFullExit()
			//case syscall.SIGUSR1:
			//	fmt.Println("usr1 signal", s)
			//case syscall.SIGUSR2:
			//	fmt.Println("usr2 signal", s)
			default:
				fmt.Println("other signal", s)
			}
		}
	}()

	fmt.Println("Program Start...")
	sum := 0
	//模拟计算消耗
	for {
		sum++
		fmt.Println("sum:", sum)
		time.Sleep(time.Second)
	}
}
func GraceFullExit() {
	fmt.Println("Start Exit...")
	fmt.Println("Execute Clean...")
	fmt.Println("End Exit...")
	os.Exit(0)
}
