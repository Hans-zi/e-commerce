package main

import (
	"e-commerce/internal/server/http"
	"e-commerce/pkg/dbs"
	"e-commerce/pkg/mq"
	"e-commerce/pkg/payment"
	"e-commerce/pkg/token"
	"e-commerce/pkg/utils"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay/v3"
	"gorm.io/gorm"
)

var db *gorm.DB
var maker token.Maker
var alipayClient *alipay.Client

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("can not load config")
		return
	}
	dsn := dbs.BuildDSN(config)
	db, err = dbs.InitDB(dsn)
	if err != nil {
		log.Fatal("can not connect to database")
		return
	}
	err = dbs.AutoMigrate(db)
	if err != nil {
		log.Fatal("can not migrate database")
		return
	}
	cache := dbs.NewRedisClient(config)
	alipayClient, err = payment.NewPayClient(config)
	k := mq.NewKafka()
	maker, err = token.NewJWTMaker(config.JWT.SymmetricKey, config.JWT.Duration)
	server := http.NewServer(db, config, maker, cache, alipayClient, k)
	gin.SetMode(config.Server.Mode)

	server.MapRoutes()

	server.Run()
}
