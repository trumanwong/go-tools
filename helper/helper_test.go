package helper

import (
	"net/http"
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
	resp, err := CheckHttp("https://trumanwl.com", 30*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatal("status code not 200, status: ", resp.StatusCode)
	}
}

func TestGenerateShortUrl(t *testing.T) {
	url, err := GenerateShortUrl("https://short.trumanwl.com", "https://trumanwl.com/xxx")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(url)
}
