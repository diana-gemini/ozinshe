package controllers

import (
	"net/http"
	"ozinshe/db/initializers"
	"ozinshe/internal/models"

	"github.com/gin-gonic/gin"
)

func Search(c *gin.Context) {
	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "query parameter is required",
		})
		return
	}

	var movies []models.Movie
	result := initializers.DB.Where("LOWER(name_of_project) LIKE LOWER(?) OR LOWER(category) LIKE LOWER(?)", "%"+query+"%", "%"+query+"%").Find(&movies)
	if err := result.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"movies": movies,
	})
}
