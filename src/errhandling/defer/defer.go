/*
资源管理
defer调用
1.确保调用在函数结束时发生
2.参数在defer语句时计算
3.defer列表为后进先出
何时使用defer调用
1.open/close
2.lock/unlock
3.printHeader/printFooter
*/
package main

import (
	"bufio"
	"fmt"
	"functional/fib"
	"os"
)

func tryDefer() {
	// defer 维护了一个栈，这个栈呢是先进后出的，所以下面打印结果是3，2，1
	defer fmt.Println(1)
	defer fmt.Println(2)
	fmt.Println(3)
	// defer 不怕中间有return
	//return
	// 甚至不怕panic
	//panic("error occurred")
	//fmt.Println(4)
	for i := 0; i < 100; i++ {
		defer fmt.Println(i)
		if i == 30 {
			panic("printed too many")
		}
	}
}

// 定义一个写文件
func writeFile(fileName string) {
	//创建一个文件
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	//defer 在函数结束，和return之前调用
	defer file.Close()

	//直接操作文件是很慢的，所以先将文件加载到文件缓冲区
	writer := bufio.NewWriter(file)
	// 操作完bufio以后再把缓冲区的内容，刷到文件中
	defer writer.Flush()
	// init fibonacci generator
	f := fib.Fibonacci()
	for i := 0; i < 20; i++ {
		// 把序列写入缓冲区 注意用的是Fprintln
		fmt.Fprintln(writer, f())
	}
}

func openFile(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		// 既然err都是实现一个Error接口，那么就可以用type assertion 类型断言方式
		// 取出接口变量的值
		if pathErr, ok := err.(*os.PathError); ok {
			fmt.Println(
				"Op is ", pathErr.Op,
				",path is ", pathErr.Path,
				",Error is ", pathErr.Err,
			)
		} else {
			fmt.Println("unknown error", err)
		}
	}
	defer file.Close()
}
func main() {
	//tryDefer()
	writeFile("fib.txt")
	// 打开一个不存在的文件出发pathError
	openFile("acb.txt")
}
