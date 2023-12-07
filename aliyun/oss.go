package aliyun

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
	"net/url"
	"strings"
)

// AliOss is a struct that holds the client and bucket information for interacting with Aliyun OSS.
type AliOss struct {
	client *oss.Client // The OSS client
	bucket *oss.Bucket // The OSS bucket
}

// NewAliOss creates a new instance of AliOss.
// It takes the endpoint, accessKeyId, accessKeySecret, and bucketName as parameters.
// It returns a pointer to the created AliOss instance and an error if any occurred during the creation.
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

// PutObject uploads a file to OSS.
// It takes the objectName, reader, and options as parameters.
// It returns an error if the upload operation fails.
func (a AliOss) PutObject(objectName string, reader io.Reader, options ...oss.Option) error {
	err := a.bucket.PutObject(objectName, reader, options...)
	return err
}

// ListObjects lists the files in the OSS bucket.
// It takes a prefix as a parameter and returns a list of ObjectProperties and an error if the operation fails.
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

// GetSignUrl generates a signed URL for a file in the OSS bucket.
// It takes the objectName as a parameter and returns the signed URL and an error if the operation fails.
func (a AliOss) GetSignUrl(objectName string) (string, error) {
	signedUrl, err := a.bucket.SignURL(objectName, oss.HTTPGet, 3600*8)
	if err != nil {
		return "", err
	}
	return signedUrl, nil
}

// DeleteObjects deletes files from the OSS bucket.
// It takes a list of object names as a parameter and returns an error if the delete operation fails.
func (a AliOss) DeleteObjects(objects []string) error {
	_, err := a.bucket.DeleteObjects(objects)
	if err != nil {
		return err
	}
	return nil
}

// MakePublicURL generates a public URL for a file in the OSS bucket.
// It takes the domain and key as parameters and returns the generated URL.
func (a AliOss) MakePublicURL(domain, key string) (finalUrl string) {
	domain = strings.TrimRight(domain, "/")
	srcUrl := fmt.Sprintf("%s/%s", domain, key)
	srcUri, _ := url.Parse(srcUrl)
	finalUrl = srcUri.String()
	return
}

// SignUrlRequest is a struct that holds the information needed to generate a signed URL.
type SignUrlRequest struct {
	ObjectKey    string         // The object key
	Method       oss.HTTPMethod // The HTTP method
	ExpiredInSec int64          // The expiration time in seconds
	CdnDomain    *string        // The CDN domain
	Options      []oss.Option   // The options
}

// SignUrl generates a signed URL for a file in the OSS bucket.
// It takes a SignUrlRequest as a parameter and returns the signed URL and an error if the operation fails.
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

// Move moves a file within the OSS bucket.
// It takes the source object key, destination object key, and options as parameters.
// It returns an error if the move operation fails.
func (a AliOss) Move(srcObjectKey, destObjectKey string, options ...oss.Option) error {
	_, err := a.bucket.CopyObject(srcObjectKey, destObjectKey, options...)
	if err != nil {
		return err
	}
	return a.DeleteObjects([]string{srcObjectKey})
}
