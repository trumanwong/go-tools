package validate

import "testing"

func TestCheckContainUrl(t *testing.T) {
	if !CheckContainUrl("This is url: https://www.baidu.com") {
		t.Error("CheckContainUrl test failed")
		return
	}
	if CheckContainUrl("This is not url") {
		t.Error("CheckContainUrl test failed")
		return
	}
}
