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

type GreenClient struct {
	client *green20220302.Client
}

// NewGreenClient 创建GreenClient实例
// 此处实例化的client请尽可能重复使用，避免重复建立连接，提升检测性能
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

type ImageModerationRequest struct {
	// 待检测图片链接，公网可访问的URL。
	ImageUrl *string
	// 使用OSS授权进行检测，必须同时传入ossBucketName、ossObjectName、ossRegionId。
	// 已授权OSS空间的Bucket名,示例：bucket001
	OssBucketName *string
	// 已授权OSS空间的文件名,示例：image/001.jpg
	OssObjectName *string
	// OSS Bucket所在区域，如oss-cn-hangzhou。
	OssRegionID *string
	// 风险标签, key为风险标签值，value为风险标签对应的分值，范围0-100，分值越高风险程度越高。
	// https://help.aliyun.com/document_detail/467829.html?spm=a2c4g.467828.0.0.2fd42592WnfIdf
	Labels map[string]float32
}

func (c GreenClient) ImageModeration(req *ImageModerationRequest) (err error) {
	//运行时参数设置，仅对使用了该运行时参数实例的请求有效
	runtime := &util.RuntimeOptions{}

	//构建图片检测请求。
	m := map[string]interface{}{
		//待检测数据的ID。
		"dataId": uuid.New(),
	}
	if req.ImageUrl != nil {
		m["imageUrl"] = *req.ImageUrl
	} else if req.OssBucketName != nil && req.OssObjectName != nil && req.OssRegionID != nil {
		//使用OSS授权进行检测
		m["ossBucketName"] = *req.OssBucketName
		m["ossObjectName"] = *req.OssObjectName
		m["ossRegionId"] = *req.OssRegionID
	}
	serviceParameters, _ := json.Marshal(m)
	imageModerationRequest := &green20220302.ImageModerationRequest{
		//图片检测service：内容安全控制台图片增强版规则配置的serviceCode，示例：baselineCheck
		Service:           tea.String("baselineCheck"),
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
