package tencent

import (
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	hunyuan "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/hunyuan/v20230901"
)

type HunYuanClient struct {
	client *hunyuan.Client
}

func NewHunYuanClient(secretId, secretKey, region string) (*HunYuanClient, error) {
	credential := common.NewCredential(
		secretId,
		secretKey,
	)
	// 实例化一个client选项，可选的，没有特殊需求可以跳过
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "hunyuan.tencentcloudapi.com"
	// 实例化要请求产品的client对象,clientProfile是可选的
	client, err := hunyuan.NewClient(credential, region, cpf)
	if err != nil {
		return nil, err
	}
	return &HunYuanClient{
		client: client,
	}, nil
}

func (c HunYuanClient) ChatCompletions(req *hunyuan.ChatCompletionsRequest) (*hunyuan.ChatCompletionsResponse, error) {
	response, err := c.client.ChatCompletions(req)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return nil, fmt.Errorf("an api error has returned: %s", err)
	}
	if err != nil {
		return nil, err
	}
	return response, nil
}
