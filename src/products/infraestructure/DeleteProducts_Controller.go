package infraestructure

import (
	"expresApi/src/products/application"
	"net/http"

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

	c.JSON(http.StatusOK, gin.H{"message": "Producto eliminado correctamente"})
}
