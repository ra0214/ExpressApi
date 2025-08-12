package domain

import "time"

// Category representa una categoría de productos
type Category struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	ImageURL    string    `json:"image_url" db:"image_url"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// CreateCategoryRequest representa la estructura para crear una categoría
type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	IsActive    bool   `json:"is_active"`
}

// UpdateCategoryRequest representa la estructura para actualizar una categoría
type UpdateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	IsActive    *bool  `json:"is_active"` // Usar puntero para distinguir entre false y nil
}

// ICategoryRepository define las operaciones de persistencia para categorías
type ICategoryRepository interface {
	CreateCategory(category CreateCategoryRequest) (*Category, error)
	GetAllCategories() ([]Category, error)
	GetCategoryByID(id int) (*Category, error)
	GetCategoryByName(name string) (*Category, error)
	UpdateCategory(id int, category UpdateCategoryRequest) (*Category, error)
	DeleteCategory(id int) error
	GetActiveCategories() ([]Category, error)
}
