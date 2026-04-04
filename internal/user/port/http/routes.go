package http

import (
	"e-commerce/internal/user/repository"
	"e-commerce/internal/user/service"
	"e-commerce/pkg/middleware"
	"e-commerce/pkg/token"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Routes(r *gin.RouterGroup, db *gorm.DB, maker token.Maker) {
	authMiddleware := middleware.AuthMiddleware(maker)

	userRepo := repository.NewUserRepo(db)
	userSvc := service.NewUserService(userRepo, maker)
	userHandler := NewUserHandler(userSvc)

	users := r.Group("/users")
	users.POST("", userHandler.Register)
	users.POST("/login", userHandler.Login)
	users.Use(authMiddleware)
	{
		users.GET("/me", userHandler.GetMe)
		users.POST("/change-password", userHandler.ChangePassword)
	}

	addressRepo := repository.NewAddressRepo(db)
	addressSvc := service.NewAddressService(addressRepo)
	addressHandler := NewAddressHandler(addressSvc)
	addresses := r.Group("/addresses")
	addresses.Use(authMiddleware)
	{
		addresses.POST("", addressHandler.CreateAddress)
		addresses.GET("", addressHandler.ListAddresses)
		addresses.GET("/:id", addressHandler.GetAddress)
		addresses.PUT("/:id", addressHandler.UpdateAddress)
		addresses.PUT("/:id/default", addressHandler.SetAddressDefault)
	}
}
