package engine

import (
	"log"
)

type SimpleEngine struct {
	ItemChan chan Item
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
			// 数据持久化
			go func() {
				e.ItemChan <- item
			}()
			log.Printf("Got item v: %v", item)
		}
		requests = append(requests, parserResult.Requests...)
	}
}
