package engine

import (
	"go-clawer/fetcher"
	"log"
)

type SimpleEngine struct {
}

func (e SimpleEngine) Run(seeds ...Request) {
	var requests []Request
	for _, req := range seeds {
		requests = append(requests, req)
	}

	for len(requests) > 0 {
		r := requests[0]
		requests = requests[1:]
		parserResult, err := Worker(r)
		if err != nil {
			continue
		}
		for _, item := range parserResult.Items {
			log.Printf("Got item %v", item)
		}
		requests = append(requests, parserResult.Requests...)
	}
}

func Worker(r Request) (ParserResult, error) {
	body, err := fetcher.Fetch(r.Url)
	//log.Printf("Get Url:%s", r.Url)
	if err != nil {
		log.Printf("Fetcher: error"+"fetching url %s:%v", r.Url, err)
		return ParserResult{}, nil
	}
	return r.ParserFun(body), nil
}
