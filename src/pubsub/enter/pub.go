package enter

import (
	"sync"
	"time"
)

//发布订阅（publish-and-subscribe）模型通常被简写为 pub/sub 模型。
//在这个模型中，消息生产者成为发布者（publisher），而消息消费者则成为订阅者（subscriber）,生产者和消费者是 M:N 的关系。
//在传统生产者和消费者模型中，是将消息发送到一个队列中，而发布订阅模型则是将消息发布给一个主题。

//在发布订阅模型中，每条消息都会传送给多个订阅者。
//发布者不会知道、也不关心哪一个订阅者正在接受主题消息。
//订阅者和发布者可以在运行时动态添加，是一种松散的耦合关系，这使得系统的复杂性可以随时间的推移而增长。
//在现实生活中，像天气预报之类的应用就可以应用这个并发模式。

type (
	subscriber chan interface{}         //订阅者为一个 interface类型的管道
	topicFunc  func(v interface{}) bool //主题为一个过滤器
)

//发布者struct
type Publisher struct {
	m           sync.RWMutex             //读写锁
	buffer      int                      //订阅队列的缓存大小
	timeout     time.Duration            //发布超时时间
	subscribers map[subscriber]topicFunc //所以订阅者信息
}

//设置发布超时时间和缓存队列的长度
func (p *Publisher) init(t time.Duration, buffer int) *Publisher {
	p.timeout = t
	p.buffer = buffer
	p.subscribers = make(map[subscriber]topicFunc)
	return p
}

//构建一个发布者对象
func New(t time.Duration, buffer int) *Publisher {
	return new(Publisher).init(t, buffer)
}

//添加一个新的订阅者，订阅过滤筛选后的主题
func (p *Publisher) SubscribeTopic(tf topicFunc) subscriber {
	out := make(chan interface{}, p.buffer)

	//开一个goroutine，用于注册订阅者方法，新增一个订阅
	//go func(tf topicFunc) {
	//上锁的原因是为了不要超过订阅者的消息个数
	p.m.Lock()              //拿锁
	defer p.m.Unlock()      //释放锁
	p.subscribers[out] = tf //注册订阅者方法，新增一个订阅
	//}(tf)

	return out
}

//订阅全部主题，没有过滤
func (p *Publisher) Subscribe() subscriber {
	return p.SubscribeTopic(nil)
}

//退出订阅
func (p *Publisher) Evict(s subscriber) {
	p.m.Lock()
	defer p.m.Unlock()

	delete(p.subscribers, s) //删除订阅者
	close(s)                 //关闭通道
}

//发布一个主题
func (p *Publisher) Publish(v interface{}) {
	p.m.RLock()         //拿读锁
	defer p.m.RUnlock() //释放锁

	var wg sync.WaitGroup

	for s, t := range p.subscribers {
		//发布主题，可以容忍一定的超时
		go func(s subscriber, tf topicFunc, v interface{}, wg *sync.WaitGroup) {
			wg.Add(1)       //添加一个等待
			defer wg.Done() //完成一个goroutine
			//没有订阅该主题
			if tf != nil && !tf(v) {
				return
			}

			select {

			case s <- v:

			case <-time.After(p.timeout):

			}
		}(s, t, v, &wg)

	}
	wg.Wait()
}

//关闭发布者对象，同时关闭所有的订阅者管道
func (p *Publisher) Close() {
	p.m.Lock()
	defer p.m.Unlock()
	//上锁为为了避免操作时，继续有订阅者注册
	for s := range p.subscribers {
		delete(p.subscribers, s)
		close(s)
	}
}
