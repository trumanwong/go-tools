package qiniu

import (
	"os"
	"testing"
)

func TestCdnClient_GetDomains(t *testing.T) {
	cdnClient := NewCdnClient(os.Getenv("QINIU_ACCESS_KEY"), os.Getenv("QINIU_SECRET_KEY"))
	resp, err := cdnClient.GetDomains(&GetDomainsRequest{
		Limit: 10,
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(resp.Marker)
	for _, v := range resp.Domains {
		t.Log(v.Name)
	}
}

func TestCdnClient_ModifyHttpsConf(t *testing.T) {
	cdnClient := NewCdnClient(os.Getenv("QINIU_ACCESS_KEY"), os.Getenv("QINIU_SECRET_KEY"))
	resp, err := cdnClient.ModifyHttpsConf(&ModifyHttpsConfRequest{
		Name:        os.Getenv("QINIU_DOMAIN_NAME"),
		CertID:      os.Getenv("QINIU_SSL_CERT_ID"),
		ForceHttps:  false,
		Http2Enable: true,
	})
	if err != nil {
		t.Error(err)
		return
	}
	if resp.Code != 200 {
		t.Error(resp.Error)
	}
}
