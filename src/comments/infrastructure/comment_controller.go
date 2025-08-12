package infrastructure

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"expresApi/src/comments/application"
	"expresApi/src/comments/domain"
	wsocket "expresApi/src/websocket"

	"github.com/gin-gonic/gin"
)

// CommentController maneja las peticiones HTTP relacionadas con comentarios
type CommentController struct {
	useCase *application.CommentUseCase
}

// NewCommentController crea una nueva instancia del controlador
func NewCommentController(useCase *application.CommentUseCase) *CommentController {
	return &CommentController{
		useCase: useCase,
	}
}

// CreateComment crea un nuevo comentario
func (c *CommentController) CreateComment(ctx *gin.Context) {
	var request domain.CreateCommentRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Datos inválidos",
			"details": err.Error(),
		})
		return
	}

	comment, err := c.useCase.CreateComment(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al crear comentario",
			"details": err.Error(),
		})
		return
	}

	// Enviar notificación WebSocket
	wsMessage := map[string]interface{}{
		"type":      "comment_created",
		"timestamp": time.Now().Format(time.RFC3339),
		"data": map[string]interface{}{
			"id":         comment.ID,
			"comment":    comment.Comment,
			"user_name":  comment.UserName,
			"rating":     comment.Rating,
			"product_id": comment.ProductID,
			"action":     "creado",
		},
	}

	if messageBytes, err := json.Marshal(wsMessage); err == nil {
		wsocket.BroadcastMessage(messageBytes)
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Comentario creado exitosamente",
		"data":    comment,
	})
}

// GetAllComments obtiene todos los comentarios
func (c *CommentController) GetAllComments(ctx *gin.Context) {
	comments, err := c.useCase.GetAllComments()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al obtener comentarios",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Comentarios obtenidos exitosamente",
		"data":    comments,
	})
}

// GetCommentsByProduct obtiene comentarios por producto
func (c *CommentController) GetCommentsByProduct(ctx *gin.Context) {
	productIDStr := ctx.Param("productId")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de producto inválido",
		})
		return
	}

	comments, err := c.useCase.GetCommentsByProduct(productID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al obtener comentarios",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Comentarios obtenidos exitosamente",
		"data":    comments,
	})
}

// GetCommentByID obtiene un comentario por ID
func (c *CommentController) GetCommentByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "ID inválido",
		})
		return
	}

	comment, err := c.useCase.GetCommentByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   "Comentario no encontrado",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Comentario obtenido exitosamente",
		"data":    comment,
	})
}

// DeleteComment elimina un comentario
func (c *CommentController) DeleteComment(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "ID inválido",
		})
		return
	}

	err = c.useCase.DeleteComment(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al eliminar comentario",
			"details": err.Error(),
		})
		return
	}

	// Enviar notificación WebSocket
	wsMessage := map[string]interface{}{
		"type":      "comment_deleted",
		"timestamp": time.Now().Format(time.RFC3339),
		"data": map[string]interface{}{
			"id":     id,
			"action": "eliminado",
		},
	}

	if messageBytes, err := json.Marshal(wsMessage); err == nil {
		wsocket.BroadcastMessage(messageBytes)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Comentario eliminado exitosamente",
	})
}

// GetCommentStats obtiene estadísticas de comentarios
func (c *CommentController) GetCommentStats(ctx *gin.Context) {
	productIDStr := ctx.Param("productId")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de producto inválido",
		})
		return
	}

	stats, err := c.useCase.GetCommentStats(productID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al obtener estadísticas",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Estadísticas obtenidas exitosamente",
		"data":    stats,
	})
}
