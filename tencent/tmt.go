package tencent

import (
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tmt "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tmt/v20180321"
	"log"
)

type TmtClient struct {
	client *tmt.Client
}

func NewTmtClient(secretId, secretKey, region string) (*TmtClient, error) {
	log.Println("=="+secretId+"===", "=="+secretKey+"===")
	credential := common.NewCredential(
		secretId,
		secretKey,
	)
	// 实例化一个client选项，可选的，没有特殊需求可以跳过
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "tmt.tencentcloudapi.com"
	// 实例化要请求产品的client对象,clientProfile是可选的
	client, err := tmt.NewClient(credential, region, cpf)
	if err != nil {
		return nil, err
	}
	return &TmtClient{
		client: client,
	}, nil
}

func (t TmtClient) TextTranslate(req *tmt.TextTranslateRequest) (string, error) {
	response, err := t.client.TextTranslate(req)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return "", fmt.Errorf("an api error has returned: %s", err)

	}
	if err != nil {
		return "", err
	}
	return *response.Response.TargetText, nil
}
