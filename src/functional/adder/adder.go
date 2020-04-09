/*
函数式编程vs函数指针
1.函数是一等公民：参数，变量，返回值都可以是函数
2.高阶函数：一个函数它的参数也可以是函数
3.函数->闭包
正统的函数式编程
1.不可变性：不能有状态，只有常量和函数（不能有变量，甚至连选择判断语句都不能有）
2.函数只能有一个参数

*/
package main

import (
	"fmt"
)

// 定义一个adder函数，他的返回值是一个函数
func adder() func(int) int {
	sum := 0
	//下面返回的这个函数就是一个闭包
	//一个闭包包括（函数体->(局部变量，自由变量)）
	// v参数就可以被看成一个局部变量 函数体内sum是一个自由变量（编译器会连一根线连接到函数体外的sum变量中）
	return func(v int) int {
		sum += v
		return sum
	}
}

//正统闭包实现不能有状态记录
//返回参数一个是当前加完的值，一个呢是它下一轮要执行的函数体 它是一个递归的定义
type iAdder func(int) (int, iAdder)

func adder2(base int) iAdder {
	return func(v int) (int, iAdder) {
		//base+v是当前加完的局部和，
		//adder2(base + v)拿当前的局部和，去初始化下次累加的函数体
		return base + v, adder2(base + v)
	}
}
func main() {
	//函数adder()返回值为一个函数 所以变量a是一个func(int) int格式的函数
	a := adder()
	for i := 0; i < 10; i++ {
		fmt.Printf("0+1+...%d=%d\n", i, a(i))
	}
	fmt.Println("------------------------------------------")
	// 正统函数
	b := adder2(0)
	var s int
	for i := 0; i < 10; i++ {
		//每执行一次循环，返回的b函数体的参数base都会变（每一次都是上一次的局部和），
		//也就是说每次执行返回的函数体b都是不一样的，再用这个新的函数体去执行下次循环
		//后面的b每执行一次都是用前面返回的函数体b，每次执行b函数体其实都会变
		s, b = b(i)
		fmt.Printf("0+1+...%d=%d\n", i, s)
	}
}
