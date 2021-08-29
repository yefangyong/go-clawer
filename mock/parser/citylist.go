package parser

import (
	"go-clawer/config"
	"go-clawer/engine"
	"regexp"
)

const cityListRe = `<a href="(http://localhost:8080/mock/www.zhenai.com/zhenghun/[0-9a-zA-Z]+)"[^>]*>([^<]+)</a>`

func ParserCityList(contents []byte, _ string) engine.ParserResult {
	re := regexp.MustCompile(cityListRe)
	match := re.FindAllStringSubmatch(string(contents), -1)
	result := engine.ParserResult{}
	for _, m := range match {
		result.Requests = append(result.Requests, engine.Request{Url: string(m[1]), Parser: engine.NewFuncParser(ParserCity, config.ParserCity)})
	}
	return result
}
