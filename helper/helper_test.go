package helper

import (
	"github.com/trumanwong/go-tools/crawler"
	"log"
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

func TestShuffleArray(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5}
	ShuffleArray(arr)
	t.Log(arr)
}

func TestDownloadFile(t *testing.T) {
	_, err := DownloadFile(&crawler.Request{
		Url:    "https://cdn.trumanwl.com/favicon.ico",
		Method: http.MethodGet,
	}, "favicon.ico", false)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetSSLExpireDate(t *testing.T) {
	expire, err := GetSSLExpireDate("gc1.midjourny.cn")
	if err != nil {
		t.Error(err)
	}
	log.Println(expire)
}

func TestTernary(t *testing.T) {
	if Ternary(true, 1, 2) != 1 {
		t.Fatal("Ternary error")
	}
	if Ternary(false, 1, 2) != 2 {
		t.Fatal("Ternary error")
	}
}
