package main

import (
	"crawler_distributed/config"
	"crawler_distributed/rpcsupport"
	"crawler_distributed/worker"
	"fmt"
	"testing"
	"time"
)

func TestCrawlService(t *testing.T) {
	// 端口别错了
	const host = 9000
	go rpcsupport.RunServer(&worker.CrawlService{}, fmt.Sprintf(":%d", host))
	time.Sleep(time.Second)

	client, err := rpcsupport.FactoryClient(fmt.Sprintf(":%d", host))

	if err != nil {
		panic(err)
	}
	var workerParserResult worker.ParserResult
	err = client.Call(config.WorkerServiceMethod, worker.Request{
		Url: "http://album.zhenai.com/u/1630078262",
		Parser: worker.SerializedParser{
			Name: config.ParserProfile,
			Args: "心語",
		},
	}, &workerParserResult)

	if err != nil {
		t.Errorf("rpc call err:%v", err)
	} else {
		fmt.Println(workerParserResult)
	}

}
