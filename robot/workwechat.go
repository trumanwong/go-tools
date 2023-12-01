package robot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type WorkWechatRobot struct {
	// 企业微信机器人url
	url string
}

func NewWorkWechatRobot(url string) *WorkWechatRobot {
	return &WorkWechatRobot{url: url}
}

type Level string

const (
	LevelInfo    Level = "info"
	LevelWarning Level = "warning"
	LevelError   Level = "error"
)

type SentTextRequest struct {
	Level   Level  `json:"level"`
	Content string `json:"content"`
	IsAtAll bool   `json:"is_at_all"`
}

func (robot *WorkWechatRobot) SendText(req *SentTextRequest) {
	messages, params := make([]string, 0), make(map[string]interface{})
	messages = append(messages, fmt.Sprintf("- 时间：%s", time.Now().Format("2006-01-02 15:04:05")))
	messages = append(messages, fmt.Sprintf("- Level：%s", req.Level))
	messages = append(messages, fmt.Sprintf("- 信息：%s", req.Content))
	text := make(map[string]interface{})
	text["content"] = strings.Join(messages, "\n")
	if req.IsAtAll {
		text["mentioned_list"] = []string{"@all"}
		text["mentioned_mobile_list"] = []string{"@all"}
	}
	params["text"] = text
	data, _ := json.Marshal(params)
	request, err := http.NewRequest("POST", robot.url, bytes.NewReader(data))
	if err != nil {
		log.Printf("NewRequest fail, %s\n", err)
	}
	request.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	_, err = client.Do(request)
	if err != nil {
		log.Printf("Request WeChat Api fail, %s\n", err)
	}
}
