package application

import (
	"expresApi/src/products/domain"
)

type CreateProduct struct {
	db domain.IProduct
}

func NewCreateProduct(db domain.IProduct) *CreateProduct {
	return &CreateProduct{db: db}
}

func (ct *CreateProduct) Execute(name string, description string, price float64, category string, imageURL string) error {
	// Crear el objeto producto
	product := domain.NewProduct(name, description, price, category, imageURL)

	// Guardar el producto una sola vez
	err := ct.db.SaveProduct(product.Name, product.Description, product.Price, product.Category, product.ImageURL)
	if err != nil {
		return err
	}

	return nil
}
