package engine

import (
	"crawler/fetcher"
	"crawler_distributed/rpcsupport"
	"log"
	"net/rpc"
)

func Worker(request Request) (ParserResult, error) {
	//提取网页内容
	body, err := fetcher.Fetch(request.Url)
	if err != nil {
		log.Printf("Fetcher: error fetching url:%s:%v", request.Url, err)
		return ParserResult{}, err
	}
	//解析价值信息
	return request.Parser.Parse(body, request.Url), nil
}

//crate worker connection pool
func CreateWorkerConnPool(hosts []string) chan *rpc.Client {
	var clients []*rpc.Client
	for _, host := range hosts {
		client, err := rpcsupport.FactoryClient(host)
		if err != nil {
			log.Printf("error connecting to %s,err:%v", host, err)
			continue
		}
		clients = append(clients, client)
		log.Printf("connected to %s", host)
	}
	//创建一个channel
	pool := make(chan *rpc.Client)
	// 开一个goroutine
	go func() {
		//要不停的往里添加 rpc worker client ,
		// workerHandler 才能不停的去接收，不然workerHandler拿到4个就没有了
		for {
			for _, workerClient := range clients {
				pool <- workerClient
			}
		}
	}()
	// 立即返回
	return pool
}
