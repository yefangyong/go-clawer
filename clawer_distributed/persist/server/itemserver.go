package main

import (
	"flag"
	"fmt"
	"go-clawer/clawer_distributed/persist"
	"go-clawer/clawer_distributed/rpcsupport"
	"go-clawer/config"
	"log"

	"github.com/olivere/elastic/v7"
)

var port = flag.Int("port", 0,
	"the port for me to listen on")

func main() {
	flag.Parse()
	if *port == 0 {
		fmt.Println("must specify a port")
		return
	}
	log.Fatal(ServeRpc(fmt.Sprintf(":%d", *port), config.ElasticIndex))
}

func ServeRpc(host string, index string) error {
	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		return err
	}
	return rpcsupport.ServeRpc(host, &persist.ItemSaverService{Client: client, Index: index})
}
