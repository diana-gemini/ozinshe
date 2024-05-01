package controllers

import (
	"net/http"
	"github.com/diana-gemini/ozinshe/db/initializers"
	"github.com/diana-gemini/ozinshe/internal/models"
	"github.com/diana-gemini/ozinshe/internal/validations"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateSeason(c *gin.Context) {
	movieID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	var videos []models.Video
	videosArray := c.PostFormArray("videos")

	for _, video := range videosArray {
		tempVideo := models.Video{
			Link: video,
		}
		videos = append(videos, tempVideo)
	}

	season := models.Season{
		Videos:  videos,
		MovieID: uint(movieID),
	}

	result := initializers.DB.Create(&season)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Cannot create season",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"season": season,
	})
}

func UpdateSeason(c *gin.Context) {
	movieID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	seasonID, err := strconv.Atoi(c.Param("seasonid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	if !validations.IsUniqueTwoValue("seasons", "id", "movie_id", uint(seasonID), uint(movieID)) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"Season": "The season of movie does not exist!",
		})
		return
	}

	var season models.Season
	result := initializers.DB.Where("movie_id = ?", movieID).First(&season, seasonID)
	if err := result.Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err,
		})
		return
	}
	var videos []models.Video
	videosArray := c.PostFormArray("videos")

	for _, video := range videosArray {
		tempVideo := models.Video{
			Link: video,
		}
		videos = append(videos, tempVideo)
	}

	updateSeason := models.Season{
		Videos:  videos,
		MovieID: uint(movieID),
	}

	// if err := initializers.DB.Model(&season).Association("Videos").Clear(); err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"error": err.Error(),
	// 	})
	// 	return
	// }
	
	if err := initializers.DB.Model(&season).Association("Videos").Replace(updateSeason.Videos); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	result = initializers.DB.Model(&season).Updates(&updateSeason)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"season": updateSeason,
	})
}

func DeleteSeason(c *gin.Context) {
	movieID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	seasonID, err := strconv.Atoi(c.Param("seasonid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	var season models.Season
	result := initializers.DB.Where("movie_id = ?", movieID).First(&season, seasonID)
	if err := result.Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err,
		})
		return
	}

	err = initializers.DB.Unscoped().Delete(&season).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err,
		})
		return
	}

	var videos []models.Video
	result = initializers.DB.Where("season_id = ?", seasonID).Find(&videos)
	if err := result.Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err,
		})
		return
	}

	err = initializers.DB.Unscoped().Delete(&videos).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "season and video delete successfully",
	})
}
