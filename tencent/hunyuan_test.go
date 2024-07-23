package tencent

import (
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	hunyuan "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/hunyuan/v20230901"
	"log"
	"os"
	"testing"
)

func TestHunYuanClient_ChatCompletions(t *testing.T) {
	client, err := NewHunYuanClient(os.Getenv("TencentSecretId"), os.Getenv("TencentSecretKey"), "ap-guangzhou")
	if err != nil {
		t.Error(err)
		return
	}

	req := hunyuan.NewChatCompletionsRequest()
	req.Model = common.StringPtr("hunyuan-lite")
	msg := "多块钟表成液状流淌，一块挂在树的两个分叉上，树上没有叶子已经枯死"
	req.Messages = []*hunyuan.Message{
		{
			Role:    common.StringPtr("user"),
			Content: common.StringPtr("下面我让你来充当翻译家，你的目标是将输入的内容翻译成英文，只要输出翻译后的英文结果: {" + msg + "}"),
		},
	}
	req.TopP = common.Float64Ptr(0.001)
	req.Temperature = common.Float64Ptr(0.001)
	req.Stream = common.BoolPtr(false)
	//req.SearchInfo = common.BoolPtr(true)
	//req.EnableEnhancement = common.BoolPtr(true)
	resp, err := client.ChatCompletions(req)
	if err != nil {
		t.Error(err)
		return
	}
	if resp.Response != nil {
		for _, v := range resp.Response.Choices {
			log.Println(*v.Message.Content)
		}
	} else {
		// 流式
		for v := range resp.Events {
			log.Println(string(v.Data))
		}
	}
}
