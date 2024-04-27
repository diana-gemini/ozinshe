package controllers

import (
	"net/http"
	"github.com/diana-gemini/ozinshe/db/initializers"
	"github.com/diana-gemini/ozinshe/internal/helpers"
	"github.com/diana-gemini/ozinshe/internal/models"
	"github.com/diana-gemini/ozinshe/internal/validations"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAllMovies(c *gin.Context) {
	var movies []models.Movie
	result := initializers.DB.Preload("Categories").
		Preload("Screenshots").
		Preload("Seasons", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Videos")
		}).Find(&movies)

	if err := result.Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"movie": "Record not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Movies": movies,
	})
}

func GetMovieByID(c *gin.Context) {
	authUser := helpers.GetAuthUser(c)
	userID := authUser.ID

	movieID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}

	isUserFavorite := validations.IsUniqueTwoValue("favorites", "user_id", "movie_id", userID, uint(movieID))

	var movie models.Movie
	result := initializers.DB.Preload("Categories").
		//Preload("Screenshots").
		Preload("Seasons", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Videos")
		}).
		First(&movie, movieID)

	if err := result.Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"movie": "Record not found",
		})
		return
	}

	var similarSerial []models.Movie
	similarSerialResult := initializers.DB.Preload("Categories").
		// Preload("Screenshots").
		Preload("Seasons", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Videos")
		}).Limit(limitOfMovie).
		Find(&similarSerial)

	if err := similarSerialResult.Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"movie": "Similar serial movies not found",
		})
	}

	movie.CountOfWatch++

	if err := initializers.DB.Model(&models.Movie{}).Where("id = ?", movieID).Update("count_of_watch", movie.CountOfWatch).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"movie":          movie,
		"isUserFavorite": isUserFavorite,
		"similarSerial":  similarSerial,
	})
}

func GetMovieSeriesByID(c *gin.Context) {
	movieID := c.Param("id")
	seasonID, err := strconv.Atoi(c.Params.ByName("seasonid"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err,
		})
		return
	}
	seriesID, err := strconv.Atoi(c.Params.ByName("seriesid"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err,
		})
		return
	}

	var movie models.Movie
	result := initializers.DB.Preload("Categories").
		//	Preload("Screenshots").
		Preload("Seasons", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Videos")
		}).First(&movie, movieID)

	if err := result.Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"movie": "Record not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Series": movie.Seasons[seasonID-1].Videos[seriesID-1].Link,
	})
}
