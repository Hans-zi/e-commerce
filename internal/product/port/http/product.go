package http

import (
	"e-commerce/internal/consts"
	"e-commerce/internal/product/dto"
	"e-commerce/internal/product/service"
	"e-commerce/pkg/dbs"
	"e-commerce/pkg/utils"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type ProductHandler struct {
	service service.ProductService
	cache   *redis.Client
}

func NewProductHandler(service service.ProductService, cache *redis.Client) *ProductHandler {
	return &ProductHandler{
		service: service,
		cache:   cache,
	}
}

func (h *ProductHandler) CreateProduct(ctx *gin.Context) {
	var req dto.CreateProductReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}
	product, err := h.service.Create(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	var res dto.CreateProductRes
	utils.Copy(&res.Product, product)
	ctx.JSON(http.StatusOK, utils.SuccessResponse("创建商品成功", res))
}

func (h *ProductHandler) ListProducts(ctx *gin.Context) {
	var req dto.ListProductsReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}
	var res dto.ListProductsRes
	cacheKey := ctx.Request.URL.String()
	err := dbs.Get(h.cache, cacheKey, &res)
	if err == nil {
		ctx.JSON(http.StatusOK, utils.SuccessResponse("获取列表成功", res))
		return
	}
	products, paging, err := h.service.ListProducts(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	utils.Copy(&res.Products, products)
	utils.Copy(&res.Paging, paging)
	ctx.JSON(http.StatusOK, utils.SuccessResponse("获取列表成功", res))
	_ = dbs.SetWithExpirationTime(h.cache, cacheKey, res, consts.ProductExpiredTime)
}

func (h *ProductHandler) GetProductByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		err := errors.New("missing product ID")
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}
	var res dto.Product
	err := dbs.Get(h.cache, id, &res)
	if err == nil {
		ctx.JSON(http.StatusOK, utils.SuccessResponse("获取商品成功", res))
		return
	}
	product, err := h.service.GetProductByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	utils.Copy(&res, product)
	ctx.JSON(http.StatusOK, utils.SuccessResponse("获取商品成功", res))
	_ = dbs.SetWithExpirationTime(h.cache, id, &product, consts.ProductExpiredTime)
}
