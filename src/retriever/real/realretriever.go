/*
实现者Retriever，它不需要说我实现了那个接口，
它只要实现接口里的方法就可以。
*/
package real

import (
	"net/http"
	"net/http/httputil"
	"time"
)

type Retriever struct {
	UserAgent string
	TimeOut   time.Duration
}

// go语言实现接口
// 接口的实现是隐式的
// 只要实现接口里的方法
func (r *Retriever) Get(url string) string {
	// 发送一个http get 请求
	resp, error := http.Get(url)
	if error != nil {
		panic(error)
	}
	// 解析http响应 转存response
	result, error := httputil.DumpResponse(resp, true)
	defer resp.Body.Close()

	if error != nil {
		panic(error)
	}
	return string(result)
}
