package dbs

import (
	cartModel "e-commerce/internal/cart/model"
	orderModel "e-commerce/internal/order/model"
	productModel "e-commerce/internal/product/model"
	userModel "e-commerce/internal/user/model"
	"e-commerce/pkg/utils"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Info),
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取数据库实例失败: %w", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	return db, nil
}

func BuildDSN(config utils.Config) string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		config.DB.User, config.DB.Password, config.DB.Host, config.DB.Port, config.DB.DbName)
	return dsn
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&userModel.User{},
		&productModel.Category{},
		&productModel.Product{},
		&cartModel.Cart{},
		&cartModel.CartLine{},
		&orderModel.Order{},
		&orderModel.OrderLine{},
		&orderModel.Payment{},
		&userModel.Address{})
}
