package aliyun

import (
	"fmt"
	"github.com/trumanwong/cryptogo"
	"net/url"
	"strings"
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
