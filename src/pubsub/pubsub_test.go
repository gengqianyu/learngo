package pubsub

import (
	"os"
	"os/signal"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestPublisher_Publish(t *testing.T) {

	pub := New(100*time.Millisecond, 10) //实例化发布者
	defer pub.Close()                    //关闭发布者对象，同时关闭所有订阅者管道

	//添加一个新的订阅者，订阅全部主题，没有过滤
	all := pub.Subscribe()

	//添加一个新的订阅者，订阅过滤筛选后的主题
	golang := pub.SubscribeTopic(func(v interface{}) bool {
		if s, ok := v.(string); ok {
			return strings.Contains(s, "golang")
		}
		return false
	})

	//发布两个主题
	pub.Publish("hello world!")
	pub.Publish("hello golang!")

	go func() {
		//接受对应订阅的主题
		for msg := range all {
			t.Logf("all:%s\n", msg)
		}
	}()

	go func() {
		//接受对应订阅的主题
		for msg := range golang {
			t.Logf("golang:%s\n", msg)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	t.Logf("quit(%v)\n", <-sig)
}
