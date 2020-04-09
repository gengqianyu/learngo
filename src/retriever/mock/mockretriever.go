package mock

import "fmt"

// 定义retriever结构体
type Retriever struct {
	Contents string
}

// 实现一个系统接口方法 格式化一个结构，将一个结构体当成一个字符串输出 类似php的__toString 魔术方法
func (r Retriever) String() string {
	//Sprintf打印的结果不会被输出，二是被返回字符串。可以把结果赋值变量
	return fmt.Sprintf("Retriever:{Contents=%s}", r.Contents)
}

// 实现retriever接口的方法Get
// go里面很神奇的地方，本文件内（包内）并没有出现retriever接口的引用
// 但是只要实现Get方法，那么就认为retriever（结构）实现了retriever（接口）
func (r Retriever) Get(string) string {
	return r.Contents
}

//再实现一个poster接口 因为要修改结构体的contents 所以要使用指针接收者
func (r *Retriever) Post(url string, form map[string]string) string {
	r.Contents = form["contents"]
	return "ok"
}
