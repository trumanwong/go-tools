package aliyun

import (
	cdn20180510 "github.com/alibabacloud-go/cdn-20180510/v4/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
)

type CdnClient struct {
	client *cdn20180510.Client
}

func NewCdnClient(config *openapi.Config) (*CdnClient, error) {
	client, err := cdn20180510.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &CdnClient{client: client}, nil
}

// DescribeCdnDomainDetail 查询指定加速域名的基本配置
func (c CdnClient) DescribeCdnDomainDetail(req *cdn20180510.DescribeCdnDomainDetailRequest, runtime *util.RuntimeOptions) (*cdn20180510.DescribeCdnDomainDetailResponse, error) {
	return c.client.DescribeCdnDomainDetailWithOptions(req, runtime)
}

// DescribeUserDomainsWithOptions 查询用户名下所有的域名与状态，支持域名模糊匹配过滤和域名状态过滤
func (c CdnClient) DescribeUserDomainsWithOptions(req *cdn20180510.DescribeUserDomainsRequest, runtime *util.RuntimeOptions) (*cdn20180510.DescribeUserDomainsResponse, error) {
	return c.client.DescribeUserDomainsWithOptions(req, runtime)
}

// DescribeDomainCertificateInfoWithOptions 查询域名证书信息
func (c CdnClient) DescribeDomainCertificateInfoWithOptions(req *cdn20180510.DescribeDomainCertificateInfoRequest, runtime *util.RuntimeOptions) (*cdn20180510.DescribeDomainCertificateInfoResponse, error) {
	return c.client.DescribeDomainCertificateInfoWithOptions(req, runtime)
}

func (c CdnClient) SetCdnDomainSSLCertificate(req *cdn20180510.SetCdnDomainSSLCertificateRequest, runtime *util.RuntimeOptions) (*cdn20180510.SetCdnDomainSSLCertificateResponse, error) {
	return c.client.SetCdnDomainSSLCertificateWithOptions(req, runtime)
}
