package dto

type Address struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Province  string `json:"province"`
	City      string `json:"city"`
	District  string `json:"district"`
	Detail    string `json:"detail"`
	IsDefault bool   `json:"is_default"`
}

type CreateAddressReq struct {
	Name     string `json:"name" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Province string `json:"province" binding:"required"`
	City     string `json:"city" binding:"required"`
	District string `json:"district" binding:"required"`
}

type UpdateAddressReq struct {
	Name     string `json:"name,omitempty"`
	Phone    string `json:"phone,omitempty"`
	Province string `json:"province,omitempty"`
	City     string `json:"city,omitempty"`
	District string `json:"district,omitempty"`
	Detail   string `json:"detail,omitempty"`
}
