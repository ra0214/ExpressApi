package infraestructure

import (
	"expresApi/src/users/application"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DeleteUserController struct {
	useCase *application.DeleteUser
}

func NewDeleteUserController(useCase *application.DeleteUser) *DeleteUserController {
	return &DeleteUserController{useCase: useCase}
}

func (du_c *DeleteUserController) Execute(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de producto inválido"})
		return
	}

	err = du_c.useCase.Execute(int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error al eliminar el usuario: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Usuario eliminado correctamente"})
}