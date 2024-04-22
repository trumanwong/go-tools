package aliyun

import (
	alidns20150109 "github.com/alibabacloud-go/alidns-20150109/v4/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"os"
	"testing"
)

func TestDnsClient_DescribeDomains(t *testing.T) {
	client, err := NewDnsClient(&openapi.Config{
		AccessKeyId:     tea.String(os.Getenv("AccessKeyId")),
		AccessKeySecret: tea.String(os.Getenv("AccessKeySecret")),
		Endpoint:        tea.String(os.Getenv("Endpoint")),
	})
	if err != nil {
		t.Error(err)
		return
	}
	describeDomainsRequest := &alidns20150109.DescribeDomainsRequest{}
	runtime := &util.RuntimeOptions{}
	resp, err := client.DescribeDomains(describeDomainsRequest, runtime)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(*resp.Body.TotalCount)
	for _, domain := range resp.Body.Domains.Domain {
		t.Log(*domain.DomainName, *domain.RecordCount)
	}
}
