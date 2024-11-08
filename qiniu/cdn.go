package qiniu

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/cdn"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type CdnClient struct {
	credentials *auth.Credentials
	cdnManager  *cdn.CdnManager
}

func NewCdnClient(accessKey, secretKey string) *CdnClient {
	mac := qbox.NewMac(accessKey, secretKey)
	credentials := auth.New(accessKey, secretKey)
	cdnManager := cdn.NewCdnManager(mac)
	return &CdnClient{
		credentials: credentials,
		cdnManager:  cdnManager,
	}
}

type GetDomainRequest struct {
	Name string
}

type GetDomainResponse struct {
	Name                   string            `json:"name"`
	PareDomain             string            `json:"pareDomain"`
	Type                   string            `json:"type"`
	Cname                  string            `json:"cname"`
	TestURLPath            string            `json:"testURLPath"`
	Protocol               string            `json:"protocol"`
	Platform               string            `json:"platform"`
	Product                string            `json:"product"`
	GeoCover               string            `json:"geoCover"`
	QiniuPrivate           bool              `json:"qiniuPrivate"`
	OperationType          string            `json:"operationType"`
	OperatingState         string            `json:"operatingState"`
	OperatingStateDesc     string            `json:"operatingStateDesc"`
	FreezeType             string            `json:"freezeType"`
	CreateAt               time.Time         `json:"createAt"`
	ModifyAt               time.Time         `json:"modifyAt"`
	CouldOperateBySelf     bool              `json:"couldOperateBySelf"`
	UIDIsFreezed           bool              `json:"uidIsFreezed"`
	OemMail                string            `json:"oemMail"`
	TagList                interface{}       `json:"tagList"`
	KvTagList              interface{}       `json:"kvTagList"`
	IPTypes                int               `json:"ipTypes"`
	DeliveryBucket         string            `json:"deliveryBucket"`
	DeliveryBucketType     string            `json:"deliveryBucketType"`
	DeliveryBucketFop      DeliveryBucketFop `json:"deliveryBucketFop"`
	IsSufy                 bool              `json:"isSufy"`
	IsPcdnBackup           bool              `json:"isPcdnBackup"`
	IsPcdnBackup302        bool              `json:"isPcdnBackup302"`
	Source                 Source            `json:"source"`
	Bsauth                 Bsauth            `json:"bsauth"`
	External               External          `json:"external"`
	Cache                  Cache             `json:"cache"`
	Referer                Referer           `json:"referer"`
	TimeACL                TimeACL           `json:"timeACL"`
	IPACL                  IPACL             `json:"ipACL"`
	UaACL                  interface{}       `json:"uaACL"`
	RequestHeaders         interface{}       `json:"requestHeaders"`
	ResponseHeaderControls []interface{}     `json:"responseHeaderControls"`
	HTTPS                  HTTPS             `json:"https"`
	RegisterNo             string            `json:"registerNo"`
	ConfigProcessRaTio     int               `json:"configProcessRa   tio"`
	HurryUpFreecert        bool              `json:"hurryUpFreecert"`
	HTTPSOPTime            time.Time         `json:"httpsOPTime"`
	Range                  Range             `json:"range"`
	OperTaskID             string            `json:"operTaskId"`
	OperTaskType           string            `json:"operTaskType"`
	OperTaskErrCode        int               `json:"operTaskErrCode"`
}
type DeliveryBucketFop struct {
	Enable           bool        `json:"enable"`
	SufyDeliveryHost string      `json:"sufyDeliveryHost"`
	NewStyle         interface{} `json:"newStyle"`
	DeleteStyleNames interface{} `json:"deleteStyleNames"`
	NewSeparator     interface{} `json:"newSeparator"`
}
type Range struct {
	Enable string `json:"enable"`
}
type SourceAuthInfo struct {
	OssProvider     string `json:"ossProvider"`
	AccessKeyID     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
}
type Source struct {
	SourceCname                  string         `json:"sourceCname"`
	SourceType                   string         `json:"sourceType"`
	SourceHost                   string         `json:"sourceHost"`
	TestSourceHost               string         `json:"testSourceHost"`
	SourceIPs                    []string       `json:"sourceIPs"`
	SourceDomain                 string         `json:"sourceDomain"`
	SourceQiniuBucket            string         `json:"sourceQiniuBucket"`
	SourceURLScheme              string         `json:"sourceURLScheme"`
	AdvancedSources              interface{}    `json:"advancedSources"`
	SkipCheckSource              bool           `json:"skipCheckSource"`
	TransferConf                 interface{}    `json:"transferConf"`
	SourceTimeACL                bool           `json:"sourceTimeACL"`
	SourceTimeACLKeys            interface{}    `json:"sourceTimeACLKeys"`
	MaxSourceRate                int            `json:"maxSourceRate"`
	MaxSourceConcurrency         int            `json:"maxSourceConcurrency"`
	AddRespHeader                interface{}    `json:"addRespHeader"`
	URLRewrites                  interface{}    `json:"urlRewrites"`
	SourceRetryCodes             interface{}    `json:"sourceRetryCodes"`
	FollowRedirect               bool           `json:"followRedirect"`
	Redirect30X                  interface{}    `json:"redirect30x"`
	SourceRequestHeaderControls  interface{}    `json:"sourceRequestHeaderControls"`
	SourceResponseHeaderControls interface{}    `json:"sourceResponseHeaderControls"`
	MaxSourceRatePerIDC          int            `json:"maxSourceRatePerIDC"`
	MaxSourceConcurrencyPerIDC   int            `json:"maxSourceConcurrencyPerIDC"`
	TestURLPath                  string         `json:"testURLPath"`
	SourceIgnoreParams           []interface{}  `json:"sourceIgnoreParams"`
	SourceIgnoreAllParams        bool           `json:"sourceIgnoreAllParams"`
	Range                        Range          `json:"range"`
	EnableSourceAuth             bool           `json:"enableSourceAuth"`
	SourceAuthInfo               SourceAuthInfo `json:"sourceAuthInfo"`
}
type UserAuthURLIPLimitConf struct {
	Enable   bool `json:"enable"`
	Limit    int  `json:"limit"`
	TimeSlot int  `json:"timeSlot"`
}
type UserAuthReqConf struct {
	Body                       interface{} `json:"body"`
	Header                     interface{} `json:"header"`
	Urlquery                   interface{} `json:"urlquery"`
	IncludeClientHeadersInBody bool        `json:"includeClientHeadersInBody"`
}
type UserAuthRespBodyConf struct {
	Enable                 bool        `json:"enable"`
	ContentType            string      `json:"contentType"`
	SuccessConditions      interface{} `json:"successConditions"`
	SuccessLogicalOperator string      `json:"successLogicalOperator"`
	FailureConditions      interface{} `json:"failureConditions"`
	FailureLogicalOperator string      `json:"failureLogicalOperator"`
}
type UserBsauthResultCacheConf struct {
	CacheEnable     bool        `json:"cacheEnable"`
	CacheSingleType string      `json:"cacheSingleType"`
	CacheKeyElems   interface{} `json:"cacheKeyElems"`
	CacheShareHost  string      `json:"cacheShareHost"`
	CacheDuration   int         `json:"cacheDuration"`
}
type UserAuthMatchRuleConf struct {
	Type string `json:"type"`
	Rule string `json:"rule"`
}
type Bsauth struct {
	Path                           interface{}               `json:"path"`
	Method                         string                    `json:"method"`
	Parameters                     interface{}               `json:"parameters"`
	TimeLimit                      int                       `json:"timeLimit"`
	UserAuthURL                    string                    `json:"userAuthUrl"`
	Strict                         bool                      `json:"strict"`
	Enable                         bool                      `json:"enable"`
	SuccessStatusCode              int                       `json:"successStatusCode"`
	FailureStatusCode              int                       `json:"failureStatusCode"`
	IsQiniuPrivate                 bool                      `json:"isQiniuPrivate"`
	BackSourceWithResourcePath     bool                      `json:"backSourceWithResourcePath"`
	BackSourceWithoutClientHeaders bool                      `json:"backSourceWithoutClientHeaders"`
	ResponseWithSourceAuthCode     bool                      `json:"responseWithSourceAuthCode"`
	ResponseWithSourceAuthBody     bool                      `json:"responseWithSourceAuthBody"`
	UserAuthURLIPLimitConf         UserAuthURLIPLimitConf    `json:"userAuthUrlIpLimitConf"`
	UserAuthReqConf                UserAuthReqConf           `json:"userAuthReqConf"`
	UserAuthContentType            string                    `json:"userAuthContentType"`
	UserAuthRespBodyConf           UserAuthRespBodyConf      `json:"userAuthRespBodyConf"`
	UserBsauthResultCacheConf      UserBsauthResultCacheConf `json:"userBsauthResultCacheConf"`
	UserAuthMatchRuleConf          UserAuthMatchRuleConf     `json:"userAuthMatchRuleConf"`
}
type ImageSlim struct {
	EnableImageSlim  bool          `json:"enableImageSlim"`
	PrefixImageSlims []interface{} `json:"prefixImageSlims"`
	RegexpImageSlims []interface{} `json:"regexpImageSlims"`
}
type External struct {
	EnableFop bool      `json:"enableFop"`
	ImageSlim ImageSlim `json:"imageSlim"`
}
type CacheControls struct {
	Time     int    `json:"time"`
	Timeunit int    `json:"timeunit"`
	Type     string `json:"type"`
	Rule     string `json:"rule"`
}
type Cache struct {
	CacheControls []CacheControls `json:"cacheControls"`
	IgnoreParam   bool            `json:"ignoreParam"`
	IgnoreParams  []interface{}   `json:"ignoreParams"`
	IncludeParams []interface{}   `json:"includeParams"`
}
type Referer struct {
	RefererType   string        `json:"refererType"`
	RefererValues []interface{} `json:"refererValues"`
	NullReferer   bool          `json:"nullReferer"`
}
type Verification struct {
	Name     string `json:"name"`
	Locate   string `json:"locate"`
	FailCode int    `json:"failCode"`
}
type TimeACL struct {
	Enable                bool         `json:"enable"`
	TimeACLKeys           interface{}  `json:"timeACLKeys"`
	AuthType              string       `json:"authType"`
	AuthDelta             int          `json:"authDelta"`
	SufyTimeACLKeys       interface{}  `json:"sufyTimeACLKeys"`
	SufyCallbackBody      interface{}  `json:"sufyCallbackBody"`
	CheckURL              string       `json:"checkUrl"`
	AdvanceFunctionEnable bool         `json:"advanceFunctionEnable"`
	RuleType              string       `json:"ruleType"`
	Rules                 interface{}  `json:"rules"`
	Params                interface{}  `json:"params"`
	ParamStr              string       `json:"paramStr"`
	ToLowerCase           string       `json:"toLowerCase"`
	URLEncode             string       `json:"urlEncode"`
	HashMethod            string       `json:"hashMethod"`
	Verification          Verification `json:"verification"`
}
type IPACL struct {
	IPACLType   string   `json:"ipACLType"`
	IPACLValues []string `json:"ipACLValues"`
}
type HTTPS struct {
	CertID      string `json:"certId"`
	ForceHTTPS  bool   `json:"forceHttps"`
	HTTP2Enable bool   `json:"http2Enable"`
	FreeCert    bool   `json:"freeCert"`
}

func (c CdnClient) GetDomain(request *GetDomainRequest) (*GetDomainResponse, error) {
	resp, err := requestQiniu(&Request{
		Method:      http.MethodGet,
		ApiUrl:      Host + "/domain/" + request.Name,
		Body:        nil,
		Credentials: c.credentials,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %s", err)
	}
	var response GetDomainResponse
	err = json.Unmarshal(b, &response)
	if err != nil {
		return nil, fmt.Errorf("unmarshal response body [%s] failed: %s", b, err)
	}
	return &response, nil
}

type GetDomainsRequest struct {
	// 域名前缀
	Marker string `json:"marker,omitempty"`
	// 返回的最大域名个数。1~1000, 不填默认为 10
	Limit int `json:"limit,omitempty"`
}

type GetDomainsResponse struct {
	Marker  string    `json:"marker"`
	Domains []Domains `json:"domains"`
}

type Domains struct {
	Name               string `json:"name"`
	Type               string `json:"type"`
	Cname              string `json:"cname"`
	TestURLPath        string `json:"testURLPath"`
	Platform           string `json:"platform"`
	GeoCover           string `json:"geoCover"`
	Protocol           string `json:"protocol"`
	OperatingState     string `json:"operatingState"`
	OperatingStateDesc string `json:"operatingStateDesc"`
	CreateAt           string `json:"createAt"`
	ModifyAt           string `json:"modifyAt"`
}

func (c CdnClient) GetDomains(req *GetDomainsRequest) (*GetDomainsResponse, error) {
	body, _ := json.Marshal(req)
	resp, err := requestQiniu(&Request{
		Method:      http.MethodGet,
		ApiUrl:      Host + "/domain",
		Body:        body,
		Credentials: c.credentials,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %s", err)
	}
	var response GetDomainsResponse
	err = json.Unmarshal(b, &response)
	if err != nil {
		return nil, fmt.Errorf("unmarshal response body [%s] failed: %s", b, err)
	}
	return &response, nil

}

type ModifyHttpsConfRequest struct {
	// 域名
	Name string `json:"-"`
	// 证书id，从上传或者获取证书列表里拿到证书id
	CertID string `json:"certId"`
	// 是否强制https跳转
	ForceHttps bool `json:"forceHttps"`
	// http2功能是否启用，false为关闭，true为开启
	Http2Enable bool `json:"http2Enable"`
}

type ModifyHttpsConfResponse struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

func (c CdnClient) ModifyHttpsConf(req *ModifyHttpsConfRequest) (*ModifyHttpsConfResponse, error) {
	body, _ := json.Marshal(req)
	resp, err := requestQiniu(&Request{
		Method:      http.MethodPut,
		ApiUrl:      Host + fmt.Sprintf("/domain/%s/httpsconf", req.Name),
		Body:        body,
		Credentials: c.credentials,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %s", err)
	}
	var response ModifyHttpsConfResponse
	err = json.Unmarshal(b, &response)
	if err != nil {
		return nil, fmt.Errorf("unmarshal response body [%s] failed: %s", b, err)
	}
	return &response, nil
}

type FluxRequest struct {
	// 开始日期，例如：2016-07-01
	StartDate string `json:"startDate"`
	// 结束日期，例如：2016-07-03
	EndDate string `json:"endDate"`
	// 时间粒度，取值为5min/hour/day
	Granularity string `json:"granularity"`
	// 域名列表，以 ；分割
	Domains string `json:"domains"`
}

type FluxResponse struct {
	Code  int                             `json:"code"`
	Error string                          `json:"error"`
	Time  []string                        `json:"time"`
	Data  map[string]FluxResponseDataItem `json:"data"`
}

type FluxResponseDataItem struct {
	China   []int64 `json:"china"`
	Oversea []int64 `json:"oversea"`
}

func (c CdnClient) Flux(req *FluxRequest) (*FluxResponse, error) {
	body, _ := json.Marshal(req)
	resp, err := requestQiniu(&Request{
		Method:      http.MethodPost,
		ApiUrl:      fusionHost + "/v2/tune/flux",
		Body:        body,
		Credentials: c.credentials,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %s", err)
	}
	var response FluxResponse
	err = json.Unmarshal(b, &response)
	if err != nil {
		return nil, fmt.Errorf("unmarshal response body [%s] failed: %s", b, err)
	}
	return &response, nil
}

type GetTopTrafficIpRequest struct {
	Domains   []string `json:"domains"`
	Region    string   `json:"region"`
	StartDate string   `json:"startDate"`
	EndDate   string   `json:"endDate"`
}

type GetTopTrafficIpResponse struct {
	Code  int                         `json:"code"`
	Error string                      `json:"error"`
	Data  GetTopTrafficIpResponseData `json:"data"`
}

type GetTopTrafficIpResponseData struct {
	Traffic []int64  `json:"traffic"`
	Ips     []string `json:"ips"`
}

func (c CdnClient) GetTopTrafficIp(req *GetTopTrafficIpRequest) (*GetTopTrafficIpResponse, error) {
	body, _ := json.Marshal(req)
	log.Println(string(body))
	resp, err := requestQiniu(&Request{
		Method:      http.MethodPost,
		ApiUrl:      fusionHost + "/v2/tune/loganalyze/toptrafficip",
		Body:        body,
		Credentials: c.credentials,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %s", err)
	}
	var response GetTopTrafficIpResponse
	err = json.Unmarshal(b, &response)
	if err != nil {
		return nil, fmt.Errorf("unmarshal response body [%s] failed: %s", b, err)
	}
	return &response, nil
}

type GetTopCountIpRequest struct {
	Domains   []string `json:"domains"`
	Region    string   `json:"region"`
	StartDate string   `json:"startDate"`
	EndDate   string   `json:"endDate"`
}

type GetTopCountIpResponse struct {
	Code  int                       `json:"code"`
	Error string                    `json:"error"`
	Data  GetTopCountIpResponseData `json:"data"`
}

type GetTopCountIpResponseData struct {
	Count []int64  `json:"count"`
	Ips   []string `json:"ips"`
}

func (c CdnClient) GetTopCountIp(req *GetTopCountIpRequest) (*GetTopCountIpResponse, error) {
	body, _ := json.Marshal(req)
	resp, err := requestQiniu(&Request{
		Method:      http.MethodPost,
		ApiUrl:      fusionHost + "/v2/tune/loganalyze/topcountip",
		Body:        body,
		Credentials: c.credentials,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %s", err)
	}
	var response GetTopCountIpResponse
	err = json.Unmarshal(b, &response)
	if err != nil {
		return nil, fmt.Errorf("unmarshal response body [%s] failed: %s", b, err)
	}
	return &response, nil
}

type UpdateIpACLRequest struct {
	Domain string
	IpAcl  IPACL
}

type UpdateIpACLResponse struct{}

func (c CdnClient) UpdateIpACL(req *UpdateIpACLRequest) (*UpdateIpACLResponse, error) {
	body, _ := json.Marshal(req.IpAcl)
	resp, err := requestQiniu(&Request{
		Method:      http.MethodPut,
		ApiUrl:      Host + "/domain/" + req.Domain + "/ipacl",
		Body:        body,
		Credentials: c.credentials,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %s", err)
	}
	return &UpdateIpACLResponse{}, nil
}

// GetCdnLogList 获取CDN域名访问日志的下载链接
func (c CdnClient) GetCdnLogList(day string, domains []string) (listLogResult cdn.ListLogResult, err error) {
	return c.cdnManager.GetCdnLogList(day, domains)
}

// GetFluxData 方法用来获取域名访问流量数据
//
//	StartDate	string		必须	开始日期，例如：2016-07-01
//	EndDate		string		必须	结束日期，例如：2016-07-03
//	Granularity	string		必须	粒度，取值：5min ／ hour ／day
//	Domains		[]string	必须	域名列表
//	Opts                            非必须   可选项
func (c CdnClient) GetFluxData(startDate, endDate, granularity string, domainList []string, opts ...cdn.FluxOption) (cdn.TrafficResp, error) {
	return c.cdnManager.GetFluxData(startDate, endDate, granularity, domainList, opts...)
}

type CdnAccessLog struct {
	ClientIp string
	// 命中信息
	StatusHit string
	// 响应时间
	RespTime float64
	// 请求时间
	ReqTime string
	// 请求方法
	Method string
	// 请求URL
	Url string
	// 请求协议
	Protocol string
	// 请求状态码
	StatusCode int64
	// 	响应大小
	RespSize int64
	// HTTP请求头中的Referer。
	Referer string
	// 请求中 User-Agent 头部的值。如果请求不包含该头部，该字段的值是 -
	UserAgent string
}

func (c CdnClient) AnalyzeCdnAccessLog(logPath string, handler func(...interface{}) error) error {
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
					pattern := `(\S+) (\S+) (\d+) \[(.*?)\] "(\S+) (.*?) (\S+)" (\d+) (\d+) "(.*?)" "(.*?)"`
					re := regexp.MustCompile(pattern)
					// Find the matches
					info := re.FindStringSubmatch(str)
					if len(info) != 12 {
						return fmt.Errorf("invalid log format, expect 12 fields, get %d fields: %s", len(info), str)
					}

					respTime, err := strconv.ParseFloat(info[3], 64)
					if err != nil {
						return fmt.Errorf("invalid respTime: %s, expect float: %s", info[3], str)
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
						ClientIp:   info[1],
						StatusHit:  info[2],
						RespTime:   respTime,
						ReqTime:    info[4],
						Method:     info[5],
						Url:        info[6],
						Protocol:   info[7],
						StatusCode: statusCode,
						RespSize:   respSize,
						Referer:    info[10],
						UserAgent:  info[11],
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
