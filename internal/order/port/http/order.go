package http

import (
	"e-commerce/internal/consts"
	"e-commerce/internal/order/dto"
	"e-commerce/internal/order/service"
	"e-commerce/pkg/token"
	"e-commerce/pkg/utils"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	service service.OrderService
}

func NewOrderHandler(service service.OrderService) *OrderHandler {
	return &OrderHandler{
		service: service,
	}
}

func (h *OrderHandler) PlaceOrder(ctx *gin.Context) {
	var req dto.PlaceOrderReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}
	payload, _ := ctx.Get(consts.AuthorizationPayloadKey)
	userID := payload.(*token.Payload).UserID

	order, err := h.service.PlaceOrder(userID, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	var res dto.Order
	utils.Copy(&res, order)
	ctx.JSON(http.StatusOK, utils.SuccessResponse("创建订单成功", res))
}

func (h *OrderHandler) GetOrderByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		err := errors.New("missing order id")
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}
	payload, _ := ctx.Get(consts.AuthorizationPayloadKey)
	userID := payload.(*token.Payload).UserID

	order, err := h.service.GetOrderByID(id, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	var res dto.Order
	utils.Copy(&res, order)
	ctx.JSON(http.StatusOK, utils.SuccessResponse("获取订单成功", res))

}

func (h *OrderHandler) CancelOrder(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		err := errors.New("missing order id")
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}
	payload, _ := ctx.Get(consts.AuthorizationPayloadKey)
	userID := payload.(*token.Payload).UserID

	err := h.service.CancelOrder(id, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, utils.SuccessResponse("取消订单成功", nil))

}
func (h *OrderHandler) DeleteOrder(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		err := errors.New("missing order id")
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}
	payload, _ := ctx.Get(consts.AuthorizationPayloadKey)
	userID := payload.(*token.Payload).UserID

	err := h.service.DeleteOrder(id, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, utils.SuccessResponse("删除订单成功", nil))

}
