package infraestructure

import (
	"expresApi/src/products/application"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ViewProductController struct {
	useCase *application.ViewProduct
}

func NewViewProductController(useCase *application.ViewProduct) *ViewProductController {
	return &ViewProductController{useCase: useCase}
}

func (et_c *ViewProductController) Execute(c *gin.Context) {
	products, err := et_c.useCase.Execute()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los productos", "detalles": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}
