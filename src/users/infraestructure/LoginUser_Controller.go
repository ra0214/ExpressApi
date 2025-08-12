package infraestructure

import (
	"expresApi/src/users/application"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginUserController struct {
	useCase *application.LoginUser
}

func NewLoginUserController(useCase *application.LoginUser) *LoginUserController {
	return &LoginUserController{useCase: useCase}
}

type LoginRequestBody struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

func (lc *LoginUserController) Execute(c *gin.Context) {
	var body LoginRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error al leer el JSON"})
		return
	}

	user, err := lc.useCase.Execute(body.UserName, body.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales inválidas"})
		return
	}

	// No devolver la contraseña en la respuesta
	user.Password = ""
	c.JSON(http.StatusOK, gin.H{
		"message": "Login exitoso",
		"user":    user,
	})
}
