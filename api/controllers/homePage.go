package controllers

import (
	"net/http"

	"github.com/diana-gemini/ozinshe/db/initializers"
	"github.com/diana-gemini/ozinshe/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var limitOfMovie = 5

// Home godoc
// @Summary Home
// @Security ApiKeyAuth
// @Tags main-page-controller
// @ID home
// @Accept  json
// @Produce  json
// @Success 200 {integer} integer 1
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /home [get]
func Home(c *gin.Context) {
	var trendMovies []models.Movie
	trendResult := initializers.DB.Preload("Categories").
		Preload("Screenshots").
		Preload("Seasons", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Videos")
		}).Order("count_of_watch desc").
		Limit(limitOfMovie).
		Find(&trendMovies)

	if err := trendResult.Error; err != nil {
		NewErrorResponse(c, http.StatusNotFound, "trend movies not found")
	}

	var newMovies []models.Movie
	newResult := initializers.DB.Preload("Categories").
		Preload("Screenshots").
		Preload("Seasons", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Videos")
		}).Order("created_at desc").
		Limit(limitOfMovie).
		Find(&newMovies)

	if err := newResult.Error; err != nil {
		NewErrorResponse(c, http.StatusNotFound, "new projects not found")
	}

	var types models.Type
	typesResult := initializers.DB.Where("type_name = ?", "Serial").Find(&types)
	if err := typesResult.Error; err != nil {
		NewErrorResponse(c, http.StatusNotFound, "serial type not found")
	}

	var telehikayaMovies []models.Movie
	typesResult = initializers.DB.
		Preload("Categories").
		Preload("Screenshots").
		Preload("Seasons", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Videos")
		}).Where("type_id = ?", types.ID).
		Order("created_at desc").
		Limit(limitOfMovie).
		Find(&telehikayaMovies)
	if err := typesResult.Error; err != nil {
		NewErrorResponse(c, http.StatusNotFound, "telehikaya not found")
	}

	hororMovies, errhororMovies := getMoviesByCategory("Horor")
	if errhororMovies != nil {
		NewErrorResponse(c, http.StatusNotFound, "horor not found")
	}

	animeMovies, errAnimeMovies := getMoviesByCategory("Anime")
	if errAnimeMovies != nil {
		NewErrorResponse(c, http.StatusNotFound, "anime not found")
	}

	var category []models.Category
	categoryResult := initializers.DB.Find(&category)
	if err := categoryResult.Error; err != nil {
		NewErrorResponse(c, http.StatusNotFound, "category not found")
	}

	var ageCategory []models.AgeCategory
	ageCategoryResult := initializers.DB.Find(&ageCategory)
	if err := ageCategoryResult.Error; err != nil {
		NewErrorResponse(c, http.StatusNotFound, "age category not found")
	}

	c.JSON(http.StatusOK, gin.H{
		"TrendMovies": trendMovies,
		"NewMovies":   newMovies,
		"Telehikaya":  telehikayaMovies,
		"Horor":       hororMovies,
		"Anime":       animeMovies,
		"Category":    category,
		"AgeCategory": ageCategory,
	})
}

func getMoviesByCategory(categoryName string) ([]models.Movie, error) {
	var movies []models.Movie
	result := initializers.DB.
		Joins("JOIN movie_category ON movies.id = movie_category.movie_id").
		Joins("JOIN categories ON movie_category.category_id = categories.id").
		Where("categories.category_name = ?", categoryName).
		Preload("Categories").
		Preload("Screenshots").
		Preload("Seasons", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Videos")
		}).Limit(limitOfMovie).
		Find(&movies)

	if err := result.Error; err != nil {
		return nil, err
	}

	return movies, nil
}
