package persist

import (
	"context"
	"crawler/engine"
	"errors"
	"log"

	"github.com/olivere/elastic"
)

func ItemServer(index string) (chan engine.Item, error) {
	out := make(chan engine.Item)
	//must turn off sniff in docker
	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		return out, err
	}
	go func() {
		itemCount := 0
		for {
			item := <-out
			itemCount++
			log.Printf("Item Server Got item: count:%d,val:%v", itemCount, item)
			err := Save(client, index, item)
			if err != nil {
				log.Printf("Item Server Saving item err: item:%v,val:%v", item, err)
			}
		}
	}()
	return out, nil
}

// save data to elasticSearch
func Save(client *elastic.Client, index string, item engine.Item) error {

	if item.Type == "" {
		return errors.New("must supply Type")
	}

	indexServer := client.Index().
		Index(index).
		Type(item.Type)

	if item.Id != "" {
		indexServer.Id(item.Id)
	}
	_, err := indexServer.
		BodyJson(item).
		Do(context.Background())
	if err != nil {
		return err
	}
	// 打印结构体
	//fmt.Printf("%+v", rep)
	return nil
}
