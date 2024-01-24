package image

import (
	"os"
	"testing"
)

func TestReSave(t *testing.T) {
	f, err := os.Open(os.Getenv("TEST_FILE_PATH"))
	if err != nil {
		t.Fatal(err)
	}
	err = ReSave(f, os.Getenv("TEST_SAVE_PATH"))
	if err != nil {
		t.Fatal(err)
	}
}
