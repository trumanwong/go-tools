package aliyun

import (
	"log"
	"os"
	"testing"
)

func TestDocMindClient_SubmitConvertPdfToWordJobAdvance(t *testing.T) {
	filename := "C:\\Users\\Admin\\Downloads\\example.pdf"
	f, err := os.Open(filename)
	if err != nil {
		t.Errorf("open file error: %s", err.Error())
	}
	defer f.Close()
	client, err := NewDocMindClient(os.Getenv("AccessKeyId"), os.Getenv("AccessKeySecret"), os.Getenv("Endpoint"))
	if err != nil {
		t.Errorf("new client error: %s", err.Error())
		return
	}
	resp, err := client.SubmitConvertPdfToWordJobAdvance("example.pdf", f)
	if err != nil {
		t.Errorf("submit job error: %s", err.Error())
		return
	}
	t.Logf("response requestId: %s", *resp.Body.Data.Id)
}

func TestDocMindClient_GetDocumentConvertResult(t *testing.T) {
	client, err := NewDocMindClient(os.Getenv("AccessKeyId"), os.Getenv("AccessKeySecret"), os.Getenv("Endpoint"))
	if err != nil {
		t.Errorf("new client error: %s", err.Error())
		return
	}
	resp, err := client.GetDocumentConvertResult(os.Getenv("RequestId"))
	if err != nil {
		t.Errorf("get result error: %s", err.Error())
		return
	}
	log.Println(resp.Body)
}
