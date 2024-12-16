package volcenginecloud

import (
	"fmt"
	"github.com/trumanwong/go-tools/crawler"
	"github.com/trumanwong/go-tools/helper"
	"github.com/trumanwong/go-tools/trans"
	"github.com/volcengine/volc-sdk-golang/service/cdn"
	"log"
	"os"
	"testing"
	"time"
)

func TestCdnClient_ListDomains(t *testing.T) {
	cdnClient := NewCdnClient(os.Getenv("VOLC_ACCESS_KEY"), os.Getenv("VOLC_SECRET_KEY"))
	resp, err := cdnClient.ListDomains(&cdn.ListCdnDomainsRequest{
		PageNum:  trans.Int64(1),
		PageSize: trans.Int64(100),
	})
	if err != nil {
		t.Error(err)
		return
	}
	log.Println(resp.Result.Total)
	for _, item := range resp.Result.Data {
		response, err := cdnClient.DescribeCdnConfig(&cdn.DescribeCdnConfigRequest{
			Domain: item.Domain,
		})
		if err != nil {
			t.Error(err)
			return
		}
		log.Println(*response.Result.DomainConfig.Status)
		log.Println(*response.Result.DomainConfig.HTTPS.Switch)
		log.Println(*response.Result.DomainConfig.HTTPS.CertInfo.ExpireTime)
	}
}

func TestCdnClient_DescribeCdnData(t *testing.T) {
	cdnClient := NewCdnClient(os.Getenv("VOLC_ACCESS_KEY"), os.Getenv("VOLC_SECRET_KEY"))
	resp, err := cdnClient.DescribeCdnData(&cdn.DescribeCdnDataRequest{
		StartTime: time.Now().AddDate(0, 0, -30).Unix(),
		EndTime:   time.Now().Unix(),
		Interval:  trans.String("day"),
		Metric:    "flux",
		Domain:    trans.String(os.Getenv("VOLC_DOMAIN")),
	})
	if err != nil {
		t.Error(err)
		return
	}
	for _, resource := range resp.Result.Resources {
		log.Println(resource.Name)
		for _, metric := range resource.Metrics {
			for _, val := range metric.Values {
				log.Println(time.Unix(val.Timestamp, 0).Format(time.DateOnly), helper.FormatByte(val.Value, 1024))
			}
		}
	}
}

func TestCdnClient_DescribeCdnConfig(t *testing.T) {
	cdnClient := NewCdnClient(os.Getenv("VOLC_ACCESS_KEY"), os.Getenv("VOLC_SECRET_KEY"))
	resp, err := cdnClient.DescribeCdnConfig(&cdn.DescribeCdnConfigRequest{
		Domain: os.Getenv("VOLC_DOMAIN"),
	})
	if err != nil {
		t.Error(err)
		return
	}
	if resp.Result.DomainConfig.IpAccessRule != nil {
		if resp.Result.DomainConfig.IpAccessRule.SharedConfig != nil {
			log.Println(*resp.Result.DomainConfig.IpAccessRule.SharedConfig.ConfigName)
		}
		log.Println(*resp.Result.DomainConfig.IpAccessRule.RuleType)
		log.Println(*resp.Result.DomainConfig.IpAccessRule.Switch)
		log.Println(resp.Result.DomainConfig.IpAccessRule.Ip)
	}
}

func TestCdnClient_UpdateCdnConfig(t *testing.T) {
	cdnClient := NewCdnClient(os.Getenv("VOLC_ACCESS_KEY"), os.Getenv("VOLC_SECRET_KEY"))
	resp, err := cdnClient.UpdateCdnConfig(&cdn.UpdateCdnConfigRequest{
		Domain: trans.String(os.Getenv("VOLC_DOMAIN")),
		IpAccessRule: &cdn.IpAccessRule{
			Switch:   trans.Bool(true),
			RuleType: trans.String("deny"),
			Ip:       []string{"111.199.100.167", "183.197.158.221", "111.199.229.206"},
		},
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(resp.ResponseMetadata.RequestId)
}

func TestCdnClient_DescribeDistrictSummary(t *testing.T) {
	cdnClient := NewCdnClient(os.Getenv("VOLC_ACCESS_KEY"), os.Getenv("VOLC_SECRET_KEY"))
	resp, err := cdnClient.DescribeDistrictSummary(&cdn.DescribeDistrictSummaryRequest{
		Metric:    "requests",
		StartTime: time.Now().AddDate(0, 0, -1).Unix(),
		EndTime:   time.Now().Unix(),
	})

	if err != nil {
		t.Error(err)
		return
	}

	for _, item := range resp.Result.MetricDataList {
		log.Println(item.Metric, item.Value)
	}
}

func TestCdnClient_DescribeStatisticalRanking(t *testing.T) {
	cdnClient := NewCdnClient(os.Getenv("VOLC_ACCESS_KEY"), os.Getenv("VOLC_SECRET_KEY"))
	resp, err := cdnClient.DescribeStatisticalRanking(&cdn.DescribeStatisticalRankingRequest{
		Domain:    os.Getenv("VOLC_DOMAIN"),
		Item:      "clientip",
		Metric:    "requests",
		StartTime: time.Now().AddDate(0, 0, -1).Unix(),
		EndTime:   time.Now().Unix(),
	})
	if err != nil {
		t.Error(err)
		return
	}
	for _, item := range resp.Result.RankingDataList {
		log.Println(item.ItemKey, item.Value)
	}
}

func TestCdnClient_DescribeCdnAccessLog(t *testing.T) {
	cdnClient := NewCdnClient(os.Getenv("VOLC_ACCESS_KEY"), os.Getenv("VOLC_SECRET_KEY"))
	intervalTime := helper.GetIntervalTime(&helper.GetIntervalTimeRequest{
		Time: time.Now(),
		Type: helper.Day,
		Num:  -1,
	})
	resp, err := cdnClient.DescribeCdnAccessLog(&cdn.DescribeCdnAccessLogRequest{
		Domain:    os.Getenv("VOLC_DOMAIN"),
		StartTime: intervalTime.StartAt.Unix(),
		EndTime:   intervalTime.EndAt.Unix(),
		PageSize:  trans.Int64(100),
	})
	if err != nil {
		t.Error(err)
		return
	}
	for _, item := range resp.Result.DomainLogDetails {
		_, err = helper.DownloadFile(&crawler.Request{
			Url: item.LogPath,
		}, fmt.Sprintf("temp/%s_%d_%d.gz", os.Getenv("VOLC_DOMAIN"), item.StartTime, item.EndTime), true)
		if err != nil {
			t.Error(err)
			return
		}
	}
}

func TestCdnClient_AnalyzeCdnAccessLog(t *testing.T) {
	cdnClient := NewCdnClient(os.Getenv("VOLC_ACCESS_KEY"), os.Getenv("VOLC_SECRET_KEY"))
	err := cdnClient.AnalyzeCdnAccessLog("temp", func(accessLog interface{}) error {
		log.Println(accessLog)
		return nil
	})
	if err != nil {
		t.Error(err)
	}
}

func TestGenTypeAUrl(t *testing.T) {
	link, err := GenTypeAUrl(
		os.Getenv("CDN_URL"),
		os.Getenv("CDN_TIMESTAMP_KEY"),
		os.Getenv("CDN_SIGN_NAME"),
		"0",
		time.Now().Add(1*time.Minute).Unix(),
	)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(link)
}

func TestGenTypeBUrl(t *testing.T) {
	link, err := GenTypeBUrl(
		os.Getenv("CDN_URL"),
		os.Getenv("CDN_TIMESTAMP_KEY"),
		time.Now().Add(-29*time.Minute).Unix(),
	)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(link)
}

func TestGenTypeCUrl(t *testing.T) {
	link, err := GenTypeCUrl(
		os.Getenv("CDN_URL"),
		os.Getenv("CDN_TIMESTAMP_KEY"),
		time.Now().Add(-29*time.Minute).Unix(),
	)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(link)
}

func TestGenTypeDUrl(t *testing.T) {
	link, err := GenTypeDUrl(
		os.Getenv("CDN_URL"),
		os.Getenv("CDN_TIMESTAMP_KEY"),
		os.Getenv("CDN_SIGN_NAME"),
		os.Getenv("CDN_T_NAME"),
		time.Now().Add(-29*time.Minute).Unix(),
		10,
	)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(link)
}

func TestGenTypeEUrl(t *testing.T) {
	link, err := GenTypeEUrl(
		os.Getenv("CDN_URL"),
		os.Getenv("CDN_TIMESTAMP_KEY"),
		os.Getenv("CDN_SIGN_NAME"),
		os.Getenv("CDN_T_NAME"),
		time.Now().Add(-29*time.Minute).Unix(),
		10,
	)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(link)
}

func TestCdnClient_ListCertInfo(t *testing.T) {
	cdnClient := NewCdnClient(os.Getenv("VOLC_ACCESS_KEY"), os.Getenv("VOLC_SECRET_KEY"))
	resp, err := cdnClient.ListCertInfo(&cdn.ListCertInfoRequest{
		Source: "cdn_cert_hosting",
		Name:   trans.String(os.Getenv("VOLC_DOMAIN")),
	})
	if err != nil {
		t.Error(err)
	}
	for _, v := range resp.Result.CertInfo {
		t.Log(v.CertId, v.ExpireTime)
	}
}
