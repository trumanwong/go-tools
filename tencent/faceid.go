package tencent

import (
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	faceid "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/faceid/v20180301"
)

type FaceIDClient struct {
	client *faceid.Client
}

func NewFaceIDClient(secretId, secretKey string) (*FaceIDClient, error) {
	credential := common.NewCredential(
		secretId,
		secretKey,
	)
	// 实例化一个client选项，可选的，没有特殊需求可以跳过
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "faceid.tencentcloudapi.com"
	client, err := faceid.NewClient(credential, "", cpf)
	if err != nil {
		return nil, err
	}
	return &FaceIDClient{
		client: client,
	}, nil
}

func (c FaceIDClient) DectectAuth(request *faceid.DetectAuthRequest) (*faceid.DetectAuthResponse, error) {
	resp, err := c.client.DetectAuth(request)
	if err != nil {
		if _, ok := err.(*errors.TencentCloudSDKError); ok {
			return nil, fmt.Errorf("an api error has returned: %s", err)
		}
		return nil, err
	}
	return resp, nil
}
