package alipay

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay"
)

type Pay struct {
	client *alipay.Client
}

func NewAliPay(appID, privateKey, notifyUrl, returnUrl string, isProd bool) (*Pay, error) {
	var client, err = alipay.NewClient(appID, privateKey, isProd)
	if err != nil {
		return nil, err
	}
	client.SetSignType(alipay.RSA2).
		//SetLocation(alipay.LocationShanghai).
		SetNotifyUrl(notifyUrl).
		SetReturnUrl(returnUrl)

	return &Pay{client: client}, nil
}

func (p Pay) TradePrecreate(ctx context.Context, subject, outTradeNo, amount, returnUrl, notifyUrl string) (*string, error) {
	bm := gopay.BodyMap{}
	bm.Set("subject", subject).
		Set("out_trade_no", outTradeNo).
		Set("total_amount", amount).
		Set("notify_url", notifyUrl).
		Set("return_url", returnUrl)
	resp, err := p.client.TradePrecreate(ctx, bm)
	if err != nil {
		return nil, err
	}
	return &resp.Response.QrCode, nil
}

func (p Pay) TradeWapPay(ctx context.Context, subject, outTradeNo, amount, returnUrl, notifyUrl string) (*string, error) {
	bm := gopay.BodyMap{}
	bm.Set("subject", subject).
		Set("out_trade_no", outTradeNo).
		Set("total_amount", amount).
		Set("return_url", returnUrl).
		Set("notify_url", notifyUrl)
	payUrl, err := p.client.TradeWapPay(ctx, bm)
	if err != nil {
		return nil, err
	}
	return &payUrl, nil
}

func (p Pay) Refund(ctx context.Context, outTradeNo, refundAmount, refundReason string) error {
	bm := gopay.BodyMap{}
	bm.Set("out_trade_no", outTradeNo).
		Set("refund_amount", refundAmount).
		Set("refund_reason", refundReason)
	resp, err := p.client.TradeRefund(ctx, bm)
	if err != nil {
		return errors.New(fmt.Sprintf("支付宝退款失败：%s", err))
	}
	if resp.Response.Code != "10000" {
		return errors.New(fmt.Sprintf("支付宝退款失败, msg: %s, code: %s, subMsg: %s, subCode: %s",
			resp.Response.Msg, resp.Response.Code, resp.Response.SubMsg, resp.Response.SubCode))
	}
	return nil
}
