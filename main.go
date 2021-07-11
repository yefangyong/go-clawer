package main

import (
	"go-clawer/engine"
	"go-clawer/mock/parser"
	"go-clawer/scheduler"
)

func main() {
	//	engine.SimpleEngine{}.Run(engine.Request{Url: "http://localhost:8080/mock/www.zhenai.com/zhenghun", ParserFun: parser.ParserCityList})
	concurrentEngine := engine.ConcurrentEngine{
		Scheduler:   &scheduler.QueueScheduler{},
		WorkerCount: 10,
	}
	concurrentEngine.Run(engine.Request{Url: "http://localhost:8080/mock/www.zhenai.com/zhenghun", ParserFun: parser.ParserCityList})

}
