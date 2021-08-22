package main

import (
	persist "go-clawer/clawer_distributed/persist/client"
	"go-clawer/engine"
	"go-clawer/mock/parser"
	"go-clawer/scheduler"
)

func main() {
	itemChan, err := persist.ItemSaver(":1234")
	if err != nil {
		panic(err)
	}

	//simpleEngine := engine.SimpleEngine{
	//	ItemChan: itemChan,
	//}
	//simpleEngine.Run(engine.Request{Url: "http://localhost:8080/mock/www.zhenai.com/zhenghun", ParserFun: parser.ParserCityList})
	concurrentEngine := engine.ConcurrentEngine{
		Scheduler:   &scheduler.QueueScheduler{},
		WorkerCount: 10,
		ItemChan:    itemChan,
	}
	concurrentEngine.Run(engine.Request{Url: "http://localhost:8080/mock/www.zhenai.com/zhenghun", ParserFun: parser.ParserCityList})

}
