/*
defined fetcher (网页提取器)
*/
package fetcher

import (
	"bufio"
	"crawler_distributed/config"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"golang.org/x/text/encoding/unicode"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
)

//速率限制器
var rateLimiter = time.Tick(time.Second / config.Qps)

func Fetch(url string) ([]byte, error) {
	log.Printf("fetching %s", url)
	<-rateLimiter
	// 自定义控制http request
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:69.0) Gecko/20100101 Firefox/69.0")

	// create an http client，custom 语义是自定义
	customClient := http.Client{
		Transport: nil,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// 打印一下重定向的地址
			//fmt.Println("redirect:", req)
			return nil
		},
		Jar:     nil,
		Timeout: 0,
	}
	// custom client send a request of http get
	resp, err := customClient.Do(request)
	//resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	// defer close response body
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("resp:%v", resp)
		return nil, fmt.Errorf("wrong status code:%d", resp.StatusCode)
	}

	//将response body转成buf reader
	bodyReader := bufio.NewReader(resp.Body)
	// 检测response 编码
	encoding := determineEncoding(bodyReader)
	// 将response body 转成utf8  (simplifiedChinese.GBK)
	utf8Reader := transform.NewReader(bodyReader, encoding.NewEncoder())
	//读取内容,并返回
	return ioutil.ReadAll(utf8Reader)
}

//检测 reader encoding
func determineEncoding(reader *bufio.Reader) encoding.Encoding {
	//buf reader read 1024 byte
	bytes, err := reader.Peek(1024)
	// 如果出错返回默认编码
	if err != nil {
		log.Printf("fetcher encoding error:%v", err)
		return unicode.UTF8
	}

	// determine encoding 判断编码
	encoding, _, _ := charset.DetermineEncoding(bytes, "")

	return encoding
}
