package qiniu

import (
	"github.com/trumanwong/go-tools/helper"
	"log"
	"os"
	"testing"
	"time"
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

func TestCdnClient_GetDomain(t *testing.T) {
	cdnClient := NewCdnClient(os.Getenv("QINIU_ACCESS_KEY"), os.Getenv("QINIU_SECRET_KEY"))
	resp, err := cdnClient.GetDomain(&GetDomainRequest{
		Name: os.Getenv("QINIU_DOMAIN_NAME"),
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(resp.HTTPS.CertID)
}

func TestCdnClient_Flux(t *testing.T) {
	cdnClient := NewCdnClient(os.Getenv("QINIU_ACCESS_KEY"), os.Getenv("QINIU_SECRET_KEY"))
	resp, err := cdnClient.Flux(&FluxRequest{
		StartDate:   time.Now().AddDate(0, 0, -30).Format("2006-01-02"),
		EndDate:     time.Now().Format("2006-01-02"),
		Granularity: "day",
		Domains:     os.Getenv("QINIU_DOMAIN_NAME"),
	})
	if err != nil {
		t.Error(err)
		return
	}
	var total int64
	for _, v := range resp.Data {
		for _, vv := range v.China {
			total += vv
		}

		for _, vv := range v.Oversea {
			total += vv
		}
	}
	log.Println(helper.FormatByte(float64(total), 1024))
}
