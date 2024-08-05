package qiniu

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth"
	"io"
	"net/http"
)

type SslClient struct {
	credentials *auth.Credentials
}

func NewSslClient(accessKey, secretKey string) *SslClient {
	credentials := auth.New(accessKey, secretKey)
	return &SslClient{
		credentials: credentials,
	}
}

type CreateSslCertRequest struct {
	// 证书名称
	Name       string `json:"name"`
	CommonName string `json:"common_name"`
	Pri        string `json:"pri"`
	Ca         string `json:"ca"`
}

type CreateSslCertResponse struct {
	// 证书ID
	CertID string `json:"certID"`
	Code   int    `json:"code"`
	Error  string `json:"error"`
}

// CreateSslCert 上传证书
func (s SslClient) CreateSslCert(ctx context.Context, request *CreateSslCertRequest) (*CreateSslCertResponse, error) {
	body, _ := json.Marshal(request)
	resp, err := requestQiniu(&Request{
		Method:      http.MethodPost,
		ApiUrl:      Host + "/sslcert",
		Body:        body,
		Credentials: s.credentials,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var response CreateSslCertResponse
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %s", err)
	}
	err = json.Unmarshal(b, &response)
	if err != nil {
		return nil, fmt.Errorf("unmarshal response body failed: %s", err)
	}
	if response.Code != 200 {
		return nil, fmt.Errorf("create ssl cert failed: %s", response.Error)
	}
	return &response, nil
}

type GetSslCertRequest struct {
	CertID string
}

type GetSslCertResponse struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
	Cert  Cert   `json:"cert"`
}

type Cert struct {
	Certid string `json:"certid"`
	// 证书名称
	Name string `json:"name"`
	// 通用名称
	CommonName string `json:"common_name"`
	// DNS域名
	DnsNames []string `json:"dnsnames"`
	// 生效时间
	NotBefore int64 `json:"not_before"`
	// 过期时间
	NotAfter int64 `json:"not_after"`
	// 证书私钥
	Pri string `json:"pri"`
	// 证书内容
	Ca               string `json:"ca"`
	UID              int    `json:"uid"`
	CreateTime       int    `json:"create_time"`
	Orderid          string `json:"orderid"`
	ProductShortName string `json:"product_short_name"`
	ProductType      string `json:"product_type"`
	CertType         string `json:"cert_type"`
	Encrypt          string `json:"encrypt"`
	EncryptParameter string `json:"encryptParameter"`
	Enable           bool   `json:"enable"`
	ChildOrderID     string `json:"child_order_id"`
	State            string `json:"state"`
	AutoRenew        bool   `json:"auto_renew"`
	Renewable        bool   `json:"renewable"`
}

// GetSslCert 用户获取单个证书的接口
func (s SslClient) GetSslCert(ctx context.Context, request *GetSslCertRequest) (*GetSslCertResponse, error) {
	resp, err := requestQiniu(&Request{
		Method:      http.MethodGet,
		ApiUrl:      Host + "/sslcert/" + request.CertID,
		Body:        nil,
		Credentials: s.credentials,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var response GetSslCertResponse
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %s", err)
	}
	err = json.Unmarshal(b, &response)
	if err != nil {
		return nil, fmt.Errorf("unmarshal response body failed: %s", err)
	}
	if response.Code != 200 {
		return nil, fmt.Errorf("get ssl cert failed: %s", response.Error)
	}
	return &response, nil
}
