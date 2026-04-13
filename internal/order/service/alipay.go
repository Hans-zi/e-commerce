package service

import (
	"context"
	"e-commerce/internal/consts"
	"e-commerce/internal/order/dto"
	"e-commerce/internal/order/model"
	"e-commerce/internal/order/repository"
	"e-commerce/pkg/mq"
	"e-commerce/pkg/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/bytedance/gopkg/util/logger"

	"github.com/smartwalle/alipay/v3"
)

type AlipayService interface {
	CreatePayment(orderID, userID string) (*model.Payment, error)
	CallBack(values url.Values) (*dto.CallBackRes, error)
}

type alipayService struct {
	paymentRepo repository.PaymentRepository
	orderRepo   repository.OrderRepository
	productRepo repository.ProductRepository
	alipay      *alipay.Client
	k           *mq.Kafka
}

func NewAlipayService(
	paymentRepo repository.PaymentRepository,
	orderRepo repository.OrderRepository,
	productRepo repository.ProductRepository,
	alipay *alipay.Client,
	k *mq.Kafka) AlipayService {

	s := &alipayService{
		paymentRepo: paymentRepo,
		orderRepo:   orderRepo,
		productRepo: productRepo,
		alipay:      alipay,
		k:           k,
	}
	go s.k.StartConsumeMessages(consts.TOPIC_REFUND_ORDER,
		consts.GROUP_ID_REFUND,
		s.Refund)
	return s
}

func (s *alipayService) CreatePayment(orderID, userID string) (*model.Payment, error) {
	//获取订单
	order, err := s.orderRepo.GetUserOrderByID(orderID, userID)
	if err != nil {
		logger.Errorf("CreatePayment.GetOrder fail: %s", err)
		return nil, err
	}

	//判断是否可支付
	if order.Status != consts.ORDER_STATUS_PENDING {
		err = errors.New("订单已支付或取消")
		logger.Errorf("CreatePayment.GetOrder fail: %s", err)
		return nil, err
	}

	var p alipay.TradePagePay
	p.Subject = fmt.Sprintf("订单支付-%s", orderID)
	p.TotalAmount = strconv.FormatFloat(order.Total, 'f', 2, 64)
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"
	p.OutTradeNo = orderID
	p.NotifyURL = consts.SERVER_DOMAIN + "/api/payment/alipay/notify"
	p.ReturnURL = ""

	url, err := s.alipay.TradePagePay(p)
	if err != nil {
		return nil, err
	}

	var payment model.Payment
	payment.OrderID = orderID
	payment.Method = consts.PAYMENT_METHOD_ALIPAY
	payment.Url = url.String()

	err = s.paymentRepo.Create(&payment)
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (s *alipayService) CallBack(values url.Values) (*dto.CallBackRes, error) {

	err := s.alipay.VerifySign(context.Background(), values)
	if err != nil {
		logger.Errorf("验签失败: %s", err)
		return nil, err
	}

	notification := &alipay.Notification{
		OutTradeNo:  values.Get("out_trade_no"),
		TradeNo:     values.Get("trade_no"),
		TotalAmount: values.Get("total_amount"),
		TradeStatus: alipay.TradeStatus(values.Get("trade_status")),
		GmtPayment:  values.Get("gmt_payment"),
	}

	if notification.TradeStatus != "TRADE_SUCCESS" {
		err := errors.New("交易支付失败")
		return nil, err

	}

	var res dto.CallBackRes
	//获取订单
	orderID := notification.OutTradeNo
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		logger.Errorf("AlipayCallback.GetOrder: %s", err)
		return nil, err
	}
	//更新订单
	order.Status = consts.ORDER_STATUS_PAID
	var paidAt = time.Now()
	order.PaidAt = &paidAt
	err = s.orderRepo.Update(order)
	if err != nil {
		logger.Errorf("AlipayCallback.UpdateOrder: %s", err)
		return nil, err
	}
	//获取支付
	payment, err := s.paymentRepo.GetByOrderID(orderID)
	if err != nil {
		logger.Errorf("AlipayCallback.GetPayment: %s", err)
		return nil, err
	}
	//更新支付
	payment.Status = consts.PAYMENT_STATUS_PAID
	payment.TransactionID = notification.TradeNo
	payment.Amount, _ = strconv.ParseFloat(notification.TotalAmount, 32)
	err = s.paymentRepo.Update(payment)
	if err != nil {
		logger.Errorf("AlipayCallback.UpdatePayment: %s", err)
		return nil, err
	}
	utils.Copy(&res, notification)
	return &res, nil
}

func (s *alipayService) AutoCancel(data []byte) error {
	logger.Infof("订单支付超时，自动取消")
	var msg dto.AutoCancelMsg
	err := json.Unmarshal(data, &msg)
	if err != nil {
		logger.Errorf("AutoCancel.Unmarshal fail: %s", err)
		return err
	}

	//获取支付
	payment, err := s.paymentRepo.GetByOrderID(msg.OrderID)
	if err != nil {
		return err
	}

	//确认是否能取消
	if payment.Status == consts.PAYMENT_STATUS_PAID {
		err = errors.New("订单已付款")
		return nil
	}

	order, err := s.orderRepo.GetUserOrderByID(msg.OrderID, msg.UserID)
	if err != nil {
		return err
	}
	//取消订单
	order.Status = consts.ORDER_STATUS_CANCELED

	//恢复库存
	err = s.productRepo.RestoreStockBatch(order.Lines)
	if err != nil {
		logger.Errorf("CancelOrder.RestoreStock fail: %s", err)
		return err
	}

	err = s.orderRepo.Update(order)
	if err != nil {
		logger.Errorf("CancelOrder.UpdateOrder fail: %s", err)
		return err
	}
	return nil
}

func (s *alipayService) Refund(data []byte) error {
	var msg dto.RefundMsg
	err := json.Unmarshal(data, &msg)
	if err != nil {
		logger.Errorf("Refund.Unmarshal fail: %s", err)
		return err
	}
	payment, err := s.paymentRepo.GetByOrderID(msg.OrderID)
	if err != nil {
		logger.Errorf("Refund.GetPayment fail: %s", err)
		return err
	}
	var req alipay.TradeRefund
	req.TradeNo = payment.TransactionID
	req.OutTradeNo = payment.OrderID
	req.RefundAmount = strconv.FormatFloat(payment.Amount, 'f', 2, 32)
	req.RefundReason = "测试退款"
	req.OutRequestNo = payment.OrderID + "_refund" // 幂等单号

	//调用支付宝退款
	resp, err := s.alipay.TradeRefund(context.Background(), req)
	if err != nil {
		logger.Errorf("Refund.TradeRefund fail: %s", err)
		return err
	}

	//验签 + 判断是否退款成功
	if !resp.IsSuccess() {
		logger.Errorf("Refund failed, code=%s msg=%s", resp.Code, resp.Msg)
		return errors.New("refund failed from alipay")
	}

	//更新支付状态为已退款
	payment.Status = consts.PAYMENT_STATUS_REFUNDED
	if err = s.paymentRepo.Update(payment); err != nil {
		logger.Errorf("Refund.UpdatePayment fail: %s", err)
		return err
	}

	logger.Infof("订单 %s 退款成功", payment.OrderID)
	return nil
}
