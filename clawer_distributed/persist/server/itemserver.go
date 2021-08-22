package main

import (
	"go-clawer/clawer_distributed/persist"
	"go-clawer/clawer_distributed/rpcsupport"
	"log"

	"github.com/olivere/elastic/v7"
)

func main() {
	log.Fatal(ServeRpc(":1234", "test1"))
}

func ServeRpc(host string, index string) error {
	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		return err
	}
	return rpcsupport.ServeRpc(host, &persist.ItemSaverService{Client: client, Index: index})
}
