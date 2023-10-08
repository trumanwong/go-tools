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

type PutFileRequest struct {
	// 上传桶
	Bucket string
	// 上传文件key
	Key string
	// 上传文件内容
	Data []byte
	// 文件上传的上传策略
	PutPolicy storage.PutPolicy
	// 文件上传，资源管理等配置
	Config *storage.Config
	// 表示分片上传
	Extra *storage.PutExtra
}

type PutRet struct {
	Key          string `json:"key"`
	Hash         string `json:"hash"`
	Fsize        int64  `json:"fsize"`
	PersistentID string `json:"persistentId"`
}

// PutFile 上传文件
func (c *Client) PutFile(ctx context.Context, req *PutFileRequest) (*PutRet, error) {
	token := c.GetUploadToken(req.PutPolicy)

	formUploader := storage.NewFormUploader(req.Config)

	ret := PutRet{}
	err := formUploader.Put(ctx, &ret, token, req.Key, bytes.NewReader(req.Data), int64(len(req.Data)), req.Extra)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (c *Client) RefreshUrls(urlsToRefresh []string) error {
	_, err := c.cdnManager.RefreshUrls(urlsToRefresh)
	return err
}

// GetUploadToken 获取上传token
func (c *Client) GetUploadToken(putPolicy storage.PutPolicy) string {
	return putPolicy.UploadToken(c.mac)
}

// Delete 删除指定文件
func (c *Client) Delete(bucket, key string) error {
	return c.bucketManager.Delete(bucket, key)
}

type BatchFilesType string

const (
	BatchFilesTypeMove BatchFilesType = "move"
	BatchFilesTypeCopy BatchFilesType = "copy"
)

type BatchFilesRequest struct {
	// 操作类型，move/copy
	Type BatchFilesType
	// 需要移动/复制的源文件key
	SrcKeys []string
	// 移动/复制到的目标文件key
	DstKeys []string
	// 源文件所在桶
	SrcBucket string
	// 目标文件所在桶
	DestBucket string
	// 是否强制覆盖
	Force bool
}

// BatchFiles 批量移动/复制文件
func (c *Client) BatchFiles(req *BatchFilesRequest) error {
	if len(req.SrcKeys) != len(req.DstKeys) {
		return errors.New("srcKeys length must equal destKeys length")
	}
	operations := make([]string, len(req.SrcKeys))
	for i, v := range req.DstKeys {
		switch req.Type {
		case BatchFilesTypeCopy:
			operations[i] = storage.URICopy(req.SrcBucket, req.SrcKeys[i], req.DestBucket, v, req.Force)
		case BatchFilesTypeMove:
			operations[i] = storage.URIMove(req.SrcBucket, req.SrcKeys[i], req.DestBucket, v, req.Force)
		}
	}
	var rets []storage.BatchOpRet
	var err error
	rets, err = c.bucketManager.Batch(operations)
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
