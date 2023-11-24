package aliyun

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
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
	err := a.bucket.PutObject(objectName, reader)
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
