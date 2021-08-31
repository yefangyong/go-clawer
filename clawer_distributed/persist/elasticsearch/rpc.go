package elasticsearch

import (
	"go-clawer/engine"
	"go-clawer/presist/elasticsearch"
	"log"

	"github.com/olivere/elastic/v7"
)

type ItemSaverService struct {
	Client *elastic.Client
	Index  string
}

func (s ItemSaverService) Save(item engine.Item, result *string) error {
	itemId, err := elasticsearch.Save(s.Client, s.Index, item)
	log.Printf("Item %v saved.", item)
	if err != nil {
		log.Printf("Item save error ,item:%v,error:%v", item, err)
	} else {
		log.Printf("Item save success ,itemId:%v", itemId)
		*result = "ok"
	}
	return err
}
