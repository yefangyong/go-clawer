package parser

import (
	"go-clawer/engine"
	"regexp"
)
//<a href="http://localhost:8080/mock/album.zhenai.com/u/1489297721356860737">混過也愛過小气鬼</a>
const cityRe = `<a href="(.*album\.zhenai\.com/u/[0-9]+)"[^>]*>([^<]+)</a>`

func ParserCity(contents []byte) engine.ParserResult {
	re := regexp.MustCompile(cityRe)
	match := re.FindAllStringSubmatch(string(contents), -1)
	result := engine.ParserResult{}
	for _, m := range match {
		result.Items = append(result.Items, "User "+string(m[2]))
		result.Requests = append(result.Requests, engine.Request{Url: string(m[1]), ParserFun: ParserProfile})
	}
	return result
}