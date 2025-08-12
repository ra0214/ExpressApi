package application

import (
	"expresApi/src/products/domain"
)

type ViewProduct struct {
	db domain.IProduct
}

func NewViewProduct(db domain.IProduct) *ViewProduct {
	return &ViewProduct{db: db}
}

func (vt ViewProduct) Execute() ([]domain.Product, error) {
	return vt.db.GetAll()
}