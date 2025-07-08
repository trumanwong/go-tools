package aliyun

import (
	"fmt"
	"os"
	"testing"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	ecs20140526 "github.com/alibabacloud-go/ecs-20140526/v7/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

func TestDescribeSecurityGroupsWithOptions(t *testing.T) {
	config := &openapi.Config{
		AccessKeyId:     tea.String(os.Getenv("AccessKeyId")),
		AccessKeySecret: tea.String(os.Getenv("AccessKeySecret")),
		RegionId:        tea.String("cn-shenzhen"),
	}
	ecsClient, err := NewEcsClient(config)
	if err != nil {
		t.Fatalf("failed to create ECS client: %v", err)
	}

	// Create a request
	req := &ecs20140526.DescribeSecurityGroupsRequest{
		RegionId:   tea.String("cn-shenzhen"),
		MaxResults: tea.Int32(100),
	}

	// Call the method
	resp, err := ecsClient.DescribeSecurityGroupsWithOptions(req, &util.RuntimeOptions{})
	if err != nil {
		t.Fatalf("failed to describe security groups: %v", err)
	}

	fmt.Println(len(resp.Body.SecurityGroups.SecurityGroup))
	for _, securityGroup := range resp.Body.SecurityGroups.SecurityGroup {
		fmt.Println(
			*securityGroup.SecurityGroupId,
			*securityGroup.SecurityGroupName,
			*securityGroup.Description,
		)
	}
}

func TestDescribeSecurityGroupAttributeWithOptions(t *testing.T) {
	config := &openapi.Config{
		AccessKeyId:     tea.String(os.Getenv("AccessKeyId")),
		AccessKeySecret: tea.String(os.Getenv("AccessKeySecret")),
		RegionId:        tea.String("cn-shenzhen"),
	}
	ecsClient, err := NewEcsClient(config)
	if err != nil {
		t.Fatalf("failed to create ECS client: %v", err)
	}

	// Create a request
	req := &ecs20140526.DescribeSecurityGroupAttributeRequest{
		RegionId:        tea.String("cn-shenzhen"),
		SecurityGroupId: tea.String(os.Getenv("SecurityGroupId")), // Replace with your security group ID
	}

	// Call the method
	resp, err := ecsClient.DescribeSecurityGroupAttributeWithOptions(req, &util.RuntimeOptions{})
	if err != nil {
		t.Fatalf("failed to describe security group attribute: %v", err)
	}

	for _, permission := range resp.Body.Permissions.Permission {
		fmt.Println(
			// 安全组规则ID
			*permission.SecurityGroupRuleId,
			// 安全组描述信息
			*permission.Description,
			// 端口范围
			*permission.PortRange,
			// 源端ip地址
			*permission.SourceCidrIp,
		)
	}
}

func TestModifySecurityGroupRuleWithOptions(t *testing.T) {
	config := &openapi.Config{
		AccessKeyId:     tea.String(os.Getenv("AccessKeyId")),
		AccessKeySecret: tea.String(os.Getenv("AccessKeySecret")),
		RegionId:        tea.String("cn-shenzhen"),
	}
	ecsClient, err := NewEcsClient(config)
	if err != nil {
		t.Fatalf("failed to create ECS client: %v", err)
	}
	// Create a request
	req := &ecs20140526.ModifySecurityGroupRuleRequest{
		RegionId:            tea.String("cn-shenzhen"),
		SecurityGroupId:     tea.String(os.Getenv("SecurityGroupId")),     // Replace with your security group ID
		SecurityGroupRuleId: tea.String(os.Getenv("SecurityGroupRuleId")), // Replace with your security group rule ID
		SourceCidrIp:        tea.String(os.Getenv("SourceCidrIp")),        // Replace with your source CIDR IP
	}
	// Call the method
	resp, err := ecsClient.ModifySecurityGroupRuleWithOptions(req, &util.RuntimeOptions{})
	if err != nil {
		t.Fatalf("failed to modify security group rule: %v", err)
	}
	fmt.Println("Security group rule modified successfully:", resp)
}
