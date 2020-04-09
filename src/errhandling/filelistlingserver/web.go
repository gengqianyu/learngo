/*
如何实现统一的错误处理逻辑
整体思想：
定义一个包装器，将一个业务函数作为参数输入，执行业务逻辑，返回错误
包装器函数处理返回的err，将业务逻辑和错误分开。输出包装以后的type

error vs panic
意料之中的使用error，如打不开文件
意料之外的使用panic，如数组越界

//使用net/http/pprof 工具查看/分析服务性能
// 查看dug ：interview localhost:8888/debug/pprof
// 使用工具 tool ：go tool pprof http://localhost:8888/debug/pprof/profile
*/
package main

import (
	"errhandling/filelistlingserver/filelisting"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
)

// 定义 appHandler 类型 用业务逻辑处理 如果出错返回一个error
type appHandler func(writer http.ResponseWriter, request *http.Request) error

// 定义http.HandleFunc中的第二个参数handler的类型
type httpHandler func(http.ResponseWriter, *http.Request)

// 包装器
// 包装函数只是 将文件处理handler->filelisting.HandleFileList函数体作为参数，输入进去，执行一下，包装成httpHandler再输出出来
// 文件处理handler->filelisting.HandleFileList只管做正确的事情,如果有错误就返回来，让这个包装器去做处理
// 函数式编程，errWrapper 参数为一个文件handler函数，经过业务逻辑，包装一个httpHandler函数输出
func errWrapper(handler appHandler) httpHandler {
	return func(writer http.ResponseWriter, request *http.Request) {
		// defer 保护 用于处理未知panic错误（也就是说不要把程序的panic错误暴露）
		defer func() {
			// 如果recover不等于nil 则去处理错误
			if r := recover(); r != nil {
				log.Printf("panic:%v", r)
				http.Error(
					writer,
					http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
		}()
		//panic(123) //测试以上defer专用 打开就可以

		//handler去做业务逻辑，出来的错误在这里统一处理
		err := handler(writer, request)
		if err != nil {
			// 服务端打印log
			log.Printf("error occurred handling request:%s", err.Error())
			// 处理用户错误 user error
			if userErr, ok := err.(userError); ok { // 使用了一个类型断言 type assertion
				// 服务端打印 user error
				log.Printf("user error:%s", userErr.Message())
				// 往response中写一个user error
				http.Error(writer, userErr.Message(), http.StatusBadRequest)
				return
			}
			// 处理系统错误 system error
			code := http.StatusOK
			switch {
			case os.IsNotExist(err): //未找到
				code = http.StatusNotFound
			case os.IsPermission(err): //没权限
				code = http.StatusForbidden
			default:
				code = http.StatusInternalServerError
			}
			http.Error(
				writer,
				http.StatusText(code),
				code)
			return
		}
	}
}

// 定义一个用户错误接口约束
type userError interface {
	error
	Message() string
}

func main() {

	handler := errWrapper(filelisting.HandleFileList)
	// 定义一个http处理函数
	http.HandleFunc("/", handler)

	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		panic(err)
	}

}
