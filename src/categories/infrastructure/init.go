package infrastructure

import (
	"expresApi/src/categories/application"
	"log"

	"github.com/gin-gonic/gin"
)

// InitCategories inicializa el módulo de categorías
func InitCategories(r *gin.Engine) {
	log.Println("Inicializando módulo de categorías...")

	// Crear grupo de rutas API
	apiGroup := r.Group("/api/v1")

	// Configurar rutas de categorías
	SetupCategoryRoutes(apiGroup)

	log.Println("Módulo de categorías inicializado correctamente")
}

// GetCategoryDependencies retorna las dependencias para usar en otros módulos
func GetCategoryDependencies() (*application.CategoryUseCase, error) {
	repository := NewMySQLCategoryRepository()
	useCase := application.NewCategoryUseCase(repository)
	return useCase, nil
}
