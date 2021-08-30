package main

import (
	"go-clawer/config"
	"go-clawer/engine"
	"go-clawer/mock/parser"
	"go-clawer/presist/elasticsearch"
	"go-clawer/scheduler"
)

func main() {
	itemChan, err := elasticsearch.ItemSaver(config.ElasticIndex)
	if err != nil {
		panic(err)
	}

	//simpleEngine := engine.SimpleEngine{
	//	ItemChan: itemChan,
	//}
	//simpleEngine.Run(engine.Request{Url: "http://localhost:8080/mock/www.zhenai.com/zhenghun", ParserFun: parser.ParserCityList})
	concurrentEngine := engine.ConcurrentEngine{
		Scheduler:        &scheduler.SimpleScheduler{},
		WorkerCount:      10,
		ItemChan:         itemChan,
		RequestProcessor: engine.Worker,
	}
	concurrentEngine.Run(engine.Request{Url: "http://localhost:8080/mock/www.zhenai.com/zhenghun",
		Parser: engine.NewFuncParser(parser.ParserCityList, "ParserCityList")})

}
