// cmd/main.go
package main

import (
	"log"
	"os"

	"merch-shop/internal/database"
	"merch-shop/internal/handlers"
	"merch-shop/internal/middleware"
	"merch-shop/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	repo := repository.NewPostgresRepository(db)

	router := gin.Default()

	router.POST("/api/auth", handlers.AuthHandler(repo))

	authGroup := router.Group("/api")
	authGroup.Use(middleware.JWTAuthMiddleware())
	{
		authGroup.GET("/info", handlers.InfoHandler(repo))
		authGroup.POST("/sendCoin", handlers.SendCoinHandler(repo))
		authGroup.GET("/buy/:item", handlers.BuyHandler(repo))
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on port %s", port)
	router.Run(":" + port)
}
