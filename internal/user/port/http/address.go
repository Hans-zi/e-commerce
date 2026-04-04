package http

import (
	"e-commerce/internal/consts"
	"e-commerce/internal/user/dto"
	"e-commerce/internal/user/service"
	"e-commerce/pkg/token"
	"e-commerce/pkg/utils"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AddressHandler struct {
	service service.AddressService
}

func NewAddressHandler(service service.AddressService) *AddressHandler {
	return &AddressHandler{
		service: service,
	}
}

func (h *AddressHandler) CreateAddress(ctx *gin.Context) {
	var req dto.CreateAddressReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}
	payload, _ := ctx.Get(consts.AuthorizationPayloadKey)

	userID := payload.(*token.Payload).UserID
	address, err := h.service.CreateAddress(userID, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	var res dto.Address
	utils.Copy(&res, address)
	ctx.JSON(http.StatusOK, utils.SuccessResponse("创建地址成功", res))
}

func (h *AddressHandler) ListAddresses(ctx *gin.Context) {
	payload, _ := ctx.Get(consts.AuthorizationPayloadKey)

	userID := payload.(*token.Payload).UserID

	addresses, err := h.service.ListAddresses(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	var res []*dto.Address
	utils.Copy(&res, &addresses)
	ctx.JSON(http.StatusOK, utils.SuccessResponse("获取地址列表成功", &res))
}

func (h *AddressHandler) GetAddress(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		err := errors.New("missing address ID")
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}
	payload, _ := ctx.Get(consts.AuthorizationPayloadKey)

	userID := payload.(*token.Payload).UserID

	address, err := h.service.GetAddressByID(id, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	var res dto.Address
	utils.Copy(&res, address)
	ctx.JSON(http.StatusOK, utils.SuccessResponse("获取地址成功", res))
}

func (h *AddressHandler) UpdateAddress(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		err := errors.New("missing address ID")
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}
	var req *dto.UpdateAddressReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}
	payload, _ := ctx.Get(consts.AuthorizationPayloadKey)

	userID := payload.(*token.Payload).UserID

	address, err := h.service.UpdateAddressByID(id, userID, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	var res dto.Address
	utils.Copy(&res, address)
	ctx.JSON(http.StatusOK, utils.SuccessResponse("更新地址成功", res))
}

func (h *AddressHandler) SetAddressDefault(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		err := errors.New("missing address ID")
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}
	payload, _ := ctx.Get(consts.AuthorizationPayloadKey)

	userID := payload.(*token.Payload).UserID

	err := h.service.SetAddressDefault(id, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, utils.SuccessResponse("设置默认成功", nil))
}
