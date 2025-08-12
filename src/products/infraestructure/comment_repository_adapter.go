package infraestructure

import (
	"database/sql"
	commentInfra "expresApi/src/comments/infrastructure"
)

// CommentRepositoryAdapter adapta el repositorio de comentarios para usar en productos
type CommentRepositoryAdapter struct {
	repo *commentInfra.MySQLCommentRepository
}

// NewCommentRepositoryAdapter crea una nueva instancia del adaptador
func NewCommentRepositoryAdapter(db *sql.DB) *CommentRepositoryAdapter {
	repo := commentInfra.NewMySQLCommentRepository(db)
	return &CommentRepositoryAdapter{
		repo: repo.(*commentInfra.MySQLCommentRepository),
	}
}

// DeleteByProductID elimina todos los comentarios de un producto
func (c *CommentRepositoryAdapter) DeleteByProductID(productID int) error {
	return c.repo.DeleteByProductID(productID)
}
