package aliyun

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/trumanwong/cryptogo"
)

type CdnAuth struct {
	// 鉴权Key
	Key string
}

func NewCdnAuth(key string) *CdnAuth {
	return &CdnAuth{Key: key}
}

// AuthA 鉴权方式A
// rawUrl: 原始URL
// expireTime: 链接到期时间，unix时间戳
// rand 随机数，建议使用uuid
func (c CdnAuth) AuthA(rawUrl string, expireTime int64, rand string) (string, error) {
	urlPath, err := url.Parse(rawUrl)
	if err != nil {
		return "", err
	}
	encodePath := strings.ReplaceAll(url.QueryEscape(urlPath.Path), "%2F", "/")
	authKey := fmt.Sprintf("%s-%d-%s-0-%s", encodePath, expireTime, rand, c.Key)
	mdHash := cryptogo.MD5([]byte(authKey))
	return urlPath.Scheme + "://" + urlPath.Host + encodePath + "?auth_key=" + fmt.Sprintf("%d-%s-0-%s", expireTime, rand, mdHash), nil
}

func (c CdnAuth) AuthB(rawUrl string, expireTime int64) (string, error) {
	urlPath, err := url.Parse(rawUrl)
	if err != nil {
		return "", err
	}
	scheme := urlPath.Scheme
	if scheme == "" {
		scheme = "http"
	}
	host := urlPath.Host
	path := urlPath.Path
	if path == "" {
		path = "/"
	}
	encodePath := strings.ReplaceAll(url.QueryEscape(path), "%2F", "/")
	args := ""
	if urlPath.RawQuery != "" {
		args = "?" + urlPath.RawQuery
	}
	// 转换时间戳为"YYYYmmDDHHMM"格式
	nexp := time.Unix(expireTime, 0).Format("200601021504")
	sstring := c.Key + nexp + encodePath
	hashvalue := cryptogo.MD5([]byte(sstring))
	// 拼接最终URL
	finalUrl := fmt.Sprintf("%s://%s/%s/%s%s%s", scheme, host, nexp, hashvalue, encodePath, args)
	return finalUrl, nil
}
