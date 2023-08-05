package helper

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseTime(t *testing.T) {
	parseTime, err := ParseTime("")
	assert.NotNil(t, err)
	assert.Nil(t, parseTime)
}
