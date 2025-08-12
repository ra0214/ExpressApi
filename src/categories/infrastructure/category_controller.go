package infrastructure

import (
	"encoding/json"
	"expresApi/src/categories/application"
	"expresApi/src/categories/domain"
	wsocket "expresApi/src/websocket"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// CategoryController maneja las peticiones HTTP para categorías
type CategoryController struct {
	useCase *application.CategoryUseCase
}

// NewCategoryController crea una nueva instancia del controlador
func NewCategoryController(useCase *application.CategoryUseCase) *CategoryController {
	return &CategoryController{
		useCase: useCase,
	}
}

// CreateCategory crea una nueva categoría
func (c *CategoryController) CreateCategory(ctx *gin.Context) {
	var request domain.CreateCategoryRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Datos inválidos",
			"details": err.Error(),
		})
		return
	}

	// Establecer valor por defecto para IsActive si no se proporciona
	if request.IsActive == false {
		request.IsActive = true
	}

	category, err := c.useCase.CreateCategory(request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error al crear categoría",
			"details": err.Error(),
		})
		return
	}

	// Enviar notificación WebSocket
	wsMessage := map[string]interface{}{
		"type":      "category_created",
		"timestamp": time.Now().Format(time.RFC3339),
		"data": map[string]interface{}{
			"id":          category.ID,
			"name":        category.Name,
			"description": category.Description,
			"image_url":   category.ImageURL,
			"is_active":   category.IsActive,
			"action":      "creada",
		},
	}

	if messageBytes, err := json.Marshal(wsMessage); err == nil {
		wsocket.BroadcastMessage(messageBytes)
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Categoría creada exitosamente",
		"data":    category,
	})
}

// GetAllCategories obtiene todas las categorías
func (c *CategoryController) GetAllCategories(ctx *gin.Context) {
	categories, err := c.useCase.GetAllCategories()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al obtener categorías",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Categorías obtenidas exitosamente",
		"data":    categories,
	})
}

// GetActiveCategories obtiene solo las categorías activas
func (c *CategoryController) GetActiveCategories(ctx *gin.Context) {
	categories, err := c.useCase.GetActiveCategories()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al obtener categorías activas",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Categorías activas obtenidas exitosamente",
		"data":    categories,
	})
}

// GetCategoryByID obtiene una categoría por ID
func (c *CategoryController) GetCategoryByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de categoría inválido",
		})
		return
	}

	category, err := c.useCase.GetCategoryByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   "Categoría no encontrada",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Categoría obtenida exitosamente",
		"data":    category,
	})
}

// UpdateCategory actualiza una categoría existente
func (c *CategoryController) UpdateCategory(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de categoría inválido",
		})
		return
	}

	var request domain.UpdateCategoryRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Datos inválidos",
			"details": err.Error(),
		})
		return
	}

	category, err := c.useCase.UpdateCategory(id, request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error al actualizar categoría",
			"details": err.Error(),
		})
		return
	}

	// Enviar notificación WebSocket
	wsMessage := map[string]interface{}{
		"type":      "category_updated",
		"timestamp": time.Now().Format(time.RFC3339),
		"data": map[string]interface{}{
			"id":          category.ID,
			"name":        category.Name,
			"description": category.Description,
			"image_url":   category.ImageURL,
			"is_active":   category.IsActive,
			"action":      "actualizada",
		},
	}

	if messageBytes, err := json.Marshal(wsMessage); err == nil {
		wsocket.BroadcastMessage(messageBytes)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Categoría actualizada exitosamente",
		"data":    category,
	})
}

// DeleteCategory elimina una categoría (soft delete)
func (c *CategoryController) DeleteCategory(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de categoría inválido",
		})
		return
	}

	err = c.useCase.DeleteCategory(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error al eliminar categoría",
			"details": err.Error(),
		})
		return
	}

	// Enviar notificación WebSocket
	wsMessage := map[string]interface{}{
		"type":      "category_deleted",
		"timestamp": time.Now().Format(time.RFC3339),
		"data": map[string]interface{}{
			"id":     id,
			"action": "eliminada",
		},
	}

	if messageBytes, err := json.Marshal(wsMessage); err == nil {
		wsocket.BroadcastMessage(messageBytes)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Categoría eliminada exitosamente",
	})
}
