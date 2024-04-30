package controllers

import (
	"net/http"
	"strconv"

	"github.com/diana-gemini/ozinshe/db/initializers"
	"github.com/diana-gemini/ozinshe/internal/helpers"
	"github.com/diana-gemini/ozinshe/internal/models"
	"github.com/diana-gemini/ozinshe/internal/validations"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetAllMovies godoc
// @Summary GetAllMovies
// @Security ApiKeyAuth
// @Tags movie-controller
// @ID get-all-movies
// @Accept  json
// @Produce  json
// @Success 200 {integer} integer 1
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /all [get]
func GetAllMovies(c *gin.Context) {
	var movies []models.Movie
	result := initializers.DB.Preload("Categories").
		Preload("Screenshots").
		Preload("Seasons", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Videos")
		}).Find(&movies)

	if err := result.Error; err != nil {
		NewErrorResponse(c, http.StatusNotFound, "movies not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Movies": movies,
	})
}

// GetMovieByID godoc
// @Summary GetMovieByID
// @Security ApiKeyAuth
// @Tags movie-controller
// @ID get-movie-by-id
// @Accept  json
// @Produce  json
// @Success 200 {integer} integer 1
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /movie/:id [get]
func GetMovieByID(c *gin.Context) {
	authUser := helpers.GetAuthUser(c)
	userID := authUser.ID

	movieID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "cannot convert to int")
	}

	isUserFavorite := validations.IsUniqueTwoValue("favorites", "user_id", "movie_id", userID, uint(movieID))

	var movie models.Movie
	result := initializers.DB.Preload("Categories").
		Preload("Screenshots").
		Preload("Seasons", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Videos")
		}).
		First(&movie, movieID)

	if err := result.Error; err != nil {
		NewErrorResponse(c, http.StatusNotFound, "movie not found")
		return
	}

	var similarSerial []models.Movie
	similarSerialResult := initializers.DB.Preload("Categories").
		Preload("Screenshots").
		Preload("Seasons", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Videos")
		}).Limit(limitOfMovie).
		Find(&similarSerial)

	if err := similarSerialResult.Error; err != nil {
		NewErrorResponse(c, http.StatusNotFound, "similar serial movie not found")
	}

	movie.CountOfWatch++

	if err := initializers.DB.Model(&models.Movie{}).Where("id = ?", movieID).Update("count_of_watch", movie.CountOfWatch).Error; err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"movie":          movie,
		"isUserFavorite": isUserFavorite,
		"similarSerial":  similarSerial,
	})
}

// GetMovieSeriesByID godoc
// @Summary GetMovieSeriesByID
// @Security ApiKeyAuth
// @Tags movie-controller
// @ID get-movie-series-by-id
// @Accept  json
// @Produce  json
// @Param movieid path string true "movieid"
// @Param seasonid path string true "seasonid"
// @Param seriesid path string true "seriesid"
// @Success 200 {integer} integer 1
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /movie/:id/series/:seasonid/:seriesid [get]
func GetMovieSeriesByID(c *gin.Context) {
	movieID := c.Param("id")

	seasonID, err := strconv.Atoi(c.Params.ByName("seasonid"))
	if err != nil {
		NewErrorResponse(c, http.StatusNotFound, "season id not found")
		return
	}

	seriesID, err := strconv.Atoi(c.Params.ByName("seriesid"))
	if err != nil {
		NewErrorResponse(c, http.StatusNotFound, "series id not found")
		return
	}

	var movie models.Movie
	result := initializers.DB.Preload("Categories").
		Preload("Screenshots").
		Preload("Seasons", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Videos")
		}).First(&movie, movieID)

	if err := result.Error; err != nil {
		NewErrorResponse(c, http.StatusNotFound, "movie not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Series": movie.Seasons[seasonID-1].Videos[seriesID-1].Link,
	})
}
