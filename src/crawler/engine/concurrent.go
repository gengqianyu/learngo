/*
并发版 concurrent
*/
package engine

// 定义一个并发结构体
type Concurrent struct {
	Scheduler     Scheduler     //调度器
	WorkerNum     int           //worker 的数量
	ItemChan      chan Item     //将parser出的item 放进此通道 让item server save
	WorkerHandler WorkerHandler //do worker
}

type WorkerHandler func(Request) (ParserResult, error)

// 定义一个调度器接口
type Scheduler interface {
	ReadyNotifier
	RequestReady(Request)
	GetWorkerChan() chan Request
	Run()
}

// 此接口和上方接口是一个组合关系
type ReadyNotifier interface {
	WorkReady(chan Request)
}

func (engine *Concurrent) Run(seeds ...Request) {
	// create channel out 用于 worker send ParserResult to engine
	out := make(chan ParserResult)
	// 启动scheduler服务
	engine.Scheduler.Run()

	// 创建多个worker去做事情
	for i := 0; i < engine.WorkerNum; i++ {
		engine.CreateWorker(engine.Scheduler.GetWorkerChan(), out, engine.Scheduler)
	}

	// 调度器向MasterWorkerChan提交请求
	for _, seed := range seeds {
		engine.Scheduler.RequestReady(seed)
	}
	//起一个服务，只做分发不做数据处理
	for {
		// engine接收一个ParserResult message, worker 通过 out 通道发过来的
		parserResult := <-out
		// 处理item
		for _, item := range parserResult.Items {
			//开一个goroutine 把item发往itemChan通道里，让itemServer 去消费
			go func() {
				engine.ItemChan <- item
			}()
		}
		// 处理request
		for _, request := range parserResult.Requests {
			engine.Scheduler.RequestReady(request)
		}
	}

}

// 创建一个worker
func (engine *Concurrent) CreateWorker(in chan Request, out chan ParserResult, notifier ReadyNotifier) {
	go func() {
		for {
			// 通知scheduler worker channel ready
			notifier.WorkReady(in)
			// received request message ,scheduler通过in通道发来的
			request := <-in
			// do worker
			// stand-alone 打击并发版
			parserResult, err := engine.WorkerHandler(request)
			if err != nil {
				continue
			}
			// 将parserResult发送进out通道中 让engine去接收
			out <- parserResult
		}
	}()
}

var visitedUrls = make(map[string]bool)

// 判断是否访问过
func isDuplicate(url string) bool {
	if visitedUrls[url] {
		return true
	}
	visitedUrls[url] = true
	return false
}
