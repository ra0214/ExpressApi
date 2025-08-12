package infraestructure

import (
	"encoding/json"
	"expresApi/src/products/application"
	wsocket "expresApi/src/websocket"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type UpdateProductController struct {
	useCase *application.UpdateProduct
}

func NewUpdateProductController(useCase *application.UpdateProduct) *UpdateProductController {
	return &UpdateProductController{useCase: useCase}
}

type UpdateRequestBody struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
	ImageURL    string  `json:"image_url"`
}

func (u *UpdateProductController) Execute(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID es requerido"})
		return
	}

	var body UpdateRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error al leer el JSON", "detalles": err.Error()})
		return
	}

	err := u.useCase.Execute(id, body.Name, body.Description, body.Price, body.Category, body.ImageURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el producto", "detalles": err.Error()})
		return
	}

	// Enviar notificación WebSocket
	wsMessage := map[string]interface{}{
		"type":      "product_update",
		"timestamp": time.Now().Format(time.RFC3339),
		"data": map[string]interface{}{
			"id":          id,
			"name":        body.Name,
			"description": body.Description,
			"price":       body.Price,
			"category":    body.Category,
			"image_url":   body.ImageURL,
			"action":      "actualizado",
		},
	}

	if messageBytes, err := json.Marshal(wsMessage); err == nil {
		wsocket.BroadcastMessage(messageBytes)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Producto actualizado correctamente"})
}
