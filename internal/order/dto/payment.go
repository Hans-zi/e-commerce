package dto

type Payment struct {
	OrderID string  `json:"order_id"`
	Amount  float64 `json:"amount"`
	Method  string  `json:"method"` // alipay, wechat
	Status  string  `json:"status"`
	Url     string  `json:"url"`
}

// AlipayNotify 支付宝异步通知结构体
type CallBackRes struct {
	OutTradeNo  string `json:"out_trade_no"` // 你的订单号
	TradeNo     string `json:"trade_no"`     // 支付宝交易号
	TotalAmount string `json:"total_amount"` // 支付金额
	TradeStatus string `json:"trade_status"` // 交易状态
	GmtPayment  string `json:"gmt_payment"`  // 支付时间
}
