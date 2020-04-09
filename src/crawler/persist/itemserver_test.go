package persist

import (
	"context"
	"crawler/Model"
	"crawler/engine"
	"encoding/json"
	"testing"

	"github.com/olivere/elastic"
)

func TestSave(t *testing.T) {
	expected := engine.Item{
		Url:  "http://album.zhenai.com/u/16742041818",
		Id:   "16742041818",
		Type: "zhenai",
		Payload: Model.Profile{
			Name:       "耿乾钰",
			Height:     170,
			Age:        30,
			Income:     "3000-5000元",
			Marriage:   "离异",
			Occupation: "大专",
			Birthplace: "石家庄",
			BasicInfo:  "未婚,29岁,魔羯座(12.22-01.19),158cm,工作地:上海青浦区,月收入:8千-1.2万,大学本科,",
			DetailInfo: "籍贯:北京,租房,",
		},
	}
	const index = "dating_test"
	client, err := elastic.NewClient(elastic.SetSniff(false))
	// save expected item
	err = Save(client, index, expected)
	if err != nil {
		panic(err)
	}
	// fetch saved item
	result, err := client.Get().
		Index(index).
		Type(expected.Type).
		Id(expected.Id).
		Do(context.Background())
	if err != nil {
		panic(err)
	}
	//t.Logf("%+v", result)
	t.Logf("%T,%v", result.Source, result.Source)
	var actual engine.Item
	// 将json数据解析到结构体中。
	// Unmarshal将传入的对象键与Marshal使用的键（结构字段名称或其标记）进行匹配
	err = json.Unmarshal([]byte(*result.Source), &actual)
	if err != nil {
		panic(err)
	}
	// 将json map转成profile type
	actual.Payload, _ = Model.FromJsonObj(actual.Payload)

	// verify result
	if expected != actual {
		t.Errorf("expected %v;got %v", expected, actual)
	}
}
