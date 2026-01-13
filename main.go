package main

import (
	"cinema-go/config"
	"cinema-go/controllers"
	"cinema-go/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
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

	// run database migrations
	if err := config.RunMigrations(db); err != nil {
		log.Fatal("Failed to run migrations: ", err)
	}

	// intiailize Gin router
	r := gin.Default()

	// initialize controller
	cinemaController := controllers.NewCinemaController(db)

	// setup routes
	routes.SetupCinemaRoutes(r, cinemaController)

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
