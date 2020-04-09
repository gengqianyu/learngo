package main

import (
	"fmt"
	"time"
)

func main() {
	//example 1
	//defer的执行顺序是倒序的，先进后出
	//for i := 0; i < 5; i++ {
	//	defer fmt.Println(i)
	//}
	//fmt.Println("test defer")

	//example 2
	//defer 传入的函数不是在退出代码块的作用域时执行的，它是在当前函数和方法返回之前被调用
	//注意这里的代码块(code block)和函数方法的区别,
	//一对大括号就能定义一个代码块。大括号内的一些数据只能在这个括号内使用，所以是有作用域的概念。
	// output#
	//code block ends
	//main ends
	//defer runs (如果是以代码块为准，拿应该是第二个输出，显然不是)

	//{
	//	defer fmt.Println("defer runs")
	//	fmt.Println("code block ends")
	//}
	//
	//fmt.Println("main ends")

	//example 3
	//	与计算参数
	//go语言中所有的函数传参都是传值的，defer 虽然是关键字，但是也继承了这个特性

	startedAt := time.Now()
	//不会得到预期的结果1秒
	//调用defer关键字会立刻对函数中引用的外部参数进行拷贝，
	//后面表达式time.Since(startedAt)的结果不是在main函数退出之前计算的，
	//而是在 defer 关键字调用时计算的，最终导致输出0s
	//defer fmt.Println(time.Since(startedAt))
	//下面拷贝的是函数的指针
	defer func() { fmt.Println(time.Since(startedAt)) }()
	time.Sleep(time.Second)
}
