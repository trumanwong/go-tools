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

func (robot *WorkWechatRobot) SendText(level, content string) {
	messages, params := make([]string, 0), make(map[string]interface{})
	messages = append(messages, fmt.Sprintf("- 时间：%s", time.Now().Format("2006-01-02 15:04:05")))
	messages = append(messages, fmt.Sprintf("- Level：%s", level))
	messages = append(messages, fmt.Sprintf("- 信息：%s", content))
	markdown := make(map[string]interface{})
	markdown["title"] = "通知"
	markdown["content"] = strings.Join(messages, "\n")
	params["timestamp"] = time.Now().Unix()
	params["msgtype"] = "markdown"
	params["markdown"] = markdown
	if level == "error" {
		params["mentioned_mobile_list"] = []string{"@all"}
	}
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
