package application

import (
	"expresApi/src/products/domain"
)

type UpdateProduct struct {
	repo domain.IProduct
}

func NewUpdateProduct(repo domain.IProduct) *UpdateProduct {
	return &UpdateProduct{repo: repo}
}

func (u *UpdateProduct) Execute(id string, name string, description string, price float64, category string, imageURL string) error {
	return u.repo.Update(id, name, description, price, category, imageURL)
}
