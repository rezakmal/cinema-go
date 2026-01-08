package main

import (
	"cinema-go/config"
	"cinema-go/controllers"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	defer db.Close()

	// intiailize Gin router
	r := gin.Default()

	// initialize controller
	cinemaController := controllers.NewCinemaController(db)

	// routes
	r.POST("/cinema", cinemaController.CreateCinema)

	// start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s ", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
