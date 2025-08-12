package infrastructure

import (
	"expresApi/src/categories/application"

	"github.com/gin-gonic/gin"
)

// SetupCategoryRoutes configura las rutas para categorías
func SetupCategoryRoutes(r *gin.RouterGroup) {
	// Crear dependencias
	repository := NewMySQLCategoryRepository()
	useCase := application.NewCategoryUseCase(repository)
	controller := NewCategoryController(useCase)

	// Configurar rutas
	categoryGroup := r.Group("/categories")
	{
		// GET /api/v1/categories - Obtener todas las categorías
		categoryGroup.GET("", controller.GetAllCategories)

		// GET /api/v1/categories/active - Obtener solo categorías activas
		categoryGroup.GET("/active", controller.GetActiveCategories)

		// GET /api/v1/categories/:id - Obtener categoría por ID
		categoryGroup.GET("/:id", controller.GetCategoryByID)

		// POST /api/v1/categories - Crear nueva categoría
		categoryGroup.POST("", controller.CreateCategory)

		// PUT /api/v1/categories/:id - Actualizar categoría
		categoryGroup.PUT("/:id", controller.UpdateCategory)

		// DELETE /api/v1/categories/:id - Eliminar categoría (soft delete)
		categoryGroup.DELETE("/:id", controller.DeleteCategory)
	}
}

// RegisterCategoryRoutes registra las rutas en un grupo existente
func RegisterCategoryRoutes(r *gin.RouterGroup) {
	SetupCategoryRoutes(r)
}
