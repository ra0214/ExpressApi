package main

import (
	categoryInfra "expresApi/src/categories/infrastructure"
	commentInfra "expresApi/src/comments/infrastructure"
	"expresApi/src/config"
	"expresApi/src/config/middleware"
	"expresApi/src/products/infraestructure"
	userInfra "expresApi/src/users/infraestructure"
	wsocket "expresApi/src/websocket"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Crear router principal
	r := gin.Default()

	// Agregar middlewares
	r.Use(middleware.NewCorsMiddleware())

	// Obtener conexión a la base de datos
	dbConfig := config.GetDBPool()
	if dbConfig.Err != "" {
		log.Fatalf("Error al configurar el pool de conexiones: %v", dbConfig.Err)
	}

	// Inicializar repositorios MySQL
	productRepo := infraestructure.NewMySQL()
	userRepo := userInfra.NewMySQL()

	// Configurar rutas de usuarios
	userGroup := r.Group("/api/v1")
	userInfra.RegisterRoutes(userGroup, userRepo)

	// Configurar rutas de productos
	productGroup := r.Group("/api/v1")
	infraestructure.RegisterRoutes(productGroup, productRepo)

	// Configurar rutas de comentarios
	commentInfra.InitComments(r, dbConfig.DB)

	// Configurar rutas de categorías
	categoryGroup := r.Group("/api/v1")
	categoryInfra.RegisterCategoryRoutes(categoryGroup)

	// Configurar WebSocket
	go wsocket.WSHub.Run()
	r.GET("/ws", wsocket.HandleWebSocket)

	// Configurar servidor
	r.SetTrustedProxies([]string{"127.0.0.1"})

	// Iniciar servidor
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
