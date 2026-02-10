package aliyun

import (
	"os"
	"testing"

	cr20181201 "github.com/alibabacloud-go/cr-20181201/v3/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/trumanwong/go-tools/trans"
)

func TestListInstanceWithOptions(t *testing.T) {
	client, err := NewCrClient(&openapi.Config{
		AccessKeyId:     trans.String(os.Getenv("AccessKeyId")),
		AccessKeySecret: trans.String(os.Getenv("AccessKeySecret")),
		Endpoint:        trans.String(os.Getenv("Endpoint")),
	})

	if err != nil {
		t.Fatal(err)
	}
	resp, err := client.ListInstanceWithOptions(&cr20181201.ListInstanceRequest{}, &util.RuntimeOptions{})
	if err != nil {
		t.Fatal(err)
	}
	for _, instance := range resp.Body.GetInstances() {
		t.Log(instance.InstanceId, instance.InstanceName)
	}
}
