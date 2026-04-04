package http

import (
	"e-commerce/internal/product/dto"
	"e-commerce/internal/product/service"
	"e-commerce/pkg/utils"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	service service.CategoryService
}

func NewCategoryHandler(categoryService service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		service: categoryService,
	}
}

func (h *CategoryHandler) Create(ctx *gin.Context) {
	var req dto.CreateCategoryReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	category, err := h.service.Create(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	var res dto.CreateCategoryRes
	utils.Copy(&res.Category, &category)
	ctx.JSON(http.StatusOK, utils.SuccessResponse("创建成功", res))
}

func (h *CategoryHandler) ListCategories(ctx *gin.Context) {
	categories, err := h.service.ListCategories()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	var res dto.ListCategoriesRes
	utils.Copy(&res.Categories, categories)

	ctx.JSON(http.StatusOK, utils.SuccessResponse("获取成功", res))
}

func (h *CategoryHandler) GetCategoryByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		err := errors.New("missing category ID")
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	category, err := h.service.GetCategoryByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	var res dto.Category
	utils.Copy(&res, category)
	ctx.JSON(http.StatusOK, utils.SuccessResponse("获取分类成功", res))
}
