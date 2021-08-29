package engine

import (
	"go-clawer/fetcher"
	"log"
)

func Worker(r Request) (ParserResult, error) {
	body, err := fetcher.Fetch(r.Url)
	log.Printf("Get Url:%s", r.Url)
	if err != nil {
		log.Printf("Fetcher: error"+"fetching url %s:%v", r.Url, err)
		return ParserResult{}, nil
	}
	return r.Parser.Parse(body, r.Url), nil
}
