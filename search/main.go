package main

import (
	"encoding/json"
	"fmt"
	"go-clawer/utils"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"golang.org/x/text/encoding/simplifiedchinese"
)

const baiduUrl = "https://www.baidu.com/su?&wd=%s&p=3&cb=BaiduSuggestion.callbacks.give1628576397062&json=1&t=%s"
const OldSearchUrl = "https://so.2345.com/index/search.php?wd=%s&t=7.10&ver=v2.0&charset=utf-8&channel=ziyou"
const SearchEtUrl = "https://so.2345.com/searchEt?wd=%s&cb=T.adZone.callback&t=7.10&ver=v2.0&charset=utf-8&channel=ziyou"
const SearchUrl = "http://index-api.2345.com/search/search?keyword=%s&channel=1&baiduKeyword=%s&cb=T.adZone.callback"

type OldSearchData struct {
	title    string
	subTitle string
}

type NewSearchData struct {
	title    string
	subTitle string
}

func main() {
	// 读取 excel 文件
	f, err := excelize.OpenFile("/Users/yfy/opt/case/go-clawer/statics/search2.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	out := count()
	// Get all the rows in the Sheet1.
	rows := f.GetRows("Sheet1")
	for _, row := range rows {
		for _, keyword := range row {
			// 获取百度联想词，老的搜索联想接口，新的搜索联想接口，对比数据
			if keyword != "" {
				// 开启一个协程，去获取数据
				go getSearchResult(keyword, out)
				time.Sleep(time.Millisecond * 50)
			}
		}
	}
	time.Sleep(time.Hour)
}

type SearchResult struct {
	NewHitCount  int // 新接口命中数
	OldHitCount  int // 老接口命中数
	SameHitCount int // 新老接口命中，且数据一致
	DiffHitCount int // 新老接口都命中，但是数据不一致
	NewDiffCount int // 新接口命中，老接口没有命中
	OldDiffCount int // 旧接口命中，新接口没有命中
}

func getSearchResult(keyword string, out chan SearchResult) {
	oldHitCount := 0  // 新接口命中数
	newHitCount := 0  // 老接口命中数
	sameHitCount := 0 // 新老接口命中，且数据一致
	diffHitCount := 0 // 新老接口都命中，但是数据不一致
	newDiffCount := 0 // 新接口命中，老接口没有命中
	oldDiffCount := 0 // 旧接口命中，新接口没有命中
	newSearchData := &NewSearchData{}
	oldSearchData := &OldSearchData{}
	baiduKeyword, err := getBaiduKeyword(keyword)
	fmt.Printf("百度联想词：%s\n", baiduKeyword)
	if err != nil {
		baiduKeyword = ""
	}

	// 新的搜索逻辑
	newSearchData, err = getSearchData(keyword, baiduKeyword)
	if err == nil && newSearchData != nil {
		newHitCount = newHitCount + 1
	}

	// 老的搜索联想逻辑
	if baiduKeyword != "" {
		// 根据联想词进行搜索，如果结果为null，则根据词根进行搜索
		oldSearchData, err = getOldSearchData(baiduKeyword)
		if oldSearchData != nil && err == nil {
			oldHitCount = oldHitCount + 1
		} else {
			oldSearchData, err = getSearchEt(keyword)
			if oldSearchData != nil && err == nil {
				oldHitCount = oldHitCount + 1
			}
		}
	} else {
		oldSearchData, err = getSearchEt(keyword)
		if oldSearchData != nil && err == nil {
			oldHitCount = oldHitCount + 1
		}
	}

	if oldSearchData != nil && newSearchData != nil {
		if oldSearchData.title == newSearchData.title && oldSearchData.subTitle == newSearchData.subTitle {
			fmt.Printf("新老接口数据一致：%v\n", oldSearchData)
			sameHitCount = sameHitCount + 1
		} else {
			fmt.Printf("老的数据：%v\n", oldSearchData)
			fmt.Printf("新的数据：%v\n", newSearchData)
			diffHitCount = diffHitCount + 1
		}
	}

	if oldSearchData == nil && newSearchData != nil {
		newDiffCount = newDiffCount + 1
	}

	if oldSearchData != nil && newSearchData == nil {
		oldDiffCount = oldDiffCount + 1
	}
	go func() {
		out <- SearchResult{
			NewHitCount:  newHitCount,
			OldHitCount:  oldHitCount,
			DiffHitCount: diffHitCount,
			NewDiffCount: newDiffCount,
			SameHitCount: sameHitCount,
			OldDiffCount: oldDiffCount,
		}
	}()
}

func getSearchData(keyword string, baiduKeyword string) (*NewSearchData, error) {
	urlPath := fmt.Sprintf(SearchUrl, url.QueryEscape(keyword), url.QueryEscape(baiduKeyword))
	res, err := utils.HttpGet(urlPath)
	if err != nil {
		return nil, err
	}
	searchData := make(map[string]string, 0)
	_ = json.Unmarshal([]byte(utils.JsonPToJson(string(res))), &searchData)
	if len(searchData) == 0 {
		return nil, nil
	}
	newSearchData := &NewSearchData{}
	newSearchData.title = searchData["title"]
	newSearchData.subTitle = searchData["sub_title"]
	return newSearchData, nil
}

func getSearchEt(keyword string) (*OldSearchData, error) {
	urlPath := fmt.Sprintf(SearchEtUrl, url.QueryEscape(keyword))
	res, err := utils.HttpGet(urlPath)
	if err != nil {
		return nil, err
	}
	result, err := simplifiedchinese.GBK.NewDecoder().Bytes([]byte(res))
	if err != nil {
		return nil, err
	}
	searchData := make(map[string]string, 0)
	_ = json.Unmarshal([]byte(utils.JsonPToJson(string(result))), &searchData)
	if len(searchData) == 0 {
		return nil, nil
	}
	oldSearchData := &OldSearchData{}
	oldSearchData.title = searchData["title"]
	oldSearchData.subTitle = searchData["subtitle"]
	return oldSearchData, nil
}

func getOldSearchData(baiduKeyword string) (*OldSearchData, error) {
	urlPath := fmt.Sprintf(OldSearchUrl, url.QueryEscape(baiduKeyword))
	res, err := utils.HttpGet(urlPath)
	if err != nil {
		return nil, err
	}
	result, err := simplifiedchinese.GBK.NewDecoder().Bytes([]byte(res))
	if err != nil {
		return nil, err
	}
	searchData := make(map[string]string, 0)
	_ = json.Unmarshal([]byte(utils.JsonPToJson(string(result))), &searchData)
	if len(searchData) == 0 {
		return nil, nil
	}
	oldSearchData := &OldSearchData{}
	oldSearchData.title = searchData["w"]
	oldSearchData.subTitle = searchData["t"]
	return oldSearchData, nil
}

// 获取百度联想词
func getBaiduKeyword(keyword string) (string, error) {
	urlPath := fmt.Sprintf(baiduUrl, strings.Trim(keyword, " "), strconv.FormatInt(time.Now().Unix(), 10))
	res, err := utils.HttpGet(urlPath)
	if err != nil {
		return "", err
	}

	result, err := simplifiedchinese.GBK.NewDecoder().Bytes([]byte(res))
	if err != nil {
		return "", err
	}

	baiduData := make(map[string]interface{}, 0)
	if len(string(result)) == 0 {
		return "", nil
	}
	err = json.Unmarshal([]byte(utils.JsonPToJson(string(result))), &baiduData)
	if err != nil {
		return "", err
	}
	baiduDataS := baiduData["s"].([]interface{})
	baiduKeyword := ""
	if len(baiduDataS) != 0 {
		baiduKeyword = baiduDataS[0].(string)
	}
	return baiduKeyword, nil
}

func count() chan SearchResult {
	out := make(chan SearchResult)
	go func() {
		oldHitCount := 0  // 新接口命中数
		newHitCount := 0  // 老接口命中数
		sameHitCount := 0 // 新老接口命中，且数据一致
		diffHitCount := 0 // 新老接口都命中，但是数据不一致
		newDiffCount := 0 // 新接口命中，老接口没有命中
		oldDiffCount := 0 // 旧接口命中，新接口没有命中
		itemCount := 0
		for {
			item := <-out
			itemCount++
			oldHitCount = oldHitCount + item.OldHitCount
			newHitCount = newHitCount + item.NewHitCount
			sameHitCount = sameHitCount + item.SameHitCount
			diffHitCount = diffHitCount + item.DiffHitCount
			newDiffCount = newDiffCount + item.NewDiffCount
			oldDiffCount = oldDiffCount + item.OldDiffCount
			log.Printf("get Item: #%d\n新接口命中数:%d\n老接口命中数:%d\n新老接口命中，且数据一致数:%d\n新老接口都命中，但是数据不一致:%d\n新接口命中，老接口没有命中:%d\n旧接口命中，新接口没有命中:%d\n",
				itemCount, newHitCount, oldHitCount, sameHitCount, diffHitCount, newDiffCount, oldDiffCount)
		}
	}()
	return out
}
