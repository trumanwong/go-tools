package aliyun

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dcdn20180115 "github.com/alibabacloud-go/dcdn-20180115/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
)

type DCdnClient struct {
	client *dcdn20180115.Client
}

func NewDCdnClient(config *openapi.Config) (*DCdnClient, error) {
	client, err := dcdn20180115.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &DCdnClient{client: client}, nil
}

// DescribeDcdnDomainDetailWithOptions 获取指定加速域名配置的基本信息
func (c DCdnClient) DescribeDcdnDomainDetailWithOptions(req *dcdn20180115.DescribeDcdnDomainDetailRequest, runtime *util.RuntimeOptions) (*dcdn20180115.DescribeDcdnDomainDetailResponse, error) {
	return c.client.DescribeDcdnDomainDetailWithOptions(req, runtime)
}

// DescribeDcdnUserDomainsWithOptions 查询用户名下所有的DCDN域名，支持域名模糊匹配过滤和域名状态过滤
func (c DCdnClient) DescribeDcdnUserDomainsWithOptions(req *dcdn20180115.DescribeDcdnUserDomainsRequest, runtime *util.RuntimeOptions) (*dcdn20180115.DescribeDcdnUserDomainsResponse, error) {
	return c.client.DescribeDcdnUserDomainsWithOptions(req, runtime)
}

// DescribeDcdnDomainCertificateInfoWithOptions 查询指定域名证书信息
func (c DCdnClient) DescribeDcdnDomainCertificateInfoWithOptions(req *dcdn20180115.DescribeDcdnDomainCertificateInfoRequest, runtime *util.RuntimeOptions) (*dcdn20180115.DescribeDcdnDomainCertificateInfoResponse, error) {
	return c.client.DescribeDcdnDomainCertificateInfoWithOptions(req, runtime)
}

func (c DCdnClient) SetDcdnDomainSSLCertificate(req *dcdn20180115.SetDcdnDomainSSLCertificateRequest, runtime *util.RuntimeOptions) (*dcdn20180115.SetDcdnDomainSSLCertificateResponse, error) {
	return c.client.SetDcdnDomainSSLCertificateWithOptions(req, runtime)
}
