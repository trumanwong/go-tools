package crawler

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

// mockTransport 模拟 HTTP 请求，不发起真实网络调用
type mockTransport struct {
	roundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.roundTripFunc(req)
}

func newMockTransport(fn func(req *http.Request) (*http.Response, error)) *mockTransport {
	return &mockTransport{roundTripFunc: fn}
}

func okResponse() *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader("ok")),
	}
}

func TestSend_Get(t *testing.T) {
	transport := newMockTransport(func(req *http.Request) (*http.Response, error) {
		if req.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", req.Method)
		}
		return okResponse(), nil
	})

	resp, err := Send(&Request{
		Url:       "http://example.com",
		Method:    http.MethodGet,
		Transport: transport,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestSend_WithHeaders(t *testing.T) {
	transport := newMockTransport(func(req *http.Request) (*http.Response, error) {
		if req.Header.Get("X-Custom") != "test-value" {
			t.Errorf("expected header X-Custom=test-value, got %s", req.Header.Get("X-Custom"))
		}
		return okResponse(), nil
	})

	resp, err := Send(&Request{
		Url:       "http://example.com",
		Method:    http.MethodGet,
		Headers:   map[string]string{"X-Custom": "test-value"},
		Transport: transport,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
}

func TestSend_WithHost(t *testing.T) {
	transport := newMockTransport(func(req *http.Request) (*http.Response, error) {
		if req.Host != "custom.host.com" {
			t.Errorf("expected host custom.host.com, got %s", req.Host)
		}
		if req.Header.Get("Host") != "custom.host.com" {
			t.Errorf("expected Host header custom.host.com, got %s", req.Header.Get("Host"))
		}
		return okResponse(), nil
	})

	host := "custom.host.com"
	resp, err := Send(&Request{
		Url:       "http://example.com",
		Method:    http.MethodGet,
		Host:      &host,
		Transport: transport,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
}

func TestSend_WithBasicAuth(t *testing.T) {
	transport := newMockTransport(func(req *http.Request) (*http.Response, error) {
		user, pass, ok := req.BasicAuth()
		if !ok || user != "admin" || pass != "secret" {
			t.Errorf("expected basic auth admin:secret, got %s:%s (ok=%v)", user, pass, ok)
		}
		return okResponse(), nil
	})

	resp, err := Send(&Request{
		Url:    "http://example.com",
		Method: http.MethodGet,
		BasicAuth: &BasicAuth{
			Username: "admin",
			Password: "secret",
		},
		Transport: transport,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
}

func TestSend_WithPostForm(t *testing.T) {
	transport := newMockTransport(func(req *http.Request) (*http.Response, error) {
		if req.PostForm == nil {
			t.Error("expected PostForm to be set")
		}
		if req.PostForm.Get("name") != "test" {
			t.Errorf("expected PostForm name=test, got %s", req.PostForm.Get("name"))
		}
		if req.PostForm.Get("age") != "18" {
			t.Errorf("expected PostForm age=18, got %s", req.PostForm.Get("age"))
		}
		return okResponse(), nil
	})

	resp, err := Send(&Request{
		Url:    "http://example.com",
		Method: http.MethodPost,
		Body:   strings.NewReader("name=test&age=18"),
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
		PostForm: url.Values{
			"name": {"test"},
			"age":  {"18"},
		},
		Transport: transport,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestSend_PostFormNil(t *testing.T) {
	transport := newMockTransport(func(req *http.Request) (*http.Response, error) {
		if req.PostForm != nil {
			t.Error("expected PostForm to be nil")
		}
		return okResponse(), nil
	})

	resp, err := Send(&Request{
		Url:       "http://example.com",
		Method:    http.MethodGet,
		Transport: transport,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
}

func TestSend_InvalidUrl(t *testing.T) {
	_, err := Send(&Request{
		Url:    "://invalid",
		Method: http.MethodGet,
	})
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}
