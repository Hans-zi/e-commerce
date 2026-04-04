package http

import (
	"e-commerce/internal/cart/repository"
	"e-commerce/internal/cart/service"
	"e-commerce/pkg/middleware"
	"e-commerce/pkg/token"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Routes(r *gin.RouterGroup, db *gorm.DB, maker token.Maker) {
	authMiddleware := middleware.AuthMiddleware(maker)

	cartRepo := repository.NewCartRepository(db)
	cartSvc := service.NewCartService(cartRepo)
	cartHandler := NewCartHandler(cartSvc)

	cart := r.Group("/cart")
	cart.Use(authMiddleware)
	{
		cart.GET("", cartHandler.GetMyCart)
		cart.POST("/add", cartHandler.AddProduct)
		cart.POST("/remove", cartHandler.RemoveProduct)
	}
}
