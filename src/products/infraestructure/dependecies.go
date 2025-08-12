package infraestructure

import (
	"log"
)

func Init() {
	mysqlRepo := NewMySQL()

	router := SetupRouter(mysqlRepo)
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
