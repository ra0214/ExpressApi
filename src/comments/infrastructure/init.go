package infrastructure

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

// InitComments inicializa el módulo de comentarios
func InitComments(router *gin.Engine, db *sql.DB) {
	// Crear dependencias
	dependencies := NewCommentDependencies(db)

	// Configurar rutas
	SetupCommentRoutes(router, dependencies.Controller)
}
