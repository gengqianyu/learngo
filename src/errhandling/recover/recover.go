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
	"fmt"
)

func tryRecover() {
	//panic 只会触发当前 Goroutine 的 defer；
	//recover 只有在 defer 中调用才会生效；
	//panic 允许在 defer 中嵌套多次调用；
	defer func() {
		//能获取 panic 的值
		//如果无法处理可以重新 panic

		r := recover() // 返回任意类型
		// 获取panic的值 type assertion
		if err, ok := r.(error); ok {
			fmt.Println("Error occurred:", err)
		} else {
			// 无法处理重新 panic，main 函数中的 defer 成本 recover 到了这个错误并打印
			panic(fmt.Sprintf("I don't know what to do %v", r))
		}
	}() //匿名函数自调用

	// 停止当前函数执行
	//panic 能够改变程序的控制流，调用 panic 后会立刻停止执行当前函数的剩余代码，并在当前 Goroutine 中递归执行调用方的 defer；
	// 如果没有遇见 recover，程序退出
	//panic(errors.New("this is an error")) //可以panic任意类型
	panic(123)
}

func main() {
	// 接收二次 panic
	defer func() {
		fmt.Println("in main")
		r := recover()
		fmt.Println(r)
	}()

	//panic 只会触发当前 Goroutine 的延迟函数调用
	//当我们运行这段代码时会发现 main 函数中的 defer 语句并没有执行，执行的只有当前 Goroutine 中的 defer。
	//go func() {
	//	defer println("in goroutine")
	//	panic("")
	//}()

	tryRecover()

}
