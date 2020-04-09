package client

import (
	"crawler/engine"
	"crawler_distributed/config"
	"crawler_distributed/rpcsupport"
	"log"
)

func ItemServer(host string) (chan engine.Item, error) {

	client, err := rpcsupport.FactoryClient(host)
	if err != nil {
		return nil, err
	}
	out := make(chan engine.Item)
	go func() {
		itemCount := 0
		for {
			//channel的特点，等到了就执行，等不到就不执行
			//received an item from engine
			item := <-out
			itemCount++
			log.Printf("Item Server Got item: count:%d,val:%v", itemCount, item)
			// call rpc to save item
			result := ""
			err = client.Call(config.ItemServiceMethod, item, &result)
			if err != nil {
				log.Printf("Item Server:error saving item %v,%v", item, err)
			}
		}

	}()
	return out, nil
}
