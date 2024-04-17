package controllers

import (
	"net/http"
	"ozinshe/db/initializers"
	"ozinshe/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetTrends(c *gin.Context) {
	var movies []models.Movie
	result := initializers.DB.Preload("Categories").
		Preload("Screenshots").
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
	result := initializers.DB.Preload("Categories").
		Preload("Screenshots").
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
	result := initializers.DB.Where("type_name = ?", "Serial").Find(&types)

	if err := result.Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"Type": "Type not found",
		})
		return
	}

	result = initializers.DB.
		Preload("Categories").
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

func Horor(c *gin.Context) {
	var movies []models.Movie
	result := initializers.DB.
		Joins("JOIN movie_category ON movies.id = movie_category.movie_id").
		Joins("JOIN categories ON movie_category.category_id = categories.id").
		Where("categories.category_name = ?", "Horor").
		Preload("Categories").
		Preload("Screenshots").
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
		"Horor": movies,
	})
}

func Anime(c *gin.Context) {
	var movies []models.Movie
	result := initializers.DB.
		Joins("JOIN movie_category ON movies.id = movie_category.movie_id").
		Joins("JOIN categories ON movie_category.category_id = categories.id").
		Where("categories.category_name = ?", "Anime").
		Preload("Categories").
		Preload("Screenshots").
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
		"Anime": movies,
	})

}
