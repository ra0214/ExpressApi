package infraestructure

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// ToggleUserStatusController maneja la solicitud HTTP para cambiar el estado de un usuario
func ToggleUserStatusController(c *gin.Context) {
	// Obtener el ID del usuario desde los parámetros de la URL
	idParam := c.Param("id")

	// Convertir el ID a int32
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "ID de usuario inválido",
		})
		return
	}

	// Crear instancia del repositorio
	repo := NewMySQL()

	// Llamar al método para cambiar el estado
	newStatus, err := repo.ToggleUserStatus(int32(id))
	if err != nil {
		c.JSON(500, gin.H{
			"error":   "Error al cambiar el estado del usuario",
			"details": err.Error(),
		})
		return
	}

	// Responder con el nuevo estado
	c.JSON(200, gin.H{
		"message":    "Estado del usuario actualizado correctamente",
		"user_id":    id,
		"new_status": newStatus,
	})
}
