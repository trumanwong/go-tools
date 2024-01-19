package comfyui

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"testing"
)

func TestServer_UploadImage(t *testing.T) {
	s := NewServer(os.Getenv("COMFY_HOST"))
	f, err := os.Open(os.Getenv("TEST_FILE_PATH"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	part, err := writer.CreateFormFile("image", filepath.Base(os.Getenv("TEST_FILE_PATH")))
	if err != nil {
		t.Fatal(err)
	}
	_, err = io.Copy(part, f)
	if err != nil {
		t.Fatal(err)
	}
	if err = writer.Close(); err != nil {
		t.Fatal(err)
	}
	_, err = s.UploadImage(payload, writer)
	if err != nil {
		t.Fatal(err)
	}
}

func TestServer_Prompt(t *testing.T) {
	s := NewServer(os.Getenv("COMFY_HOST"))
	prompt := make(map[string]interface{})
	extraData := make(map[string]interface{})
	if err := json.Unmarshal([]byte(os.Getenv("TEST_PROMPT")), &prompt); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal([]byte(os.Getenv("TEST_EXTRA_DATA")), &extraData); err != nil {
		t.Fatal(err)
	}
	resp, err := s.Prompt(uuid.New().String(), prompt, extraData)
	if err != nil {
		t.Fatal(err)
	}
	if resp.PromptID == "" {
		t.Fatal("prompt id is empty")
	}
	log.Println(resp.PromptID)
}

func TestServer_History(t *testing.T) {
	s := NewServer(os.Getenv("COMFY_HOST"))
	resp, err := s.History("c6070c4a-cd1d-4c9f-8578-9542a8ea7356", []string{"40"}...)
	if err != nil {
		t.Fatal(err)
	}
	if len(resp) == 0 {
		t.Fatal("images is empty")
	}
	for _, image := range resp {
		for _, img := range image.Images {
			log.Println(img.Filename)
		}
	}
}
