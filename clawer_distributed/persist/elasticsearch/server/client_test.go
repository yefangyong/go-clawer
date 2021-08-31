package main

import (
	"go-clawer/clawer_distributed/rpcsupport"
	"go-clawer/engine"
	"go-clawer/models"
	"testing"
	"time"
)

func TestItemSaver(t *testing.T) {
	// start ItemServer
	const host = "127.0.0.1:1234"
	go ServeRpc(host, "test")
	time.Sleep(time.Second)
	// start ItemClient
	client, err := rpcsupport.NewClient(host)
	if err != nil {
		panic(err)
	}
	item := engine.Item{
		Url: "http://album.zhenai.com/u/108906739",
		Id:  "108906739",
		PayLoad: models.Profiles{
			Age:        34,
			Height:     162,
			Weight:     57,
			Income:     "3001-5000元",
			Gender:     "女",
			Name:       "安静的雪",
			Xinzuo:     "牡羊座",
			Occupation: "人事/行政",
			Marriage:   "离异",
			House:      "已购房",
			Hokou:      "山东菏泽",
			Education:  "大学本科",
			Car:        "未购车",
		},
	}

	result := ""
	// call save
	err = client.Call("ItemSaverService.Save", item, &result)

	if err != nil || result != "ok" {
		t.Errorf("result: %s; err: %s",
			result, err)
	}
}
