package volcenginecloud

import (
	"bytes"
	"context"
	"crypto"
	"crypto/md5"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"

	"github.com/volcengine/ve-tos-golang-sdk/v2/tos"
)

type OssClient struct {
	client *tos.ClientV2
}

func NewOssClient(accessKey, secretKey, endpoint, region string) (*OssClient, error) {
	client, err := tos.NewClientV2(
		endpoint,
		tos.WithCredentials(tos.NewStaticCredentials(accessKey, secretKey)),
		tos.WithRegion(region),
	)
	if err != nil {
		return nil, err
	}
	return &OssClient{client: client}, nil
}

// HeadObject 查询对象元数据信息（判断对象是否存在）
func (c OssClient) HeadObject(ctx context.Context, req *tos.HeadObjectV2Input) (*tos.HeadObjectV2Output, error) {
	return c.client.HeadObjectV2(ctx, req)
}

// DeleteObject 删除对象
func (c OssClient) DeleteObject(ctx context.Context, req *tos.DeleteObjectV2Input) (*tos.DeleteObjectV2Output, error) {
	return c.client.DeleteObjectV2(ctx, req)
}

// DeleteObject 批量删除对象
func (c OssClient) DeleteMultiObjects(ctx context.Context, req *tos.DeleteMultiObjectsInput) (*tos.DeleteMultiObjectsOutput, error) {
	return c.client.DeleteMultiObjects(ctx, req)
}

// RenameObject 重命名对象
func (c OssClient) RenameObject(ctx context.Context, req *tos.RenameObjectInput) (*tos.RenameObjectOutput, error) {
	return c.client.RenameObject(ctx, req)
}

// RestoreObject 恢复对象
func (c OssClient) RestoreObject(ctx context.Context, req *tos.RestoreObjectInput) (*tos.RestoreObjectOutput, error) {
	return c.client.RestoreObject(ctx, req)
}

// ListObjectsType2 列举对象
func (c OssClient) ListObjectsType2(ctx context.Context, req *tos.ListObjectsType2Input) (*tos.ListObjectsType2Output, error) {
	return c.client.ListObjectsType2(ctx, req)
}

// PutObject 普通上传
func (c OssClient) PutObject(ctx context.Context, req *tos.PutObjectV2Input) (*tos.PutObjectV2Output, error) {
	return c.client.PutObjectV2(ctx, req)
}

// AppendObject 追加上传
func (c OssClient) AppendObject(ctx context.Context, req *tos.AppendObjectV2Input) (*tos.AppendObjectV2Output, error) {
	return c.client.AppendObjectV2(ctx, req)
}

// PreSignedPostSignature Post表单预签名
func (c OssClient) PreSignedPostSignature(ctx context.Context, req *tos.PreSingedPostSignatureInput) (*tos.PreSingedPostSignatureOutput, error) {
	return c.client.PreSignedPostSignature(ctx, req)
}

// PreSignedURL 普通预签名上传
func (c OssClient) PreSignedURL(req *tos.PreSignedURLInput) (*tos.PreSignedURLOutput, error) {
	return c.client.PreSignedURL(req)
}

const TosAuthorization = "Authorization"
const XTosPublicKeyUrl = "x-tos-pub-key-url"

func calcStringToSign(req *http.Request, body []byte) []byte {
	buf := bytes.NewBuffer(nil)
	buf.WriteString(req.URL.Path)

	query := req.URL.Query()
	if len(query) > 0 {
		buf.WriteByte('?')
		keys := make([]string, 0, len(query))
		for key := range query {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		for _, key := range keys {
			values := query[key]
			for _, value := range values {
				buf.WriteString(key)
				buf.WriteByte('=')
				buf.WriteString(value)
				buf.WriteByte('&')
			}
		}

		buf.Truncate(buf.Len() - 1)
	}

	buf.WriteByte('\n')
	buf.Write(body)
	return buf.Bytes()
}

func getSignature(r *http.Request) ([]byte, error) {
	authorization := r.Header.Get(TosAuthorization)
	if authorization == "" {
		return nil, nil
	}

	return base64.StdEncoding.DecodeString(authorization)
}

func getPublicKey(r *http.Request) ([]byte, error) {
	var bytePublicKey []byte

	publicKeyURLBase64 := r.Header.Get(XTosPublicKeyUrl)
	if publicKeyURLBase64 == "" {
		return bytePublicKey, errors.New("no x-tos-pub-key-url field in request header ")
	}
	publicKeyURL, _ := base64.StdEncoding.DecodeString(publicKeyURLBase64)

	resp, err := http.Get(string(publicKeyURL))
	if err != nil {
		return bytePublicKey, err
	}

	bytePublicKey, err = io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return bytePublicKey, err
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, errors.New("get public key failed")
	}

	return bytePublicKey, nil
}

func getContentAndSignMD5(r *http.Request) ([]byte, []byte, error) {
	bodyContent, err := io.ReadAll(r.Body)
	_ = r.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	stringToSign := calcStringToSign(r, bodyContent)

	authMd5 := md5.New()
	authMd5.Write(stringToSign)
	signMd5 := authMd5.Sum(nil)
	return signMd5, bodyContent, nil
}

func verifySignature(publicKey []byte, signMd5 []byte, sign []byte) error {
	pubBlock, _ := pem.Decode(publicKey)
	if pubBlock == nil {
		return errors.New("invalid public key pem")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(pubBlock.Bytes)
	if err != nil || pubInterface == nil {
		return errors.New("invalid public key pem")
	}

	pub := pubInterface.(*rsa.PublicKey)

	err = rsa.VerifyPKCS1v15(pub, crypto.MD5, signMd5, sign)
	if err != nil {
		return err
	}
	return nil
}

// GetUploadCallbackParams 校验上传回调并返回body
func GetUploadCallbackParams(r *http.Request) ([]byte, error) {
	// 获取签名
	sign, err := getSignature(r)
	if err != nil {
		return nil, fmt.Errorf("get callback signature err: %s", err.Error())
	}

	// 获取公共签名，该结果可以缓存
	publicKey, err := getPublicKey(r)
	if err != nil {
		return nil, fmt.Errorf("get public key err: %s", err.Error())
	}

	// 读取 body 并计算签名
	signMd5, body, err := getContentAndSignMD5(r)
	if err != nil {
		return nil, fmt.Errorf("get content and sign md5 err: %s", err.Error())
	}

	// 验证签名
	if err = verifySignature(publicKey, signMd5, sign); err != nil {
		return nil, fmt.Errorf("verify signature err: %s", err.Error())
	}

	return body, nil
}
