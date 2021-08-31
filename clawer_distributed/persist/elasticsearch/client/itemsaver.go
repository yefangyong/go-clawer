package persist

import (
	"fmt"
	"go-clawer/clawer_distributed/rpcsupport"
	"go-clawer/engine"
	"log"
)

func ItemSaver(host string) (chan engine.Item, error) {
	out := make(chan engine.Item)
	client, err := rpcsupport.NewClient(host)
	if err != nil {
		return nil, err
	}
	go func() {
		itemCount := 0
		for {
			item := <-out
			log.Printf("Item Saver: got item "+
				"#%d: %v", itemCount, item)
			itemCount++
			result := ""
			err := client.Call("ItemSaverService.Save", item, &result)
			if err != nil {
				fmt.Printf("save error, saving item:%v,error:%v", item, err)
			}
		}

	}()
	return out, nil
}
