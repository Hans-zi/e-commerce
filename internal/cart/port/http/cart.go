package http

import (
	"e-commerce/internal/cart/dto"
	"e-commerce/internal/cart/service"
	"e-commerce/internal/consts"
	"e-commerce/pkg/token"
	"e-commerce/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	service service.CartService
}

func NewCartHandler(service service.CartService) *CartHandler {
	return &CartHandler{
		service: service,
	}
}

func (h *CartHandler) AddProduct(ctx *gin.Context) {
	var req dto.AddProductReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}
	payload, _ := ctx.Get(consts.AuthorizationPayloadKey)
	userID := payload.(*token.Payload).UserID

	cart, err := h.service.AddProduct(userID, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	var res dto.Cart
	utils.Copy(&res, cart)
	ctx.JSON(http.StatusOK, utils.SuccessResponse("添加成功", res))
}

func (h *CartHandler) GetMyCart(ctx *gin.Context) {
	payload, _ := ctx.Get(consts.AuthorizationPayloadKey)
	userID := payload.(*token.Payload).UserID

	cart, err := h.service.GetCartByUserID(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	var res dto.Cart
	utils.Copy(&res, cart)
	ctx.JSON(http.StatusOK, utils.SuccessResponse("获取成功", res))
}

func (h *CartHandler) RemoveProduct(ctx *gin.Context) {
	var req dto.RemoveProductReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}
	payload, _ := ctx.Get(consts.AuthorizationPayloadKey)
	userID := payload.(*token.Payload).UserID

	cart, err := h.service.RemoveProduct(userID, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	var res dto.Cart
	utils.Copy(&res, cart)
	ctx.JSON(http.StatusOK, utils.SuccessResponse("删除成功", res))
}
