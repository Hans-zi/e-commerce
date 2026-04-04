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

type PaymentHandler struct {
	alipaySvc service.AlipayService
}

func NewPaymentHandler(alipaySvc service.AlipayService) *PaymentHandler {
	return &PaymentHandler{
		alipaySvc: alipaySvc,
	}
}

func (h *PaymentHandler) CreatePayment(ctx *gin.Context) {
	orderID := ctx.Param("id")
	if orderID == "" {
		err := errors.New("missing order id")
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	payload, _ := ctx.Get(consts.AuthorizationPayloadKey)
	userID := payload.(*token.Payload).UserID
	payment, err := h.alipaySvc.CreatePayment(orderID, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	var res dto.Payment
	utils.Copy(&res, payment)
	ctx.JSON(http.StatusOK, utils.SuccessResponse("", res))
}

func (h *PaymentHandler) AlipayCallBack(ctx *gin.Context) {
	if err := ctx.Request.ParseForm(); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	res, err := h.alipaySvc.CallBack(ctx.Request.Form)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, utils.SuccessResponse("", res))
}
