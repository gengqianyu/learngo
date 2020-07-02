//许多将unsafe.Pointer指针转为原生数字，然后再转回为unsafe.Pointer类型指针的操作也是不安全的。
//比如下面的例子需要将变量x的地址加上b字段地址偏移量转化为*int16类型指针，然后通过该指针更新x.b：

package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

func Float64bits(f float64) uint64 {
	fmt.Println(reflect.TypeOf(unsafe.Pointer(&f)))            //unsafe.Pointer
	fmt.Println(reflect.TypeOf((*uint64)(unsafe.Pointer(&f)))) //*uint64<br>　　　　　//(*uint64)(&f)  //这种类型转换语法是无效的
	return *(*uint64)(unsafe.Pointer(&f))                      //在变量前加*，是取值
	//代码  *(*uint64)(unsafe.Pointer(&f)) 我们进行分解：
	//1. unsafe.Pointer(&f) 转换变量指针f成Pointer类型
	//2. (*uint64)(unsafe.Pointer(&f))  将上述的Pointer类型转成uint6l的指针类型
	//3. 在2的结果(指针类型)前加上*,就是获取该指针的变量值
}

func main() {
	//一个普通的T类型指针可以被转化为unsafe.Pointer类型指针，并且一个unsafe.Pointer类型指针也可以被转回普通的指针，被转回普通的指针类型并不需要和原始的T类型相同。
	//说明：
	//unsafe.Pointer 转换的变量类型，一定是指针类型；
	//& 取址，* 取值；

	//unsafe.Pointer是特别定义的一种指针类型（译注：类似C语言中的void类型的指针），它可以包含任意类型变量的地址。
	//（1）任何类型的指针都可以被转化为unsafe.Pointer
	//（2）unsafe.Pointer可以被转化为任何类型的指针
	//（3）uintptr可以被转化为unsafe.Pointer
	//（4）unsafe.Pointer可以被转化为uintptr
	fmt.Printf("%#016x\n", Float64bits(1.0)) // "0x3ff0000000000000"

	v1 := uint(12)
	v2 := int(12)

	fmt.Println(reflect.TypeOf(v1)) //uint
	fmt.Println(reflect.TypeOf(v2)) //int

	fmt.Println(reflect.TypeOf(&v1)) //*uint
	fmt.Println(reflect.TypeOf(&v2)) //*int

	p := &v1

	//两个变量的类型不同,不能赋值
	//p = &v2 //报错，cannot use &v2 (type *int) as type *uint in assignment，类型不同，不可赋值
	//p = (*uint)(&v2)  //这种类型转换语法也是无效的
	fmt.Println(reflect.TypeOf(p)) // *unit

	//type uintptr uintptr
	//uintptr是golang的内置类型，是能存储指针的整型，在64位平台上底层的数据类型是，
	//
	//typedef unsigned long long int  uint64;
	//typedef uint64          uintptr;
	//
	//一个unsafe.Pointer指针也可以被转化为uintptr类型，然后保存到指针型数值变量中（注：这只是和当前指针相同的一个数字值，并不是一个指针），
	//然后用以做必要的指针数值运算。（uintptr是一个无符号的整型数，足以保存一个地址）**这种转换虽然也是可逆的，
	//但是将uintptr转为unsafe.Pointer指针可能会破坏类型系统，因为并不是所有的数字都是有效的内存地址。
	var x struct {
		a bool
		b int16
		c []int
	}

	/**
	  unsafe.Offsetof 函数的参数必须是一个字段 x.f, 然后返回 f 字段相对于 x 起始地址的偏移量, 包括可能的空洞.
	*/

	//  指针的运算
	//  uintptr(unsafe.Pointer(&x)) + unsafe.Offsetof(x.b)
	// 和 pb := &x.b 等价 很像操作系统的页内寻址，x的起始地址+偏移地址 就是 x.b的地址位置
	pb := (*int16)(unsafe.Pointer(uintptr(unsafe.Pointer(&x)) + unsafe.Offsetof(x.b)))
	*pb = 42
	fmt.Println(x.b) // "42"
}

//上面的写法尽管很繁琐，但在这里并不是一件坏事，因为这些功能应该很谨慎地使用。
//不要试图引入一个uintptr类型的临时变量，因为它可能会破坏代码的安全性（注：这是真正可以体会unsafe包为何不安全的例子）。

// NOTE: subtly incorrect! 下面这段代码就是错误的
//tmp := uintptr(unsafe.Pointer(&x)) + unsafe.Offsetof(x.b)
//pb := (*int16)(unsafe.Pointer(tmp))
//*pb = 42
//产生错误的原因很微妙。
//**有时候垃圾回收器会移动一些变量以降低内存碎片等问题。这类垃圾回收器被称为移动GC。
//当一个变量被移动，所有的保存改变量旧地址的指针必须同时被更新为变量移动后的新地址。
//从垃圾收集器的视角来看，一个unsafe.Pointer是一个指向变量的指针，因此当变量被移动是对应的指针也必须被更新；
//但是uintptr类型的临时变量只是一个普通的数字，所以其值不应该被改变。
//上面错误的代码因为引入一个非指针的临时变量tmp，导致垃圾收集器无法正确识别这个是一个指向变量x的指针。
//当第二个语句执行时，变量x可能已经被转移，这时候临时变量tmp也就不再是现在的&x.b地址。
//**第三个向之前无效地址空间的赋值语句将彻底摧毁整个程序！
