package tencent

import (
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	tmt "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tmt/v20180321"
	"os"
	"testing"
)

func TestTmtClient_TextTranslate(t *testing.T) {
	client, err := NewTmtClient(os.Getenv("TENCENT_SECRET_ID"), os.Getenv("TENCENT_SECRET_KEY"), "ap-guangzhou")
	if err != nil {
		t.Error(err)
		return
	}
	b, err := os.ReadFile(os.Getenv("SOURCE_FILE_PATH"))
	if err != nil {
		t.Error(err)
		return
	}
	req := tmt.NewTextTranslateRequest()
	req.SourceText = common.StringPtr(string(b))
	req.Source = common.StringPtr("zh")
	req.Target = common.StringPtr("jp")
	req.ProjectId = common.Int64Ptr(0)
	resp, err := client.TextTranslate(req)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(resp)
}
