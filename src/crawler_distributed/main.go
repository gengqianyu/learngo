/*
get package text and html
gopm get -g -v golang.org/x/text
gopm get -g -v golang.org/x/net/html

获取城市名称和连接
1.使用css选择器
2.使用xpath
3.使用正则表达式：本课程仅使用正则表达式
*/
package main

import (
	"crawler/engine"
	"crawler/scheduler"
	"crawler/zhenai/parser"
	"crawler_distributed/config"
	itemClient "crawler_distributed/persist/client"
	workerClient "crawler_distributed/worker/client"
	"flag"
	"fmt"
	"strings"
)

// 定义命令行参数
var (
	itemServerHost    = flag.String("item_server_host", "", `connection to  host of itemServer (example --item_server_host ":1234")`)
	workerServerHosts = flag.String("worker_server_hosts", "", `used to create worker connection pool (example --worker_server_hosts ":9000,:9001,:9002,:9003")`)
)

func main() {
	//解析命令行参数
	flag.Parse()

	if *itemServerHost == "" {
		fmt.Printf("must specify item server host")
		return
	}
	if *workerServerHosts == "" {
		fmt.Printf("must specify worker server hosts")
		return
	}
	//单worker版
	//engine.Simple{}.Run(engine.Request{
	//	Url:    "http://www.zhenai.com/zhenghun",
	//	Parser: parser.CityList,
	//})

	const index = "dating_profile"
	// 多worker版
	// 耦合版
	//itemChan, err := persist.ItemServer(index)
	// 分布式版
	itemChan, err := itemClient.ItemServer(*itemServerHost)
	if err != nil {
		panic(err)
	}

	//create worker connection pool
	hosts := strings.Split(*workerServerHosts, ",")
	workerPool := engine.CreateWorkerConnPool(hosts)
	// 创建分布式版 worker handler client
	workerHandler, err := workerClient.CreateWorkerHandler(workerPool)
	if err != nil {
		panic(err)
	}

	e := engine.Concurrent{
		//Scheduler: &scheduler.Simple{},
		Scheduler: &scheduler.QueuedScheduler{},
		WorkerNum: 20,
		ItemChan:  itemChan,
		//WorkerHandler: engine.Worker,//单机版
		WorkerHandler: workerHandler,
	}
	e.Run(engine.Request{
		Url:    "http://www.zhenai.com/zhenghun",
		Parser: engine.NewFuncParser(parser.CityList, config.ParserCityList),
	})

}
