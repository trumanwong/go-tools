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
			Content:       "This is a test content--test ttt --version 6.0 --turbo",
			DisableParams: []string{"turbo"},
		},
		{
			Content: "This is a test content --aspect test --ar 1 --invalid 11",
		},
		{
			Content: "This is a test content --aspect test --ar 1 --chaos 11 --iw 1.0 --quality 0.5 --repeat 10 --seed 100",
		},
		{
			Content: "This is a test content --aspect test --ar 1 --chaos 11 --sref https://google.com::100 https://google.com --cref https://www.test.com",
		},
		{
			Content: "This is a test content --aspect test --ar 1 --chaos 11 --sref https://google.com::100 --cref",
		},
		{
			Content: "This is a test content --aspect test --ar 1 --chaos 11 --sref --cref",
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
