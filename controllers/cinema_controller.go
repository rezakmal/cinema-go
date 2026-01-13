package controllers

import (
	"cinema-go/models"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
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
			INSERT INTO cinema (name, location, rating, created_at, updated_at)
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

func (bc *CinemaController) GetCinemas(c *gin.Context) {
	query := `SELECT id, name, location, rating, created_at, updated_at FROM cinema ORDER BY created_at DESC`

	rows, err := bc.db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch cinemas",
			"details": err.Error(),
		})
		return
	}
	defer rows.Close()

	var cinemas []models.Cinema
	for rows.Next() {
		var cinema models.Cinema
		err := rows.Scan(
			&cinema.ID,
			&cinema.Name,
			&cinema.Location,
			&cinema.Rating,
			&cinema.CreatedAt,
			&cinema.UpdatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to scan cinema data",
				"details": err.Error(),
			})
			return
		}
		cinemas = append(cinemas, cinema)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error occurred while reading cinemas",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cinemas retrieved successfully",
		"data":    cinemas,
		"count":   len(cinemas),
	})
}

func (bc *CinemaController) GetCinemaByID(c *gin.Context) {
	id := c.Param("id")

	// Convert id to int
	var cinemaID int
	_, err := fmt.Sscanf(id, "%d", &cinemaID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid cinema ID format",
		})
		return
	}

	query := `SELECT id, name, location, rating, created_at, updated_at FROM cinema WHERE id = $1`

	var cinema models.Cinema
	err = bc.db.QueryRow(query, cinemaID).Scan(
		&cinema.ID,
		&cinema.Name,
		&cinema.Location,
		&cinema.Rating,
		&cinema.CreatedAt,
		&cinema.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Cinema not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch cinema",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cinema retrieved successfully",
		"data":    cinema,
	})
}

func (bc *CinemaController) UpdateCinema(c *gin.Context) {
	id := c.Param("id")

	// Convert id to int
	var cinemaID int
	_, err := fmt.Sscanf(id, "%d", &cinemaID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid cinema ID format",
		})
		return
	}

	var req models.UpdateCinemaRequest

	// bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// validate rating range if provided
	if req.Rating != nil && (*req.Rating < 0 || *req.Rating > 5) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Rating should be between 0.0 to 5.0",
		})
		return
	}

	// check if cinema exists
	var existingCinema models.Cinema
	checkQuery := `SELECT id, name, location, rating, created_at, updated_at FROM cinema WHERE id = $1`
	err = bc.db.QueryRow(checkQuery, cinemaID).Scan(
		&existingCinema.ID,
		&existingCinema.Name,
		&existingCinema.Location,
		&existingCinema.Rating,
		&existingCinema.CreatedAt,
		&existingCinema.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Cinema not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to check cinema existence",
			"details": err.Error(),
		})
		return
	}

	// build dynamic update query based on provided fields
	var updateFields []string
	var args []interface{}
	argCount := 1

	// check which fields are provided and add to update
	if req.Name != nil {
		updateFields = append(updateFields, fmt.Sprintf("name = $%d", argCount))
		args = append(args, *req.Name)
		argCount++
	}

	if req.Location != nil {
		updateFields = append(updateFields, fmt.Sprintf("location = $%d", argCount))
		args = append(args, *req.Location)
		argCount++
	}

	if req.Rating != nil {
		updateFields = append(updateFields, fmt.Sprintf("rating = $%d", argCount))
		args = append(args, *req.Rating)
		argCount++
	}

	updateFields = append(updateFields, "updated_at = CURRENT_TIMESTAMP")

	args = append(args, cinemaID)

	if len(updateFields) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No valid fields to update",
		})
		return
	}

	updateQuery := fmt.Sprintf(`
		UPDATE cinema
		SET %s
		WHERE id = $%d
		RETURNING id, name, location, rating, created_at, updated_at
	`, strings.Join(updateFields, ", "), len(args))

	var updatedCinema models.Cinema
	err = bc.db.QueryRow(updateQuery, args...).Scan(
		&updatedCinema.ID,
		&updatedCinema.Name,
		&updatedCinema.Location,
		&updatedCinema.Rating,
		&updatedCinema.CreatedAt,
		&updatedCinema.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update cinema",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cinema updated successfully",
		"data":    updatedCinema,
	})
}
