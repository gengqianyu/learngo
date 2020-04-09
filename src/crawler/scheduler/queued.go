package scheduler

import "crawler/engine"

type QueuedScheduler struct {
	// channel of engine send request message to scheduler
	requestChan chan engine.Request
	//scheduler 到 worker群 会有100个channel，我们把这100channel都灌进masterChan
	masterChan chan chan engine.Request
	// 定义两个队列
	requestQ    []engine.Request
	workerChanQ []chan engine.Request
}

func (s *QueuedScheduler) GetWorkerChan() chan engine.Request {
	return make(chan engine.Request)
}

func (s *QueuedScheduler) RequestReady(r engine.Request) {
	s.requestChan <- r
}

//如果有worker ready 就把这个worker的channel当作一个message，发送给masterChan
func (s *QueuedScheduler) WorkReady(workerChan chan engine.Request) {
	s.masterChan <- workerChan
}

//起一个服务
func (s *QueuedScheduler) Run() {
	s.requestChan = make(chan engine.Request)
	s.masterChan = make(chan chan engine.Request)
	go func() {
		for {
			//不要定义在for，这是利用nil channel 不被select到的特点，
			// 放到外面就第二次循环就不为nil channel 了
			var activeRequest engine.Request
			var activeWorkerChan chan engine.Request
			// 如果两个队列都有值，那么我们(scheduler)就挑一个workerChan把request发给worker
			if len(s.requestQ) > 0 && len(s.workerChanQ) > 0 {
				activeRequest = s.requestQ[0]
				activeWorkerChan = s.workerChanQ[0]
			}

			select {
			case r := <-s.requestChan: //从requestChan收一个消息，放入队列
				s.requestQ = append(s.requestQ, r)
			case w := <-s.masterChan: // 从masterChan收一个消息，放入队列
				s.workerChanQ = append(s.workerChanQ, w)
			case activeWorkerChan <- activeRequest: // 将request放入worker channel
				s.requestQ = s.requestQ[1:]
				s.workerChanQ = s.workerChanQ[1:]
			}
		}
	}()
}
