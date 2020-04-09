package controller

import (
	"context"
	"crawler/engine"
	"crawler/frontend/model"
	"crawler/frontend/view"
	"log"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/olivere/elastic"
)

type SearchHandler struct {
	View   view.View
	Client *elastic.Client
}

var host = "http://localhost:9200/"

func FactorySearchHandler(template string) *SearchHandler {
	client, err := elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL(host),
		elastic.SetHealthcheckInterval(10*time.Second),
		elastic.SetGzip(true),
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
	)
	if err != nil {
		panic(err)
	}
	return &SearchHandler{
		View:   view.FactoryView(template),
		Client: client,
	}
}

//implement Handler interface
func (s *SearchHandler) ServeHTTP(w http.ResponseWriter, rep *http.Request) {
	// input
	q := strings.TrimSpace(rep.FormValue("q"))
	from, err := strconv.Atoi(rep.FormValue("from"))
	if err != nil {
		from = 0
	}
	// 将格式化后的字符串输出到 http response writer
	//fmt.Fprintf(w, "q=%s,from=%d", q, from)

	var searchResult model.SearchResult
	searchResult, err = s.Execute(q, from)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	// 渲染模板 写入http response writer
	err = s.View.Render(w, searchResult)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (s *SearchHandler) Execute(q string, from int) (model.SearchResult, error) {
	var searchResult model.SearchResult
	// search 10 item from elastic search
	result, err := s.Client.Search().
		Index("dating_profile").
		Query(elastic.NewQueryStringQuery(RewriteQueryString(q))).
		Pretty(true).
		From(from).Size(10).
		Do(context.Background())
	if err != nil {
		return searchResult, err
	}
	//defer s.client.Stop()
	searchResult.Query = q
	searchResult.Hits = int(result.TotalHits())
	searchResult.Start = from

	// 反射原理
	for _, v := range result.Each(reflect.TypeOf(engine.Item{})) {
		//	// 将interface{} 转成 engine.Item type
		item := v.(engine.Item)
		searchResult.Items = append(searchResult.Items, item)
	}

	searchResult.PrevFrom = searchResult.Start - len(searchResult.Items)
	searchResult.NextFrom = searchResult.Start + len(searchResult.Items)

	return searchResult, nil
}

func RewriteQueryString(q string) string {
	re := regexp.MustCompile("([A-Za-z]+):")
	return re.ReplaceAllString(q, "Payload.$1:")
}
