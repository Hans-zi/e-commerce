package dto

import (
	"e-commerce/pkg/utils"
	"time"
)

type Product struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	Status      string    `json:"status"`
	SalesCount  int       `json:"sales_count"`
	ViewCount   int       `json:"view_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Images      []string  `json:"images,omitempty"`

	CategoryID string `json:"category_id,omitempty"`
}

type CreateProductReq struct {
	Name        string   `json:"name" binding:"required"`
	Slug        string   `json:"slug"`
	Description string   `json:"description"`
	Price       float64  `json:"price" binding:"required"`
	Stock       int      `json:"stock"`
	Images      []string `json:"images,omitempty"`
	CategoryID  string   `json:"category_id,omitempty"`
}

type CreateProductRes struct {
	Product Product `json:"product"`
}

type ListProductsReq struct {
	Page       int64   `json:"page" `
	PageSize   int64   `json:"page_size"`
	Name       string  `json:"name,omitempty"`
	Slug       string  `json:"slug,omitempty" `
	MaxPrice   float64 `json:"max_price"`
	MinPrice   float64 `json:"min_price"`
	CategoryID string  `json:"category_id,omitempty"`
	OrderBy    string  `json:"order_by"`
	OrderDesc  bool    `json:"-"`
}

type ListProductsRes struct {
	Products []Product        `json:"products"`
	Paging   utils.Pagination `json:"paging"`
}

type GetProductReq struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type GetProductRes struct {
	Product Product `json:"product"`
}
