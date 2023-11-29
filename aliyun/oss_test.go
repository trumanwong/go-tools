package aliyun

import (
	"net/http"
	"os"
	"testing"
)

func TestAliOss_PutObject(t *testing.T) {
	client, err := NewAliOss(os.Getenv("OssEndpoint"), os.Getenv("AccessKeyId"), os.Getenv("AccessKeySecret"), os.Getenv("Bucket"))
	if err != nil {
		t.Error(err)
		return
	}
	resp, err := http.Get("https://imgmil.gmw.cn/attachement/jpg/site2/20231129/f44d305ea6dd26d0d59d42.jpg")
	if err != nil {
		t.Error(err)
		return
	}
	err = client.PutObject("oaidalleapiprodscus/test2.png", resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
}
