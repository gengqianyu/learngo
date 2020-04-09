package engine

import (
	"log"
)

type Simple struct {
}

// 运行函数 参数seeds 可以是很多个request
func (s *Simple) Run(seeds ...Request) {
	// 定义一个请求队列
	var requests []Request
	//  将请求种子放入队列
	for _, seed := range seeds {
		requests = append(requests, seed)
	}
	// 如果请求队列不为空就执行请求
	for len(requests) > 0 {
		request := requests[0]
		requests = requests[1:]
		parserResult, err := Worker(request)
		if err != nil {
			continue
		}
		// 将requests切片展开，放进队列
		requests = append(requests, parserResult.Requests...)

		for _, item := range parserResult.Items {
			log.Printf("Got item %v", item)
		}

	}
}
