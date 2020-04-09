/*
首先要知道 Go 语言在传递参数时都是传值的。
对于 &Cat{} 来说，这意味着拷贝一个新的 &Cat{} 指针，这个指针与原来的指针指向一个相同并且唯一的结构体，所以编译器可以隐式的对变量解引用（dereference）获取指针指向的结构体；
对于 Cat{} 来说，这意味着 Quack 方法会接受一个全新的 Cat{}，因为方法的参数是 *Cat，编译器不会无中生有创建一个新的指针；即使编译器可以创建新指针，这个指针指向的也不是最初调用该方法的结构体；
*/
package main

import "fmt"

//使用值接收者的方法既可以通过值调用，也可以通过指针调用。
type S struct {
	data string
}

//hte focus of this example is this read function
//值接收者的方法，既可以用值调用，也可以用指针调用
func (s S) Read() string {
	return s.data
}

//只能通过指针调用
func (s *S) Write(str string) {
	s.data = str
}

func main() {
	sValues := map[int]S{1: {data: "A"}}
	//你只能通过值来调用read
	data := sValues[1].Read()
	fmt.Printf("data:%s\n", data)
	//不能编译通过
	//sValues[1].Write("b")
	//下面才能通过
	sptr := &S{data: "d"}
	sptr.Write("B")

	Sptrs := map[int]*S{1: &S{data: "B"}}
	Sptrs[1].Write("C")
	data = Sptrs[1].Read()
	fmt.Printf("data:%s\n", data)

	// 函数的类型检查
	//将 *RPCError 类型的变量赋值给 error 类型的变量 rpcErr；
	var rpcErr error = NewRPCError(400, "unknown err")
	//将 *RPCError 类型的变量 rpcErr 传递给签名中参数类型为 error 的 AsErr 函数；
	err := AsError(rpcErr)
	fmt.Println(err)
}

// 定义一个错误类型结构体
type RPCError struct {
	Code    int64
	Message string
}

//在 Go 中：实现接口的所有方法就隐式的实现了接口；RPCError实现了error接口，因为实现了error接口的Error方法
func (e *RPCError) Error() string {
	return fmt.Sprintf("%s,code=%d", e.Message, e.Code)
}

// NewRPCError Create a RPCError
func NewRPCError(code int64, msg string) error {
	//将 *RPCError 类型的变量从函数签名的返回值类型为 error 的 NewRPCError 函数中返回；
	return &RPCError{
		Code:    code,
		Message: msg,
	}
}

func AsError(err error) error {
	return err
}
