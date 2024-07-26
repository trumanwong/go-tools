package tencent

import (
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	faceid "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/faceid/v20180301"
	"log"
	"os"
	"testing"
)

func TestFaceIDClient_DectectAuth(t *testing.T) {
	client, err := NewFaceIDClient(os.Getenv("TencentSecretId"), os.Getenv("TencentSecretKey"))
	if err != nil {
		t.Error(err)
		return
	}
	request := faceid.NewDetectAuthRequest()
	request.RuleId = common.StringPtr(os.Getenv("TencentFaceIDRuleId"))
	request.RedirectUrl = common.StringPtr(os.Getenv("TencentFaceIDRedirectUrl"))
	resp, err := client.DectectAuth(request)
	if err != nil {
		t.Error(err)
		return
	}
	log.Println(*resp.Response.Url)
}
