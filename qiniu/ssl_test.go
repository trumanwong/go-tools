package qiniu

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestSslClient_CreateSslCert(t *testing.T) {
	sslClient := NewSslClient(os.Getenv("QINIU_ACCESS_KEY"), os.Getenv("QINIU_SECRET_KEY"))
	request := &CreateSslCertRequest{
		Name:       os.Getenv("QINIU_SSL_NAME"),
		CommonName: os.Getenv("QINIU_SSL_COMMON_NAME"),
		Pri:        os.Getenv("QINIU_SSL_PRI"),
		Ca:         os.Getenv("QINIU_SSL_CA"),
	}
	response, err := sslClient.CreateSslCert(nil, request)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(response)
}

func TestSslClient_GetSslCert(t *testing.T) {
	sslClient := NewSslClient(os.Getenv("QINIU_ACCESS_KEY"), os.Getenv("QINIU_SECRET_KEY"))
	response, err := sslClient.GetSslCert(context.Background(), &GetSslCertRequest{
		CertID: os.Getenv("QINIU_CERT_ID"),
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(response.Cert.NotAfter < time.Now().Add(time.Hour*24*60).Unix())
}
