/*
show http client

*/
package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
)

func main() {
	// 自定义控制http request
	request, err := http.NewRequest(http.MethodGet, "http://www.imooc.com", nil)
	if err != nil {
		panic(err)
	}
	// add an request header 模拟手机客户端
	request.Header.Add("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 12_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/12.0 Mobile/15A372 Safari/604.1")

	// create an http client，custom 语义是自定义
	customClient := http.Client{
		Transport: nil,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// 打印一下重定向的地址
			fmt.Println("redirect:", req)
			return nil
		},
		Jar:     nil,
		Timeout: 0,
	}
	// custom client send a request of http get
	resp, err := customClient.Do(request)

	// client send a http get request
	//resp, err := http.DefaultClient.Do(request)

	// send a http get request，return a http response
	//resp, err := http.Get("http://www.imooc.com")
	if err != nil {
		// panic 语义是惊慌的
		panic(err)
	}
	// defer 语义是推迟
	defer resp.Body.Close()

	// util 是跑龙套的意思，http龙套库 dumps 语义是存储 存储http响应 第二个参数是代表是否存储 response body
	result, err := httputil.DumpResponse(resp, true)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", result)
}
