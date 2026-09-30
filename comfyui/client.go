package comfyui

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"
	"strings"
)

// Client 通过 Comfy 公网 API 操作单台机器；调用方负责调度与重试。
type Client struct {
	BaseURL string
	HTTP    *http.Client
	Token   string
}

type Artifact struct {
	Filename  string `json:"filename"`
	Subfolder string `json:"subfolder"`
	Type      string `json:"type"`
}

type HistoryResult struct {
	Status  string
	Outputs map[string][]Artifact
}

func (c Client) request(ctx context.Context, method, route, contentType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, strings.TrimRight(c.BaseURL, "/")+route, body)
	if err != nil {
		return nil, err
	}
	if c.Token != "" {
		req.Header.Set("Authorization", c.Token)
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	client := c.HTTP
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		resp.Body.Close()
		return nil, fmt.Errorf("Comfy %s returned HTTP %d", route, resp.StatusCode)
	}
	return resp, nil
}

// UploadInput 上传不超过 maxBytes 的图片或音频到 Comfy input 目录。
func (c Client) UploadInput(ctx context.Context, filename string, input io.Reader, maxBytes int64) (string, error) {
	if path.Base(filename) != filename || filename == "" || maxBytes <= 0 {
		return "", errors.New("invalid Comfy input filename or size limit")
	}
	var payload bytes.Buffer
	writer := multipart.NewWriter(&payload)
	part, err := writer.CreateFormFile("image", filename)
	if err != nil {
		return "", err
	}
	n, err := io.Copy(part, io.LimitReader(input, maxBytes+1))
	if err != nil {
		return "", err
	}
	if n > maxBytes {
		return "", errors.New("Comfy input exceeds size limit")
	}
	if err := writer.Close(); err != nil {
		return "", err
	}
	resp, err := c.request(ctx, http.MethodPost, "/upload/image", writer.FormDataContentType(), &payload)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var result UploadImageResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if result.Name == "" {
		return "", errors.New("Comfy upload returned no filename")
	}
	return result.Name, nil
}

// Submit 发送 API 格式节点图并返回用于绑定机器的 prompt ID。
func (c Client) Submit(ctx context.Context, graph any) (string, error) {
	data, err := json.Marshal(map[string]any{"prompt": graph})
	if err != nil {
		return "", err
	}
	resp, err := c.request(ctx, http.MethodPost, "/prompt", "application/json", bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var result PromptResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if result.PromptID == "" {
		return "", errors.New("Comfy rejected prompt")
	}
	return result.PromptID, nil
}

// QueryHistory 区分排队、完成和失败，兼容 SaveVideo 的 images 产物。
func (c Client) QueryHistory(ctx context.Context, promptID string) (HistoryResult, error) {
	if promptID == "" || strings.ContainsAny(promptID, "/?#") {
		return HistoryResult{}, errors.New("invalid prompt ID")
	}
	resp, err := c.request(ctx, http.MethodGet, "/history/"+url.PathEscape(promptID), "", nil)
	if err != nil {
		return HistoryResult{}, err
	}
	defer resp.Body.Close()
	var history map[string]struct {
		Status struct {
			StatusStr string `json:"status_str"`
		} `json:"status"`
		Outputs map[string]struct {
			Images []Artifact `json:"images"`
			Videos []Artifact `json:"videos"`
			Gifs   []Artifact `json:"gifs"`
		} `json:"outputs"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&history); err != nil {
		return HistoryResult{}, err
	}
	item, ok := history[promptID]
	if !ok {
		return HistoryResult{Status: "pending"}, nil
	}
	result := HistoryResult{Status: item.Status.StatusStr, Outputs: make(map[string][]Artifact)}
	for key, output := range item.Outputs {
		result.Outputs[key] = append(append(output.Images, output.Videos...), output.Gifs...)
	}
	return result, nil
}

// DownloadArtifact 将 /view 响应流式写入调用方提供的 writer。
func (c Client) DownloadArtifact(ctx context.Context, artifact Artifact, output io.Writer) error {
	if artifact.Filename == "" {
		return errors.New("missing artifact filename")
	}
	query := url.Values{"filename": {artifact.Filename}, "subfolder": {artifact.Subfolder}, "type": {artifact.Type}}
	resp, err := c.request(ctx, http.MethodGet, "/view?"+query.Encode(), "", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(output, resp.Body)
	return err
}

// Health 检查目标 Comfy 实例及 LTX latent 节点是否可用。
func (c Client) Health(ctx context.Context) error {
	for _, route := range []string{"/system_stats", "/object_info/LoadLatent"} {
		resp, err := c.request(ctx, http.MethodGet, route, "", nil)
		if err != nil {
			return err
		}
		resp.Body.Close()
	}
	return nil
}
