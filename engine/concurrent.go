package engine

import "fmt"

type ConcurrentEngine struct {
	Scheduler   Scheduler
	WorkerCount int
}

type Scheduler interface {
	ConfigureMasterWorkerChan(chan Request)
	Submit(request Request)
	WorkerReady(chan Request)
	Run()
}

func (e *ConcurrentEngine) Run(seeds ...Request) {
	out := make(chan ParserResult)
	e.Scheduler.Run()
	for i := 0; i < e.WorkerCount; i++ {
		createWorker(out,e.Scheduler)
	}
	for _, r := range seeds {
		e.Scheduler.Submit(r)
	}
	for {
		result := <-out
		for _, item := range result.Items {
			fmt.Printf("Get item:%v\n", item)
		}
		for _, request := range result.Requests {
			e.Scheduler.Submit(request)
		}
	}
}

func createWorker(out chan ParserResult, s Scheduler) {
	in := make(chan Request)
	go func() {
		for {
			s.WorkerReady(in)
			request := <-in
			result, err := Worker(request)
			if err != nil {
				continue
			}
			out <- result
		}
	}()
}
