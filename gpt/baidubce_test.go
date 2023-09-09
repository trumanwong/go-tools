package gpt

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetBaidubceAccessToken(t *testing.T) {
	_, err := GetBaidubceAccessToken(os.Getenv("BAIDU_API_KEY"), os.Getenv("BAIDU_SECRET_KEY"))
	assert.Equal(t, nil, err)
}
