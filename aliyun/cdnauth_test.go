package aliyun

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestCdnAuth_AuthB(t *testing.T) {
	cdnAuth := NewCdnAuth("aliyuncdnexp1234")
	url, err := cdnAuth.AuthB("http://domain.example.com/4/44/44c0909bcfc20a01afaf256ca99a8b8b.mp3", 1439596800)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, "http://domain.example.com/201508150800/9044548ef1527deadafa49a890a377f0/4/44/44c0909bcfc20a01afaf256ca99a8b8b.mp3", url)
}
