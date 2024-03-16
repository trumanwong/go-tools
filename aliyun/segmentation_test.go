package aliyun

import (
	imageseg20191230 "github.com/alibabacloud-go/imageseg-20191230/v3/client"
	"github.com/alibabacloud-go/tea/tea"
	"log"
	"os"
	"testing"
)

func TestSegmentationClient_SegmentCommodity(t *testing.T) {
	client, err := NewSegmentationClient(os.Getenv("AccessKeyId"), os.Getenv("AccessKeySecret"), "imageseg.cn-shanghai.aliyuncs.com")
	if err != nil {
		t.Error(err)
		return
	}
	resp, err := client.SegmentCommodity(&imageseg20191230.SegmentCommodityRequest{
		ImageURL:   tea.String("http://viapi-test.oss-cn-shanghai.aliyuncs.com/viapi-3.0domepic/imageseg/SegmentCommodity/SegmentCommodity1.jpg"),
		ReturnForm: nil,
	})
	if err != nil {
		t.Error(err)
		return
	}
	log.Println(*resp.Body.Data.ImageURL)
}
