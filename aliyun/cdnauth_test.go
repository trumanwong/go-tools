package aliyun

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCdnAuth_AuthA(t *testing.T) {
	cdnAuth := NewCdnAuth("aliyuncdnexp1234")
	url, err := cdnAuth.AuthA("http://domain.example.com/video/standard/test.mp4", 1444435200, "0")
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, "http://domain.example.com/video/standard/test.mp4?auth_key=1444435200-0-0-23bf85053008f5c0e791667a313e28ce", url)
}
