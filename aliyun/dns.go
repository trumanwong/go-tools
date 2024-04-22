package aliyun

import (
	alidns20150109 "github.com/alibabacloud-go/alidns-20150109/v4/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
)

type DnsClient struct {
	client *alidns20150109.Client
}

func NewDnsClient(config *openapi.Config) (*DnsClient, error) {
	client, err := alidns20150109.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &DnsClient{client: client}, nil
}

// DescribeDomains 获取域名列表
func (d DnsClient) DescribeDomains(describeDomainsRequest *alidns20150109.DescribeDomainsRequest, runtime *util.RuntimeOptions) (*alidns20150109.DescribeDomainsResponse, error) {
	return d.client.DescribeDomainsWithOptions(describeDomainsRequest, runtime)
}
