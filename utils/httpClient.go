package utils

import (
	"io/ioutil"
	"net/http"
	"strings"
)

type HttpCallback func(resp string, err error)

func HttpRequest(method string, url string, params string, headers map[string]string) (string, error) {
	client := &http.Client{}

	req, err := http.NewRequest(method, url, strings.NewReader(params))
	if err != nil {
		return "", err
	}

	for headerKey, headerVal := range headers {
		req.Header.Set(headerKey, headerVal)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func HttpGet(url string) (string, error) {
	header := make(map[string]string, 0)
	return HttpRequest(http.MethodGet, url, "", header)
}

func HttpGetAsync(url string, callback HttpCallback) {
	go func() {
		callback(HttpGet(url))
	}()
}

func HttpPost(url string, params map[string]string) (string, error) {
	header := make(map[string]string, 0)
	header["Content-Type"] = "application/x-www-form-urlencoded"
	paramsStr := MapToString(params, "&", "=")
	return HttpRequest(http.MethodPost, url, paramsStr, header)
}

func HttpPostAsync(url string, params map[string]string, callback HttpCallback) {
	go func() {
		callback(HttpPost(url, params))
	}()
}
