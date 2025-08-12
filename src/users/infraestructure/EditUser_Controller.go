package infraestructure

import (
	"expresApi/src/users/application"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EditUserController struct {
	useCase *application.EditUser
}

func NewEditUserController(useCase *application.EditUser) *EditUserController {
	return &EditUserController{useCase: useCase}
}

func (eu_c *EditUserController) Execute(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuario inválido"})
		return
	}

	var body struct {
		ID          int32 `json:"id"`
	    UserName	string `json:"userName"`
	    Email 		string `json:"email"`
	    Password	string `json:"password"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error al leer los datos"})
		return
	}

	err = eu_c.useCase.Execute(int32(id), body.UserName, body.Email, body.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el usuario", "detalles": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Usuario actualizado correctamente"})
}