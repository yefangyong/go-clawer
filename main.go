package main

import (
	"go-clawer/engine"
	"go-clawer/mock/parser"
	"go-clawer/presist"
	"go-clawer/scheduler"
)

func main() {
	//	engine.SimpleEngine{}.Run(engine.Request{Url: "http://localhost:8080/mock/www.zhenai.com/zhenghun", ParserFun: parser.ParserCityList})
	itemChan, err := presist.ItemSaver()
	if err != nil {
		panic(err)
	}
	concurrentEngine := engine.ConcurrentEngine{
		Scheduler:   &scheduler.QueueScheduler{},
		WorkerCount: 10,
		ItemChan:    itemChan,
	}
	concurrentEngine.Run(engine.Request{Url: "http://localhost:8080/mock/www.zhenai.com/zhenghun", ParserFun: parser.ParserCityList})

}
