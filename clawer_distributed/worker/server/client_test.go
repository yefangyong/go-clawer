package main

import (
	"fmt"
	"go-clawer/clawer_distributed/rpcsupport"
	"go-clawer/clawer_distributed/worker"
	"go-clawer/config"
	"testing"
)

func TestWorker(t *testing.T) {
	// 创建一个客户端
	const host = "127.0.0.1:2345"
	client, err := rpcsupport.NewClient(host)
	if err != nil {
		t.Error(err)
	}
	req := worker.Request{
		Url: "http://localhost:8080/mock/album.zhenai.com/u/6121809463065467146",
		Parser: worker.SerializedParse{
			Name: config.ParseProfile,
			Args: "不良骚年深碍",
		},
	}
	var result worker.ParserResult
	err = client.Call("CrawlService.Process", req, &result)
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println(result)
	}
}
