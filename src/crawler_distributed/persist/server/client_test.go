package main

import (
	"crawler/Model"
	"crawler/engine"
	"crawler_distributed/config"
	"crawler_distributed/rpcsupport"
	"testing"
	"time"
)

func TestStartUpRpcServer(t *testing.T) {
	host := ":12345"
	index := "dating_test"

	go StartUpRpcServer(host, index)
	time.Sleep(time.Second)

	client, err := rpcsupport.FactoryClient(host)

	if err != nil {
		panic(err)
	}
	args := engine.Item{
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
	result := ""
	err = client.Call(config.ItemServiceMethod, args, &result)
	if err != nil || result != "ok" {
		t.Errorf("err:%+v,result:%s", err, result)
	}
}
