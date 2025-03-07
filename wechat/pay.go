package wechat

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/h5"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/native"
	"github.com/wechatpay-apiv3/wechatpay-go/services/refunddomestic"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
	"time"
)

type WeChatPay struct {
	client *core.Client
}

// NewWeChatPay
// mchCertificateSerialNumber 商户证书序列号
// mchAPIv3Key 商户APIv3密钥
// apiClientKeyPath 商户私钥路径
func NewWeChatPay(ctx context.Context, mchID, mchCertificateSerialNumber, mchAPIv3Key, apiClientKeyPath string) (*WeChatPay, error) {
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath(apiClientKeyPath)
	if err != nil {
		return nil, err
	}

	// 使用商户私钥等初始化 client，并使它具有自动定时获取微信支付平台证书的能力
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(mchID, mchCertificateSerialNumber, mchPrivateKey, mchAPIv3Key),
	}
	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return &WeChatPay{client: client}, nil
}

// NativePrePay Native支付
func (p WeChatPay) NativePrePay(ctx context.Context, appId, mchId, title, tradeNo, notifyUrl string, amount int64) (*string, error) {
	svc := native.NativeApiService{Client: p.client}
	resp, _, err := svc.Prepay(ctx, native.PrepayRequest{
		Appid:       core.String(appId),
		Mchid:       core.String(mchId),
		Description: core.String(title),
		OutTradeNo:  core.String(tradeNo),
		NotifyUrl:   core.String(notifyUrl),
		Amount: &native.Amount{
			Total:    core.Int64(amount),
			Currency: core.String("CNY"),
		},
	})
	if err != nil {
		return nil, err
	}
	return resp.CodeUrl, nil
}

// H5PrePay h5支付
func (p WeChatPay) H5PrePay(ctx context.Context, appId, mchId, title, tradeNo, notifyUrl string, amount int64, payerClientIP *string) (*string, error) {
	svc := h5.H5ApiService{Client: p.client}
	resp, _, err := svc.Prepay(ctx, h5.PrepayRequest{
		Appid:       core.String(appId),
		Mchid:       core.String(mchId),
		Description: core.String(title),
		OutTradeNo:  core.String(tradeNo),
		NotifyUrl:   core.String(notifyUrl),
		Amount: &h5.Amount{
			Total:    core.Int64(amount),
			Currency: core.String("CNY"),
		},
		SceneInfo: &h5.SceneInfo{
			PayerClientIp: payerClientIP,
			H5Info: &h5.H5Info{
				Type: core.String("H5"),
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return resp.H5Url, nil
}

func (p WeChatPay) JsApiPay(ctx context.Context, appId, mchId, title, tradeNo, notifyUrl, openId string, amount int64) (*string, error) {
	svc := jsapi.JsapiApiService{Client: p.client}
	resp, _, err := svc.PrepayWithRequestPayment(ctx,
		jsapi.PrepayRequest{
			Appid:       core.String(appId),
			Mchid:       core.String(mchId),
			Description: core.String(title),
			OutTradeNo:  core.String(tradeNo),
			NotifyUrl:   core.String(notifyUrl),
			Amount: &jsapi.Amount{
				Total: core.Int64(amount),
			},
			Payer: &jsapi.Payer{
				Openid: core.String(openId),
			},
		},
	)
	if err != nil {
		return nil, err
	}
	b, _ := json.Marshal(&resp)
	result := string(b)
	return &result, nil
}

func (p WeChatPay) Refund(ctx context.Context, title, tradeNo, refundNo, reason string, refundAmount, total, refundQuantity int64) error {
	svc := refunddomestic.RefundsApiService{Client: p.client}
	resp, _, err := svc.Create(
		ctx,
		refunddomestic.CreateRequest{
			// SubMchid:    core.String(mchId),
			OutTradeNo:  core.String(tradeNo),
			OutRefundNo: core.String(refundNo),
			Reason:      core.String(reason),
			Amount: &refunddomestic.AmountReq{
				Currency: core.String("CNY"),
				//From: []refunddomestic.FundsFromItem{{
				//	Account: refunddomestic.ACCOUNT_AVAILABLE.Ptr(),
				//	Amount:  core.Int64(refundAmount),
				//}},
				Refund: core.Int64(refundAmount),
				Total:  core.Int64(total),
			},
			GoodsDetail: []refunddomestic.GoodsDetail{{
				GoodsName:       core.String(title),
				MerchantGoodsId: core.String(tradeNo),
				RefundAmount:    core.Int64(refundAmount),
				RefundQuantity:  core.Int64(refundQuantity),
				UnitPrice:       core.Int64(total),
			}},
		},
	)
	if err != nil {
		return fmt.Errorf("微信退款失败: %s", err)
	}
	if resp.Status == nil {
		return errors.New("微信退款状态异常: nil")
	}
	if *resp.Status != refunddomestic.STATUS_SUCCESS && *resp.Status != refunddomestic.STATUS_PROCESSING {
		return fmt.Errorf("微信退款状态异常: %s", *resp.Status)
	}
	return nil
}

func (p WeChatPay) QueryRefund(ctx context.Context, refundNo string) (*time.Time, error) {
	svc := refunddomestic.RefundsApiService{Client: p.client}
	resp, _, err := svc.QueryByOutRefundNo(
		ctx,
		refunddomestic.QueryByOutRefundNoRequest{
			OutRefundNo: core.String(refundNo),
		},
	)
	if err != nil {
		return nil, err
	}
	if resp.Status == nil {
		return nil, errors.New("查询退款状态失败")
	}
	if *resp.Status != refunddomestic.STATUS_SUCCESS {
		return nil, fmt.Errorf("查询退款状态失败，状态为%s", *resp.Status)
	}
	return resp.SuccessTime, nil
}
