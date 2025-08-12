package application

import (
	"expresApi/src/comments/domain"
)

// CommentUseCase contiene la lógica de negocio para comentarios
type CommentUseCase struct {
	repository domain.CommentRepository
}

// NewCommentUseCase crea una nueva instancia del caso de uso
func NewCommentUseCase(repository domain.CommentRepository) *CommentUseCase {
	return &CommentUseCase{
		repository: repository,
	}
}

// CreateComment crea un nuevo comentario
func (uc *CommentUseCase) CreateComment(req domain.CreateCommentRequest) (*domain.Comment, error) {
	// Validaciones de negocio aquí si es necesario

	// Crear comentario
	return uc.repository.Create(req)
}

// GetCommentsByProduct obtiene comentarios por producto
func (uc *CommentUseCase) GetCommentsByProduct(productID int) ([]domain.Comment, error) {
	return uc.repository.GetByProductID(productID)
}

// GetCommentByID obtiene un comentario por ID
func (uc *CommentUseCase) GetCommentByID(id int) (*domain.Comment, error) {
	return uc.repository.GetByID(id)
}

// DeleteComment elimina un comentario
func (uc *CommentUseCase) DeleteComment(id int) error {
	return uc.repository.Delete(id)
}

// GetAllComments obtiene todos los comentarios
func (uc *CommentUseCase) GetAllComments() ([]domain.Comment, error) {
	return uc.repository.GetAll()
}

// GetCommentStats obtiene estadísticas de comentarios
func (uc *CommentUseCase) GetCommentStats(productID int) (*domain.CommentStats, error) {
	return uc.repository.GetStats(productID)
}
