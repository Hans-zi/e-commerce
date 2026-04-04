package http

import (
	"database/sql"
	"e-commerce/internal/consts"
	"e-commerce/internal/user/dto"
	"e-commerce/internal/user/service"
	"e-commerce/pkg/token"
	"e-commerce/pkg/utils"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{service: userService}
}

func (h *UserHandler) Register(ctx *gin.Context) {
	var req dto.RegisterReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	user, err := h.service.Register(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	var res dto.RegisterRes
	utils.Copy(&res.User, &user)

	ctx.JSON(http.StatusOK, res)
}

func (h *UserHandler) Login(ctx *gin.Context) {
	var req dto.LoginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}
	user, jtoken, err := h.service.Login(&req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, utils.ErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	var res dto.LoginRes
	utils.Copy(&res.User, user)
	res.Token = jtoken
	ctx.JSON(http.StatusOK, utils.SuccessResponse("用户登录成功", res))
}

func (h *UserHandler) GetMe(ctx *gin.Context) {
	payload, _ := ctx.Get(consts.AuthorizationPayloadKey)

	id := payload.(*token.Payload).UserID

	user, err := h.service.GetUserById(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	var res dto.User
	utils.Copy(&res, user)
	ctx.JSON(http.StatusOK, utils.SuccessResponse("获取信息成功", res))
}

func (h *UserHandler) ChangePassword(ctx *gin.Context) {
	var req dto.ChangePasswordReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}
	payload, _ := ctx.Get(consts.AuthorizationPayloadKey)
	id := payload.(*token.Payload).UserID

	user, err := h.service.ChangePassword(id, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	var res dto.User
	utils.Copy(&res, user)
	ctx.JSON(http.StatusOK, utils.SuccessResponse("修改密码成功", res))
}
