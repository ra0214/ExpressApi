package domain

import "time"

// Comment representa un comentario de producto
type Comment struct {
	ID        int       `json:"id" db:"id"`
	Comment   string    `json:"comment" db:"comment"`
	UserName  string    `json:"user_name" db:"user_name"`
	Rating    int       `json:"rating" db:"rating"`
	ProductID int       `json:"product_id" db:"product_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// CreateCommentRequest estructura para crear un comentario
type CreateCommentRequest struct {
	Comment   string `json:"comment" validate:"required,min=5,max=500"`
	UserName  string `json:"user_name" validate:"required,min=2,max=100"`
	Rating    int    `json:"rating" validate:"required,min=1,max=5"`
	ProductID int    `json:"product_id" validate:"required,min=1"`
}

// CommentRepository interfaz para el repositorio de comentarios
type CommentRepository interface {
	Create(comment CreateCommentRequest) (*Comment, error)
	GetByProductID(productID int) ([]Comment, error)
	GetByID(id int) (*Comment, error)
	Delete(id int) error
	DeleteByProductID(productID int) error // Nuevo método
	GetAll() ([]Comment, error)
	GetStats(productID int) (*CommentStats, error)
}

// CommentStats estadísticas de comentarios por producto
type CommentStats struct {
	ProductID     int         `json:"product_id"`
	TotalComments int         `json:"total_comments"`
	AverageRating float64     `json:"average_rating"`
	RatingCounts  map[int]int `json:"rating_counts"`
}
