package main

import (
	"go-clawer/engine"
	"go-clawer/mock/parser"
)

func main() {
	engine.Run(engine.Request{Url: "http://localhost:8080/mock/www.zhenai.com/zhenghun", ParserFun: parser.ParserCityList})
}
