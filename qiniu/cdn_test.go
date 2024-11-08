package qiniu

import (
	"fmt"
	"github.com/trumanwong/go-tools/crawler"
	"github.com/trumanwong/go-tools/helper"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
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
	t.Log(resp.IPACL.IPACLValues)
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

func TestCdnClient_GetTopTrafficIp(t *testing.T) {
	cdnClient := NewCdnClient(os.Getenv("QINIU_ACCESS_KEY"), os.Getenv("QINIU_SECRET_KEY"))
	resp, err := cdnClient.GetTopTrafficIp(&GetTopTrafficIpRequest{
		StartDate: time.Now().Format("2006-01-02"),
		EndDate:   time.Now().Format("2006-01-02"),
		Region:    "global",
		Domains:   []string{os.Getenv("QINIU_DOMAIN_NAME")},
	})
	if err != nil {
		t.Error(err)
		return
	}
	for i, v := range resp.Data.Ips {
		t.Log(v, helper.FormatByte(float64(resp.Data.Traffic[i]), 1024))
	}
}

func TestCdnClient_GetTopCountIp(t *testing.T) {
	cdnClient := NewCdnClient(os.Getenv("QINIU_ACCESS_KEY"), os.Getenv("QINIU_SECRET_KEY"))
	resp, err := cdnClient.GetTopCountIp(&GetTopCountIpRequest{
		StartDate: time.Now().Format("2006-01-02"),
		EndDate:   time.Now().Format("2006-01-02"),
		Region:    "global",
		Domains:   []string{os.Getenv("QINIU_DOMAIN_NAME")},
	})
	if err != nil {
		t.Error(err)
		return
	}
	for i, v := range resp.Data.Ips {
		t.Log(v, resp.Data.Count[i])
	}
}

func TestCdnClient_UpdateIpACL(t *testing.T) {
	cdnClient := NewCdnClient(os.Getenv("QINIU_ACCESS_KEY"), os.Getenv("QINIU_SECRET_KEY"))
	_, err := cdnClient.UpdateIpACL(&UpdateIpACLRequest{
		Domain: os.Getenv("QINIU_DOMAIN_NAME"),
		IpAcl: IPACL{
			IPACLType: "black",
			IPACLValues: []string{
				"39.144.0.75",
			},
		},
	})
	if err != nil {
		t.Error(err)
		return
	}
}

func TestCdnClient_GetCdnLogList(t *testing.T) {
	cdnClient := NewCdnClient(os.Getenv("QINIU_ACCESS_KEY"), os.Getenv("QINIU_SECRET_KEY"))
	logListResult, err := cdnClient.GetCdnLogList(time.Now().AddDate(0, 0, -1).Format(time.DateOnly), []string{os.Getenv("QINIU_DOMAIN")})
	if err != nil {
		t.Error(err)
		return
	}
	for domain, logList := range logListResult.Data {
		for _, v := range logList {
			u, _ := url.Parse(v.URL)
			savePath := fmt.Sprintf("temp/%s/%s", domain, filepath.Base(u.Path))
			_, err = helper.DownloadFile(&crawler.Request{
				Url:    v.URL,
				Method: http.MethodGet,
			}, savePath, false)
			if err != nil {
				t.Error(err)
				break
			}
		}
	}
}

func TestCdnClient_AnalyzeCdnAccessLog(t *testing.T) {
	cdnClient := NewCdnClient(os.Getenv("QINIU_ACCESS_KEY"), os.Getenv("QINIU_SECRET_KEY"))
	err := cdnClient.AnalyzeCdnAccessLog(os.Getenv("QINIU_LOG_PATH"), func(accessLog interface{}) error {
		return nil
	})
	if err != nil {
		t.Error(err)
	}
}
