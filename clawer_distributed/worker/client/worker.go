package client

import (
	"go-clawer/clawer_distributed/rpcsupport"
	"go-clawer/clawer_distributed/worker"
	"go-clawer/config"
	"go-clawer/engine"
)

func CreateProcess() (engine.Processor, error) {
	client, err := rpcsupport.NewClient("127.0.0.1:2345")
	if err != nil {
		return nil, err
	}
	return func(r engine.Request) (engine.ParserResult, error) {
		sReq := worker.SerializedRequest(r)
		var result worker.ParserResult
		err := client.Call(config.CrawlServiceRpc, sReq, &result)
		if err != nil {
			return engine.ParserResult{}, nil
		}
		return worker.DeserializedResult(result), nil
	}, nil
}
