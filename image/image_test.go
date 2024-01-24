package image

import (
	"os"
	"testing"
)

func TestReSave(t *testing.T) {
	err := ReSave(os.Getenv("TEST_FILE_PATH"), os.Getenv("TEST_SAVE_PATH"))
	if err != nil {
		t.Fatal(err)
	}
}
