package aliyun

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/live"
	"os"
	"testing"
)

func TestNewLiveClient(t *testing.T) {
	client, err := NewLiveClient(os.Getenv("ALIYUN_REGION_ID"), os.Getenv("ALIYUN_ACCESS_KEY_ID"), os.Getenv("ALIYUN_ACCESS_KEY_SECRET"))
	if err != nil {
		t.Errorf("NewLiveClient() error = %v", err)
		return
	}

	describeLiveUserDomainsReq := live.CreateDescribeLiveUserDomainsRequest()
	describeLiveUserDomainsResp, err := client.DescribeLiveUserDomains(describeLiveUserDomainsReq)
	if err != nil {
		t.Errorf("DescribeLiveUserDomains() error = %v", err)
		return
	}
	for _, domain := range describeLiveUserDomainsResp.Domains.PageData {
		describeLiveDomainCertificateInfoReq := live.CreateDescribeLiveDomainCertificateInfoRequest()
		describeLiveDomainCertificateInfoReq.DomainName = domain.DomainName
		describeLiveDomainCertificateInfoResp, err := client.DescribeLiveDomainCertificateInfo(describeLiveDomainCertificateInfoReq)
		if err != nil {
			t.Errorf("DescribeLiveDomainCertificateInfo() error = %v", err)
			return
		}

		for _, v := range describeLiveDomainCertificateInfoResp.CertInfos.CertInfo {
			t.Logf("[%s] cert status: %s, expire: %v", domain.DomainName, v.SSLProtocol, v.CertExpireTime)
		}
	}
}
