package engine

import (
	"fmt"
	"go-clawer/models"
)

type ConcurrentEngine struct {
	Scheduler   Scheduler
	WorkerCount int
	ItemChan chan Item
}

type Scheduler interface {
	WorkChan() chan Request
	Submit(request Request)
	ReadyNotify
	Run()
}

type ReadyNotify interface {
	WorkerReady(chan Request)
}

func (e *ConcurrentEngine) Run(seeds ...Request) {
	out := make(chan ParserResult)
	e.Scheduler.Run()
	for i := 0; i < e.WorkerCount; i++ {
		createWorker(e.Scheduler.WorkChan(), out, e.Scheduler)
	}
	for _, r := range seeds {
		if isDuplicate(r.Url) {
			fmt.Printf("Duplicate request:%v\n", r)
			continue
		}
		e.Scheduler.Submit(r)
	}

	for {
		result := <-out
		for _, item := range result.Items {
			if _, ok := item.PayLoad.(models.Profiles); ok {
				go func() {
					e.ItemChan <- item
				}()
			}
		}
		for _, request := range result.Requests {
			if isDuplicate(request.Url) {
				fmt.Printf("Duplicate request:%v\n", request)
				continue
			}
			e.Scheduler.Submit(request)
		}
	}
}

// 去重
var visitUrls = make(map[string]bool)

func isDuplicate(url string) bool {
	if visitUrls[url] {
		return true
	}
	visitUrls[url] = true
	return false
}

func createWorker(in chan Request, out chan ParserResult, notify ReadyNotify) {
	go func() {
		for {
			notify.WorkerReady(in)
			request := <-in
			result, err := Worker(request)
			if err != nil {
				continue
			}
			out <- result
		}
	}()
}
