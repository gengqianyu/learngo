/*
panic
停止当前函数执行
一直向上返回，执行每一层的defer
如果没有遇见recover，程序退出

recover
仅在defer调用中使用
能获取panic的值
如果无法处理可以重新panic
*/
package main

import (
	"errors"
	"fmt"
)

func tryRecover() {
	defer func() {
		//recover 仅在defer调用中使用
		//能获取panic的值
		//如果无法处理可以重新panic

		r := recover() // 返回任意类型
		// 获取panic的值 type assertion
		if err, ok := r.(error); ok {
			fmt.Println("Error occurred:", err)
		} else {
			// 无法处理重新panic
			panic(fmt.Sprintf("I don't know what to do %v", r))
		}
	}() //匿名函数自调用

	// 停止当前函数执行
	// 一直向上返回，执行每一层的defer
	// 如果没有遇见recover，程序退出
	panic(errors.New("this is an error")) //可以panic任意类型
	//panic(123)
}
func main() {
	tryRecover()
}
