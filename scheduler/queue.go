package scheduler

import "go-clawer/engine"

type QueueScheduler struct {
	requestChan chan engine.Request
	workerChan  chan chan engine.Request
}

func (q *QueueScheduler) WorkerReady(requests chan engine.Request) {
	q.workerChan <- requests
}

func (q *QueueScheduler) ConfigureMasterWorkerChan(requests chan engine.Request) {
	panic("implement me")
}

func (q *QueueScheduler) Submit(request engine.Request) {
	q.requestChan <- request
}

func (q *QueueScheduler) Run() {
	q.requestChan = make(chan engine.Request)
	q.workerChan = make(chan chan engine.Request)
	go func() {
		var RequestQ []engine.Request
		var WorkerQ []chan engine.Request
		for {
			var activeRequest engine.Request
			var activeWorker chan engine.Request
			if len(RequestQ) > 0 && len(WorkerQ) > 0 {
				activeRequest = RequestQ[0]
				activeWorker = WorkerQ[0]
			}
			select {
			case r := <-q.requestChan:
				RequestQ = append(RequestQ, r)
			case w := <-q.workerChan:
				WorkerQ = append(WorkerQ, w)
			case activeWorker <- activeRequest:
				RequestQ = RequestQ[1:]
				WorkerQ = WorkerQ[1:]
			}
		}
	}()
}
