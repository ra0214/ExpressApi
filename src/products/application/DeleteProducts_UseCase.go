package application

import (
	"expresApi/src/products/domain"
	"fmt"
	"strconv"
)

type DeleteProduct struct {
	repo        domain.IProduct
	commentRepo CommentRepository // Añadimos dependencia a comentarios
}

// CommentRepository interfaz para operaciones de comentarios
type CommentRepository interface {
	DeleteByProductID(productID int) error
}

func NewDeleteProduct(repo domain.IProduct, commentRepo CommentRepository) *DeleteProduct {
	return &DeleteProduct{
		repo:        repo,
		commentRepo: commentRepo,
	}
}

func (d *DeleteProduct) Execute(id string) error {
	// Convertir string a int para el ID del producto
	productID, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("ID de producto inválido: %v", err)
	}

	// Primero eliminar todos los comentarios del producto
	if d.commentRepo != nil {
		err := d.commentRepo.DeleteByProductID(productID)
		if err != nil {
			// Log el error pero no fallas, continuamos con la eliminación del producto
			fmt.Printf("Advertencia: No se pudieron eliminar algunos comentarios del producto %d: %v\n", productID, err)
		}
	}

	// Luego eliminar el producto
	return d.repo.Delete(id)
}
