package controllers

import (
	"fmt"
	"net/http"
	"ozinshe/db/initializers"
	"ozinshe/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

var limitOfProjectToDisplay = 5

func GetAllMovies(c *gin.Context) {
	var movies []models.Movie
	result := initializers.DB.Find(&movies)

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
	result := initializers.DB.First(&movie, id)

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
	result := initializers.DB.First(&movie, movieID)

	if err := result.Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"movie": "Record not found",
		})
		return
	}

	seriesLink := fmt.Sprintf("https://www.youtube.com/watch?v=%v", movie.CountOfSeason[seasonID-1].LinkOfSeries[seriesID-1])

	c.JSON(http.StatusOK, gin.H{
		"Series": seriesLink,
	})
}

func GetTrends(c *gin.Context) {
	var movies []models.Movie
	result := initializers.DB.Order("count_of_watch desc").Limit(limitOfProjectToDisplay).Find(&movies)

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

func GetAllTrends(c *gin.Context) {
	var movies []models.Movie
	result := initializers.DB.Order("count_of_watch desc").Find(&movies)

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
	result := initializers.DB.Order("created_at desc").Limit(limitOfProjectToDisplay).Find(&movies)

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

func GetAllNewprojects(c *gin.Context) {
	var movies []models.Movie
	result := initializers.DB.Order("created_at desc").Find(&movies)

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
	result := initializers.DB.Where("type_of_project = ?", "Serial").
		Order("created_at desc").
		Limit(limitOfProjectToDisplay).
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

func GetAllTelehikaya(c *gin.Context) {
	var movies []models.Movie
	result := initializers.DB.Where("type_of_project = ?", "Serial").
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
