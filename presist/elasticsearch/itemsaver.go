package elasticsearch

import (
	"context"
	"fmt"
	"go-clawer/engine"
	"log"

	"github.com/olivere/elastic/v7"
)

func ItemSaver(index string) (chan engine.Item, error) {
	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		return nil, err
	}
	out := make(chan engine.Item)
	go func() {
		itemCount := 0
		for {
			item := <-out
			log.Printf("Item Saver: got item "+
				"#%d: %v", itemCount, item)
			itemCount++
			_, err := Save(client, index, item)
			if err != nil {
				fmt.Printf("save error, saving item:%v,error:%v", item, err)
			}
		}

	}()
	return out, nil
}

func Save(client *elastic.Client, index string, item engine.Item) (string, error) {
	indexClient := client.Index().Index(index)
	if item.Id != "" {
		indexClient = indexClient.Id(item.Id)
	}
	resp, err := indexClient.BodyJson(item).Do(context.Background())
	if err != nil {
		return "", err
	}
	return resp.Id, nil
}
