package routes

import (
	"cinema-go/controllers"

	"github.com/gin-gonic/gin"
)

func SetupCinemaRoutes(router *gin.Engine, cinemaController *controllers.CinemaController) {
	// Cinema routes
	cinema := router.Group("/cinema")
	{
		cinema.POST("", cinemaController.CreateCinema)       // POST /cinema
		cinema.GET("", cinemaController.GetCinemas)          // GET /cinema
		cinema.GET("/:id", cinemaController.GetCinemaByID)   // GET /cinema/:id
		cinema.PUT("/:id", cinemaController.UpdateCinema)    // PUT /cinema/:id
		cinema.DELETE("/:id", cinemaController.DeleteCinema) // DELETE /cinema/:id
	}
}
