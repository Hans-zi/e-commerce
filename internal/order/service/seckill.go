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

	"github.com/bytedance/gopkg/util/logger"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type SeckillService interface {
	Seckill(userID string, req *dto.SeckillReq) (*model.Order, error)
	SeckillConsume(data []byte) error
	SeckillSchedule() error
}

type seckillService struct {
	orderRepo   repository.OrderRepository
	productRepo repository.ProductRepository
	cache       *redis.Client
	k           *mq.Kafka
}

func NewSeckillService(
	orderRepo repository.OrderRepository,
	productRepo repository.ProductRepository,
	cache *redis.Client,
	k *mq.Kafka) SeckillService {

	s := &seckillService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
		cache:       cache,
		k:           k,
	}
	go s.k.StartConsumeMessages(consts.TOPIC_SECKILL,
		consts.GROUP_ID_SECKILL,
		s.SeckillConsume)
	return s
}

func (s *seckillService) Seckill(userID string, req *dto.SeckillReq) (*model.Order, error) {
	var line model.OrderLine
	utils.Copy(&line, req)

	exists, err := s.cache.Exists(context.Background(), utils.GenerateSeckillOrderCode(userID, line.ProductID)).Result()
	if err != nil || exists == 1 {
		return nil, errors.New("请勿重复参与秒杀")
	}
	var order model.Order
	order.UserID = userID
	order.ID = uuid.New().String()
	var deductStockScript = redis.NewScript(`
		local stockKey = KEYS[1]
		local num = tonumber(ARGV[1])
		local stock = tonumber(redis.call("get", stockKey) or "0")
		
		if stock < num then
			return 0
		end
		redis.call("decrby", stockKey, num)
		return 1
`)

	key := utils.GenerateSeckillProductCode(line.ProductID)
	res, err := deductStockScript.Run(
		context.Background(),
		s.cache,
		[]string{key},
		line.Quantity).Int()
	if err != nil {
		logger.Errorf("Seckill.DecrStock fail: %s", err)
		return nil, err
	}

	if res == 1 {
		kMsg := dto.SeckillOrderAsyncMsg{
			OrderID:   order.ID,
			UserID:    userID,
			ProductID: line.ProductID,
			Quantity:  line.Quantity,
		}
		err = s.k.ProduceMessage(consts.TOPIC_SECKILL, kMsg)
		if err != nil {
			logger.Errorf("Seckill.ProduceMessage fail: %s", err)
			return nil, err
		}
	}

	return &order, nil
}

func (s *seckillService) SeckillConsume(data []byte) error {
	var msg dto.SeckillOrderAsyncMsg
	err := json.Unmarshal(data, &msg)

	if err != nil {
		logger.Errorf("Seckill.Unmarshal fail: %s", err)
		return err
	}

	order, _ := s.orderRepo.GetByID(msg.OrderID)
	if order != nil {
		return nil // 已处理，直接跳过
	}

	order = &model.Order{
		ID:     msg.OrderID,
		UserID: msg.UserID,
	}
	//扣减库存
	err = s.productRepo.DecreaseStock(msg.ProductID, msg.Quantity)
	if err != nil {
		logger.Errorf("Seckill.DecreaseStock fail: %s", err)
		return err
	}

	product, err := s.productRepo.GetProductById(msg.ProductID)

	if err != nil {
		logger.Errorf("Seckill.GetProduct fail: %s", err)
		return err
	}
	var line model.OrderLine
	utils.Copy(&line, &msg)
	order.Total = product.Price * float64(line.Quantity)

	if err = s.orderRepo.CreateOrder(order, []*model.OrderLine{&line}); err != nil {
		logger.Errorf("Seckill.CreateOrder fail: %s", err)
		return err
	}
	return nil
}

func (s *seckillService) SeckillSchedule() error {
	products, err := s.productRepo.ListSeckillProduct()
	if err != nil {
		logger.Errorf("Schedule.ListProducts fail: %s", err)
		return err
	}

	for _, p := range products {
		err = s.cache.Set(context.Background(),
			utils.GenerateSeckillProductCode(p.ID),
			p.Stock,
			consts.ProductExpiredTime).Err()
		if err != nil {
			return err
		}
	}
	return nil
}
