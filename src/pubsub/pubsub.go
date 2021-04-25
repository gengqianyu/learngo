package pubsub

//发布订阅（publish-and-subscribe）模型通常被简写为 pub/sub 模型。
//在这个模型中，消息生产者成为发布者（publisher），而消息消费者则成为订阅者（subscriber）,生产者和消费者是 M:N 的关系。
//在传统生产者和消费者模型中，是将消息发送到一个队列中，而发布订阅模型则是将消息发布给一个主题。

//在发布订阅模型中，每条消息都会传送给多个订阅者。
//发布者不会知道、也不关心哪一个订阅者正在接受主题消息。
//订阅者和发布者可以在运行时动态添加，是一种松散的耦合关系，这使得系统的复杂性可以随时间的推移而增长。
//在现实生活中，像天气预报之类的应用就可以应用这个并发模式。

//发布者
type Publisher struct {
}

func (p *Publisher) init() *Publisher {
	return p
}

//构建一个发布者实例
func New() *Publisher {
	return new(Publisher).init()
}

//发布者发布一个主题
func (p *Publisher) Publish(topic interface{}) {

}

//交换机
type Exchange struct {
	subscribers map[string][]subscriber
}

//订阅者
type subscriber chan interface{}
