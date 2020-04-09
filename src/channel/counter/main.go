package main

import (
	"fmt"
	"time"
)

func Count(ch chan int, i int) {
	fmt.Println("Counting ", i)
	time.Sleep(5 * time.Second)
	ch <- i
	close(ch)
}

func main() {
	// 声明一个长度为10类型为chan int的slice
	chs := make([]chan int, 10)

	for i := 0; i < 10; i++ {
		// 初始化chan int切片
		chs[i] = make(chan int)
		go Count(chs[i], i)
	}
	// 程序不会在这里阻塞，你会发现第一个打出的就是"阻塞"。而是等到收的时候阻塞
	for _, ch := range chs {
		fmt.Println("阻塞")
		// 如果没人给当前ch发消息程序会阻塞在这里等待goroutine(count函数体)向channel发消息
		// 在count函数中加个sleep模拟延时添加消息
		v := <-ch
		fmt.Printf("counter %d\n", v)
	}
}
