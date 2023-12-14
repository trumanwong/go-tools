package aliyun

import (
	"encoding/json"
	"errors"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	green20220302 "github.com/alibabacloud-go/green-20220302/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/google/uuid"
	"net/http"
)

// GreenClient is a struct that holds the client information for interacting with Alibaba Cloud's Green service.
type GreenClient struct {
	client *green20220302.Client // The Green service client
}

// NewGreenClient creates a new instance of GreenClient.
// It takes the accessKeyId, accessKeySecret, and endpoint as parameters.
// It returns a pointer to the created GreenClient instance and an error if any occurred during the creation.
// The created client should be reused as much as possible to avoid repeatedly establishing connections and improve detection performance.
func NewGreenClient(accessKeyId, accessKeySecret, endpoint string) (*GreenClient, error) {
	config := &openapi.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
		Endpoint:        tea.String(endpoint),
	}
	client, err := green20220302.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &GreenClient{client: client}, nil
}

// ImageModerationRequest is a struct that holds the information needed for image moderation.
type ImageModerationRequest struct {
	// The URL of the image to be checked, which must be publicly accessible.
	ImageUrl *string
	// To use OSS authorization for detection, ossBucketName, ossObjectName, and ossRegionId must be passed in at the same time.
	// The name of the authorized OSS space, for example: bucket001
	OssBucketName *string
	// The file name of the authorized OSS space, for example: image/001.jpg
	OssObjectName *string
	// The region where the OSS Bucket is located, such as oss-cn-hangzhou.
	OssRegionID *string
	// Risk labels, where the key is the value of the risk label and the value is the corresponding score, ranging from 0-100. The higher the score, the higher the risk level.
	// https://help.aliyun.com/document_detail/467829.html?spm=a2c4g.467828.0.0.2fd42592WnfIdf
	Labels map[string]float32
	// 检测规则，可选值：baselineCheck（基础版检测）
	Service *string
}

// ImageModeration performs image moderation.
// It takes an ImageModerationRequest as a parameter and returns an error if the moderation operation fails.
func (c GreenClient) ImageModeration(req *ImageModerationRequest) (err error) {
	// Runtime parameter settings, only effective for requests using this runtime parameter instance
	runtime := &util.RuntimeOptions{}

	// Build the image detection request.
	m := map[string]interface{}{
		// The ID of the data to be checked.
		"dataId": uuid.New(),
	}
	if req.ImageUrl != nil {
		m["imageUrl"] = *req.ImageUrl
	} else if req.OssBucketName != nil && req.OssObjectName != nil && req.OssRegionID != nil {
		// Use OSS authorization for detection
		m["ossBucketName"] = *req.OssBucketName
		m["ossObjectName"] = *req.OssObjectName
		m["ossRegionId"] = *req.OssRegionID
	}
	serviceParameters, _ := json.Marshal(m)
	service := tea.String("baselineCheck")
	if req.Service != nil {
		service = req.Service
	}
	imageModerationRequest := &green20220302.ImageModerationRequest{
		// Image detection service: the serviceCode configured by the content security console image enhanced version rules, example: baselineCheck
		Service:           service,
		ServiceParameters: tea.String(string(serviceParameters)),
	}

	response, err := c.client.ImageModerationWithOptions(imageModerationRequest, runtime)
	if err != nil {
		return err
	}
	statusCode := tea.IntValue(tea.ToInt(response.StatusCode))
	body := response.Body
	imageModerationResponseData := body.Data
	if statusCode == http.StatusOK {
		if tea.IntValue(tea.ToInt(body.Code)) == 200 {
			result := imageModerationResponseData.Result
			for i := 0; i < len(result); i++ {
				label := tea.StringValue(result[i].Label)
				confidence := tea.Float32Value(result[i].Confidence)
				if _, ok := req.Labels[label]; ok && confidence > req.Labels[label] {
					return errors.New("invalid image. label:" + label + ", confidence:" + tea.ToString(confidence))
				}
			}
			return nil
		}
		return errors.New("image moderation not success. status" + fmt.Sprintf("%d", tea.IntValue(tea.ToInt(body.Code))))
	}
	return errors.New("image moderation failed. statusCode:" + tea.ToString(statusCode))
}
