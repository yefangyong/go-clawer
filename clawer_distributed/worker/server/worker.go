package main

import (
	"go-clawer/clawer_distributed/rpcsupport"
	"go-clawer/clawer_distributed/worker"
	"log"
)

func main() {
	log.Fatal(rpcsupport.ServeRpc("127.0.0.1:2345", worker.CrawlService{}))
}
