package product

type CreateProductRequest struct {
	Name        string   `json:"name" validate:"required"`
	Description string   `json:"description" validate:"required"`
	Price       float32  `json:"price" validate:"required,number,gt=0"`
	Images      []string `json:"images"`
}

type UpdateProductRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float32  `json:"price" validate:"omitempty,number,gt=0"`
	Images      []string `json:"images"`
}
