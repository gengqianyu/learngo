/*
iface and eface
Go 语言使用 iface 结构体表示第一种接口，使用 eface 结构体表示第二种空接口，两种接口虽然都使用 interface 声明，
但是由于后者在 Go 语言中非常常见，所以在实现时使用了特殊的类型。

需要注意的是，与 C 语言中的 void * 不同，interface{} 类型不是任意类型，
如果我们将类型转换成了 interface{} 类型，这边变量在运行期间的类型也发生了变化，获取变量类型时就会得到 interface{}。
*/
package main

import "fmt"

func main() {
	type test struct{}
	v := test{}
	Print(v)
}

//上述函数不接受任意类型的参数，只接受 interface{} 类型的值，
//在调用 Print 函数时会对参数 v 进行类型转换，将原来的 Test 类型转换成 interface{} 类型。
func Print(i interface{}) {
	fmt.Println(i)
}
