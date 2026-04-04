package http

import (
	"e-commerce/pkg/mq"
	"e-commerce/pkg/token"
	"e-commerce/pkg/utils"
	"fmt"

	cartHttp "e-commerce/internal/cart/port/http"
	orderHttp "e-commerce/internal/order/port/http"
	productHttp "e-commerce/internal/product/port/http"
	userHttp "e-commerce/internal/user/port/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/smartwalle/alipay/v3"
	"gorm.io/gorm"
)

type Server struct {
	r            *gin.Engine
	db           *gorm.DB
	cfg          utils.Config
	maker        token.Maker
	cache        *redis.Client
	alipayClient *alipay.Client
	k            *mq.Kafka
}

func NewServer(
	db *gorm.DB,
	cfg utils.Config,
	maker token.Maker,
	cache *redis.Client,
	alipayClient *alipay.Client,
	k *mq.Kafka) *Server {
	return &Server{
		r:            gin.Default(),
		db:           db,
		cfg:          cfg,
		maker:        maker,
		cache:        cache,
		alipayClient: alipayClient,
		k:            k,
	}
}

func (s *Server) MapRoutes() error {
	api := s.r.Group("/api")
	userHttp.Routes(api, s.db, s.maker)
	productHttp.Routes(api, s.db, s.maker, s.cache)
	cartHttp.Routes(api, s.db, s.maker)
	orderHttp.Routes(api, s.db, s.maker, s.alipayClient, s.cache, s.k)
	return nil
}

func (s *Server) Run() {
	s.r.Run(fmt.Sprintf("%s:%d", s.cfg.Server.Host, s.cfg.Server.Port))
}
