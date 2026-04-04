package dto

import "e-commerce/internal/product/model"

type Category struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}
type CreateCategoryReq struct {
	Name        string `json:"name" binding:"required"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

type CreateCategoryRes struct {
	Category Category
}

type ListCategoriesRes struct {
	Categories []model.Category
}

type GetCategoryReq struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type GetCategoryRes struct {
	Category Category `json:"category"`
}
