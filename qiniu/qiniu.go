package qiniu

import (
	"bytes"
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/cdn"
	"github.com/qiniu/go-sdk/v7/sms/rpc"
	"github.com/qiniu/go-sdk/v7/storage"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const Host = "https://api.qiniu.com"

// Client is a struct that represents a Qiniu client.
// It contains a Mac object for authentication, a CdnManager for CDN operations, and a BucketManager for bucket operations.
type Client struct {
	// mac is a Mac object that contains the access key and secret key for authentication.
	mac *qbox.Mac
	// cdnManager is a CdnManager object that provides methods for CDN operations.
	cdnManager *cdn.CdnManager
	// bucketManager is a BucketManager object that provides methods for bucket operations.
	bucketManager *storage.BucketManager
}

// NewClient is a function that creates a new Qiniu client.
// It takes an access key and a secret key as parameters, and returns a pointer to a Client object.
// The function creates a new Mac object with the access key and secret key, a new CdnManager with the Mac object,
// and a new BucketManager with the Mac object and a Config object that enables HTTPS.
// The function then creates a new Client object with the Mac object, the CdnManager, and the BucketManager, and returns a pointer to it.
func NewClient(accessKey, secretKey string) *Client {
	// Create a new Mac object with the access key and secret key.
	mac := qbox.NewMac(accessKey, secretKey)
	// Create a new CdnManager with the Mac object.
	cdnManager := cdn.NewCdnManager(mac)
	// Create a new Config object that enables HTTPS.
	cfg := storage.Config{
		// 是否使用https域名进行资源管理
		UseHTTPS: true,
	}
	// Create a new BucketManager with the Mac object and the Config object.
	bucketManager := storage.NewBucketManager(mac, &cfg)
	// Create a new Client object with the Mac object, the CdnManager, and the BucketManager.
	return &Client{mac: mac, cdnManager: cdnManager, bucketManager: bucketManager}
}

type PutFileRequest struct {
	// 上传桶
	Bucket string
	// 上传文件key
	Key string
	// 上传文件内容，Data和LocalFile必须有一个不为nil
	Data []byte
	// 本地文件路径
	LocalFile *string
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
func (c Client) PutFile(ctx context.Context, req *PutFileRequest) (*PutRet, error) {
	token := c.GetUploadToken(req.PutPolicy)

	formUploader := storage.NewFormUploader(req.Config)

	ret := PutRet{}
	var err error
	if req.Data != nil {
		err = formUploader.Put(ctx, &ret, token, req.Key, bytes.NewReader(req.Data), int64(len(req.Data)), req.Extra)
	} else if req.LocalFile != nil {
		err = formUploader.PutFile(ctx, &ret, token, req.Key, *req.LocalFile, req.Extra)
	} else {
		return nil, errors.New("data and filePath can't both be nil")
	}
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

// RefreshUrls is a method of Client that refreshes the cache of the specified URLs.
// It takes a slice of strings representing the URLs to refresh as a parameter.
// The method calls the RefreshUrls method of the CdnManager object in the Client, passing the URLs to refresh.
// The method returns an error if the RefreshUrls method of the CdnManager returns an error.
func (c Client) RefreshUrls(urlsToRefresh []string) error {
	_, err := c.cdnManager.RefreshUrls(urlsToRefresh)
	return err
}

// GetUploadToken is a method of Client that gets an upload token.
// It takes a PutPolicy object as a parameter.
// The method calls the UploadToken method of the PutPolicy object, passing the Mac object in the Client.
// The method returns the upload token as a string.
func (c Client) GetUploadToken(putPolicy storage.PutPolicy) string {
	return putPolicy.UploadToken(c.mac)
}

// Delete is a method of Client that deletes a specified file.
// It takes a string representing the bucket where the file is located and a string representing the key of the file as parameters.
// The method calls the Delete method of the BucketManager object in the Client, passing the bucket and the key.
// The method returns an error if the Delete method of the BucketManager returns an error.
func (c Client) Delete(bucket, key string) error {
	return c.bucketManager.Delete(bucket, key)
}

type BatchFilesType string

const (
	BatchFilesTypeMove      BatchFilesType = "move"
	BatchFilesTypeCopy      BatchFilesType = "copy"
	BatchFilesTypeRestoreAr                = "restoreAr"
)

type BatchFilesRequest struct {
	// 操作类型，move/copy/restoreAr
	Type BatchFilesType
	// 需要移动/复制/解冻的源文件key
	SrcKeys []string
	// 移动/复制到的目标文件key
	DstKeys []string
	// 源文件所在桶
	SrcBucket *string
	// 目标文件所在桶
	DestBucket *string
	// 解冻文件天数
	AfterDay int
	// 是否强制覆盖
	Force bool
}

// BatchFiles 批量移动/复制文件
func (c Client) BatchFiles(req *BatchFilesRequest) error {
	if len(req.SrcKeys) != len(req.DstKeys) {
		return errors.New("srcKeys length must equal destKeys length")
	}
	operations := make([]string, len(req.SrcKeys))
	for i, v := range req.DstKeys {
		switch req.Type {
		case BatchFilesTypeCopy:
			operations[i] = storage.URICopy(*req.SrcBucket, req.SrcKeys[i], *req.DestBucket, v, req.Force)
		case BatchFilesTypeMove:
			operations[i] = storage.URIMove(*req.SrcBucket, req.SrcKeys[i], *req.DestBucket, v, req.Force)
		case BatchFilesTypeRestoreAr:
			operations[i] = storage.URIRestoreAr(*req.SrcBucket, req.SrcKeys[i], req.AfterDay)
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
func (c Client) ListFiles(bucket, prefix, delimiter, marker string, limit int) ([]storage.ListItem, []string, string, bool, error) {
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
func (c Client) MakePublicUrl(domain, key string) string {
	return storage.MakePublicURL(domain, key)
}

func (c Client) GetTimestampSignUrl(urlPath url.URL, secretKey string, expiration time.Duration) string {
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

// VerifyCallback 验证上传回调请求是否来自七牛
func (c Client) VerifyCallback(req *http.Request) (bool, error) {
	return qbox.VerifyCallback(c.mac, req)
}

func requestQiniu(method, apiUrl string, credentials *auth.Credentials, body []byte) (*http.Response, error) {
	var reader *bytes.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	}
	req, err := http.NewRequest(method, Host+apiUrl, reader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.ContentLength = int64(len(body))
	err = credentials.AddToken(auth.TokenQiniu, req)
	client := http.Client{}
	return client.Do(req)
}
