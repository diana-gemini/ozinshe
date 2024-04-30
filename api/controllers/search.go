package controllers

import (
	"net/http"

	"github.com/diana-gemini/ozinshe/db/initializers"
	"github.com/diana-gemini/ozinshe/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Search godoc
// @Summary Search
// @Security ApiKeyAuth
// @Tags search-controller
// @ID search
// @Accept json
// @Produce json
// @Param search query string true "search param received in the URL"
// @Success 200 {integer} integer 1
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /search [get]
func Search(c *gin.Context) {
	search := c.Query("search")
	if search == "" {
		NewErrorResponse(c, http.StatusBadRequest, "search not found in URL")
		return
	}

	var movies []models.Movie

	result := initializers.DB.Preload("Categories").
		Preload("Screenshots").
		Preload("Seasons", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Videos")
		}).Where("LOWER(name_of_project) LIKE LOWER(?) ", "%"+search+"%").Find(&movies)
	if err := result.Error; err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"movies": movies,
	})
}
