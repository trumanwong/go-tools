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
	resp, err := s.History("d7281953-9746-4990-91e0-f0e99dbafdef", []string{"58"}...)
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
		for _, img := range image.Gifs {
			log.Println(img.Filename)
			log.Println(img.Format)
		}
	}
}

func TestServer_QueueIsRunning(t *testing.T) {
	s := NewServer(os.Getenv("COMFY_HOST"))
	resp, err := s.QueueIsRunning("78cc8638-160c-468e-8675-f462d61ca4d8")
	if err != nil {
		t.Fatal(err)
	}
	log.Println(resp)
}

func TestServer_Cancel(t *testing.T) {
	s := NewServer(os.Getenv("COMFY_HOST"))
	err := s.Cancel([]string{"78cc8638-160c-468e-8675-f462d61ca4d8"}...)
	if err != nil {
		t.Fatal(err)
	}
}
