package infraestructure

import (
	"expresApi/src/users/application"
	"expresApi/src/users/domain"

	"github.com/gin-gonic/gin"
)

func SetupRouter(repo domain.IUser) *gin.Engine {
	r := gin.Default()

	createUser := application.NewCreateUser(repo)
	createUserController := NewCreateUserController(createUser)

	viewUser := application.NewViewUser(repo)
	viewUserController := NewViewUserController(viewUser)

	editUserUseCase := application.NewEditUser(repo)
	editUserController := NewEditUserController(editUserUseCase)

	deleteUserUseCase := application.NewDeleteUser(repo)
	deleteUserController := NewDeleteUserController(deleteUserUseCase)

	loginUser := application.NewLoginUser(repo)
	loginUserController := NewLoginUserController(loginUser)

	r.POST("/user", createUserController.Execute)
	r.GET("/user", viewUserController.Execute)
	r.PUT("/user/:id", editUserController.Execute)
	r.DELETE("/user/:id", deleteUserController.Execute)
	r.POST("/login", loginUserController.Execute)

	return r
}

// Nueva función para registrar rutas en un grupo
func RegisterRoutes(r *gin.RouterGroup, repo domain.IUser) {
	// Crear instancia de MySQL para acceder al método ToggleUserStatusController
	mysqlRepo := repo.(*MySQL)

	createUser := application.NewCreateUser(repo)
	createUserController := NewCreateUserController(createUser)

	viewUser := application.NewViewUser(repo)
	viewUserController := NewViewUserController(viewUser)

	editUserUseCase := application.NewEditUser(repo)
	editUserController := NewEditUserController(editUserUseCase)

	deleteUserUseCase := application.NewDeleteUser(repo)
	deleteUserController := NewDeleteUserController(deleteUserUseCase)

	loginUser := application.NewLoginUser(repo)
	loginUserController := NewLoginUserController(loginUser)

	r.POST("/users", createUserController.Execute)
	r.GET("/users", viewUserController.Execute)
	r.PUT("/users/:id", editUserController.Execute)
	r.DELETE("/users/:id", deleteUserController.Execute)
	r.POST("/login", loginUserController.Execute)
	r.PUT("/users/:id/toggle-status", mysqlRepo.ToggleUserStatusController)
	// También agregar la función independiente como alternativa
	r.PATCH("/users/:id/toggle-status", ToggleUserStatusController)
}
