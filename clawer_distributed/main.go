package main

import (
	"flag"
	"fmt"
	elasticsearch "go-clawer/clawer_distributed/persist/elasticsearch/client"
	"go-clawer/clawer_distributed/rpcsupport"
	"go-clawer/clawer_distributed/worker/client"
	"go-clawer/config"
	"go-clawer/engine"
	"go-clawer/mock/parser"
	"go-clawer/scheduler"
	"net/rpc"
	"strings"
)

var (
	itemSaverHost = flag.String(
		"itemsaver_host", "", "itemsaver host")

	workerHosts = flag.String(
		"worker_hosts", "",
		"worker hosts (comma separated)")
)

func main() {
	flag.Parse()
	itemChan, err := elasticsearch.ItemSaver(*itemSaverHost)
	if err != nil {
		panic(err)
	}

	rpcClient := CreateClientPool(strings.Split(*workerHosts, ","))
	RequestProcessor := client.CreateProcess(rpcClient)

	//simpleEngine := engine.SimpleEngine{
	//	ItemChan: itemChan,
	//}
	//simpleEngine.Run(engine.Request{Url: "http://localhost:8080/mock/www.zhenai.com/zhenghun", ParserFun: parser.ParserCityList})
	concurrentEngine := engine.ConcurrentEngine{
		Scheduler:        &scheduler.QueueScheduler{},
		WorkerCount:      10,
		ItemChan:         itemChan,
		RequestProcessor: RequestProcessor,
	}
	concurrentEngine.Run(engine.Request{Url: "http://localhost:8080/mock/www.zhenai.com/zhenghun", Parser: engine.NewFuncParser(parser.ParserCityList, config.ParserCityList)})
}

// 创建worker
func CreateClientPool(hosts []string) chan *rpc.Client {
	out := make(chan *rpc.Client)
	var clients []*rpc.Client
	for _, h := range hosts {
		newClient, err := rpcsupport.NewClient(h)
		if err != nil {
			fmt.Printf("connect to host:%s failed, error:%v\n", h, err)
		}
		clients = append(clients, newClient)
	}
	go func() {
		for {
			for _, c := range clients {
				out <- c
			}
		}
	}()
	return out
}
