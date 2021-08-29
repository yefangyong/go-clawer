package worker

import "go-clawer/engine"

type CrawlService struct {
}

func (CrawlService) Process(request Request, result *ParserResult) error {
	engineReq, err := DeserializedRequest(request)
	if err != nil {
		return err
	}
	engineResult, err := engine.Worker(engineReq)
	if err != nil {
		return err
	}
	*result = SerializedResult(engineResult)
	return nil
}
