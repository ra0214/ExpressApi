package infrastructure

import (
	"github.com/gin-gonic/gin"
)

// SetupCommentRoutes configura las rutas para comentarios
func SetupCommentRoutes(router *gin.Engine, controller *CommentController) {
	commentGroup := router.Group("/api/comments")
	{
		// Crear comentario
		commentGroup.POST("/", controller.CreateComment)

		// Obtener todos los comentarios
		commentGroup.GET("/", controller.GetAllComments)

		// Obtener comentario por ID
		commentGroup.GET("/:id", controller.GetCommentByID)

		// Obtener comentarios por producto
		commentGroup.GET("/product/:productId", controller.GetCommentsByProduct)

		// Obtener estadísticas por producto
		commentGroup.GET("/product/:productId/stats", controller.GetCommentStats)

		// Eliminar comentario
		commentGroup.DELETE("/:id", controller.DeleteComment)
	}
}
