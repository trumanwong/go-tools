package aliyun

import "github.com/aliyun/alibaba-cloud-sdk-go/services/live"

type LiveClient struct {
	client *live.Client
}

func NewLiveClient(regionId, accessKeyId, accessKeySecret string) (*LiveClient, error) {
	client, err := live.NewClientWithAccessKey(regionId, accessKeyId, accessKeySecret)
	if err != nil {
		return nil, err
	}
	return &LiveClient{client: client}, nil
}

func (l LiveClient) DescribeLiveUserDomains(req *live.DescribeLiveUserDomainsRequest) (*live.DescribeLiveUserDomainsResponse, error) {
	return l.client.DescribeLiveUserDomains(req)
}

func (l LiveClient) DescribeLiveCertificateDetail(req *live.DescribeLiveCertificateDetailRequest) (*live.DescribeLiveCertificateDetailResponse, error) {
	return l.client.DescribeLiveCertificateDetail(req)
}

func (l LiveClient) DescribeLiveDomainCertificateInfo(req *live.DescribeLiveDomainCertificateInfoRequest) (*live.DescribeLiveDomainCertificateInfoResponse, error) {
	return l.client.DescribeLiveDomainCertificateInfo(req)
}

func (l LiveClient) SetLiveDomainCertificate(req *live.SetLiveDomainCertificateRequest) (*live.SetLiveDomainCertificateResponse, error) {
	return l.client.SetLiveDomainCertificate(req)
}
