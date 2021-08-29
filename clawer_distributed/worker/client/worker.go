package client

import (
	"go-clawer/clawer_distributed/worker"
	"go-clawer/config"
	"go-clawer/engine"
	"net/rpc"
)

func CreateProcess(client chan *rpc.Client) engine.Processor {
	return func(r engine.Request) (engine.ParserResult, error) {
		sReq := worker.SerializedRequest(r)
		var result worker.ParserResult
		c := <-client
		err := c.Call(config.CrawlServiceRpc, sReq, &result)
		if err != nil {
			return engine.ParserResult{}, nil
		}
		return worker.DeserializedResult(result), nil
	}
}
