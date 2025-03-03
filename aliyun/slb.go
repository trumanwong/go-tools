package aliyun

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	slb20140515 "github.com/alibabacloud-go/slb-20140515/v4/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

type SLBClient struct {
	client *slb20140515.Client
}

func NewSLBClient(accessKeyId, accessKeySecret, endpoint string) (*SLBClient, error) {
	client, err := slb20140515.NewClient(&openapi.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
		Endpoint:        tea.String(endpoint),
	})
	if err != nil {
		return nil, err
	}
	return &SLBClient{client: client}, nil
}

// DescribeLoadBalancersWithOptions 查询已创建的负载均衡实例
func (s *SLBClient) DescribeLoadBalancersWithOptions(req *slb20140515.DescribeLoadBalancersRequest, runtime *util.RuntimeOptions) (*slb20140515.DescribeLoadBalancersResponse, error) {
	return s.client.DescribeLoadBalancersWithOptions(req, runtime)
}

func (s *SLBClient) DescribeLoadBalancerHTTPSListenerAttributeWithOptions(req *slb20140515.DescribeLoadBalancerHTTPSListenerAttributeRequest, runtime *util.RuntimeOptions) (*slb20140515.DescribeLoadBalancerHTTPSListenerAttributeResponse, error) {
	return s.client.DescribeLoadBalancerHTTPSListenerAttributeWithOptions(req, runtime)
}

// DescribeCACertificatesWithOptions 查询CA证书列表
func (s *SLBClient) DescribeCACertificatesWithOptions(req *slb20140515.DescribeCACertificatesRequest, runtime *util.RuntimeOptions) (*slb20140515.DescribeCACertificatesResponse, error) {
	return s.client.DescribeCACertificatesWithOptions(req, runtime)
}

func (s *SLBClient) DescribeServerCertificatesWithOptions(req *slb20140515.DescribeServerCertificatesRequest, runtime *util.RuntimeOptions) (*slb20140515.DescribeServerCertificatesResponse, error) {
	return s.client.DescribeServerCertificatesWithOptions(req, runtime)
}

// UploadServerCertificateWithOptions 上传服务器证书
func (s *SLBClient) UploadServerCertificateWithOptions(req *slb20140515.UploadServerCertificateRequest, runtime *util.RuntimeOptions) (*slb20140515.UploadServerCertificateResponse, error) {
	return s.client.UploadServerCertificateWithOptions(req, runtime)
}

// SetLoadBalancerHTTPSListenerAttributeWithOptions 修改HTTPS监听的配置
func (s *SLBClient) SetLoadBalancerHTTPSListenerAttributeWithOptions(req *slb20140515.SetLoadBalancerHTTPSListenerAttributeRequest, runtime *util.RuntimeOptions) (*slb20140515.SetLoadBalancerHTTPSListenerAttributeResponse, error) {
	return s.client.SetLoadBalancerHTTPSListenerAttributeWithOptions(req, runtime)
}
