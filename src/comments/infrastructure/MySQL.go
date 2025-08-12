package infrastructure

import (
	"database/sql"
	"expresApi/src/comments/domain"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLCommentRepository struct {
	db *sql.DB
}

func NewMySQLCommentRepository(db *sql.DB) domain.CommentRepository {
	return &MySQLCommentRepository{
		db: db,
	}
}

// Create crea un nuevo comentario
func (r *MySQLCommentRepository) Create(req domain.CreateCommentRequest) (*domain.Comment, error) {
	query := `
		INSERT INTO comments (comment, user_name, rating, product_id, created_at) 
		VALUES (?, ?, ?, ?, NOW())
	`

	result, err := r.db.Exec(query, req.Comment, req.UserName, req.Rating, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("error creating comment: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert id: %w", err)
	}

	// Obtener el comentario creado
	comment, err := r.GetByID(int(id))
	if err != nil {
		return nil, fmt.Errorf("error retrieving created comment: %w", err)
	}

	return comment, nil
}

// GetByProductID obtiene todos los comentarios de un producto
func (r *MySQLCommentRepository) GetByProductID(productID int) ([]domain.Comment, error) {
	query := `
		SELECT id, comment, user_name, rating, product_id, 
		       DATE_FORMAT(created_at, '%Y-%m-%d %H:%i:%s') as created_at_formatted
		FROM comments 
		WHERE product_id = ? 
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, productID)
	if err != nil {
		return nil, fmt.Errorf("error querying comments: %w", err)
	}
	defer rows.Close()

	var comments []domain.Comment
	for rows.Next() {
		var comment domain.Comment
		var createdAtStr string

		err := rows.Scan(
			&comment.ID,
			&comment.Comment,
			&comment.UserName,
			&comment.Rating,
			&comment.ProductID,
			&createdAtStr,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning comment: %w", err)
		}

		// Convertir string a time.Time con manejo de errores mejorado
		if createdAtStr != "" {
			createdAt, err := time.Parse("2006-01-02 15:04:05", createdAtStr)
			if err != nil {
				// Intentar otros formatos comunes
				createdAt, err = time.Parse("2006-01-02T15:04:05Z", createdAtStr)
				if err != nil {
					// Si todo falla, usar fecha actual
					createdAt = time.Now()
				}
			}
			comment.CreatedAt = createdAt
		} else {
			comment.CreatedAt = time.Now()
		}

		comments = append(comments, comment)
	}

	return comments, nil
}

// GetAll obtiene todos los comentarios
func (r *MySQLCommentRepository) GetAll() ([]domain.Comment, error) {
	query := `
		SELECT id, comment, user_name, rating, product_id, 
		       DATE_FORMAT(created_at, '%Y-%m-%d %H:%i:%s') as created_at_formatted
		FROM comments 
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying all comments: %w", err)
	}
	defer rows.Close()

	var comments []domain.Comment
	for rows.Next() {
		var comment domain.Comment
		var createdAtStr string

		err := rows.Scan(
			&comment.ID,
			&comment.Comment,
			&comment.UserName,
			&comment.Rating,
			&comment.ProductID,
			&createdAtStr,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning comment: %w", err)
		}

		// Convertir string a time.Time con manejo de errores mejorado
		if createdAtStr != "" {
			createdAt, err := time.Parse("2006-01-02 15:04:05", createdAtStr)
			if err != nil {
				// Intentar otros formatos comunes
				createdAt, err = time.Parse("2006-01-02T15:04:05Z", createdAtStr)
				if err != nil {
					// Si todo falla, usar fecha actual
					createdAt = time.Now()
				}
			}
			comment.CreatedAt = createdAt
		} else {
			comment.CreatedAt = time.Now()
		}

		comments = append(comments, comment)
	}

	return comments, nil
}

// GetByID obtiene un comentario por su ID
func (r *MySQLCommentRepository) GetByID(id int) (*domain.Comment, error) {
	query := `
		SELECT id, comment, user_name, rating, product_id, 
		       DATE_FORMAT(created_at, '%Y-%m-%d %H:%i:%s') as created_at_formatted
		FROM comments 
		WHERE id = ?
	`

	var comment domain.Comment
	var createdAtStr string

	err := r.db.QueryRow(query, id).Scan(
		&comment.ID,
		&comment.Comment,
		&comment.UserName,
		&comment.Rating,
		&comment.ProductID,
		&createdAtStr,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("comment not found")
		}
		return nil, fmt.Errorf("error getting comment: %w", err)
	}

	// Convertir string a time.Time con manejo de errores mejorado
	if createdAtStr != "" {
		createdAt, err := time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			// Intentar otros formatos comunes
			createdAt, err = time.Parse("2006-01-02T15:04:05Z", createdAtStr)
			if err != nil {
				// Si todo falla, usar fecha actual
				createdAt = time.Now()
			}
		}
		comment.CreatedAt = createdAt
	} else {
		comment.CreatedAt = time.Now()
	}

	return &comment, nil
}

// Delete elimina un comentario
func (r *MySQLCommentRepository) Delete(id int) error {
	query := `DELETE FROM comments WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting comment: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("comment not found")
	}

	return nil
}

// DeleteByProductID elimina todos los comentarios de un producto
func (r *MySQLCommentRepository) DeleteByProductID(productID int) error {
	query := `DELETE FROM comments WHERE product_id = ?`

	result, err := r.db.Exec(query, productID)
	if err != nil {
		return fmt.Errorf("error deleting comments for product %d: %w", productID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	// No es un error si no hay comentarios para eliminar
	fmt.Printf("Eliminados %d comentarios del producto %d\n", rowsAffected, productID)
	return nil
}

// GetStats obtiene estadísticas de comentarios para un producto
func (r *MySQLCommentRepository) GetStats(productID int) (*domain.CommentStats, error) {
	// Obtener total de comentarios y promedio de rating
	mainQuery := `
		SELECT 
			COUNT(*) as total_comments,
			COALESCE(AVG(rating), 0) as average_rating
		FROM comments 
		WHERE product_id = ?
	`

	var stats domain.CommentStats
	err := r.db.QueryRow(mainQuery, productID).Scan(
		&stats.TotalComments,
		&stats.AverageRating,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting comment stats: %w", err)
	}

	stats.ProductID = productID

	// Obtener conteo por rating
	ratingQuery := `
		SELECT rating, COUNT(*) as count
		FROM comments 
		WHERE product_id = ?
		GROUP BY rating
	`

	rows, err := r.db.Query(ratingQuery, productID)
	if err != nil {
		return nil, fmt.Errorf("error getting rating counts: %w", err)
	}
	defer rows.Close()

	stats.RatingCounts = make(map[int]int)
	for rows.Next() {
		var rating, count int
		err := rows.Scan(&rating, &count)
		if err != nil {
			return nil, fmt.Errorf("error scanning rating count: %w", err)
		}
		stats.RatingCounts[rating] = count
	}

	return &stats, nil
}
