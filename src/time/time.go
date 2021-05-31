package main

import (
	"fmt"
	"time"
)

func main() {
	t := time.Now()
	//格式化当前时间
	fmt.Println(t.Format("2006-01-02 15:04:05"))
	//时间戳
	fmt.Println(t.Unix())

	//	时间戳转时间
	fmt.Println(time.Unix(t.Unix(), 0).Format("2006-01-02 15:04:05"))

	// 程序执行时间
	defer func() {
		fmt.Println(time.Since(t))
	}()
	time.Sleep(1 * time.Second)
}
