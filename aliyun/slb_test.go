package aliyun

import (
	"fmt"
	slb20140515 "github.com/alibabacloud-go/slb-20140515/v4/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"os"
	"testing"
)

func TestNewSLBClient(t *testing.T) {
	client, err := NewSLBClient(os.Getenv("ALIYUN_ACCESS_KEY_ID"), os.Getenv("ALIYUN_ACCESS_KEY_SECRET"), "slb.cn-shenzhen.aliyuncs.com")
	if err != nil {
		t.Fatal(err)
	}

	runtime := &util.RuntimeOptions{}
	resp, err := client.DescribeLoadBalancersWithOptions(&slb20140515.DescribeLoadBalancersRequest{}, runtime)
	if err != nil {
		t.Fatal(err)
	}
	for _, balance := range resp.Body.LoadBalancers.LoadBalancer {
		if *balance.LoadBalancerStatus == "active" {
			balancerHttpsListener, err := client.DescribeLoadBalancerHTTPSListenerAttributeWithOptions(&slb20140515.DescribeLoadBalancerHTTPSListenerAttributeRequest{
				LoadBalancerId: balance.LoadBalancerId,
				ListenerPort:   tea.Int32(443),
			}, runtime)
			if err != nil {
				t.Fatal(err)
			}
			caCertificates, err := client.DescribeServerCertificatesWithOptions(
				&slb20140515.DescribeServerCertificatesRequest{
					RegionId:            tea.String("cn-shenzhen"),
					ServerCertificateId: balancerHttpsListener.Body.ServerCertificateId,
				}, runtime)
			if err != nil {
				t.Fatal(err)
			}
			for _, certificate := range caCertificates.Body.ServerCertificates.ServerCertificate {
				fmt.Println(*certificate.ExpireTime, *certificate.ExpireTimeStamp)
			}
		}
	}
}
