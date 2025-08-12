package infraestructure

import (
	"expresApi/src/products/application"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateProductController struct {
	useCase *application.CreateProduct
}

func NewCreateProductController(useCase *application.CreateProduct) *CreateProductController {
	return &CreateProductController{useCase: useCase}
}

type RequestBody struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
	ImageURL    string  `json:"image_url"`
}

func (ct_c *CreateProductController) Execute(c *gin.Context) {
	var body RequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error al leer el JSON", "detalles": err.Error()})
		return
	}

	err := ct_c.useCase.Execute(body.Name, body.Description, body.Price, body.Category, body.ImageURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar el producto", "detalles": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Producto agregado correctamente"})
}
