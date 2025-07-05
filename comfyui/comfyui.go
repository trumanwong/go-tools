package comfyui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/trumanwong/go-tools/crawler"
)

type Server struct {
	host string
}

// 上传图片
const uploadImageApi = "/upload/image"

// 提交任务
const promptApi = "/prompt"

// 查询任务状态
const historyApi = "/history/%s"

// 查询队列
const queueApi = "/queue"

func NewServer(host string) *Server {
	return &Server{
		host: host,
	}
}

type UploadImageResponse struct {
	Name      string `json:"name"`
	SubFolder string `json:"subfolder"`
	Type      string `json:"type"`
}

// UploadImage 上传图片
func (s Server) UploadImage(payload *bytes.Buffer, writer *multipart.Writer) (*UploadImageResponse, error) {
	resp, err := crawler.Send(&crawler.Request{
		Url:    s.host + uploadImageApi,
		Method: http.MethodPost,
		Headers: map[string]string{
			"Content-Type": writer.FormDataContentType(),
		},
		Body: payload,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body error: %s, body: %s", err.Error(), body)
	}
	var result UploadImageResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, fmt.Errorf("unmarshal error: %s, body: %s", err.Error(), body)
	}
	return &result, nil
}

type PromptResponse struct {
	NodeErrors any    `json:"node_errors"`
	Number     int    `json:"number"`
	PromptID   string `json:"prompt_id"`
}

// Prompt 加入队列
func (s Server) Prompt(clientId string, prompt map[string]any, extraData map[string]any) (*PromptResponse, error) {
	payload, err := json.Marshal(map[string]any{
		"client_id": clientId,
		"prompt":    prompt,
		"extraData": extraData,
	})
	if err != nil {
		return nil, err
	}
	resp, err := crawler.Send(&crawler.Request{
		Url:    s.host + promptApi,
		Method: http.MethodPost,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: bytes.NewBuffer(payload),
	})

	if err != nil {
		return nil, fmt.Errorf("send request error: %s", err.Error())
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body error: %s", err.Error())
	}
	var result PromptResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, fmt.Errorf("decode response error: %s, body: %s", err.Error(), body)
	}
	return &result, nil
}

type Image struct {
	KeyType   string `json:"key_type"`
	Filename  string `json:"filename"`
	SubFolder string `json:"subfolder"`
	Type      string `json:"type"`
	Format    string `json:"format"`
}

type ImageResult struct {
	Images []*Image `json:"images"`
	Text   []string `json:"text"`
	Key    string   `json:"key"`
}

// History 查询任务状态
// promptId 任务id
// keys 要取的结果（outputs）的key
func (s Server) History(promptId string, timeout *time.Duration, keys ...string) ([]*ImageResult, map[string]any, error) {
	req := &crawler.Request{
		Url:    s.host + fmt.Sprintf(historyApi, promptId),
		Method: http.MethodGet,
	}
	if timeout != nil {
		req.Timeout = *timeout
	}
	resp, err := crawler.Send(req)
	if err != nil {
		return nil, nil, err
	}
	m := make(map[string]any)
	err = json.NewDecoder(resp.Body).Decode(&m)
	if err != nil {
		return nil, nil, err
	}
	if _, ok := m[promptId]; !ok {
		return nil, nil, fmt.Errorf("promptId: %s not found", promptId)
	}
	// 取出任务详情
	task := m[promptId].(map[string]any)
	var status map[string]any
	if _, ok := task["status"]; !ok {
		return nil, nil, fmt.Errorf("status not found")
	} else {
		status = task["status"].(map[string]any)
		if !status["completed"].(bool) || status["status_str"] != "success" {
			return nil, status, fmt.Errorf("status is not success")
		}
	}
	if _, ok := task["outputs"]; !ok {
		return nil, status, fmt.Errorf("outputs not found")
	}
	outputs := task["outputs"].(map[string]any)
	result := make([]*ImageResult, 0)
	for _, key := range keys {
		if _, ok := outputs[key]; !ok {
			return nil, status, fmt.Errorf("key: %s not found", key)
		}
		keyImages := outputs[key].(map[string]any)
		images := make([]*Image, 0)
		text := make([]string, 0)
		var b []byte
		if _, ok := keyImages["images"]; ok {
			b, _ = json.Marshal(keyImages["images"])
			err = json.Unmarshal(b, &images)
			for i, _ := range images {
				images[i].KeyType = "images"
			}
		} else if _, ok := keyImages["gifs"]; ok {
			b, _ = json.Marshal(keyImages["gifs"])
			err = json.Unmarshal(b, &images)
			for i, _ := range images {
				images[i].KeyType = "gifs"
			}
		} else if _, ok := keyImages["text"]; ok {
			b, _ = json.Marshal(keyImages["text"])
			err = json.Unmarshal(b, &text)
			for i, _ := range images {
				images[i].KeyType = "text"
			}
		} else if _, ok := keyImages["string"]; ok {
			b, _ = json.Marshal(keyImages["string"])
			err = json.Unmarshal(b, &text)
			for i, _ := range images {
				images[i].KeyType = "string"
			}
		} else {
			return nil, status, fmt.Errorf("key images/gifs/text is not exists")
		}
		if err != nil {
			return nil, status, fmt.Errorf("key: %s unmarshal error, images: %s", key, b)
		}
		result = append(result, &ImageResult{
			Images: images,
			Text:   text,
			Key:    key,
		})
	}
	return result, status, nil
}

func (s Server) QueueIsRunning(promptId string) (bool, error) {
	resp, err := crawler.Send(&crawler.Request{
		Url:    s.host + queueApi,
		Method: http.MethodGet,
	})
	if err != nil {
		return false, err
	}
	m := make(map[string]any)
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("read body error: %s", err.Error())
	}
	err = json.Unmarshal(b, &m)
	if err != nil {
		return false, fmt.Errorf("unmarshal error: %s, body: %s", err.Error(), b)
	}
	if _, ok := m["queue_running"]; !ok {
		return false, fmt.Errorf("queue_running not found, body: %s", b)
	}
	for _, item := range m["queue_running"].([]any) {
		arr := item.([]any)
		if arr[1].(string) == promptId {
			return true, nil
		}
	}
	return false, nil
}

func (s Server) Cancel(promptId ...string) error {
	body, _ := json.Marshal(map[string]any{
		"delete": promptId,
	})
	resp, err := crawler.Send(&crawler.Request{
		Url:     s.host + queueApi,
		Method:  http.MethodPost,
		Body:    bytes.NewReader(body),
		Timeout: 60 * time.Second,
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
