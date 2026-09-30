package comfyui

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClientWorkflow(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "token" {
			t.Error("missing auth")
		}
		switch r.URL.Path {
		case "/upload/image":
			if err := r.ParseMultipartForm(1024); err != nil {
				t.Error(err)
			}
			w.Write([]byte(`{"name":"input.png"}`))
		case "/prompt":
			w.Write([]byte(`{"prompt_id":"task"}`))
		case "/history/task":
			w.Write([]byte(`{"task":{"status":{"status_str":"success"},"outputs":{"55":{"images":[{"filename":"out.mp4","type":"output"}]}}}}`))
		case "/view":
			w.Write([]byte("video"))
		case "/system_stats", "/object_info/LoadLatent":
			w.Write([]byte(`{}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	client := Client{BaseURL: server.URL, Token: "token"}
	ctx := context.Background()
	if err := client.Health(ctx); err != nil {
		t.Fatal(err)
	}
	name, err := client.UploadInput(ctx, "input.png", strings.NewReader("image"), 1024)
	if err != nil || name != "input.png" {
		t.Fatalf("upload: %s %v", name, err)
	}
	id, err := client.Submit(ctx, map[string]any{"1": map[string]any{"class_type": "LoadLatent"}})
	if err != nil || id != "task" {
		t.Fatalf("submit: %s %v", id, err)
	}
	history, err := client.QueryHistory(ctx, id)
	if err != nil || history.Status != "success" || len(history.Outputs["55"]) != 1 {
		t.Fatalf("history: %+v %v", history, err)
	}
	var output bytes.Buffer
	if err := client.DownloadArtifact(ctx, history.Outputs["55"][0], &output); err != nil || output.String() != "video" {
		t.Fatalf("download: %s %v", output.String(), err)
	}
}

func TestClientRejectsOversizedInput(t *testing.T) {
	_, err := (Client{}).UploadInput(context.Background(), "input.png", strings.NewReader("too big"), 2)
	if err == nil {
		t.Fatal("oversized upload accepted")
	}
}
