package crawler

import (
	"github.com/trumanwong/cryptogo"
	"io"
	"net/http"
	"time"
)

func GenerateCookie() string {
	mUID := cryptogo.MD5ToUpper(time.Now().Format(time.DateTime))
	gUID := cryptogo.MD5ToUpper(mUID)
	dOB := time.Now().Format("20060102150405")
	sID := cryptogo.MD5ToUpper(dOB)
	result := "MUID=" + mUID + "; SNRHOP=I=&TS=; SRCHD=AF=NOFORM; SRCHHPGUSR=SRCHLANG=zh-Hans; SRCHUID=V=2&GUID=" + gUID + "&dmnchg=1; SRCHUSR=DOB=" + dOB + "; SUID=M; _EDGE_S=F=1&SID=" + sID + "; _EDGE_V=1; _SS=SID=" + sID + "; MUIDB=" + mUID
	return result
}

// Request 发送请求
// Request send request
// url: 请求地址
// method: 请求方法
// headers: 请求头
// body: 请求体
func Request(url, method string, headers map[string]string, body io.Reader) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, body)

	if err != nil {
		return nil, err
	}

	if _, ok := headers["Cookie"]; !ok {
		req.Header.Add("Cookie", GenerateCookie())
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36 Edg/115.0.1901.200")
	req.Header.Set("Referer", "https://developer.microsoft.com/")
	req.Header.Set("Sec-Ch-Ua", "\"Not/A)Brand\";v=\"99\", \"Microsoft Edge\";v=\"115\", \"Chromium\";v=\"115\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"Windows\"")
	req.Header.Set("Sec-Ch-Ua-Platform-Version", "\"15.0.0\"")
	req.Header.Set("Sec-Ch-Ua-Arch", "\"x64\"")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	return client.Do(req)
}
