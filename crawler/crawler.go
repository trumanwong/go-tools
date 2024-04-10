package crawler

import (
	"io"
	"net/http"
	"time"
)

type Request struct {
	// 请求地址
	Url string
	// 请求方法
	Method string
	// 请求头
	Headers map[string]string
	// 请求体
	Body io.Reader
	// 请求代理
	Transport http.RoundTripper
	// 请求超时时间
	Timeout   time.Duration
	BasicAuth *BasicAuth
}

type BasicAuth struct {
	Username string
	Password string
}

// Send 发送请求
func Send(request *Request) (*http.Response, error) {
	client := &http.Client{
		Transport: request.Transport,
		Timeout:   request.Timeout,
	}
	req, err := http.NewRequest(request.Method, request.Url, request.Body)
	if err != nil {
		return nil, err
	}
	if request.BasicAuth != nil {
		req.SetBasicAuth(request.BasicAuth.Username, request.BasicAuth.Password)
	}
	if request.Headers != nil {
		for k, v := range request.Headers {
			req.Header.Set(k, v)
		}
	}

	return client.Do(req)
}
