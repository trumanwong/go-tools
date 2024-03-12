package midjourney

import (
	"log"
	"testing"
)

func TestGetPromptAndParameters(t *testing.T) {
	requests := []*GetPromptAndParametersRequest{
		{
			Content: "This is a test content",
		},
		{
			Content: "This is a test content--test ttt --version 5.4",
		},
		{
			Content: "This is a test content --aspect test --ar 1 --invalid 11",
		},
	}
	for _, req := range requests {
		resp, err := GetPromptAndParameters(req)
		if err != nil {
			t.Errorf("GetPromptAndParameters() error = %v", err)
		}
		log.Println(resp)
	}
}
