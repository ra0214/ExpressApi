package application

import (
	"errors"
	"expresApi/src/categories/domain"
	"strings"
)

// CategoryUseCase maneja la lógica de negocio para categorías
type CategoryUseCase struct {
	repository domain.ICategoryRepository
}

// NewCategoryUseCase crea una nueva instancia del caso de uso
func NewCategoryUseCase(repository domain.ICategoryRepository) *CategoryUseCase {
	return &CategoryUseCase{
		repository: repository,
	}
}

// CreateCategory crea una nueva categoría
func (uc *CategoryUseCase) CreateCategory(request domain.CreateCategoryRequest) (*domain.Category, error) {
	// Validaciones de negocio
	if strings.TrimSpace(request.Name) == "" {
		return nil, errors.New("el nombre de la categoría es requerido")
	}

	// Verificar que no exista una categoría con el mismo nombre
	existingCategory, _ := uc.repository.GetCategoryByName(request.Name)
	if existingCategory != nil {
		return nil, errors.New("ya existe una categoría con ese nombre")
	}

	// Normalizar el nombre (primera letra mayúscula)
	request.Name = strings.Title(strings.ToLower(strings.TrimSpace(request.Name)))

	return uc.repository.CreateCategory(request)
}

// GetAllCategories obtiene todas las categorías
func (uc *CategoryUseCase) GetAllCategories() ([]domain.Category, error) {
	return uc.repository.GetAllCategories()
}

// GetActiveCategories obtiene solo las categorías activas
func (uc *CategoryUseCase) GetActiveCategories() ([]domain.Category, error) {
	return uc.repository.GetActiveCategories()
}

// GetCategoryByID obtiene una categoría por ID
func (uc *CategoryUseCase) GetCategoryByID(id int) (*domain.Category, error) {
	if id <= 0 {
		return nil, errors.New("ID de categoría inválido")
	}

	category, err := uc.repository.GetCategoryByID(id)
	if err != nil {
		return nil, err
	}

	if category == nil {
		return nil, errors.New("categoría no encontrada")
	}

	return category, nil
}

// UpdateCategory actualiza una categoría existente
func (uc *CategoryUseCase) UpdateCategory(id int, request domain.UpdateCategoryRequest) (*domain.Category, error) {
	if id <= 0 {
		return nil, errors.New("ID de categoría inválido")
	}

	// Verificar que la categoría existe
	existingCategory, err := uc.repository.GetCategoryByID(id)
	if err != nil {
		return nil, err
	}
	if existingCategory == nil {
		return nil, errors.New("categoría no encontrada")
	}

	// Si se está actualizando el nombre, verificar que no exista otra categoría con ese nombre
	if request.Name != "" && request.Name != existingCategory.Name {
		duplicateCategory, _ := uc.repository.GetCategoryByName(request.Name)
		if duplicateCategory != nil && duplicateCategory.ID != id {
			return nil, errors.New("ya existe una categoría con ese nombre")
		}
		// Normalizar el nombre
		request.Name = strings.Title(strings.ToLower(strings.TrimSpace(request.Name)))
	}

	return uc.repository.UpdateCategory(id, request)
}

// DeleteCategory elimina una categoría (soft delete)
func (uc *CategoryUseCase) DeleteCategory(id int) error {
	if id <= 0 {
		return errors.New("ID de categoría inválido")
	}

	// Verificar que la categoría existe
	existingCategory, err := uc.repository.GetCategoryByID(id)
	if err != nil {
		return err
	}
	if existingCategory == nil {
		return errors.New("categoría no encontrada")
	}

	return uc.repository.DeleteCategory(id)
}
