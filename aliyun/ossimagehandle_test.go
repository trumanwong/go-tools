package aliyun

import (
	"os"
	"testing"
)

func TestGetOssImageInfo(t *testing.T) {
	resp, err := GetOssImageInfo(os.Getenv("OSS_IMAGE_URL"))
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(resp)
}
