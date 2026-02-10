package aliyun

import (
	cr20181201 "github.com/alibabacloud-go/cr-20181201/v3/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
)

type CrClient struct {
	client *cr20181201.Client
}

func NewCrClient(config *openapi.Config) (*CrClient, error) {
	client, err := cr20181201.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &CrClient{client: client}, nil
}

func (c CrClient) ListInstanceWithOptions(req *cr20181201.ListInstanceRequest, runtime *util.RuntimeOptions) (*cr20181201.ListInstanceResponse, error) {
	return c.client.ListInstanceWithOptions(req, runtime)
}
