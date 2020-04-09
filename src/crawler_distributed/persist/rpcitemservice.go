package persist

import (
	"crawler/engine"
	"crawler/persist"
	"log"

	"github.com/olivere/elastic"
)

type ItemService struct {
	Client *elastic.Client
	Index  string
}

func (s *ItemService) Save(args engine.Item, result *string) error {
	err := persist.Save(s.Client, s.Index, args)
	log.Printf("Item %v saved", args)
	if err == nil {
		*result = "ok"
	} else {
		log.Printf("Error saving Item item:%v,err:%v", args, err)
	}
	return err
}
