package main

import "fmt"

type TestStruct struct{}

func NilOrNot(i interface{}) bool {
	return i == nil
}

func main() {
	var t *TestStruct
	fmt.Println(t == nil)    //true
	fmt.Println(NilOrNot(t)) //false
	//调用 NilOrNot 函数时发生了隐式的类型转换，除了向方法传入参数之外，变量的赋值也会触发隐式类型转换。
	//在类型转换时，*TestStruct 类型会转换成 interface{} 类型，转换后的变量不仅包含转换前的变量，还包含变量的类型信息 TestStruct。
	//所以转换后的变量与 nil 不相等。再次说明了Go 语言的interface{}接口类型不是任意类型
}
