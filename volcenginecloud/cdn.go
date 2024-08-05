package volcenginecloud

import (
	"github.com/volcengine/volc-sdk-golang/service/cdn"
)

type CdnClient struct {
	instance *cdn.CDN
}

func NewCdnClient(accessKey, secretKey string) *CdnClient {
	instance := cdn.NewInstance()
	instance.Client.SetAccessKey(accessKey)
	instance.Client.SetSecretKey(secretKey)
	return &CdnClient{instance: instance}
}

// ListDomains 获取加速域名列表
func (c CdnClient) ListDomains(req *cdn.ListCdnDomainsRequest, options ...cdn.OptionArg) (*cdn.ListCdnDomainsResponse, error) {
	return c.instance.ListCdnDomains(req, options...)
}

// DescribeCdnConfig 获取加速域名配置
func (c CdnClient) DescribeCdnConfig(req *cdn.DescribeCdnConfigRequest) (*cdn.DescribeCdnConfigResponse, error) {
	return c.instance.DescribeCdnConfig(req)
}

// UpdateCdnConfig 修改加速域名配置
func (c CdnClient) UpdateCdnConfig(req *cdn.UpdateCdnConfigRequest) (*cdn.UpdateCdnConfigResponse, error) {
	return c.instance.UpdateCdnConfig(req)

}

// AddCertificate 上传证书
func (c CdnClient) AddCertificate(req *cdn.AddCdnCertificateRequest) (*cdn.AddCdnCertificateResponse, error) {
	return c.instance.AddCdnCertificate(req)
}

// BatchDeployCert 关联证书与加速域名
func (c CdnClient) BatchDeployCert(req *cdn.BatchDeployCertRequest) (*cdn.BatchDeployCertResponse, error) {
	return c.instance.BatchDeployCert(req)
}

// SubmitRefreshTask 提交刷新任务
// 节流限制：您每秒最多可以发送 20 个请求。
//
//	默认情况下，每个火山引擎主账号的刷新额度如下：
//	每天最多刷新 1,000 个 URL。
//	每天最多刷新 50 个目录。
//	每个任务可以包含 100 个 URL 或目录。
func (c CdnClient) SubmitRefreshTask(req *cdn.SubmitRefreshTaskRequest) (*cdn.SubmitRefreshTaskResponse, error) {
	return c.instance.SubmitRefreshTask(req)
}

// SubmitPreloadTask 提交预热任务
// 节流限制：您每秒最多可以发送 20 个请求。
//
//	默认情况下，每个火山引擎主账号的任务额度如下：
//	每天最多预热 1,000 个 URL。
//	每个任务最多预热 100 个 URL。
func (c CdnClient) SubmitPreloadTask(req *cdn.SubmitPreloadTaskRequest) (*cdn.SubmitPreloadTaskResponse, error) {
	return c.instance.SubmitPreloadTask(req)
}

// DescribeContentTasks 获取刷新与预热任务列表
func (c CdnClient) DescribeContentTasks(req *cdn.DescribeContentTasksRequest) (*cdn.DescribeContentTasksResponse, error) {
	return c.instance.DescribeContentTasks(req)
}

// DescribeContentQuota 获取刷新、预热、封禁、解封的配额
func (c CdnClient) DescribeContentQuota(options ...cdn.OptionArg) (*cdn.DescribeContentQuotaResponse, error) {
	return c.instance.DescribeContentQuota(options...)
}

// DescribeCdnAccessLog 获取访问日志的下载链接
func (c CdnClient) DescribeCdnAccessLog(req *cdn.DescribeCdnAccessLogRequest, options ...cdn.OptionArg) (*cdn.DescribeCdnAccessLogResponse, error) {
	return c.instance.DescribeCdnAccessLog(req, options...)
}
