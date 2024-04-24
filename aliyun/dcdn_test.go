package aliyun

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dcdn20180115 "github.com/alibabacloud-go/dcdn-20180115/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"os"
	"testing"
)

func TestDCdnClient_DescribeDcdnUserDomainsWithOptions(t *testing.T) {
	client, err := NewDCdnClient(&openapi.Config{
		AccessKeyId:     tea.String(os.Getenv("AccessKeyId")),
		AccessKeySecret: tea.String(os.Getenv("AccessKeySecret")),
		Endpoint:        tea.String(os.Getenv("Endpoint")),
	})
	if err != nil {
		t.Error(err)
		return
	}
	resp, err := client.DescribeDcdnUserDomainsWithOptions(&dcdn20180115.DescribeDcdnUserDomainsRequest{
		PageSize:   tea.Int32(100),
		PageNumber: tea.Int32(1),
	}, &util.RuntimeOptions{})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(*resp.Body.TotalCount)
	for _, v := range resp.Body.Domains.PageData {
		t.Log(*v.DomainName, *v.DomainStatus, *v.SSLProtocol)
	}
}

func TestDCdnClient_DescribeDcdnDomainCertificateInfoWithOptions(t *testing.T) {
	client, err := NewDCdnClient(&openapi.Config{
		AccessKeyId:     tea.String(os.Getenv("AccessKeyId")),
		AccessKeySecret: tea.String(os.Getenv("AccessKeySecret")),
		Endpoint:        tea.String(os.Getenv("Endpoint")),
	})
	if err != nil {
		t.Error(err)
		return
	}

	resp, err := client.DescribeDcdnDomainCertificateInfoWithOptions(&dcdn20180115.DescribeDcdnDomainCertificateInfoRequest{
		DomainName: tea.String(os.Getenv("DomainName")),
	}, &util.RuntimeOptions{})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(len(resp.Body.CertInfos.CertInfo))
	t.Log(*resp.Body.CertInfos.CertInfo[0].CertExpireTime)
}

func TestDCdnClient_SetDcdnDomainSSLCertificate(t *testing.T) {
	client, err := NewDCdnClient(&openapi.Config{
		AccessKeyId:     tea.String(os.Getenv("AccessKeyId")),
		AccessKeySecret: tea.String(os.Getenv("AccessKeySecret")),
		Endpoint:        tea.String(os.Getenv("Endpoint")),
	})
	if err != nil {
		t.Error(err)
		return
	}

	resp, err := client.SetDcdnDomainSSLCertificate(&dcdn20180115.SetDcdnDomainSSLCertificateRequest{
		DomainName:  tea.String(os.Getenv("DomainName")),
		CertName:    tea.String(os.Getenv("CertName")),
		SSLProtocol: tea.String("on"),
		SSLPub:      tea.String(os.Getenv("SSLPub")),
		SSLPri:      tea.String(os.Getenv("SSLPri")),
	}, &util.RuntimeOptions{})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(*resp.Body.RequestId)
}
