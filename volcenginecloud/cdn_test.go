package volcenginecloud

import (
	"github.com/trumanwong/go-tools/trans"
	"github.com/volcengine/volc-sdk-golang/service/cdn"
	"log"
	"os"
	"testing"
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
