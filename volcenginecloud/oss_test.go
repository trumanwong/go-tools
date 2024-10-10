package volcenginecloud

import (
	"context"
	"encoding/base64"
	"github.com/volcengine/ve-tos-golang-sdk/v2/tos"
	"log"
	"os"
	"testing"
)

func TestOssClient_PreSignedPostSignature(t *testing.T) {
	client, err := NewOssClient(os.Getenv("VOLC_ACCESS_KEY"), os.Getenv("VOLC_SECRET_KEY"), os.Getenv("VOLC_OSS_ENDPOINT"), os.Getenv("VOLC_REGION"))
	if err != nil {
		t.Error(err)
		return
	}

	callback := `
    {
       "callbackBody" : "{\"bucket\": ${bucket}, \"object\": ${object}, \"key1\": ${x:key1}}",
       "callbackBodyType" : "application/json",
       "callbackUrl" : "` + os.Getenv("VOLC_OSS_CALLBACK_URL") + `"
    }`

	callbackVar := `
    {
       "x:key1" : "ceshi"
    }`

	log.Println(callback)
	log.Println(callbackVar)
	output, err := client.PreSignedPostSignature(context.Background(), &tos.PreSingedPostSignatureInput{
		Bucket:  os.Getenv("VOLC_OSS_BUCKET"),
		Key:     "test.jpg",
		Expires: 3600,
		Conditions: []tos.PostSignatureCondition{
			{
				Key:   "x-tos-callback",
				Value: base64.StdEncoding.EncodeToString([]byte(callback)),
			},
			{
				Key:   "x-tos-callback-var",
				Value: base64.StdEncoding.EncodeToString([]byte(callbackVar)),
			},
		},
	})

	log.Println(base64.StdEncoding.EncodeToString([]byte(callback)))
	log.Println(base64.StdEncoding.EncodeToString([]byte(callbackVar)))
	if err != nil {
		t.Error(err)
		return
	}
	log.Println(output.OriginPolicy)
	log.Println("policy: ", output.Policy)
	log.Println("x-tos-algorithm", output.Algorithm)
	log.Println("x-tos-credential", output.Credential)
	log.Println("x-tos-date", output.Date)
	log.Println("x-tos-signature", output.Signature)
}
