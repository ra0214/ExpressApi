package infrastructure

import (
	"database/sql"
	"expresApi/src/categories/domain"
	"expresApi/src/config"
	"fmt"
	"log"
)

// MySQLCategoryRepository implementa ICategoryRepository para MySQL
type MySQLCategoryRepository struct {
	conn *config.Conn_MySQL
}

// NewMySQLCategoryRepository crea una nueva instancia del repositorio
func NewMySQLCategoryRepository() domain.ICategoryRepository {
	conn := config.GetDBPool()
	if conn.Err != "" {
		log.Fatalf("Error al configurar el pool de conexiones: %v", conn.Err)
	}
	return &MySQLCategoryRepository{conn: conn}
}

// CreateCategory crea una nueva categoría en la base de datos
func (r *MySQLCategoryRepository) CreateCategory(request domain.CreateCategoryRequest) (*domain.Category, error) {
	query := `
		INSERT INTO categories (name, description, image_url, is_active) 
		VALUES (?, ?, ?, ?)
	`

	result, err := r.conn.ExecutePreparedQuery(
		query,
		request.Name,
		request.Description,
		request.ImageURL,
		request.IsActive,
	)
	if err != nil {
		return nil, fmt.Errorf("error al crear la categoría: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error al obtener el ID de la categoría creada: %v", err)
	}

	// Obtener la categoría creada
	return r.GetCategoryByID(int(id))
}

// GetAllCategories obtiene todas las categorías
func (r *MySQLCategoryRepository) GetAllCategories() ([]domain.Category, error) {
	query := `
		SELECT id, name, description, image_url, is_active, created_at, updated_at 
		FROM categories 
		ORDER BY name ASC
	`

	rows, err := r.conn.FetchRows(query)
	if err != nil {
		return nil, fmt.Errorf("error al obtener las categorías: %v", err)
	}
	defer rows.Close()

	return r.scanCategories(rows)
}

// GetActiveCategories obtiene solo las categorías activas
func (r *MySQLCategoryRepository) GetActiveCategories() ([]domain.Category, error) {
	query := `
		SELECT id, name, description, image_url, is_active, created_at, updated_at 
		FROM categories 
		WHERE is_active = true 
		ORDER BY name ASC
	`

	rows, err := r.conn.FetchRows(query)
	if err != nil {
		return nil, fmt.Errorf("error al obtener las categorías activas: %v", err)
	}
	defer rows.Close()

	return r.scanCategories(rows)
}

// GetCategoryByID obtiene una categoría por su ID
func (r *MySQLCategoryRepository) GetCategoryByID(id int) (*domain.Category, error) {
	query := `
		SELECT id, name, description, image_url, is_active, created_at, updated_at 
		FROM categories 
		WHERE id = ?
	`

	rows, err := r.conn.FetchRows(query, id)
	if err != nil {
		return nil, fmt.Errorf("error al obtener la categoría: %v", err)
	}
	defer rows.Close()

	categories, err := r.scanCategories(rows)
	if err != nil {
		return nil, err
	}

	if len(categories) == 0 {
		return nil, nil
	}

	return &categories[0], nil
}

// GetCategoryByName obtiene una categoría por su nombre
func (r *MySQLCategoryRepository) GetCategoryByName(name string) (*domain.Category, error) {
	query := `
		SELECT id, name, description, image_url, is_active, created_at, updated_at 
		FROM categories 
		WHERE name = ?
	`

	rows, err := r.conn.FetchRows(query, name)
	if err != nil {
		return nil, fmt.Errorf("error al obtener la categoría por nombre: %v", err)
	}
	defer rows.Close()

	categories, err := r.scanCategories(rows)
	if err != nil {
		return nil, err
	}

	if len(categories) == 0 {
		return nil, nil
	}

	return &categories[0], nil
}

// UpdateCategory actualiza una categoría existente
func (r *MySQLCategoryRepository) UpdateCategory(id int, request domain.UpdateCategoryRequest) (*domain.Category, error) {
	// Construir la consulta dinámicamente basada en los campos a actualizar
	setParts := []string{}
	args := []interface{}{}

	if request.Name != "" {
		setParts = append(setParts, "name = ?")
		args = append(args, request.Name)
	}
	if request.Description != "" {
		setParts = append(setParts, "description = ?")
		args = append(args, request.Description)
	}
	if request.ImageURL != "" {
		setParts = append(setParts, "image_url = ?")
		args = append(args, request.ImageURL)
	}
	if request.IsActive != nil {
		setParts = append(setParts, "is_active = ?")
		args = append(args, *request.IsActive)
	}

	if len(setParts) == 0 {
		return r.GetCategoryByID(id) // No hay cambios, devolver la categoría actual
	}

	query := fmt.Sprintf(`
		UPDATE categories 
		SET %s, updated_at = CURRENT_TIMESTAMP 
		WHERE id = ?
	`, fmt.Sprintf("%s", setParts[0]))

	// Agregar las demás partes del SET si existen
	for i := 1; i < len(setParts); i++ {
		query = fmt.Sprintf(`
			UPDATE categories 
			SET %s, updated_at = CURRENT_TIMESTAMP 
			WHERE id = ?
		`, fmt.Sprintf("%s, %s", setParts[0], setParts[i]))
	}

	// Reconstruir la consulta correctamente
	query = fmt.Sprintf(`
		UPDATE categories 
		SET %s, updated_at = CURRENT_TIMESTAMP 
		WHERE id = ?
	`, fmt.Sprintf("%s", setParts[0]))

	for i := 1; i < len(setParts); i++ {
		query = fmt.Sprintf(`
			UPDATE categories 
			SET %s, %s, updated_at = CURRENT_TIMESTAMP 
			WHERE id = ?
		`, setParts[0], setParts[i])
	}

	// Simplificar: construir la consulta de manera más directa
	setClause := ""
	for i, part := range setParts {
		if i > 0 {
			setClause += ", "
		}
		setClause += part
	}

	query = fmt.Sprintf(`
		UPDATE categories 
		SET %s, updated_at = CURRENT_TIMESTAMP 
		WHERE id = ?
	`, setClause)

	args = append(args, id)

	_, err := r.conn.ExecutePreparedQuery(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error al actualizar la categoría: %v", err)
	}

	// Obtener la categoría actualizada
	return r.GetCategoryByID(id)
}

// DeleteCategory elimina una categoría (soft delete)
func (r *MySQLCategoryRepository) DeleteCategory(id int) error {
	query := `UPDATE categories SET is_active = false WHERE id = ?`

	_, err := r.conn.ExecutePreparedQuery(query, id)
	if err != nil {
		return fmt.Errorf("error al eliminar la categoría: %v", err)
	}

	return nil
}

// scanCategories convierte las filas de la base de datos en estructuras Category
func (r *MySQLCategoryRepository) scanCategories(rows *sql.Rows) ([]domain.Category, error) {
	var categories []domain.Category

	for rows.Next() {
		var category domain.Category
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Description,
			&category.ImageURL,
			&category.IsActive,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear la categoría: %v", err)
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error en las filas: %v", err)
	}

	return categories, nil
}
