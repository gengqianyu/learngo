package main

import (
	"fmt"
	"reflect"
)

func main() {
	author := "draven"
	//可以将一个普通的变量转换成『反射』包中提供的 Type 和 Value，随后就可以使用反射包中的方法对它们进行复杂的操作。
	//reflect.TypeOf 和 reflect.ValueOf 能够获取 Go 语言中的变量对应的反射对象。
	//一旦获取了反射对象，我们就能得到跟当前类型相关数据和操作，并可以使用这些运行时获取的结构执行方法。
	//获取变量反射对象的类型
	//从interface{}变量可以反射出反射对象，TypeOf和ValueOf入参是interface{}，反射出反射对象，入参时会做类型转换，string类型to interface{}
	//因为类型转换这一过程一般都是隐式的，所以我不太需要关心它，只有在我们需要将反射对象转换回基本类型时才需要显式的转换操作。
	reflectObjType := reflect.TypeOf(author)
	fmt.Println(reflectObjType)
	//获取变量反射对象的值
	fmt.Println(reflect.ValueOf(author))

	// 从反射对象转换回基本类型
	//获取反射对象
	a := 1
	reflectObjValue := reflect.ValueOf(a)
	//从反射对象转到接口,接口也是有类型的
	a = reflectObjValue.Interface().(int)
	fmt.Printf("%T,%d\n", a, a)

	i := 1
	//调用 reflect.ValueOf 函数获取变量指针；
	v := reflect.ValueOf(&i)
	//调用 reflect.Value.Elem 方法获取指针指向的变量；
	//用 reflect.Value.SetInt 方法更新变量的值：
	v.Elem().SetInt(10)
	fmt.Println(i)
}
