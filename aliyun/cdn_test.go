package aliyun

import (
	"fmt"
	cdn20180510 "github.com/alibabacloud-go/cdn-20180510/v4/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/trumanwong/go-tools/crawler"
	"github.com/trumanwong/go-tools/helper"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

func TestCdnClient_DescribeUserDomainsWithOptions(t *testing.T) {
	client, err := NewCdnClient(&openapi.Config{
		AccessKeyId:     tea.String(os.Getenv("AccessKeyId")),
		AccessKeySecret: tea.String(os.Getenv("AccessKeySecret")),
		Endpoint:        tea.String(os.Getenv("Endpoint")),
	})
	if err != nil {
		t.Error(err)
		return
	}
	resp, err := client.DescribeUserDomainsWithOptions(&cdn20180510.DescribeUserDomainsRequest{
		PageSize:   tea.Int32(100),
		PageNumber: tea.Int32(1),
	}, &util.RuntimeOptions{})
	if err != nil {
		t.Error(err)
		return
	}
	for _, v := range resp.Body.Domains.PageData {
		log.Println(*v.DomainName, *v.DomainStatus, *v.SslProtocol)
	}
}

func TestCdnClient_DescribeDomainCertificateInfoWithOptions(t *testing.T) {
	client, err := NewCdnClient(&openapi.Config{
		AccessKeyId:     tea.String(os.Getenv("AccessKeyId")),
		AccessKeySecret: tea.String(os.Getenv("AccessKeySecret")),
		Endpoint:        tea.String(os.Getenv("Endpoint")),
	})
	if err != nil {
		t.Error(err)
		return
	}

	resp, err := client.DescribeDomainCertificateInfoWithOptions(&cdn20180510.DescribeDomainCertificateInfoRequest{
		DomainName: tea.String(os.Getenv("DomainName")),
	}, &util.RuntimeOptions{})
	if err != nil {
		t.Error(err)
		return
	}
	log.Println(len(resp.Body.CertInfos.CertInfo))
	log.Println(*resp.Body.CertInfos.CertInfo[0].CertExpireTime < time.Now().AddDate(0, 0, 60).Format(time.RFC3339))
}

func TestCdnClient_SetCdnDomainSSLCertificate(t *testing.T) {
	client, err := NewCdnClient(&openapi.Config{
		AccessKeyId:     tea.String(os.Getenv("AccessKeyId")),
		AccessKeySecret: tea.String(os.Getenv("AccessKeySecret")),
		Endpoint:        tea.String(os.Getenv("Endpoint")),
	})
	if err != nil {
		t.Error(err)
		return
	}
	resp, err := client.SetCdnDomainSSLCertificate(&cdn20180510.SetCdnDomainSSLCertificateRequest{
		DomainName:  tea.String(os.Getenv("DomainName")),
		SSLProtocol: tea.String("on"),
		SSLPub:      tea.String(os.Getenv("SSL_CERT")),
		SSLPri:      tea.String(os.Getenv("SSL_KEY")),
	}, &util.RuntimeOptions{})
	if err != nil {
		t.Error(err)
		return
	}
	log.Println(*resp.Body.RequestId)
}

func TestCdnClient_DescribeDomainUsageData(t *testing.T) {
	client, err := NewCdnClient(&openapi.Config{
		AccessKeyId:     tea.String(os.Getenv("AccessKeyId")),
		AccessKeySecret: tea.String(os.Getenv("AccessKeySecret")),
		Endpoint:        tea.String(os.Getenv("Endpoint")),
	})
	if err != nil {
		t.Error(err)
		return
	}

	resp, err := client.DescribeDomainUsageData(&cdn20180510.DescribeDomainUsageDataRequest{
		DomainName: tea.String(os.Getenv("DomainName")),
		StartTime:  tea.String(time.Now().UTC().AddDate(0, 0, -30).Format("2006-01-02T15:04:05Z")),
		EndTime:    tea.String(time.Now().UTC().Format("2006-01-02T15:04:05Z")),
		Field:      tea.String("traf"),
		Interval:   tea.String("86400"),
	}, &util.RuntimeOptions{})
	if err != nil {
		t.Error(err)
		return
	}
	var total int64
	for _, v := range resp.Body.UsageDataPerInterval.DataModule {
		flux, _ := strconv.ParseInt(*v.Value, 10, 64)
		total += flux
	}
	log.Println(helper.FormatByte(float64(total), 1024))
}

func TestCdnClient_DescribeCdnDomainLogs(t *testing.T) {
	client, err := NewCdnClient(&openapi.Config{
		AccessKeyId:     tea.String(os.Getenv("AccessKeyId")),
		AccessKeySecret: tea.String(os.Getenv("AccessKeySecret")),
	})
	if err != nil {
		t.Error(err)
		return
	}

	intervalTime := helper.GetIntervalTime(&helper.GetIntervalTimeRequest{
		Time: time.Now(),
		Type: helper.Day,
		Num:  -1,
	})
	log.Println(intervalTime.StartAt.UTC().Format("2006-01-02T15:04:05Z"))
	log.Println(intervalTime.EndAt.UTC().Format("2006-01-02T15:04:05Z"))
	resp, err := client.DescribeCdnDomainLogs(&cdn20180510.DescribeCdnDomainLogsRequest{
		DomainName: tea.String(os.Getenv("DomainName")),
		StartTime:  tea.String(intervalTime.StartAt.UTC().Format("2006-01-02T15:04:05Z")),
		EndTime:    tea.String(intervalTime.EndAt.UTC().Format("2006-01-02T15:04:05Z")),
	}, &util.RuntimeOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	for _, v := range resp.Body.DomainLogDetails.DomainLogDetail {
		log.Println(*v.LogCount)
		for _, detail := range v.LogInfos.LogInfoDetail {
			log.Println(*detail.LogSize, *detail.LogPath)
			link := fmt.Sprintf("https://%s", *detail.LogPath)
			u, _ := url.Parse(link)
			fileName := filepath.Base(u.Path)
			_, err = helper.DownloadFile(&crawler.Request{
				Url: link,
			}, fmt.Sprintf("temp/%s/%s", time.Now().Format(time.DateOnly), fileName), true)
			if err != nil {
				t.Error(err)
				return
			}
		}
	}
}

func TestCdnClient_AnalyzeCdnAccessLog(t *testing.T) {
	client, err := NewCdnClient(&openapi.Config{
		AccessKeyId:     tea.String(os.Getenv("AccessKeyId")),
		AccessKeySecret: tea.String(os.Getenv("AccessKeySecret")),
	})
	if err != nil {
		t.Error(err)
		return
	}
	err = client.AnalyzeCdnAccessLog(os.Getenv("LogPath"), func(accessLog *CdnAccessLog) error {
		return nil
	})
	if err != nil {
		t.Error(err)
	}
}
