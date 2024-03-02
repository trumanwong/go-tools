package helper

import (
	"testing"
	"time"
)

func TestCheckPort(t *testing.T) {
	err := CheckPort("127.0.0.1", "80", 3*time.Second)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCheckHttp(t *testing.T) {
	err := CheckHttp("https://trumanwl.com", 30*time.Second)
	if err != nil {
		t.Fatal(err)
	}
}
