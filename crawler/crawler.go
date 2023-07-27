package crawler

import (
	"io"
	"net/http"
)

// Request 发送请求
// Request send request
// url: 请求地址
// method: 请求方法
// headers: 请求头
// body: 请求体
func Request(url, method string, headers map[string]string, body io.Reader) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, body)

	if err != nil {
		return nil, err
	}

	if headers != nil {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return io.ReadAll(res.Body)
}
