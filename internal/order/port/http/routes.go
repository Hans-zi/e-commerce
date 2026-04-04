package http

import (
	"e-commerce/internal/cart/repository"
	orderRepository "e-commerce/internal/order/repository"
	"e-commerce/internal/order/service"
	"e-commerce/pkg/middleware"
	"e-commerce/pkg/mq"
	"e-commerce/pkg/token"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/smartwalle/alipay/v3"
	"gorm.io/gorm"
)

func Routes(r *gin.RouterGroup,
	db *gorm.DB,
	maker token.Maker,
	alipayClient *alipay.Client,
	cache *redis.Client,
	k *mq.Kafka) {
	authMiddleware := middleware.AuthMiddleware(maker)

	productRepo := orderRepository.NewProductRepository(db)
	orderRepo := orderRepository.NewOrderRepo(db)
	cartRepo := repository.NewCartRepository(db)
	orderSvc := service.NewOrderService(orderRepo, productRepo, cartRepo, k)
	orderhandler := NewOrderHandler(orderSvc)

	orders := r.Group("/orders")
	orders.Use(authMiddleware)
	{
		orders.POST("", orderhandler.PlaceOrder)
		orders.GET("/:id", orderhandler.GetOrderByID)
		orders.PUT("/:id/cancel", orderhandler.CancelOrder)
		orders.DELETE("/:id", orderhandler.DeleteOrder)
	}

	paymentRepo := orderRepository.NewPaymentRepository(db)
	alipaySvc := service.NewAlipayService(paymentRepo, orderRepo, productRepo, alipayClient, k)
	alipayHandler := NewPaymentHandler(alipaySvc)

	alipays := r.Group("/payment/alipay")

	alipays.Use(authMiddleware)
	{
		alipays.POST("/:id", alipayHandler.CreatePayment)
	}

	callback := r.Group("/payment")
	{
		callback.POST("/alipay/notify", alipayHandler.AlipayCallBack)
	}

	seckillSvc := service.NewSeckillService(orderRepo, productRepo, cache, k)
	seckillHandler := NewSeckillHandler(seckillSvc)
	seckill := r.Group("/seckill")

	seckill.Use(authMiddleware)
	{
		seckill.POST("", seckillHandler.Seckill)
		seckill.GET("/schedule", seckillHandler.SeckillSchedule)
	}

}
