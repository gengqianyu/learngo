package main

import (
	"crawler/frontend/controller"
	"flag"
	"fmt"
	"net/http"
)

// 定义一个命令行参数
var port = flag.String("port", "", `port for http server (example --port ":8888")`)

func main() {
	//解析命令行参数
	flag.Parse()
	if *port == "" {
		fmt.Printf("must specify a port")
		return
	}
	// http 文件系统
	http.Handle("/", http.FileServer(http.Dir("view/")))

	// http 处理host/search 的请求
	http.Handle("/search", controller.FactorySearchHandler("view/template.html"))
	// 启动http服务并监听8888端口
	err := http.ListenAndServe(*port, nil)
	if err != nil {
		panic(err)
	}
}
