package infrastructure

import (
	"database/sql"
	"expresApi/src/comments/application"
	"expresApi/src/comments/domain"
)

// CommentDependencies contiene todas las dependencias del módulo de comentarios
type CommentDependencies struct {
	Repository domain.CommentRepository
	UseCase    *application.CommentUseCase
	Controller *CommentController
}

// NewCommentDependencies crea e inicializa todas las dependencias
func NewCommentDependencies(db *sql.DB) *CommentDependencies {
	// Crear repositorio
	repository := NewMySQLCommentRepository(db)

	// Crear caso de uso
	useCase := application.NewCommentUseCase(repository)

	// Crear controlador
	controller := NewCommentController(useCase)

	return &CommentDependencies{
		Repository: repository,
		UseCase:    useCase,
		Controller: controller,
	}
}
