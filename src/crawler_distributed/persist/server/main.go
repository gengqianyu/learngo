package main

import (
	"crawler_distributed/config"
	"crawler_distributed/persist"
	"crawler_distributed/rpcsupport"
	"flag"
	"fmt"
	"log"

	"github.com/olivere/elastic"
)

var port = flag.Int("port", 0, "the port for server of item to listen on (example --port 1234)")

func main() {
	flag.Parse()
	if *port == 0 {
		fmt.Printf("must specify a port")
		return
	}
	// 如果服务有问题直接退出
	log.Fatal(StartUpRpcServer(fmt.Sprintf(":%d", *port), config.ElasticIndex))
}

func StartUpRpcServer(host, index string) error {
	// new an elastic search client
	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		return err
	}
	// init rpc Item service pointer 作为rpc服务的注册接收者
	receiver := &persist.ItemService{
		Client: client,
		Index:  index,
	}
	// start up rpc server
	return rpcsupport.RunServer(receiver, host)
}
