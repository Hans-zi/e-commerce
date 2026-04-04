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
	//go s.k.StartConsumeMessages(consts.TOPIC_CANCEL_ORDER,
	//	consts.GROUP_ID_ORDER,
	//	s.AutoCancel)
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
	p.NotifyURL = "https://28a5ffd8.r17.cpolar.top/api/payment/alipay/notify"
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
	orderID := notification.OutTradeNo
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		logger.Errorf("AlipayCallback.GetOrder: %s", err)
		return nil, err
	}
	order.Status = consts.ORDER_STATUS_PAID
	err = s.orderRepo.Update(order)
	if err != nil {
		logger.Errorf("AlipayCallback.UpdateOrder: %s", err)
		return nil, err
	}
	payment, err := s.paymentRepo.GetByOrderID(orderID)
	if err != nil {
		logger.Errorf("AlipayCallback.GetPayment: %s", err)
		return nil, err
	}
	payment.Status = consts.PAYMENT_STATUS_PAID
	err = s.paymentRepo.Update(payment)
	if err != nil {
		logger.Errorf("AlipayCallback.UpdatePayment: %s", err)
		return nil, err
	}
	utils.Copy(&res, notification)
	return &res, nil
}

func (s *alipayService) AutoCancel(data []byte) error {
	t := time.NewTicker(time.Second)
	<-t.C
	t.Stop()
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
