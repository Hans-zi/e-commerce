package http

import (
	"e-commerce/internal/consts"
	"e-commerce/internal/order/dto"
	"e-commerce/internal/order/service"
	"e-commerce/pkg/token"
	"e-commerce/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SeckillHandler struct {
	service service.SeckillService
}

func NewSeckillHandler(service service.SeckillService) *SeckillHandler {
	return &SeckillHandler{
		service: service,
	}
}

func (h *SeckillHandler) Seckill(ctx *gin.Context) {
	var req dto.SeckillReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}
	payload, _ := ctx.Get(consts.AuthorizationPayloadKey)
	userID := payload.(*token.Payload).UserID

	order, err := h.service.Seckill(userID, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	var res dto.Order
	utils.Copy(&res, order)
	ctx.JSON(http.StatusOK, utils.SuccessResponse("秒杀成功", res))
}

func (h *SeckillHandler) SeckillSchedule(ctx *gin.Context) {
	err := h.service.SeckillSchedule()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("预热库存成功", nil))
}
