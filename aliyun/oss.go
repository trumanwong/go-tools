package aliyun

import (
	"crypto"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/trumanwong/go-tools/encoding"
	"hash"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// AliOss is a struct that holds the client and bucket information for interacting with Aliyun OSS.
type AliOss struct {
	client          *oss.Client // The OSS client
	bucket          *oss.Bucket // The OSS bucket
	accessKeyId     string
	accessKeySecret string
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
	return &AliOss{
		client:          client,
		bucket:          bucket,
		accessKeyId:     accessKeyId,
		accessKeySecret: accessKeySecret,
	}, nil
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

func (a AliOss) getGmtIso8601(expireEnd int64) string {
	var tokenExpire = time.Unix(expireEnd, 0).UTC().Format("2006-01-02T15:04:05Z")
	return tokenExpire
}

type ConfigStruct struct {
	Expiration string     `json:"expiration"`
	Conditions [][]string `json:"conditions"`
}

type PolicyToken struct {
	AccessKeyId string `json:"accessid"`
	Host        string `json:"host"`
	Expire      int64  `json:"expire"`
	Signature   string `json:"signature"`
	Policy      string `json:"policy"`
	Directory   string `json:"dir"`
	Callback    string `json:"callback"`
}

type CallbackParam struct {
	CallbackUrl      string `json:"callbackUrl"`
	CallbackBody     string `json:"callbackBody"`
	CallbackBodyType string `json:"callbackBodyType"`
}

// GetPolicyToken generates a policy token for uploading files to the OSS bucket.
// expireSeconds Token有效期，单位秒.
// ossHost OSS域名，"http://bucket-name.oss-cn-hangzhou.aliyuncs.com"
// directoryPrefix 上传目录前缀.
// callBackUrl 上传回调地址.
func (a AliOss) GetPolicyToken(expireSeconds int64, ossHost string, directoryPrefix string, callBackUrl string) ([]byte, error) {
	var policyToken PolicyToken
	policyToken.AccessKeyId = a.accessKeyId
	now := time.Now().Unix()
	policyToken.Host = ossHost
	policyToken.Expire = now + expireSeconds
	var tokenExpire = a.getGmtIso8601(policyToken.Expire)

	//create post policy json
	var config ConfigStruct
	config.Expiration = tokenExpire
	var condition []string
	condition = append(condition, "starts-with")
	condition = append(condition, "$key")
	condition = append(condition, directoryPrefix)
	config.Conditions = append(config.Conditions, condition)

	//calucate signature
	result, err := json.Marshal(config)
	debyte := base64.StdEncoding.EncodeToString(result)
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(a.accessKeySecret))
	io.WriteString(h, debyte)
	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))

	var callbackParam CallbackParam
	callbackParam.CallbackUrl = callBackUrl
	callbackParam.CallbackBody = "filename=${object}&size=${size}&mimeType=${mimeType}&height=${imageInfo.height}&width=${imageInfo.width}"
	callbackParam.CallbackBodyType = "application/x-www-form-urlencoded"
	callbackStr, err := json.Marshal(callbackParam)
	if err != nil {
		return nil, fmt.Errorf("callback json err: %s", err.Error())
	}

	policyToken.Signature = signedStr
	policyToken.Directory = directoryPrefix
	policyToken.Policy = debyte
	policyToken.Callback = base64.StdEncoding.EncodeToString(callbackStr)
	response, err := json.Marshal(policyToken)
	if err != nil {
		return nil, fmt.Errorf("json err: %s", err)
	}
	return response, nil
}

// VerifyCallbackSignature 校验回调请求的签名是否合法
func (a AliOss) VerifyCallbackSignature(r *http.Request) error {
	// Get PublicKey bytes
	bytePublicKey, err := a.getPublicKey(r)
	if err != nil {
		return err
	}
	// Get Authorization bytes : decode from Base64String
	byteAuthorization, err := a.getAuthorization(r)
	if err != nil {
		return err
	}
	// Get MD5 bytes from Newly Constructed Authorization String.
	byteMD5, err := a.getMD5FromNewAuthString(r)
	if err != nil {
		return err
	}
	return a.verifySignature(bytePublicKey, byteMD5, byteAuthorization)
}

// getPublicKey : Get PublicKey bytes from Request.URL
func (a AliOss) getPublicKey(r *http.Request) ([]byte, error) {
	var bytePublicKey []byte
	// get PublicKey URL
	publicKeyURLBase64 := r.Header.Get("x-oss-pub-key-url")
	if publicKeyURLBase64 == "" {
		return nil, errors.New("no x-oss-pub-key-url field in Request header ")
	}
	publicKeyURL, _ := base64.StdEncoding.DecodeString(publicKeyURLBase64)
	responsePublicKeyURL, err := http.Get(string(publicKeyURL))
	if err != nil {
		return nil, fmt.Errorf("Get PublicKey Content from URL failed : %s \n", err.Error())
	}
	bytePublicKey, err = io.ReadAll(responsePublicKeyURL.Body)
	if err != nil {
		return nil, fmt.Errorf("Read PublicKey Content from URL failed : %s \n", err.Error())
	}
	defer responsePublicKeyURL.Body.Close()
	return bytePublicKey, nil
}

// getAuthorization : decode from Base64String
func (a AliOss) getAuthorization(r *http.Request) ([]byte, error) {
	var byteAuthorization []byte
	// Get Authorization bytes : decode from Base64String
	strAuthorizationBase64 := r.Header.Get("authorization")
	if strAuthorizationBase64 == "" {
		return byteAuthorization, fmt.Errorf("Failed to get authorization field from request header. ")
	}
	byteAuthorization, _ = base64.StdEncoding.DecodeString(strAuthorizationBase64)
	return byteAuthorization, nil
}

// getMD5FromNewAuthString : Get MD5 bytes from Newly Constructed Authorization String.
func (a AliOss) getMD5FromNewAuthString(r *http.Request) ([]byte, error) {
	var byteMD5 []byte
	// Construct the New Auth String from URI+Query+Body
	bodyContent, err := io.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		return byteMD5, fmt.Errorf("Read Request Body failed : %s \n", err.Error())
	}
	strCallbackBody := string(bodyContent)
	// fmt.Printf("r.URL.RawPath={%s}, r.URL.Query()={%s}, strCallbackBody={%s}\n", r.URL.RawPath, r.URL.Query(), strCallbackBody)
	strURLPathDecode, errUnescape := encoding.UnescapePath(r.URL.Path, encoding.EncodePathSegment) //url.PathUnescape(r.URL.Path) for Golang v1.8.2+
	if errUnescape != nil {
		fmt.Printf("url.PathUnescape failed : URL.Path=%s, error=%s \n", r.URL.Path, err.Error())
		return byteMD5, errUnescape
	}

	// Generate New Auth String prepare for MD5
	strAuth := ""
	if r.URL.RawQuery == "" {
		strAuth = fmt.Sprintf("%s\n%s", strURLPathDecode, strCallbackBody)
	} else {
		strAuth = fmt.Sprintf("%s?%s\n%s", strURLPathDecode, r.URL.RawQuery, strCallbackBody)
	}
	// fmt.Printf("NewlyConstructedAuthString={%s}\n", strAuth)

	// Generate MD5 from the New Auth String
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(strAuth))
	byteMD5 = md5Ctx.Sum(nil)

	return byteMD5, nil
}

/*  VerifySignature
*   VerifySignature需要三个重要的数据信息来进行签名验证： 1>获取公钥PublicKey;  2>生成新的MD5鉴权串;  3>解码Request携带的鉴权串;
*   1>获取公钥PublicKey : 从RequestHeader的"x-oss-pub-key-url"字段中获取 URL, 读取URL链接的包含的公钥内容， 进行解码解析， 将其作为rsa.VerifyPKCS1v15的入参。
*   2>生成新的MD5鉴权串 : 把Request中的url中的path部分进行urldecode， 加上url的query部分， 再加上body， 组合之后进行MD5编码， 得到MD5鉴权字节串。
*   3>解码Request携带的鉴权串 ： 获取RequestHeader的"authorization"字段， 对其进行Base64解码，作为签名验证的鉴权对比串。
*   rsa.VerifyPKCS1v15进行签名验证，返回验证结果。
* */
func (a AliOss) verifySignature(bytePublicKey []byte, byteMd5 []byte, authorization []byte) error {
	pubBlock, _ := pem.Decode(bytePublicKey)
	if pubBlock == nil {
		return fmt.Errorf("failed to parse PEM block containing the public key")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(pubBlock.Bytes)
	if (pubInterface == nil) || (err != nil) {
		return fmt.Errorf("x509.ParsePKIXPublicKey(publicKey) failed : %s \n", err.Error())
	}
	pub := pubInterface.(*rsa.PublicKey)

	errorVerifyPKCS1v15 := rsa.VerifyPKCS1v15(pub, crypto.MD5, byteMd5, authorization)
	if errorVerifyPKCS1v15 != nil {
		//printByteArray(byteMd5, "AuthMd5(fromNewAuthString)")
		//printByteArray(bytePublicKey, "PublicKeyBase64")
		//printByteArray(authorization, "AuthorizationFromRequest")
		return fmt.Errorf("\nSignature Verification is Failed : %s \n", errorVerifyPKCS1v15.Error())
	}

	return nil
}

// PutObjectFromFile uploads a file to OSS from a local file.
func (a AliOss) PutObjectFromFile(objectName, filePath string, options ...oss.Option) error {
	err := a.bucket.PutObjectFromFile(objectName, filePath, options...)
	return err
}
