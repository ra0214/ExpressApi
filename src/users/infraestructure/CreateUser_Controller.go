package infraestructure

import (
	"expresApi/src/users/application"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateUserController struct {
	useCase *application.CreateUser
}

func NewCreateUserController(useCase *application.CreateUser) *CreateUserController {
	return &CreateUserController{useCase: useCase}
}

type RequestBody struct {
	UserName	string `json:"userName"`
	Email 		string `json:"email"`
	Password	string `json:"password"`
}

func (cu_c *CreateUserController) Execute(c *gin.Context) {
	var body RequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error al leer el JSON", "detalles": err.Error()})
		return
	}

	err := cu_c.useCase.Execute(body.UserName, body.Email, body.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al agregar el usuario", "detalles": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Usuario agregado correctamente"})
}