package qiniu

import (
	"bytes"
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/cdn"
	"github.com/qiniu/go-sdk/v7/sms/rpc"
	"github.com/qiniu/go-sdk/v7/storage"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	mac           *qbox.Mac
	cdnManager    *cdn.CdnManager
	bucketManager *storage.BucketManager
}

func NewClient(accessKey, secretKey string) *Client {
	mac := qbox.NewMac(accessKey, secretKey)
	cdnManager := cdn.NewCdnManager(mac)
	cfg := storage.Config{
		// 是否使用https域名进行资源管理
		UseHTTPS: true,
	}
	bucketManager := storage.NewBucketManager(mac, &cfg)
	return &Client{mac: mac, cdnManager: cdnManager, bucketManager: bucketManager}
}

func (c Client) PutFile(ctx context.Context, bucket, key string, data []byte) error {
	putPolicy := storage.PutPolicy{
		Scope: fmt.Sprintf("%s:%s", bucket, key),
	}
	upToken := putPolicy.UploadToken(c.mac)

	resumeUploader := storage.NewResumeUploaderV2(&storage.Config{
		Zone:          &storage.ZoneHuanan,
		UseHTTPS:      true,
		UseCdnDomains: false,
	})

	ret := storage.PutRet{}
	err := resumeUploader.Put(ctx, &ret, upToken, key, bytes.NewReader(data), int64(len(data)), &storage.RputV2Extra{})
	if err != nil {
		return err
	}
	return nil
}

func (c Client) RefreshUrls(urlsToRefresh []string) error {
	_, err := c.cdnManager.RefreshUrls(urlsToRefresh)
	return err
}

func (c Client) GetUploadToken(putPolicy storage.PutPolicy, bucket string) string {
	return putPolicy.UploadToken(c.mac)
}

func (c Client) Delete(bucket, key string) error {
	return c.bucketManager.Delete(bucket, key)
}

// MoveFiles 批量移动或重命名文件
func (c *Client) MoveFiles(moveKeys []string, srcBucket, destBucket string, force bool) error {
	moveOps := make([]string, 0, len(moveKeys))
	for _, v := range moveKeys {
		moveOps = append(moveOps, storage.URIMove(srcBucket, v, destBucket, v, force))
	}
	rets, err := c.bucketManager.Batch(moveOps)
	if err != nil {
		if _, ok := err.(*rpc.ErrorInfo); ok {
			for _, ret := range rets {
				if ret.Code != 200 {
					return err
				}
			}
		} else {
			return err
		}
	}
	return nil
}

// CopyFiles 复制文件
func (c *Client) CopyFiles(copyKeys, destKeys []string, srcBucket, destBucket string, force bool) error {
	if len(copyKeys) != len(destKeys) {
		return errors.New("copyKeys length must equal destKeys length")
	}
	copyOps := make([]string, len(copyKeys))
	for i := 0; i < len(copyKeys); i++ {
		copyOps[i] = storage.URICopy(srcBucket, copyKeys[i], destBucket, destKeys[i], force)
	}
	rets, err := c.bucketManager.Batch(copyOps)
	if err != nil {
		if _, ok := err.(*rpc.ErrorInfo); ok {
			for _, ret := range rets {
				if ret.Code != 200 {
					return err
				}
			}
		} else {
			return err
		}
	}
	return nil
}

// ListFiles 列举文件，每次最大1000
func (c *Client) ListFiles(bucket, prefix, delimiter, marker string, limit int) ([]storage.ListItem, []string, string, bool, error) {
	entries, commonPrefixes, nextMarker, hasNext, err := c.bucketManager.ListFiles(
		bucket,
		prefix,
		delimiter,
		marker,
		limit,
	)
	return entries, commonPrefixes, nextMarker, hasNext, err
}

// MakePublicUrl 公开空间访问链接
func (c *Client) MakePublicUrl(domain, key string) string {
	return storage.MakePublicURL(domain, key)
}

func (c *Client) GetTimestampSignUrl(urlPath url.URL, secretKey string, expiration time.Duration) string {
	t := fmt.Sprintf("%x", time.Now().Add(expiration).Unix())
	encodePath := strings.ReplaceAll(url.QueryEscape(urlPath.Path), "%2F", "/")
	sign := strings.ToLower(fmt.Sprintf("%x", md5.Sum([]byte(secretKey+encodePath+t))))
	attname := strings.Trim(urlPath.Query().Get("attname"), " ")
	location := urlPath.Scheme + "://" + urlPath.Host + encodePath + fmt.Sprintf("?sign=%s&t=%s", sign, t)
	if attname != "" {
		location = urlPath.Scheme + "://" + urlPath.Host + encodePath + fmt.Sprintf("?sign=%s&t=%s&attname=%s", sign, t, attname)
	}
	return location
}
