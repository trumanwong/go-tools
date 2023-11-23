package crawler

import (
	"github.com/trumanwong/cryptogo"
	"io"
	"net/http"
	"time"
)

func GenerateCookie() string {
	mUID := cryptogo.MD5ToUpper([]byte(time.Now().Format(time.DateTime)))
	gUID := cryptogo.MD5ToUpper([]byte(mUID))
	dOB := time.Now().Format("20060102150405")
	sID := cryptogo.MD5ToUpper([]byte(dOB))
	result := "MUID=" + mUID + "; SNRHOP=I=&TS=; SRCHD=AF=NOFORM; SRCHHPGUSR=SRCHLANG=zh-Hans; SRCHUID=V=2&GUID=" + gUID + "&dmnchg=1; SRCHUSR=DOB=" + dOB + "; SUID=M; _EDGE_S=F=1&SID=" + sID + "; _EDGE_V=1; _SS=SID=" + sID + "; MUIDB=" + mUID
	return result
}

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
	Timeout time.Duration
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

	if request.Headers != nil {
		for k, v := range request.Headers {
			req.Header.Set(k, v)
		}
	}

	return client.Do(req)
}
