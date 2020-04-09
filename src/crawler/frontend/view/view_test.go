package view

import (
	"crawler/Model"
	"crawler/engine"
	"crawler/frontend/model"
	"os"
	"testing"
)

func TestView_Render(t *testing.T) {
	view := FactoryView("template.html")
	page := model.SearchResult{}
	// 创建一个文件
	file, err := os.Create("template.test.html")
	if err != nil {
		panic(err)
	}
	page.Hits = 23
	item := engine.Item{
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

	for i := 0; i < 10; i++ {
		page.Items = append(page.Items, item)
	}
	// os.Stdout 直接把数据输出到屏幕
	//err := template.Execute(os.Stdout, page)
	err = view.Render(file, page)
	if err != nil {
		panic(err)
	}
}
