package service

import (
	"e-commerce/internal/cart/repository"
	"e-commerce/internal/consts"
	"e-commerce/internal/order/dto"
	"e-commerce/internal/order/model"
	orderRepo "e-commerce/internal/order/repository"
	"e-commerce/pkg/mq"
	"e-commerce/pkg/utils"
	"encoding/json"
	"errors"
	"time"

	"github.com/bytedance/gopkg/util/logger"
)

type OrderService interface {
	PlaceOrder(userID string, req *dto.PlaceOrderReq) (*model.Order, error)
	GetOrderByID(id, userID string) (*model.Order, error)
	CancelOrder(id, userID string) error
	DeleteOrder(id, userID string) error
}

type orderService struct {
	orderRepo   orderRepo.OrderRepository
	productRepo orderRepo.ProductRepository
	cartRepo    repository.CartRepository
	k           *mq.Kafka
}

func NewOrderService(
	orderRepo orderRepo.OrderRepository,
	productRepo orderRepo.ProductRepository,
	cartRepo repository.CartRepository,
	k *mq.Kafka) OrderService {

	s := &orderService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
		cartRepo:    cartRepo,
		k:           k,
	}
	go s.k.StartConsumeMessages(consts.TOPIC_CANCEL_ORDER,
		consts.GROUP_ID_ORDER,
		s.AutoCancel)
	return s
}

func (s *orderService) PlaceOrder(userID string, req *dto.PlaceOrderReq) (*model.Order, error) {

	var lines []*model.OrderLine
	utils.Copy(&lines, req.Lines)

	var total float64 = 0

	//校验扣减库存
	err := s.productRepo.DecreaseStockBatch(lines)
	if err != nil {
		logger.Errorf("PlaceOrder.DecreaseStock fail: %s", err)
		return nil, err
	}

	//计算总价
	var order model.Order
	for _, line := range lines {
		product, err := s.productRepo.GetProductById(line.ProductID)
		if err != nil {
			return nil, err
		}
		total += product.Price * float64(line.Quantity)
	}
	order.UserID = userID
	order.Total = total

	//创建订单
	err = s.orderRepo.CreateOrder(&order, lines)
	if err != nil {
		logger.Errorf("PlaceOrder.CreateOrder fail: %s", err)
		return nil, err
	}

	//清理购物车
	go func() {
		cart, err := s.cartRepo.GetCartByUserID(userID)
		if err != nil {
			logger.Errorf("PlaceOrder.GetCart fail: %s", err)
			return
		}
		for _, line := range lines {
			err = s.cartRepo.RemoveCartLine(cart.ID, line.ProductID)
			if err != nil {
				logger.Errorf("PlaceOrder.ClearCartLine fail: %s", err)
			}
		}
	}()

	go func() {
		var kmsg dto.AutoCancelMsg
		kmsg.OrderID = order.ID
		kmsg.UserID = order.UserID
		s.k.ProduceMessage(consts.TOPIC_CANCEL_ORDER, kmsg)

	}()

	return s.GetOrderByID(order.ID, userID)
}

func (s *orderService) GetOrderByID(id, userID string) (*model.Order, error) {
	order, err := s.orderRepo.GetUserOrderByID(id, userID)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (s *orderService) CancelOrder(id, userID string) error {
	//获取订单
	order, err := s.orderRepo.GetUserOrderByID(id, userID)
	if err != nil {
		return err
	}

	//确认是否能取消
	if !s.canChangeOrderStatus(order.Status, consts.ORDER_STATUS_CANCELED) {
		err = errors.New("无法取消订单")
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

func (s *orderService) DeleteOrder(id, userID string) error {
	return s.orderRepo.Delete(id, userID)
}

func (r *orderService) canChangeOrderStatus(cur, nxt string) bool {
	allowed := map[string][]string{
		consts.ORDER_STATUS_PENDING:   {consts.ORDER_STATUS_PAID, consts.ORDER_STATUS_COMPLETED, consts.ORDER_STATUS_CANCELED},
		consts.ORDER_STATUS_PAID:      {consts.ORDER_STATUS_SHIPPED, consts.ORDER_STATUS_COMPLETED, consts.ORDER_STATUS_CANCELED},
		consts.ORDER_STATUS_SHIPPED:   {consts.ORDER_STATUS_COMPLETED, consts.ORDER_STATUS_CANCELED},
		consts.ORDER_STATUS_COMPLETED: {},
		consts.ORDER_STATUS_CANCELED:  {},
	}

	for _, status := range allowed[cur] {
		if status == nxt {
			return true
		}
	}

	return false
}

func (s *orderService) AutoCancel(data []byte) error {
	t := time.NewTicker(time.Second)
	<-t.C
	t.Stop()

	var msg dto.AutoCancelMsg
	err := json.Unmarshal(data, &msg)
	if err != nil {
		logger.Errorf("AutoCancel.Unmarshal fail: %s", err)
		return err
	}

	//获取订单
	order, err := s.orderRepo.GetUserOrderByID(msg.OrderID, msg.UserID)
	if err != nil {
		return err
	}

	//确认是否能取消
	if order.Status != consts.ORDER_STATUS_PENDING {
		err = errors.New("订单已完成")
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
	logger.Infof("订单支付超时，自动取消")
	return nil
}
