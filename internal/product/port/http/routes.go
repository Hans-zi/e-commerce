package http

import (
	"e-commerce/internal/product/repository"
	"e-commerce/internal/product/service"
	"e-commerce/pkg/middleware"
	"e-commerce/pkg/token"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func Routes(r *gin.RouterGroup, db *gorm.DB, maker token.Maker, cache *redis.Client) {
	authMiddleware := middleware.AuthMiddleware(maker)
	adminOnly := middleware.AdminOnly()

	categoryRepo := repository.NewCategoryRepository(db)
	categorySvc := service.NewCategoryService(categoryRepo)
	categoryHandler := NewCategoryHandler(categorySvc)
	categories := r.Group("/categories")

	categories.GET("/", categoryHandler.ListCategories)
	categories.GET("/:id", categoryHandler.GetCategoryByID)
	categories.Use(authMiddleware, adminOnly)
	{
		categories.POST("", categoryHandler.Create)
	}

	productRepo := repository.NewProductRepository(db)
	productSvc := service.NewProductService(productRepo)
	productHandler := NewProductHandler(productSvc, cache)

	products := r.Group("/products")

	products.GET("", productHandler.ListProducts)
	products.GET("/:id", productHandler.GetProductByID)
	products.Use(authMiddleware, adminOnly)
	{
		products.POST("", productHandler.CreateProduct)
	}
}
