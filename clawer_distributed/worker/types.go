package worker

import (
	"errors"
	"fmt"
	"go-clawer/config"
	"go-clawer/engine"
	"go-clawer/mock/parser"
	"log"
)

type SerializedParse struct {
	Name string
	Args interface{}
}

type Request struct {
	Url    string
	Parser SerializedParse
}

type ParserResult struct {
	Items    []engine.Item
	Requests []Request
}

func SerializedRequest(r engine.Request) Request {
	name, args := r.Parser.Serialize()
	return Request{
		Url: r.Url,
		Parser: SerializedParse{
			Name: name,
			Args: args,
		},
	}
}

func SerializedResult(r engine.ParserResult) ParserResult {
	result := ParserResult{
		Items: r.Items,
	}
	for _, req := range r.Requests {
		result.Requests = append(result.Requests, SerializedRequest(req))
	}
	return result
}

func DeserializedRequest(r Request) (engine.Request, error) {
	deserializedParse, err := DeserializedParse(r.Parser)
	if err != nil {
		return engine.Request{}, err
	}
	return engine.Request{
		Url:    r.Url,
		Parser: deserializedParse,
	}, nil
}

func DeserializedResult(r ParserResult) engine.ParserResult {
	result := engine.ParserResult{
		Items: r.Items,
	}
	for _, req := range r.Requests {
		engineReq, err := DeserializedRequest(req)
		if err != nil {
			log.Printf("error deserializing "+
				"request: %v", err)
		}
		result.Requests = append(result.Requests, engineReq)
	}
	return result
}

func DeserializedParse(p SerializedParse) (engine.Parser, error) {
	switch p.Name {
	case config.ParserCity:
		return engine.NewFuncParser(parser.ParserCity, config.ParserCity), nil
	case config.ParserCityList:
		return engine.NewFuncParser(parser.ParserCityList, config.ParserCityList), nil
	case config.ParserProfile:
		if userName, ok := p.Args.(string); ok {
			return parser.NewProfileParse(userName), nil
		} else {
			return nil, fmt.Errorf("invalid "+
				"arg: %v", p.Args)
		}
	default:
		return nil, errors.New("unknown parse name")
	}
}
