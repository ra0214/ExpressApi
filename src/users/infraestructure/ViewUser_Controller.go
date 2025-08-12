package infraestructure

import (
	"expresApi/src/users/application"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ViewUserController struct {
	useCase *application.ViewUser
}

func NewViewUserController(useCase *application.ViewUser) *ViewUserController {
	return &ViewUserController{useCase: useCase}
}

func (eu_c *ViewUserController) Execute(c *gin.Context) {
	user, err := eu_c.useCase.Execute()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}