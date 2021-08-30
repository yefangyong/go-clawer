package mysql

import (
	"encoding/json"
	"fmt"
	"go-clawer/engine"
	"log"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

type Zhenai struct {
	Id      int
	Url     string
	Payload string
}

var mysqlMasterConfig = map[string]interface{}{
	"Username": "root",
	"Password": "root",
	"Host":     "127.0.0.1",
	"Port":     3306,
	"Database": "crawl",
}
var mysqlSlaveConfig = []map[string]interface{}{
	{
		"Username": "root",
		"Password": "root",
		"Host":     "127.0.0.1",
		"Port":     3306,
		"Database": "crawl",
	},
}

var Driver = "mysql"

// 初始化链接数据库
func initMysql() *xorm.EngineGroup {
	var dsnList []string
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", mysqlMasterConfig["Username"], mysqlMasterConfig["Password"], mysqlMasterConfig["Host"], mysqlMasterConfig["Port"], mysqlMasterConfig["Database"])
	dsnList = append(dsnList, dsn)
	for _, config := range mysqlSlaveConfig {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", config["Username"], config["Password"], config["Host"], config["Port"], config["Database"])
		dsnList = append(dsnList, dsn)
	}
	engineGroup, err := xorm.NewEngineGroup(Driver, dsnList)
	if err != nil {
		log.Fatal(err)
	}
	engineGroup.ShowSQL(true)
	return engineGroup
}

func ItemSaver() (chan engine.Item, error) {
	out := make(chan engine.Item)
	go func() {
		itemCount := 0
		engineGroup := initMysql()
		for {
			item := <-out
			log.Printf("Item Saver: got item "+
				"#%d: %v", itemCount, item)
			itemCount++
			Save(engineGroup, item)
		}
	}()
	return out, nil
}

func Save(engineGroup *xorm.EngineGroup, item engine.Item) {
	zhenai := new(Zhenai)
	zhenai.Id, _ = strconv.Atoi(item.Id)
	zhenai.Url = item.Url
	payload, _ := json.Marshal(item.PayLoad)
	zhenai.Payload = string(payload)
	res, err := engineGroup.Master().Insert(zhenai)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}
