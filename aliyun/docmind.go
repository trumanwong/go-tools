package aliyun

import (
	openClient "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/docmind-api-20220711/client"
	"github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"io"
	"strings"
)

// DocMindClient is a struct that holds a client instance for interacting with the DocMind API.
type DocMindClient struct {
	client *client.Client
}

// NewDocMindClient is a function that initializes a new DocMindClient. It takes in an accessKeyId, accessKeySecret, and endpoint as parameters.
// It returns a pointer to a DocMindClient and an error.
func NewDocMindClient(accessKeyId, accessKeySecret, endpoint string) (*DocMindClient, error) {
	config := openClient.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
		Endpoint:        tea.String(endpoint),
	}
	cli, err := client.NewClient(&config)
	if err != nil {
		return nil, err
	}
	return &DocMindClient{client: cli}, nil
}

// SubmitConvertPdfToWordJobAdvance is a method on the DocMindClient struct that submits a job to convert a PDF to a Word document.
// It takes in a fileName and a reader as parameters. The reader is expected to be a stream of the PDF file.
// It returns a pointer to a SubmitConvertPdfToWordJobResponse and an error.
func (d DocMindClient) SubmitConvertPdfToWordJobAdvance(fileName string, reader io.Reader) (*client.SubmitConvertPdfToWordJobResponse, error) {
	req := client.SubmitConvertPdfToWordJobAdvanceRequest{
		FileName:      tea.String(fileName),
		FileUrlObject: reader,
	}

	// 创建RuntimeObject实例并设置运行参数
	options := service.RuntimeOptions{}
	return d.client.SubmitConvertPdfToWordJobAdvance(&req, &options)
}

// GetDocumentConvertResult is a method on the DocMindClient struct that retrieves the result of a document conversion job.
// It takes in a requestId as a parameter, which is the ID of the job to retrieve the result for.
// It returns a pointer to a GetDocumentConvertResultResponse and an error.
// https://help.aliyun.com/zh/document-mind/developer-reference/convertpdftoword?spm=a2c4g.11186623.0.0.778453cedQWluW#fba5587e2duz1
func (d DocMindClient) GetDocumentConvertResult(requestId string) (*client.GetDocumentConvertResultResponse, error) {
	req := client.GetDocumentConvertResultRequest{Id: &requestId}
	resp, err := d.client.GetDocumentConvertResult(&req)
	if err != nil {
		return nil, err
	}
	if resp.Body != nil && resp.Body.Data != nil {
		for i, v := range resp.Body.Data {
			if v.Url != nil {
				//注意，使用go语言sdk会将返回结果中Url中的特殊字符&做Unicode转码，转码成\u0026，需要调用者手动将\u0026进行Unicode转码回&，才可以正常下载URL
				link := strings.ReplaceAll(*v.Url, `\u0026`, `&`)
				resp.Body.Data[i].Url = &link
			}
		}
	}
	return resp, nil
}
