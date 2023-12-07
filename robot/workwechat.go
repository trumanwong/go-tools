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

// WorkWechatRobot is a struct that represents a Work WeChat robot.
// It contains the URL of the Work WeChat robot.
type WorkWechatRobot struct {
	// url is the URL of the Work WeChat robot.
	url string
}

// NewWorkWechatRobot is a function that creates a new WorkWechatRobot object.
// It takes the URL of the Work WeChat robot as a parameter,
// and returns a pointer to the created WorkWechatRobot object.
func NewWorkWechatRobot(url string) *WorkWechatRobot {
	return &WorkWechatRobot{url: url}
}

// Level is a type that represents the level of a message.
type Level string

// Constants for the different levels of a message.
const (
	LevelInfo    Level = "info"
	LevelWarning Level = "warning"
	LevelError   Level = "error"
)

// SentTextRequest is a struct that represents a request to send a text message.
// It contains the level of the message, the content of the message, and a boolean indicating whether to mention all users.
type SentTextRequest struct {
	// Level is the level of the message.
	Level Level `json:"level"`
	// Content is the content of the message.
	Content string `json:"content"`
	// IsAtAll is a boolean indicating whether to mention all users.
	IsAtAll bool `json:"is_at_all"`
}

// SendText is a method of WorkWechatRobot that sends a text message.
// It takes a pointer to a SentTextRequest struct as a parameter.
// The method constructs the message, creates a POST request with the message in the body, and sends the request to the Work WeChat robot.
// If there is an error in creating the request or sending the request, the method logs the error.
func (robot *WorkWechatRobot) SendText(req *SentTextRequest) {
	// Create the message and the parameters for the request.
	messages, params := make([]string, 0), make(map[string]interface{})
	// Add the time, level, and content to the message.
	messages = append(messages, fmt.Sprintf("- 时间：%s", time.Now().Format("2006-01-02 15:04:05")))
	messages = append(messages, fmt.Sprintf("- Level：%s", req.Level))
	messages = append(messages, fmt.Sprintf("- 信息：%s", req.Content))
	// Create the text for the request.
	text := make(map[string]interface{})
	text["content"] = strings.Join(messages, "\n")
	// If IsAtAll is true, mention all users.
	if req.IsAtAll {
		text["mentioned_list"] = []string{"@all"}
		text["mentioned_mobile_list"] = []string{"@all"}
	}
	// Set the message type and the text in the parameters for the request.
	params["msgtype"] = "text"
	params["text"] = text
	// Marshal the parameters into JSON.
	data, _ := json.Marshal(params)
	// Create a POST request with the URL of the Work WeChat robot and the JSON data in the body.
	request, err := http.NewRequest("POST", robot.url, bytes.NewReader(data))
	if err != nil {
		// If there is an error in creating the request, log the error.
		log.Printf("NewRequest fail, %s\n", err)
	}
	// Set the Content-Type header to "application/json".
	request.Header.Set("Content-Type", "application/json")
	// Create a new HTTP client and send the request.
	client := http.Client{}
	_, err = client.Do(request)
	if err != nil {
		// If there is an error in sending the request, log the error.
		log.Printf("Request WeChat Api fail, %s\n", err)
	}
}
