package scheduler

import "crawler/engine"

// 定义一个简单的调度器结构体
type Simple struct {
	WorkerChan chan engine.Request //scheduler to worker channel
}

func (s *Simple) WorkReady(chan engine.Request) {
}

func (s *Simple) RequestReady(request engine.Request) {
	//防止循环等待
	go func() {
		s.WorkerChan <- request
	}()
}

func (s *Simple) GetWorkerChan() chan engine.Request {
	return s.WorkerChan
}

func (s *Simple) Run() {
	s.WorkerChan = make(chan engine.Request)
}
