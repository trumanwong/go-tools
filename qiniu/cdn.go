package qiniu

import (
	"encoding/json"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth"
	"io"
	"net/http"
)

type CdnClient struct {
	credentials *auth.Credentials
}

func NewCdnClient(accessKey, secretKey string) *CdnClient {
	credentials := auth.New(accessKey, secretKey)
	return &CdnClient{
		credentials: credentials,
	}
}

type GetDomainsRequest struct {
	// 域名前缀
	Marker string `json:"marker,omitempty"`
	// 返回的最大域名个数。1~1000, 不填默认为 10
	Limit int `json:"limit,omitempty"`
}

type GetDomainsResponse struct {
	Marker  string    `json:"marker"`
	Domains []Domains `json:"domains"`
}
type Domains struct {
	Name               string `json:"name"`
	Type               string `json:"type"`
	Cname              string `json:"cname"`
	TestURLPath        string `json:"testURLPath"`
	Platform           string `json:"platform"`
	GeoCover           string `json:"geoCover"`
	Protocol           string `json:"protocol"`
	OperatingState     string `json:"operatingState"`
	OperatingStateDesc string `json:"operatingStateDesc"`
	CreateAt           string `json:"createAt"`
	ModifyAt           string `json:"modifyAt"`
}

func (c CdnClient) GetDomains(req *GetDomainsRequest) (*GetDomainsResponse, error) {
	body, _ := json.Marshal(req)
	resp, err := requestQiniu(http.MethodGet, "/domain", c.credentials, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %s", err)
	}
	var response GetDomainsResponse
	err = json.Unmarshal(b, &response)
	if err != nil {
		return nil, fmt.Errorf("unmarshal response body failed: %s", err)
	}
	return &response, nil

}

type ModifyHttpsConfRequest struct {
	// 域名
	Name string `json:"-"`
	// 证书id，从上传或者获取证书列表里拿到证书id
	CertID string `json:"certId"`
	// 是否强制https跳转
	ForceHttps bool `json:"forceHttps"`
	// http2功能是否启用，false为关闭，true为开启
	Http2Enable bool `json:"http2Enable"`
}

type ModifyHttpsConfResponse struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

func (c CdnClient) ModifyHttpsConf(req *ModifyHttpsConfRequest) (*ModifyHttpsConfResponse, error) {
	body, _ := json.Marshal(req)
	resp, err := requestQiniu(http.MethodPut, fmt.Sprintf("/domain/%s/httpsconf", req.Name), c.credentials, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %s", err)
	}
	var response ModifyHttpsConfResponse
	err = json.Unmarshal(b, &response)
	if err != nil {
		return nil, fmt.Errorf("unmarshal response body failed: %s", err)
	}
	return &response, nil
}
