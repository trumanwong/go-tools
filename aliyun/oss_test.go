package aliyun

import (
	"log"
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

func TestAliOss_GetSignUrl(t *testing.T) {
	client, err := NewAliOss(os.Getenv("OssEndpoint"), os.Getenv("AccessKeyId"), os.Getenv("AccessKeySecret"), os.Getenv("Bucket"))
	if err != nil {
		t.Error(err)
		return
	}
	url, err := client.GetSignUrl(os.Getenv("ObjectName"), 600)
	if err != nil {
		t.Error(err)
	}
	log.Println(url)
}
