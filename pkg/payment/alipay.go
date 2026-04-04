package payment

import (
	"e-commerce/pkg/utils"

	"github.com/smartwalle/alipay/v3"
)

func NewPayClient(config utils.Config) (*alipay.Client, error) {
	client, err := alipay.New(config.AliPay.AppID, config.AliPay.PrivateKey, false)
	if err != nil {
		return nil, err
	}
	err = client.LoadAliPayPublicKey(config.AliPay.AlipayPublicKey)
	if err != nil {
		return nil, err
	}

	return client, nil
}
