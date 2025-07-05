package qiniu

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/sms"
	"github.com/qiniu/go-sdk/v7/sms/rpc"
)

type SmsClient struct {
	manager     *sms.Manager
	credentials *auth.Credentials
}

func NewSmsClient(accessKey, secretKey string) *SmsClient {
	credentials := auth.New(accessKey, secretKey)
	return &SmsClient{
		credentials: credentials,
		manager:     sms.NewManager(credentials),
	}
}

func (c SmsClient) SendMessage(signatureId, templateId string, mobiles []string, parameters map[string]any) (*string, error) {
	request := sms.MessagesRequest{
		SignatureID: signatureId,
		TemplateID:  templateId,
		Mobiles:     mobiles,
		Parameters:  parameters,
	}
	params, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", sms.Host+"/v1/message", bytes.NewReader(params))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.ContentLength = int64(len(params))
	err = c.credentials.AddToken(auth.TokenQiniu, req)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	var response sms.MessagesResponse
	if resp.StatusCode/100 == 2 || resp.StatusCode/100 == 3 {
		if resp.ContentLength > 0 {
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(b, &response)
			if err != nil {
				return nil, err
			}
			return &response.JobID, nil
		}
	}
	return nil, rpc.ResponseError(resp)
}
