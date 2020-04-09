package main

import (
	"crawler_distributed/rpcsupport"
	"crawler_distributed/worker"
	"flag"
	"fmt"
	"log"
)

// 定义命令行参数 port
var port = flag.Int("port", 0, "the port for server of worker to listen on(example --port 9000)")

func main() {
	//解析命令行参数。
	//必须在定义所有参数之后并且在程序访问参数之前调用。
	flag.Parse()
	// *读取一下指针的val
	if *port == 0 {
		fmt.Printf("must specify a port")
		return
	}
	//启动rpc worker server 有err就log fatal 写一个致命log，然后退出
	log.Fatal(rpcsupport.RunServer(&worker.CrawlService{}, fmt.Sprintf(":%d", *port)))
}
