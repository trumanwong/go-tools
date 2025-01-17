package aliyun

import (
	"bufio"
	"compress/gzip"
	"fmt"
	cdn20180510 "github.com/alibabacloud-go/cdn-20180510/v4/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type CdnClient struct {
	client *cdn20180510.Client
}

func NewCdnClient(config *openapi.Config) (*CdnClient, error) {
	client, err := cdn20180510.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &CdnClient{client: client}, nil
}

// DescribeCdnDomainDetail 查询指定加速域名的基本配置
func (c CdnClient) DescribeCdnDomainDetail(req *cdn20180510.DescribeCdnDomainDetailRequest, runtime *util.RuntimeOptions) (*cdn20180510.DescribeCdnDomainDetailResponse, error) {
	return c.client.DescribeCdnDomainDetailWithOptions(req, runtime)
}

// DescribeUserDomainsWithOptions 查询用户名下所有的域名与状态，支持域名模糊匹配过滤和域名状态过滤
func (c CdnClient) DescribeUserDomainsWithOptions(req *cdn20180510.DescribeUserDomainsRequest, runtime *util.RuntimeOptions) (*cdn20180510.DescribeUserDomainsResponse, error) {
	return c.client.DescribeUserDomainsWithOptions(req, runtime)
}

// DescribeDomainCertificateInfoWithOptions 查询域名证书信息
func (c CdnClient) DescribeDomainCertificateInfoWithOptions(req *cdn20180510.DescribeDomainCertificateInfoRequest, runtime *util.RuntimeOptions) (*cdn20180510.DescribeDomainCertificateInfoResponse, error) {
	return c.client.DescribeDomainCertificateInfoWithOptions(req, runtime)
}

func (c CdnClient) SetCdnDomainSSLCertificate(req *cdn20180510.SetCdnDomainSSLCertificateRequest, runtime *util.RuntimeOptions) (*cdn20180510.SetCdnDomainSSLCertificateResponse, error) {
	return c.client.SetCdnDomainSSLCertificateWithOptions(req, runtime)
}

func (c CdnClient) RefreshObjectCachesWithOptions(req *cdn20180510.RefreshObjectCachesRequest, runtime *util.RuntimeOptions) (*cdn20180510.RefreshObjectCachesResponse, error) {
	return c.client.RefreshObjectCachesWithOptions(req, runtime)
}

func (c CdnClient) DescribeDomainUsageData(req *cdn20180510.DescribeDomainUsageDataRequest, runtime *util.RuntimeOptions) (*cdn20180510.DescribeDomainUsageDataResponse, error) {
	return c.client.DescribeDomainUsageDataWithOptions(req, runtime)
}

func (c CdnClient) DescribeCdnDomainLogs(req *cdn20180510.DescribeCdnDomainLogsRequest, runtime *util.RuntimeOptions) (*cdn20180510.DescribeCdnDomainLogsResponse, error) {
	return c.client.DescribeCdnDomainLogsWithOptions(req, runtime)
}

// CdnAccessLog https://help.aliyun.com/zh/cdn/user-guide/download-logs?spm=a2c4g.11186623.help-menu-27099.d_2_7_0_1.40dd77db3L8Tlf
type CdnAccessLog struct {
	// 请求结束时间
	ReqTime string
	// 从用户请求携带的请求头X-Forwarded-For中提取左边第一个IP地址（即client_ip，客户端IP，如果客户端与CDN节点之间没有经过代理的话，等同于客户端与CDN节点建联使用的IP）。
	ClientIp string
	// 从用户请求携带的请求头X-Forwarded-For中提取左边第二个IP地址（即proxy_ip，代理IP，如果客户端与CDN节点之间没有经过代理的话，提取到的空值用-表示）。
	ProxyIp string
	// 请求响应时间，单位为毫秒
	RespTime float64
	// HTTP请求头中的Referer。
	Referer string
	// 请求方法
	Method string
	// 用户请求的URL链接。
	Url string
	// HTTP状态码。
	Status int64
	// 请求大小，单位为字节
	ReqSize int64
	// CDN 向用户传输的数据量，单位是 bytes。
	RespSize int64
	// 命中信息
	StatusHit string
	// 请求中 User-Agent 头部的值。如果请求不包含该头部，该字段的值是 -
	UserAgent string
	// CDN 响应中 Content-Type 头部的值。如果响应不包含该头部，该字段的值是 -。
	ContentType string
	// 建连IP地址
	IP string
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
					pattern := `\[(.*?)\] (\S+) (\S+) (\S+) "(.*?)" "(\S+) (.*?)" (\d+) (\d+) (\d+) (\S+) "(.*?)" "(.*?)" (\S+)`
					re := regexp.MustCompile(pattern)
					// Find the matches
					info := re.FindStringSubmatch(str)
					if len(info) != 15 {
						return fmt.Errorf("invalid log format, expect 15 fields, get %d fields: %s", len(info), str)
					}

					respTime, err := strconv.ParseFloat(info[4], 64)
					if err != nil {
						return fmt.Errorf("invalid respTime: %s, expect float: %s", info[4], str)
					}

					statusCode, err := strconv.ParseInt(info[8], 10, 64)
					if err != nil {
						return fmt.Errorf("invalid statusCode: %s, expect int: %s", info[8], str)
					}
					reqSize, err := strconv.ParseInt(info[9], 10, 64)
					if err != nil {
						return fmt.Errorf("invalid reqSize: %s, expect int: %s", info[9], str)
					}
					respSize, err := strconv.ParseInt(info[10], 10, 64)
					if err != nil {
						return fmt.Errorf("invalid respSize: %s, expect int: %s", info[10], str)
					}
					accessLog := CdnAccessLog{
						ReqTime:     info[1],
						ClientIp:    info[2],
						ProxyIp:     info[3],
						RespTime:    respTime,
						Referer:     info[5],
						Method:      info[6],
						Url:         info[7],
						Status:      statusCode,
						ReqSize:     reqSize,
						RespSize:    respSize,
						StatusHit:   info[11],
						UserAgent:   info[12],
						ContentType: info[13],
						IP:          info[14],
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

// DescribeDomainTopClientIpVisitWithOptions 获取ClientIP列表的排序数据
func (c CdnClient) DescribeDomainTopClientIpVisitWithOptions(req *cdn20180510.DescribeDomainTopClientIpVisitRequest, runtime *util.RuntimeOptions) (*cdn20180510.DescribeDomainTopClientIpVisitResponse, error) {
	return c.client.DescribeDomainTopClientIpVisitWithOptions(req, runtime)
}

// DescribeCdnDomainConfigsWithOptions 获取加速域名的配置信息
func (c CdnClient) DescribeCdnDomainConfigsWithOptions(req *cdn20180510.DescribeCdnDomainConfigsRequest, runtime *util.RuntimeOptions) (*cdn20180510.DescribeCdnDomainConfigsResponse, error) {
	return c.client.DescribeCdnDomainConfigsWithOptions(req, runtime)
}

type CdnConfigFunctions struct {
	FunctionArgs []CdnConfigFunctionArg `json:"functionArgs"`
	FunctionName string                 `json:"functionName"`
	ParentID     *string                `json:"parentId"`
}

type CdnConfigFunctionArg struct {
	ArgName  string `json:"argName"`
	ArgValue string `json:"argValue"`
}

func (c CdnClient) BatchSetCdnDomainConfig(req *cdn20180510.BatchSetCdnDomainConfigRequest, runtime *util.RuntimeOptions) (*cdn20180510.BatchSetCdnDomainConfigResponse, error) {
	return c.client.BatchSetCdnDomainConfigWithOptions(req, runtime)
}

func (c CdnClient) StopCdnDomainWithOptions(req *cdn20180510.StopCdnDomainRequest, runtime *util.RuntimeOptions) (*cdn20180510.StopCdnDomainResponse, error) {
	return c.client.StopCdnDomainWithOptions(req, runtime)
}

func (c CdnClient) StartCdnDomainWithOptions(req *cdn20180510.StartCdnDomainRequest, runtime *util.RuntimeOptions) (*cdn20180510.StartCdnDomainResponse, error) {
	return c.client.StartCdnDomainWithOptions(req, runtime)
}
