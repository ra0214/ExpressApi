package infraestructure

import (
	"encoding/json"
	"expresApi/src/products/application"
	wsocket "expresApi/src/websocket"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type DeleteProductController struct {
	useCase *application.DeleteProduct
}

func NewDeleteProductController(useCase *application.DeleteProduct) *DeleteProductController {
	return &DeleteProductController{useCase: useCase}
}

func (d *DeleteProductController) Execute(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID es requerido"})
		return
	}

	err := d.useCase.Execute(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar el producto", "detalles": err.Error()})
		return
	}

	// Enviar notificación WebSocket
	wsMessage := map[string]interface{}{
		"type":      "product_deleted",
		"timestamp": time.Now().Format(time.RFC3339),
		"data": map[string]interface{}{
			"id":     id,
			"action": "eliminado",
		},
	}

	if messageBytes, err := json.Marshal(wsMessage); err == nil {
		wsocket.BroadcastMessage(messageBytes)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Producto eliminado correctamente"})
}
