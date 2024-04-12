package paypal

import (
	"errors"
	"fmt"
	"github.com/trumanwong/go-tools/crawler"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	// VerifyURI Production Postback URL
	VerifyURI = "https://ipnpb.paypal.com/cgi-bin/webscr"
	// SandboxVerifyURI Sandbox Postback URL
	SandboxVerifyURI = "https://ipnpb.sandbox.paypal.com/cgi-bin/webscr"
	// Valid Response from PayPal indicating validation was successful
	Valid = "VERIFIED"
	// Invalid Response from PayPal indicating validation failed
	Invalid = "INVALID"
)

type IPN struct {
	useSandbox    bool
	useLocalCerts bool
	postData      map[string]string
}

func NewPaypalIPN(body []byte) (*IPN, error) {
	postData, err := readBody(body)
	if err != nil {
		return nil, fmt.Errorf("read body failed: %v", err)
	}

	if len(postData) == 0 {
		return nil, errors.New("missing POST data")
	}

	return &IPN{
		useSandbox:    false,
		useLocalCerts: true,
		postData:      postData,
	}, nil
}

func (p *IPN) UseSandbox() {
	p.useSandbox = true
}

func (p *IPN) GetPaypalURI() string {
	if p.useSandbox {
		return SandboxVerifyURI
	}
	return VerifyURI
}

func readBody(body []byte) (map[string]string, error) {
	var err error
	payload := make(map[string]string)
	arr := strings.Split(string(body), "&")
	for _, v := range arr {
		keyVal := strings.Split(v, "=")
		if len(keyVal) != 2 {
			continue
		}
		payload[keyVal[0]], err = url.QueryUnescape(keyVal[1])
		if err != nil {
			return nil, fmt.Errorf("could not unescape: %v", err)
		}
	}
	return payload, nil
}

func (p *IPN) VerifyIPN() (bool, error) {
	reqBody := "cmd=_notify-validate"
	for key, value := range p.postData {
		reqBody += "&" + key + "=" + url.QueryEscape(value)
	}

	resp, err := crawler.Send(&crawler.Request{
		Url:    p.GetPaypalURI(),
		Method: http.MethodPost,
		Body:   strings.NewReader(reqBody),
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
			"User-Agent":   "Go-IPN-Verification-Script",
			"Connection":   "Close",
		},
	})

	if err != nil {
		return false, fmt.Errorf("request %s fail, error: %s", p.GetPaypalURI(), err.Error())
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return false, errors.New("PayPal responded with HTTP code " + resp.Status)
	}

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	if string(resBody) == Valid {
		return true, nil
	}
	return false, fmt.Errorf("PayPal responded with %s", string(resBody))
}
