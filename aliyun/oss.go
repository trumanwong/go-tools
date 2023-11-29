package aliyun

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
	"net/url"
	"strings"
)

type AliOss struct {
	client *oss.Client
	bucket *oss.Bucket
}

func NewAliOss(endpoint, accessKeyId, accessKeySecret, bucketName string) (*AliOss, error) {
	// 获取OSSClient实例
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		return nil, err
	}
	// 获取存储空间
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return nil, err
	}
	return &AliOss{client: client, bucket: bucket}, nil
}

// PutObject put文件至OSS
func (a AliOss) PutObject(objectName string, reader io.Reader, options ...oss.Option) error {
	err := a.bucket.PutObject(objectName, reader, options...)
	return err
}

// ListObjects 列举文件
func (a AliOss) ListObjects(prefix string) ([]oss.ObjectProperties, error) {
	fileList := make([]oss.ObjectProperties, 0)
	marker := ""
	for {
		lsRes, err := a.bucket.ListObjects(oss.Marker(marker), oss.Prefix(prefix))
		if err != nil {
			return fileList, err
		}

		for _, object := range lsRes.Objects {
			fileList = append(fileList, object)
		}

		if lsRes.IsTruncated {
			marker = lsRes.NextMarker
		} else {
			break
		}
	}
	return fileList, nil
}

// GetSignUrl 获取签名链接
func (a AliOss) GetSignUrl(objectName string) (string, error) {
	signedUrl, err := a.bucket.SignURL(objectName, oss.HTTPGet, 3600*8)
	if err != nil {
		return "", err
	}
	return signedUrl, nil
}

// DeleteObjects 删除文件
func (a AliOss) DeleteObjects(objects []string) error {
	_, err := a.bucket.DeleteObjects(objects)
	if err != nil {
		return err
	}
	return nil
}

// MakePublicURL 用来生成公开空间资源下载链接，注意该方法并不会对 key 进行 escape
func (a AliOss) MakePublicURL(domain, key string) (finalUrl string) {
	domain = strings.TrimRight(domain, "/")
	srcUrl := fmt.Sprintf("%s/%s", domain, key)
	srcUri, _ := url.Parse(srcUrl)
	finalUrl = srcUri.String()
	return
}

type SignUrlRequest struct {
	ObjectKey    string
	Method       oss.HTTPMethod
	ExpiredInSec int64
	CdnDomain    *string
	Options      []oss.Option
}

func (a AliOss) SignUrl(req *SignUrlRequest) (string, error) {
	signUrl, err := a.bucket.SignURL(req.ObjectKey, req.Method, req.ExpiredInSec, req.Options...)
	if err != nil {
		return "", err
	}
	link, err := url.Parse(signUrl)
	if err != nil {
		return "", err
	}
	if req.CdnDomain != nil {
		link.Host = *req.CdnDomain
		if strings.Contains(*req.CdnDomain, "http://") || strings.Contains(*req.CdnDomain, "https://") {
			cdnDomain, err := url.Parse(*req.CdnDomain)
			if err != nil {
				return "", err
			}
			link.Scheme = cdnDomain.Scheme
			link.Host = cdnDomain.Host
		}
	}
	return link.String(), nil
}
