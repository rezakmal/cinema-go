package controllers

import (
	"cinema-go/models"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CinemaController struct {
	db *sql.DB
}

func NewCinemaController(db *sql.DB) *CinemaController {
	return &CinemaController{db: db}
}

func (bc *CinemaController) CreateCinema(c *gin.Context) {
	var req models.CreateCinemaRequest

	// bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// validate required fields
	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Name cannot be empty",
		})
		return
	}

	if req.Location == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Location cannot be empty",
		})
		return
	}

	// set default rating if not provided
	if req.Rating == 0 {
		req.Rating = 0.0
	}

	// validate rating range
	if req.Rating < 0 || req.Rating > 5 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Rating should be between 0.0 to 5.0",
		})
		return
	}

	// prepare SQL statement

	query := `
			INSER INTO cinema_db (name, location, rating, created_at, updated_at) 
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id, name, location, rating, created_at, updated_at
	`

	// execute query
	now := time.Now()
	var cinema models.Cinema
	err := bc.db.QueryRow(
		query,
		req.Name,
		req.Location,
		req.Rating,
		now,
		now,
	).Scan(
		&cinema.ID,
		&cinema.Name,
		&cinema.Location,
		&cinema.Rating,
		&cinema.CreatedAt,
		&cinema.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to save cinema data",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Cinema added",
		"data":    cinema,
	})
}
