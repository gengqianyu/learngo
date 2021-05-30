package main

import (
	"fmt"
	"time"
)

func main() {
	//example 1
	//向 defer 关键字传入的函数会在函数返回之前运行
	//defer的执行顺序是倒序的，先进后出
	//for i := 0; i < 5; i++ {
	//	defer fmt.Println(i)
	//}
	//fmt.Println("test defer")

	//example 2
	//defer 传入的函数不是在退出代码块的作用域时执行的，它是在当前函数和方法 return 之前被调用
	//注意这里的代码块(code block)和函数方法的区别,
	//一对大括号就能定义一个代码块。大括号内的一些数据只能在这个括号内使用，所以是有作用域的概念。

	//{
	//	defer fmt.Println("defer runs")
	//	fmt.Println("code block ends")
	//}
	//
	//fmt.Println("main ends")

	// output#
	//code block ends
	//main ends
	//defer runs (如果是以代码块为准，那应该是第二个输出，而不是最后输出，显然不是)

	//example 3
	//	与计算参数
	//go 语言中所有的函数传参都是传值的，defer 虽然是关键字，但是也继承了这个特性
	executeTime()
}

func executeTime() {
	// 函数开始执行的时间
	startedAt := time.Now()
	//不会得到预期的结果1秒
	//调用 defer 关键字会立刻对函数中引用的外部参数进行拷贝，
	//后面表达式 time.Since(startedAt) 的结果不是在 main 函数退出之前计算的，
	//而是在 defer 关键字调用时计算的，最终导致输出 0s
	defer fmt.Println(time.Since(startedAt))
	//虽然调用 defer 关键字时也使用值传递，但是因为拷贝的是函数指针，所以 time.Since(startedAt) 会在 main 函数返回前调用并打印出符合预期的结果。
	defer func() {
		fmt.Println(time.Since(startedAt))
	}()
	//休眠 1 秒钟
	time.Sleep(time.Second)
}
