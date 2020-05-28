/*
duck typing
像鸭子走路，像鸭子叫，长得像鸭子，那么就是鸭子。
描述事物的外部行为，而非内部结构
严格来说go属于结构化类型系统，类似duck typing

go 语言的duck typing
同时需要readable and appendable 两个interface怎么办？
同时具有python 和c++ duck typing的灵活性
  python运行时才检查有没有实现get方法
  c++是在编译的时候检查，编译时才知道传入的retriever(实现者)有没有get方法
  这样带来一个麻烦就是写代码的时候需要去注释说明接口（实现者）实现了那些方法，才能送给接收者去使用
又具有java的类型检查(java是在写程序的时候做约束 检查是否实现了接口的get方法)
  java的问题在于不能同时实现多个接口
还能具有php能实现多个接口的灵活性
  php的问题在于，无法实现多接口的约束()。

接口的定义
使用者(download)->实现者(retriever)
接口是由使用者定义的，使用者想要什么的样的约束，就去定义什么样的接口，可以单接口约束，也可以组合多接口约束
实现者并没有说我实现了接口，只是定义了一个和接口相同的方法

接口的值类型：
接口变量里有什么？
接口变量也是相当于一个结构
第一种：接口变量包含（实现者的类型，实现者的值）
第二种：接口变量包含（实现者的类型，实现者的指针->[实现者]），实现者指针指向一个实现者
接口变量的典型用法
接口变量自带指针（有时候也可以带值，通常情况下自己带指针的）
接口变量同样采用值传递，但是因为接口变量肚子里有个实现者的指针，所以几乎不需要使用接口的指针
指针接收者实现只能以指针方式使用；值接收者都可以
查看接口变量：
			表示任何类型：interface{}
			type assertion
			type switch

*/
package main

import (
	"fmt"
	"retriever/mock"
	"retriever/real"
	"time"
)

//定义接口
// 由使用者来规定，我这个retriever必须有Get方法
// *也就是说，接口由使用者来定义的
type Retriever interface {
	Get(url string) string
}

// 定义一个poster接口
type Poster interface {
	Post(url string, form map[string]string) string
}

const url = "http://www.imooc.com"

//定义使用者方法
// 此处遵循控制反转，依赖倒置原则 参数r类型约束，为Retriever接口类型，并非具体结构体
func download(r Retriever) string {
	return r.Get(url)
}

//定义一个post方法
func post(poster Poster) {
	poster.Post(url, map[string]string{
		"name":   "ccmouse",
		"course": "golang",
	})
}

//组合接口
type RetrieverPoster interface {
	Retriever
	Poster
}

//定义一个session方法
func session(s RetrieverPoster) string {
	//调用post修改结构体contents的值
	s.Post(url, map[string]string{
		"contents": "another faked imooc.com",
	})
	return s.Get(url)
}

// 使用者
func main() {
	//定义一个r接口变量，约束是一个Retriever接口类型
	var r Retriever
	// 生成retriever结构体对象并调用download方法
	// 使用值调用
	r = mock.Retriever{"this is a pake www.baidu.com"}
	// 实现的接口变量r肚子里有两个东西，（接口实现者的类型(mock.Retriever)，接口实现者的值(mock.Retriever整个结构体)）
	inspect(r)
	fmt.Println(download(r))
	// 使用指针调用
	r = &mock.Retriever{"this is a pointer"}
	// type switch方式查看接口变量
	inspect(r)
	fmt.Println(download(r))

	//使用指针调用 因为realRetriever 实现Get方法使用的是指针接收者
	r = &real.Retriever{
		UserAgent: "Mozilla/5.0",
		TimeOut:   time.Minute,
	}
	// 实现的 接口变量r 肚子里有两个东西，（接口实现者的类型(*real.Retriever)，接口实现者的指针(&real.Retriever)）
	inspect(r)
	// 实现的接口变量r肚子里有两个东西，一个是接口实现者的类型，一个是接口实现者的值/指针

	// type assertion 类型断言的方式查看接口变量
	//使用 接口变量 点一个接口类型 来拿出接口的值
	//严格版本
	typeRetriever := r.(*real.Retriever)
	fmt.Println("timeout:", typeRetriever.TimeOut)
	// 点一个错误的接口类型
	// 不严格版本加一个ok
	if mockRetriever, ok := r.(mock.Retriever); ok {
		fmt.Println(mockRetriever.Contents)
	} else {
		fmt.Println("not a mock retriever")
	}
	//try an download

	//fmt.Println(download(r))

	// try an session
	// 获取初始化mockRetriever结构体的指针s
	s := &mock.Retriever{"this is a pake www.baidu.com"}
	fmt.Println(session(s))
}

// 检查类型
func inspect(r Retriever) {
	fmt.Printf("r type: %T->value: %v\n", r, r)
	// 获取r的类型
	switch v := r.(type) {
	case mock.Retriever:
		fmt.Println("contents:", v.Contents)
	case *mock.Retriever:
		fmt.Println("contents:", v.Contents)
	case *real.Retriever:
		fmt.Println("userAgent:", v.UserAgent)
	}
	fmt.Println()
}
