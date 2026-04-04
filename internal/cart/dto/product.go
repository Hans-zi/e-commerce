package dto

type Product struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Slug        string   `json:"slug"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	Images      []string `json:"images"`
}
