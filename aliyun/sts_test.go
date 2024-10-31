package aliyun

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	sts20150401 "github.com/alibabacloud-go/sts-20150401/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"os"
	"testing"
)

func TestStsClient_AssumeRoleWithOptions(t *testing.T) {
	client, err := NewStsClient(&openapi.Config{
		AccessKeyId:     tea.String(os.Getenv("AccessKeyId")),
		AccessKeySecret: tea.String(os.Getenv("AccessKeySecret")),
		Endpoint:        tea.String(os.Getenv("Endpoint")),
	})
	if err != nil {
		t.Error(err)
		return
	}
	var req sts20150401.AssumeRoleRequest
	req.SetRoleArn(os.Getenv("RoleArn"))
	req.SetRoleSessionName("test")
	req.SetDurationSeconds(3600)
	resp, err := client.AssumeRoleWithOptions(&req, &util.RuntimeOptions{})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(*resp.Body.Credentials.AccessKeyId)
	t.Log(*resp.Body.Credentials.AccessKeySecret)
	t.Log(*resp.Body.Credentials.Expiration)
	t.Log(*resp.Body.Credentials.SecurityToken)
	t.Log(*resp.Body.RequestId)
}
