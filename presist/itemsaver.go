package presist

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"log"
)

func ItemSaver() chan interface{} {
	out := make(chan interface{})
	go func() {
		itemCount := 0
		for {
			item := <-out
			log.Printf("Item Saver: got item "+
				"#%d: %v", itemCount, item)
			itemCount++
			_, err := save(item)
			if err != nil {
				fmt.Printf("save error, saving item:%v,error:%v", item, err)
			}
		}

	}()
	return out
}

func save(item interface{}) (string, error) {
	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		return "", err
	}
	resp, err := client.Index().Index("crawler_data").BodyJson(item).Do(context.Background())
	if err != nil {
		return "", err
	}
	return resp.Id, nil
}
