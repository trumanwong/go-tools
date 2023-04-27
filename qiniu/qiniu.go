package qiniu

import (
	"bytes"
	"context"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/cdn"
	"github.com/qiniu/go-sdk/v7/storage"
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

func (c Client) GetUploadToken(bucket string) string {
	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	return putPolicy.UploadToken(c.mac)
}

func (c Client) Delete(bucket, key string) error {
	return c.bucketManager.Delete(bucket, key)
}
