package tencent

import (
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	faceid "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/faceid/v20180301"
	"log"
	"os"
	"testing"
)

func TestFaceIDClient_DetectAuth(t *testing.T) {
	client, err := NewFaceIDClient(os.Getenv("TencentSecretId"), os.Getenv("TencentSecretKey"))
	if err != nil {
		t.Error(err)
		return
	}
	request := faceid.NewDetectAuthRequest()
	request.RuleId = common.StringPtr(os.Getenv("TencentFaceIDRuleId"))
	request.RedirectUrl = common.StringPtr(os.Getenv("TencentFaceIDRedirectUrl"))
	resp, err := client.DetectAuth(request)
	if err != nil {
		t.Error(err)
		return
	}
	log.Println(*resp.Response.Url)
	log.Println(*resp.Response.BizToken)
}

func TestFaceIDClient_GetDetectInfoEnhanced(t *testing.T) {
	client, err := NewFaceIDClient(os.Getenv("TencentSecretId"), os.Getenv("TencentSecretKey"))
	if err != nil {
		t.Error(err)
		return
	}
	request := faceid.NewGetDetectInfoEnhancedRequest()
	request.RuleId = common.StringPtr(os.Getenv("TencentFaceIDRuleId"))
	request.BizToken = common.StringPtr(os.Getenv("TencentFaceIDBizToken"))
	request.InfoType = common.StringPtr("0")
	resp, err := client.GetDetectInfoEnhanced(request)
	if err != nil {
		t.Error(err)
		return
	}
	//log.Println(*resp.Response.IdCardData.OcrFront)
	log.Println(*resp.Response.Text.Name)
	log.Println(*resp.Response.Text.IdCard)
}
