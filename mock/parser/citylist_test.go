package parser

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestParserCityList(t *testing.T) {
	contents, _ := ioutil.ReadFile("citylisthtml.text")
	result := ParserCityList(contents)
	const ResultCount = 470
	if len(result.Requests) != ResultCount {
		fmt.Println("错误")
	} else {
		fmt.Println("正确")
	}
}
