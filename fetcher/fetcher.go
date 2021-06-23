package fetcher

import (
	"bufio"
	"fmt"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io/ioutil"
	"net/http"
)

func Fetch(url string) ([]byte, error) {
	resp, err := http.Get("http://localhost:8080/mock/www.zhenai.com/zhenghun")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch error status code:%d", resp.StatusCode)
	}
	bufferReader := bufio.NewReader(resp.Body)
	e := determineEncoding(bufferReader)
	urf8Reader := transform.NewReader(bufferReader, e.NewDecoder())
	return ioutil.ReadAll(urf8Reader)
}

func determineEncoding(r *bufio.Reader) encoding.Encoding {
	bytes, err := r.Peek(1024)
	if err != nil {
		//读取文档编码错误，返回默认编码：UTF-8
		return unicode.UTF8
	}
	//根据文档内容自动发现文档的编码
	e, _, _ := charset.DetermineEncoding(
		bytes,
		"",
	)
	return e
}
