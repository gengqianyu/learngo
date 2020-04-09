package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// 定义一个类型
type fib func() int

//定义一个斐波那契数列生成器
func fibonacci() fib {
	a, b := 0, 1
	return func() int {
		a, b = b, a+b
		return a
	}
}

func shengtuzi(n int, f fib) int {
	sum := 0
	for i := 1; i <= n; i++ {
		sum = f()
	}
	return sum
}

//终极实现法则
func fibonacci2(m int) int {
	if m == 0 {
		return 0
	}
	if m < 3 {
		return 1
	}
	return fibonacci2(m-1) + fibonacci2(m-2)
}

//为函数实现接口，go语言的厉害之处
//这是因为函数的接收者和普通的参数没有什么区别。只是调用上不一样
//接收者写在前面，调用的时候用f.Read，写在后面调用时用Read（f）
// Read是什么意思呢，它的意思是将下一个元素写入到p []byte中去。
func (f fib) Read(p []byte) (n int, err error) {
	//像文件读取一行一样，我们先获取一个数
	next := f()
	// 因为f（）是生成器是永远读不完的所以就加个判断终止获取
	if next > 10000 {
		return 0, io.EOF
	}
	//将这个next数转成一个字符串
	s := strconv.Itoa(next) + "\n"
	//s := fmt.Sprintf("%d\n", next)
	//生成一个string reader，让它代理去实现reader功能
	return strings.NewReader(s).Read(p)
}

//打印文件内容
func printFileContents(reader io.Reader) {
	// 生成一个扫描器
	scanner := bufio.NewScanner(reader)
	//fmt.Printf("scanner type:%T value:%v \n", scanner, scanner)
	// 逐行扫描打印文件内容
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}

func main() {
	//定义一个变量f
	f := fibonacci()
	//fmt.Println(f())
	//fmt.Println(f())
	//fmt.Println(f())
	//fmt.Println(f())
	//fmt.Println(f())
	//fmt.Println(f())
	//fmt.Println(f())
	//fmt.Println(f())
	//fmt.Println(f())
	//fmt.Println(f())

	//第10个月一共生多少
	fmt.Printf("第10个月一共生了%d对兔子", shengtuzi(10, f))
	fmt.Println()
	fmt.Println(fibonacci2(10))

	// 函数实现接口调用 f实现了reader接口那么它就可以当文件来用
	printFileContents(f)
}
