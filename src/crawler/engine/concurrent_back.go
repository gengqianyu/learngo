/*
并发版 concurrent
*/
package engine

//import (
//	"log"
//)
//
//// 定义一个并发结构体
//type Concurrent struct {
//	Scheduler Scheduler //调度器
//	WorkerNum int       //worker 的数量
//}
//
//// 定义一个调度器接口
//type Scheduler interface {
//	Submit(Request)
//	SetMasterWorkerChan(chan Request)
//}
//
//func (engine *Concurrent) Run(seeds ...Request) {
//	//定义worker的输入和输出
//	// create channel in 用于received Request from Scheduler send
//	in := make(chan Request)
//	// create channel out 用于 worker send ParserResult to engine
//	out := make(chan ParserResult)
//	// 设置调度器向worker发送request message的通道
//	engine.Scheduler.SetMasterWorkerChan(in)
//
//	// 创建多个worker去做事情
//	for i := 0; i < engine.WorkerNum; i++ {
//		CreateWorker(in, out)
//	}
//
//	// 调度器向MasterWorkerChan提交请求
//	for _, seed := range seeds {
//		engine.Scheduler.Submit(seed)
//	}
//	itemNum := 0
//	//起一个服务
//	for {
//		// engine接收一个ParserResult message, worker 通过 out 通道发过来的
//		parserResult := <-out
//		// 打印解析过的item
//		for _, item := range parserResult.Items {
//			log.Printf("Got item num:%d,val:%v", itemNum, item)
//			itemNum++
//		}
//		// 将Request 提交给调度器
//		// 处理city的时候，会有400多个request被submit 到in通道内
//		// 前10个request被10个worker从in通道内接收并处理，
//		// 此时所有worker都去网络fetching，事都还没完。
//		// 这时候再去submit 10个以后的request，这时候in通道的request却没有worker去接
//		// 所以会造成阻塞，（循环等待）submit函数添加request的时候添加不进去等在哪里
//		// 最简单的解决办法是，submit逻辑放入goroutine中执行，
//		// 这样submit开好goroutine后就立即返回执行下次添加，不会阻塞了
//		for _, request := range parserResult.Requests {
//			engine.Scheduler.Submit(request)
//		}
//	}
//
//}
//
//// 创建一个worker
//func CreateWorker(in chan Request, out chan ParserResult) {
//	go func() {
//		for {
//			// received request message ,scheduler通过in通道发来的
//			request := <-in
//			// do worker
//			parserResult, err := worker(request)
//			if err != nil {
//				continue
//			}
//			// 将parserResult发送进out通道中 让engine去接收
//			out <- parserResult
//		}
//	}()
//}
