package controllers

import (
	"fmt"
	"net/http"
	"ozinshe/db/initializers"
	"ozinshe/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAllMovies(c *gin.Context) {
	var movies []models.Movie
	result := initializers.DB.Preload("Screenshots").
		Preload("AgeCategories").
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
	id := c.Param("id")

	var movie models.Movie
	result := initializers.DB.Preload("Screenshots").
		Preload("AgeCategories").
		Preload("Seasons", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Videos")
		}).First(&movie, id)

	if err := result.Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"movie": "Record not found",
		})
		return
	}

	movie.CountOfWatch++

	if err := initializers.DB.Model(&models.Movie{}).Where("id = ?", id).Update("count_of_watch", movie.CountOfWatch).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"movie": movie,
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
	result := initializers.DB.Preload("Screenshots").
		Preload("AgeCategories").
		Preload("Seasons", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Videos")
		}).First(&movie, movieID)

	if err := result.Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"movie": "Record not found",
		})
		return
	}

	seriesLink := fmt.Sprintf("https://www.youtube.com/watch?v=%v", movie.Seasons[seasonID-1].Videos[seriesID-1].Link)

	c.JSON(http.StatusOK, gin.H{
		"Series": seriesLink,
	})
}

func GetTrends(c *gin.Context) {
	var movies []models.Movie
	result := initializers.DB.Preload("Screenshots").
		Preload("AgeCategories").
		Preload("Seasons", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Videos")
		}).Order("count_of_watch desc").Find(&movies)

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

func GetNewprojects(c *gin.Context) {
	var movies []models.Movie
	result := initializers.DB.Preload("Screenshots").
		Preload("AgeCategories").
		Preload("Seasons", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Videos")
		}).Order("created_at desc").Find(&movies)

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

func GetTelehikaya(c *gin.Context) {
	var movies []models.Movie

	var types models.Type
	result := initializers.DB.Preload("Types").Where("type_name = ?", "Serial").Find(&types)

	if err := result.Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"Type": "Type not found",
		})
		return
	}

	result = initializers.DB.
		Preload("Screenshots").
		Preload("AgeCategories").
		Preload("Seasons", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Videos")
		}).Where("type_id = ?", types.ID).
		Order("created_at desc").
		Find(&movies)

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
