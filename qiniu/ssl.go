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
	resp, err := requestQiniu(http.MethodPost, "/sslcert", s.credentials, body)
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
