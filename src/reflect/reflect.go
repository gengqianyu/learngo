package main

import (
	"fmt"
	"reflect"
)

func main() {
	author := order{1, 3}
	//可以将一个普通的变量转换成『反射』包中提供的 Type 和 Value，随后就可以使用反射包中的方法对它们进行复杂的操作。
	//reflect.TypeOf 和 reflect.ValueOf 能够获取 Go 语言中的变量对应的反射对象。
	//一旦获取了反射对象，我们就能得到跟当前类型相关数据和操作，并可以使用这些运行时获取的结构执行方法。

	//从interface{}变量可以反射出反射对象，TypeOf和ValueOf入参是interface{}，反射出反射对象，入参时会做类型转换，string类型to interface{}
	//因为类型转换这一过程一般都是隐式的，所以我不太需要关心它，只有在我们需要将反射对象转换回基本类型时才需要显式的转换操作。

	//获取变量反射对象的类型 reflect.Type
	reflectType := reflect.TypeOf(author)
	fmt.Println("type of author:", reflectType)
	//获取变量反射对象的值 reflect.Value
	reflectValue := reflect.ValueOf(author)
	fmt.Println("value of author", reflectValue)
	//反射类型、反射值的具体种类 kind函数reflect.Type和reflect.Value都能调用
	typeKind := reflectType.Kind()
	fmt.Println("具体反射类型(reflect.Type)的种类：", typeKind)
	valueKind := reflectValue.Kind()
	fmt.Println("具体反射值(reflect.Value)的种类：", valueKind)
	//获取反射对象的值(reflect.Value)结构体属性字段的数量 此方法仅在结构体类型中生效 不是struct就会panic
	numFieldReflectValue := reflectValue.NumField()
	fmt.Println("反射对象值结构体属性的数量：", numFieldReflectValue)

	//获取反射对象的类型(reflect.Type)结构体属性字段的数量 此方法仅在结构体类型中生效 不是struct就会panic
	numFieldReflectType := reflectType.NumField()
	fmt.Println("反射对象类型 结构的属性的数量：", numFieldReflectType)

	//reflect.Type.Name()获取反射对象类型名称
	//reflect.Type.Field(i).Name 获取反射对象类型，结构体属性字段名称
	//reflect.Value.Field(i) 获取反射对象值的结构体属性字段值可以用reflect.Value.Field(i).Int() reflect.Value.Field(i).String() 将值转成int或string
	for i := 0; i < numFieldReflectValue; i++ {
		fmt.Println(reflectType.Name(), reflectType.Kind(), " field:", reflectType.Field(i).Name, ":", reflectValue.Field(i).Int())
	}

	i := 123
	//用 reflect.Value.SetInt 方法更新变量的值：
	//获取指针变量反射对象的值；
	pinterReflectValue := reflect.ValueOf(&i)
	//通过指针变量反射对象的值，获取变量的反射对象；reflect.Value.Elem
	ReflectValue := pinterReflectValue.Elem()

	ReflectValue.SetInt(100)
	fmt.Println(ReflectValue)

	// 从反射对象转换回基本类型
	//获取反射对象
	a := 1
	reflectObjValue := reflect.ValueOf(a)
	//从反射对象转到接口,接口也是有类型的
	a = reflectObjValue.Interface().(int)
	fmt.Printf("%T,%d\n", a, a)

	//反射综合实现
	o := order{
		ordId:      456,
		customerId: 56,
	}
	createQuery(o)

	e := employee{
		name:    "Naveen",
		id:      565,
		address: "Coimbatore",
		salary:  90000,
		country: "India",
	}
	createQuery(e)
	i = 90
	createQuery(i)
}

//订单struct
type order struct {
	ordId      int
	customerId int
}

type employee struct {
	name    string
	id      int
	address string
	salary  int
	country string
}

func createQuery(q interface{}) {

	if reflect.ValueOf(q).Kind() == reflect.Struct {
		t := reflect.TypeOf(q)
		n := t.Name()
		query := fmt.Sprintf("insert into %s (", n)

		for i := 0; i < t.NumField(); i++ {
			if i == 0 {
				query = fmt.Sprintf("%s%s", query, t.Field(i).Name)
			} else {
				query = fmt.Sprintf("%s,%s", query, t.Field(i).Name)
			}

		}
		query = fmt.Sprintf("%s)values(", query)

		v := reflect.ValueOf(q)
		for i := 0; i < v.NumField(); i++ {
			switch v.Field(i).Kind() {
			case reflect.Int:
				if i == 0 {
					query = fmt.Sprintf("%s%d", query, v.Field(i).Int())
				} else {
					query = fmt.Sprintf("%s,%d", query, v.Field(i).Int())
				}
			case reflect.String:
				if i == 0 {
					query = fmt.Sprintf("%s\"%s\"", query, v.Field(i).String())
				} else {
					query = fmt.Sprintf("%s, \"%s\"", query, v.Field(i).String())
				}
			default:
				fmt.Println("Unsupported type")
				return
			}
		}
		query = fmt.Sprintf("%s)", query)
		fmt.Println(query)
		return

	}
	fmt.Println("unsupported type")
}
