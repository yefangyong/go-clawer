package main

import (
	"bufio"
	"fmt"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
)

func main() {
	resp, err := http.Get("http://localhost:8080/mock/www.zhenai.com/zhenghun")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Println("error: status code", resp.StatusCode)
		return
	}
	e := determineEncoding(resp.Body)
	urf8Reader := transform.NewReader(resp.Body, e.NewDecoder())
	all, err := ioutil.ReadAll(urf8Reader)
	printCityList(all)
}

func determineEncoding(r io.Reader) encoding.Encoding {
	bytes, err := bufio.NewReader(r).Peek(1024)
	if err != nil {
		panic(err)
	}
	//根据文档内容自动发现文档的编码
	e, _, _ := charset.DetermineEncoding(
		bytes,
		"",
	)
	return e
}

//<a href="http://localhost:8080/mock/www.zhenai.com/zhenghun/aomen" class="">澳门</a>
func printCityList(contents []byte) {
	re := regexp.MustCompile(`<a href="(http://localhost:8080/mock/www.zhenai.com/zhenghun/[0-9a-zA-Z]+)"[^>]*>([^<]+)</a>`)
	match := re.FindAllStringSubmatch(string(contents), -1)
	fmt.Println(len(match))
	for _, m := range match {
		fmt.Printf("City:%s,Url:%s\n",m[2], m[1])
	}
}
