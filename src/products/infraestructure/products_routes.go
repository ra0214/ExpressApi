package infraestructure

import (
	"expresApi/src/config"
	"expresApi/src/products/application"
	"expresApi/src/products/domain"
	"log"

	"github.com/gin-gonic/gin"
)

func SetupRouter(repo domain.IProduct) *gin.Engine {
	r := gin.Default()

	CreateProduct := application.NewCreateProduct(repo)
	createProductController := NewCreateProductController(CreateProduct)

	viewProduct := application.NewViewProduct(repo)
	viewProductController := NewViewProductController(viewProduct)

	r.POST("/product", createProductController.Execute)
	r.GET("/product", viewProductController.Execute)

	return r
}

// Nueva función para registrar rutas en un grupo
func RegisterRoutes(r *gin.RouterGroup, repo domain.IProduct) {
	// Obtener conexión a la base de datos para comentarios
	dbConfig := config.GetDBPool()
	if dbConfig.Err != "" {
		log.Printf("Error al configurar el pool de conexiones para comentarios: %v", dbConfig.Err)
	}

	// Crear adaptador de repositorio de comentarios
	var commentRepo application.CommentRepository
	if dbConfig.Err == "" {
		commentRepo = NewCommentRepositoryAdapter(dbConfig.DB)
	}

	CreateProduct := application.NewCreateProduct(repo)
	createProductController := NewCreateProductController(CreateProduct)

	viewProduct := application.NewViewProduct(repo)
	viewProductController := NewViewProductController(viewProduct)

	updateProduct := application.NewUpdateProduct(repo)
	updateProductController := NewUpdateProductController(updateProduct)

	deleteProduct := application.NewDeleteProduct(repo, commentRepo)
	deleteProductController := NewDeleteProductController(deleteProduct)

	r.POST("/products", createProductController.Execute)
	r.GET("/products", viewProductController.Execute)
	r.PUT("/products/:id", updateProductController.Execute)
	r.DELETE("/products/:id", deleteProductController.Execute)
}
