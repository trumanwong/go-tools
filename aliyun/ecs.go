package aliyun

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	ecs20140526 "github.com/alibabacloud-go/ecs-20140526/v7/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
)

type EcsClient struct {
	client *ecs20140526.Client
}

func NewEcsClient(config *openapi.Config) (*EcsClient, error) {
	client, _err := ecs20140526.NewClient(config)
	if _err != nil {
		return nil, _err
	}

	return &EcsClient{
		client: client,
	}, nil
}

// DescribeSecurityGroupsWithOptions 查询安全组基本信息列表
func (c EcsClient) DescribeSecurityGroupsWithOptions(req *ecs20140526.DescribeSecurityGroupsRequest, runtime *util.RuntimeOptions) (*ecs20140526.DescribeSecurityGroupsResponse, error) {
	return c.client.DescribeSecurityGroupsWithOptions(req, runtime)
}

// DescribeSecurityGroupAttributeWithOptions 获取安全组和组内规则信息
func (c EcsClient) DescribeSecurityGroupAttributeWithOptions(req *ecs20140526.DescribeSecurityGroupAttributeRequest, runtime *util.RuntimeOptions) (*ecs20140526.DescribeSecurityGroupAttributeResponse, error) {
	return c.client.DescribeSecurityGroupAttributeWithOptions(req, runtime)
}

// 修改安全组入方向规则
func (c EcsClient) ModifySecurityGroupRuleWithOptions(req *ecs20140526.ModifySecurityGroupRuleRequest, runtime *util.RuntimeOptions) (*ecs20140526.ModifySecurityGroupRuleResponse, error) {
	return c.client.ModifySecurityGroupRuleWithOptions(req, runtime)
}
