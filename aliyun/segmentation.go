package aliyun

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	imageseg20191230 "github.com/alibabacloud-go/imageseg-20191230/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

type SegmentationClient struct {
	client *imageseg20191230.Client
}

func NewSegmentationClient(accessKeyId, accessKeySecret, endpoint string) (*SegmentationClient, error) {
	config := &openapi.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
		Endpoint:        tea.String(endpoint),
	}
	client, err := imageseg20191230.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &SegmentationClient{client: client}, nil
}

// SegmentHDCommonImage 通用高清分割
func (s SegmentationClient) SegmentHDCommonImage(req *imageseg20191230.SegmentHDCommonImageRequest) (*imageseg20191230.SegmentHDCommonImageResponse, error) {
	resp, err := s.client.SegmentHDCommonImage(req)
	return resp, err
}

// SegmentCommonImage 通用分割
func (s SegmentationClient) SegmentCommonImage(req *imageseg20191230.SegmentCommonImageRequest) (*imageseg20191230.SegmentCommonImageResponse, error) {
	resp, err := s.client.SegmentCommonImage(req)
	return resp, err
}

// SegmentCommonImage 通用分割
func (s SegmentationClient) SegmentCommonImageAdvance(req *imageseg20191230.SegmentCommonImageAdvanceRequest, runtime *util.RuntimeOptions) (*imageseg20191230.SegmentCommonImageResponse, error) {
	resp, err := s.client.SegmentCommonImageAdvance(req, runtime)
	return resp, err
}

// SegmentBody 人体分割
func (s SegmentationClient) SegmentBody(req *imageseg20191230.SegmentBodyRequest) (*imageseg20191230.SegmentBodyResponse, error) {
	resp, err := s.client.SegmentBody(req)
	return resp, err
}

// SegmentHDBody 高清人体分割
func (s SegmentationClient) SegmentHDBody(req *imageseg20191230.SegmentHDBodyRequest) (*imageseg20191230.SegmentHDBodyResponse, error) {
	resp, err := s.client.SegmentHDBody(req)
	return resp, err
}

// SegmentCommodity 商品分割
func (s SegmentationClient) SegmentCommodity(req *imageseg20191230.SegmentCommodityRequest) (*imageseg20191230.SegmentCommodityResponse, error) {
	resp, err := s.client.SegmentCommodity(req)
	return resp, err
}

func (s SegmentationClient) SegmentCommodityAdvance(req *imageseg20191230.SegmentCommodityAdvanceRequest, runtime *util.RuntimeOptions) (*imageseg20191230.SegmentCommodityResponse, error) {
	return s.client.SegmentCommodityAdvance(req, runtime)
}

// SegmentBodyAdvance 人体分割
func (s *SegmentationClient) SegmentBodyAdvance(req *imageseg20191230.SegmentBodyAdvanceRequest, runtime *util.RuntimeOptions) (*imageseg20191230.SegmentBodyResponse, error) {
	resp, err := s.client.SegmentBodyAdvance(req, runtime)
	return resp, err
}

func (s *SegmentationClient) GetAsyncJobResultWithOptions(req *imageseg20191230.GetAsyncJobResultRequest, runtime *util.RuntimeOptions) (*imageseg20191230.GetAsyncJobResultResponse, error) {
	resp, err := s.client.GetAsyncJobResultWithOptions(req, runtime)
	return resp, err
}
