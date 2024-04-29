package controllers

import (
	"net/http"

	"github.com/diana-gemini/ozinshe/db/initializers"
	"github.com/diana-gemini/ozinshe/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetTrends godoc
// @Summary GetTrends
// @Security ApiKeyAuth
// @Tags main-page-controller
// @ID get-trends
// @Accept json
// @Produce json
// @Success 200 {integer} integer 1
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /trends [get]
func GetTrends(c *gin.Context) {
	var movies []models.Movie
	result := initializers.DB.Preload("Categories").
		Preload("Screenshots").
		Preload("Seasons", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Videos")
		}).Order("count_of_watch desc").Find(&movies)
	if err := result.Error; err != nil {
		newErrorResponse(c, http.StatusNotFound, "trend movies not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Movies": movies,
	})
}

// GetNewprojects godoc
// @Summary GetNewprojects
// @Security ApiKeyAuth
// @Tags main-page-controller
// @ID get-new-projects
// @Accept json
// @Produce json
// @Success 200 {integer} integer 1
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /newprojects [get]
func GetNewprojects(c *gin.Context) {
	var movies []models.Movie
	result := initializers.DB.Preload("Categories").
		Preload("Screenshots").
		Preload("Seasons", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Videos")
		}).Order("created_at desc").Find(&movies)
	if err := result.Error; err != nil {
		newErrorResponse(c, http.StatusNotFound, "new projects not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Movies": movies,
	})
}

// GetTelehikaya godoc
// @Summary GetTelehikaya
// @Security ApiKeyAuth
// @Tags main-page-controller
// @ID get-telehikaya
// @Accept json
// @Produce json
// @Success 200 {integer} integer 1
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /telehikaya [get]
func GetTelehikaya(c *gin.Context) {
	var types models.Type
	result := initializers.DB.Where("type_name = ?", "Serial").Find(&types)
	if err := result.Error; err != nil {
		newErrorResponse(c, http.StatusNotFound, "serial type not found")
		return
	}

	var movies []models.Movie
	result = initializers.DB.
		Preload("Categories").
		Preload("Screenshots").
		Preload("Seasons", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Videos")
		}).Where("type_id = ?", types.ID).
		Order("created_at desc").
		Find(&movies)

	if err := result.Error; err != nil {
		newErrorResponse(c, http.StatusNotFound, "telehikaya not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Movies": movies,
	})
}

// Horor godoc
// @Summary Horor
// @Security ApiKeyAuth
// @Tags main-page-controller
// @ID horor
// @Accept json
// @Produce json
// @Success 200 {integer} integer 1
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /horor [get]
func Horor(c *gin.Context) {
	var movies []models.Movie
	result := initializers.DB.
		Joins("JOIN movie_category ON movies.id = movie_category.movie_id").
		Joins("JOIN categories ON movie_category.category_id = categories.id").
		Where("categories.category_name = ?", "Horor").
		Preload("Categories").
		Preload("Screenshots").
		Preload("Seasons", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Videos")
		}).Find(&movies)

	if err := result.Error; err != nil {
		newErrorResponse(c, http.StatusNotFound, "horor not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Horor": movies,
	})
}

// Anime godoc
// @Summary Anime
// @Security ApiKeyAuth
// @Tags main-page-controller
// @ID anime
// @Accept json
// @Produce json
// @Success 200 {integer} integer 1
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /anime [get]
func Anime(c *gin.Context) {
	var movies []models.Movie
	result := initializers.DB.
		Joins("JOIN movie_category ON movies.id = movie_category.movie_id").
		Joins("JOIN categories ON movie_category.category_id = categories.id").
		Where("categories.category_name = ?", "Anime").
		Preload("Categories").
		Preload("Screenshots").
		Preload("Seasons", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Videos")
		}).Find(&movies)

	if err := result.Error; err != nil {
		newErrorResponse(c, http.StatusNotFound, "anime not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Anime": movies,
	})
}
