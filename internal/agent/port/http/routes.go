package http

import (
	"e-commerce/internal/agent/repository"
	"e-commerce/internal/agent/service"
	"e-commerce/pkg/middleware"
	"e-commerce/pkg/token"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Routes(r *gin.RouterGroup,
	db *gorm.DB,
	maker token.Maker) {
	authMiddleware := middleware.AuthMiddleware(maker)
	productRepo := repository.NewProductRepository(db)
	agentSvc := service.NewAgentService(productRepo)
	agentHandler := NewAgentHandler(agentSvc)

	agents := r.Group("/agents")
	agents.Use(authMiddleware)
	{
		agents.POST("/chat", agentHandler.Chat)
	}

}
