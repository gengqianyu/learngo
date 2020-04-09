package client

import (
	"crawler/engine"
	"crawler_distributed/config"
	"crawler_distributed/worker"
	"net/rpc"
)

func CreateWorkerHandler(clientChan chan *rpc.Client) (engine.WorkerHandler, error) {

	//client, err := rpcsupport.FactoryClient(fmt.Sprintf(":%d", config.RpcWorkerServerPort))
	//if err != nil {
	//	return nil, err
	//}

	// 注意下面这个坑 如果不写在return 函数内，就是现场执行。
	// 这会导致 return workerHandler 永远只用第一个rpc client
	// received an rpc client from client channel
	//rpcClient := <-clientChan
	return func(request engine.Request) (engine.ParserResult, error) {
		// received an rpc client from client channel
		rpcClient := <-clientChan

		workerRequest := worker.SerializeRequest(request)
		var result worker.ParserResult
		err := rpcClient.Call(config.WorkerServiceMethod, workerRequest, &result)
		if err != nil {
			return engine.ParserResult{}, err
		}
		return worker.DeserializeParserResult(result), nil
	}, nil
}
