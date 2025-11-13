package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"go-api/internal/db"
	"go-api/internal/handlers"
)

func main() {
	db.InitDB()

	router := gin.Default()

	api := router.Group("/api/v1")
	{
		handlers.RegisterUserRoutes(api)
	}

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}
