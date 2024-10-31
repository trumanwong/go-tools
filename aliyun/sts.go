package aliyun

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	sts20150401 "github.com/alibabacloud-go/sts-20150401/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
)

type StsClient struct {
	client *sts20150401.Client
}

func NewStsClient(config *openapi.Config) (*StsClient, error) {
	client, err := sts20150401.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &StsClient{client: client}, nil
}

func (c StsClient) AssumeRoleWithOptions(assumeRoleRequest *sts20150401.AssumeRoleRequest, runtime *util.RuntimeOptions) (*sts20150401.AssumeRoleResponse, error) {
	return c.client.AssumeRoleWithOptions(assumeRoleRequest, runtime)
}
