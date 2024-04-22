package qiniu

import (
	"os"
	"testing"
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
