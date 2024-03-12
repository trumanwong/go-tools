package midjourney

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestGetPromptAndParameters(t *testing.T) {
	requests := []*GetPromptAndParametersRequest{
		{
			Content: "This is a test content",
		},
		{
			Content: "This is a test content --aspect test --ar 1 --invalid 11",
		},
	}
	for _, req := range requests {
		resp, err := GetPromptAndParameters(req)
		if err != nil {
			assert.Errorf(t, err, "GetPromptAndParameters() error = %v", err)
		}
		log.Println(resp)
	}
}
