package http

import (
	"e-commerce/internal/agent/dto"
	"e-commerce/internal/agent/service"
	"e-commerce/internal/consts"
	"e-commerce/pkg/token"
	"e-commerce/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AgentHandler struct {
	agentSvc *service.AgentService
}

func NewAgentHandler(agentSvc *service.AgentService) *AgentHandler {
	return &AgentHandler{
		agentSvc: agentSvc,
	}
}

func (h *AgentHandler) Chat(ctx *gin.Context) {
	var req dto.ChatMsg
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}
	payload, _ := ctx.Get(consts.AuthorizationPayloadKey)
	userID := payload.(*token.Payload).UserID

	res, err := h.agentSvc.Chat(ctx, userID, req.Msg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(res.Msg, res.Data))
}
