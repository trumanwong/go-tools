package volcenginecloud

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"github.com/google/uuid"
	"github.com/trumanwong/cryptogo"
	"github.com/volcengine/volc-sdk-golang/service/cdn"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
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
func (c CdnClient) AddCertificate(req *cdn.AddCertificateRequest) (*cdn.AddCertificateResponse, error) {
	return c.instance.AddCertificate(req)
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

// DescribeCdnData 获取访问统计的细分数据
func (c CdnClient) DescribeCdnData(req *cdn.DescribeCdnDataRequest, options ...cdn.OptionArg) (*cdn.DescribeCdnDataResponse, error) {
	return c.instance.DescribeCdnData(req, options...)
}

// DescribeDistrictSummary 获取访问统计的汇总数据
func (c CdnClient) DescribeDistrictSummary(req *cdn.DescribeDistrictSummaryRequest) (*cdn.DescribeDistrictSummaryResponse, error) {
	return c.instance.DescribeDistrictSummary(req)
}

// DescribeDistrictData 获取访问统计的细分数据
func (c CdnClient) DescribeDistrictData(req *cdn.DescribeDistrictDataRequest) (*cdn.DescribeDistrictDataResponse, error) {
	return c.instance.DescribeDistrictData(req)
}

func (c CdnClient) DescribeStatisticalRanking(req *cdn.DescribeStatisticalRankingRequest) (*cdn.DescribeStatisticalRankingResponse, error) {
	return c.instance.DescribeStatisticalRanking(req)
}

// CdnAccessLog CDN访问日志
// Cdn日志字段说明：https://www.volcengine.com/docs/6454/71376
type CdnAccessLog struct {
	// 请求结束时间
	ReqTime string
	// 客户端ip
	ClientIp string
	// 请求响应时间
	RespTime float64
	// 请求使用的方法
	Method string
	// 请求 URL 中的 scheme
	Scheme string
	// 请求 URL 中的 域名
	Domain string
	// 请求 URL 中的路径和查询字符串，以斜杠（/）开头
	Url string
	// 请求使用的 HTTP 协议版本。
	HttpProtocol string
	// CDN 的响应状态码。
	StatusCode int64
	// CDN 向用户传输的数据量，单位是 bytes。
	RespSize int64
	// 请求是否命中 CDN 的缓存。
	BdStatusHit string
	// 请求中 Range 头部的值。如果请求不包含该头部，该字段的值是 -
	Range string
	// 请求中 Referer 头部的值。如果请求不包含该头部，该字段的值是 -
	Referer string
	// 请求中 User-Agent 头部的值。如果请求不包含该头部，该字段的值是 -
	UserAgent string
	// CDN 响应中 Content-Type 头部的值。如果响应不包含该头部，该字段的值是 -。
	ContentType string
}

func (c CdnClient) AnalyzeCdnAccessLog(logPath string, handler func(interface{}) error) error {
	fileInfoList, err := os.ReadDir(logPath)
	if err != nil {
		return err
	}

	for i := 0; i < len(fileInfoList); i++ {
		err = func(i int) error {
			file, err := os.Open(logPath + "/" + fileInfoList[i].Name())
			if err != nil {
				return err
			}
			defer file.Close()

			gz, err := gzip.NewReader(file)
			if err != nil {
				return err
			}
			defer gz.Close()
			scanner := bufio.NewScanner(gz)
			for scanner.Scan() {
				text := scanner.Text()
				strs := strings.Split(text, "\n")
				for _, str := range strs {
					info := strings.Split(str, "\t")
					if len(info) != 15 {
						return fmt.Errorf("invalid log format, expect 15 fields: %s", str)
					}

					respTime, err := strconv.ParseFloat(info[2], 64)
					if err != nil {
						return fmt.Errorf("invalid respTime: %s, expect float: %s", info[2], str)
					}

					statusCode, err := strconv.ParseInt(info[8], 10, 64)
					if err != nil {
						return fmt.Errorf("invalid statusCode: %s, expect int: %s", info[8], str)
					}
					respSize, err := strconv.ParseInt(info[9], 10, 64)
					if err != nil {
						return fmt.Errorf("invalid respSize: %s, expect int: %s", info[9], str)
					}
					accessLog := CdnAccessLog{
						ReqTime:      info[0],
						ClientIp:     info[1],
						RespTime:     respTime,
						Method:       info[3],
						Scheme:       info[4],
						Domain:       info[5],
						Url:          info[6],
						HttpProtocol: info[7],
						StatusCode:   statusCode,
						RespSize:     respSize,
						BdStatusHit:  info[10],
						Range:        info[11],
						Referer:      info[12],
						UserAgent:    info[13],
						ContentType:  info[14],
					}
					err = handler(&accessLog)
					if err != nil {
						return err
					}
				}
			}
			return nil
		}(i)

		if err != nil {
			return err
		}
	}
	return nil
}

// GenTypeAUrl generates a signed URL of type A.
// It takes a link, key, signName, uid, and timestamp as input parameters.
// The function parses the URL, generates a random string, and creates a hash using the provided parameters.
// It then constructs the signed URL with the generated hash and returns it.
// For more details, refer to the documentation: https://www.volcengine.com/docs/6454/99849#%E7%AD%BE%E5%90%8D%E8%AE%A1%E7%AE%97%E7%A4%BA%E4%BE%8B%E4%BB%A3%E7%A0%81
func GenTypeAUrl(link string, key string, signName string, uid string, timestamp int64) (string, error) {
	parseUrl, err := url.Parse(link)
	if err != nil {
		return "", err
	}

	randStr := strings.ReplaceAll(uuid.New().String(), "-", "")
	hash := cryptogo.MD5([]byte(fmt.Sprintf("%s-%d-%s-%s-%s", parseUrl.Path, timestamp, randStr, uid, key)))
	authArg := fmt.Sprintf("%s=%d-%s-%s-%s", signName, timestamp, randStr, uid, hash)
	parseUrl.RawQuery = fmt.Sprintf("%s&%s", parseUrl.RawQuery, authArg)
	return parseUrl.String(), nil
}

// splitUrl splits a URL into its components using a regular expression.
// It returns a slice of strings containing the matched components.
func splitUrl(url string) []string {
	reg := regexp.MustCompile("^(http://|https://)?([^/?]+)(/[^?]*)?(\\?.*)?$")
	return reg.FindStringSubmatch(url)
}

// GenTypeBUrl generates a signed URL of type B.
// It takes a link, key, and timestamp as input parameters.
// The function parses the URL, formats the timestamp, and creates a hash using the provided parameters.
// It then constructs the signed URL with the generated hash and returns it.
func GenTypeBUrl(link string, key string, timestamp int64) (string, error) {
	parseUrl, err := url.Parse(link)
	if err != nil {
		return "", err
	}
	tsStr := time.Unix(timestamp, 0).Format("200601021504")
	hash := cryptogo.MD5([]byte(fmt.Sprintf("%s%s%s", key, tsStr, parseUrl.Path)))
	return fmt.Sprintf(
		"%s://%s/%s/%s%s%s",
		parseUrl.Scheme,
		parseUrl.Host,
		tsStr,
		hash,
		parseUrl.Path,
		parseUrl.RawQuery,
	), nil
}

// GenTypeCUrl generates a signed URL of type C.
// It takes a link, key, and timestamp as input parameters.
// The function parses the URL, converts the timestamp to a hexadecimal string, and creates a hash using the provided parameters.
// It then constructs the signed URL with the generated hash and returns it.
func GenTypeCUrl(link string, key string, timestamp int64) (string, error) {
	parseUrl, err := url.Parse(link)
	if err != nil {
		return "", err
	}
	tsStr := strconv.FormatInt(timestamp, 16)
	hash := cryptogo.MD5([]byte(fmt.Sprintf("%s%s%s", key, parseUrl.Path, tsStr)))
	return fmt.Sprintf(
		"%s://%s/%s/%s%s%s",
		parseUrl.Scheme,
		parseUrl.Host,
		hash,
		tsStr,
		parseUrl.Path,
		parseUrl.RawQuery,
	), nil
}

// GenTypeDUrl generates a signed URL of type D.
// It takes a link, key, signName, timeName, timestamp, and base as input parameters.
// The function parses the URL, converts the timestamp to a string in the specified base, and creates a hash using the provided parameters.
// It then constructs the signed URL with the generated hash and returns it.
func GenTypeDUrl(link, key, signName, timeName string, timestamp int64, base int) (string, error) {
	parseUrl, err := url.Parse(link)
	if err != nil {
		return "", err
	}
	tsStr := strconv.FormatInt(timestamp, base)
	hash := cryptogo.MD5([]byte(fmt.Sprintf("%s%s%s", key, parseUrl.Path, tsStr)))
	authArg := fmt.Sprintf("%s=%s&%s=%s", signName, hash, timeName, tsStr)
	parseUrl.RawQuery = fmt.Sprintf("%s&%s", parseUrl.RawQuery, authArg)
	return parseUrl.String(), nil
}

// GenTypeEUrl generates a signed URL of type E by custom rule (e.g., key+domain+uri+timestamp).
// It takes a link, key, signName, tsName, timestamp, and base as input parameters.
// The function parses the URL, converts the timestamp to a string in the specified base, and creates a hash using the provided parameters.
// It then constructs the signed URL with the generated hash and returns it.
func GenTypeEUrl(link, key, signName, tsName string, timestamp int64, base int) (string, error) {
	parseUrl, err := url.Parse(link)
	if err != nil {
		return "", err
	}
	tsStr := strconv.FormatInt(timestamp, base)
	hash := cryptogo.MD5([]byte(fmt.Sprintf("%s%s%s%s", key, parseUrl.Host, parseUrl.Path, tsStr)))
	authArg := fmt.Sprintf("%s=%s&%s=%s", signName, hash, tsName, tsStr)
	parseUrl.RawQuery = fmt.Sprintf("%s&%s", parseUrl.RawQuery, authArg)
	return parseUrl.String(), nil
}
