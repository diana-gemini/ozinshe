package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/diana-gemini/ozinshe/db/initializers"
	"github.com/diana-gemini/ozinshe/internal/models"
	"github.com/diana-gemini/ozinshe/internal/validations"

	"github.com/gin-gonic/gin"
)

type NewSeason struct {
	Videos []string `form:"videos" binding:"required" example:"video link"`
}

// CreateSeason godoc
// @Summary CreateSeason
// @Security ApiKeyAuth
// @Tags admin-movie-season-controller
// @ID create-season
// @Accept  multipart/form-data
// @Produce  json
// @Param id path integer true "movieID"
// @Param newSeason formData NewSeason true "newSeason"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /admin/movie/{id}/season/create [post]
func CreateSeason(c *gin.Context) {
	movieID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "cannot convert to int")
		return
	}

	if !validations.IsExistValue("movies", "id", movieID) {
		NewErrorResponse(c, http.StatusBadRequest, "cannot find movie")
		return
	}

	var newSeason NewSeason
	if err := c.ShouldBind(&newSeason); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	videosArray := strings.Split(newSeason.Videos[0], ",")

	var videos []models.Video

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
		NewErrorResponse(c, http.StatusInternalServerError, "cannot create season")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"season": season,
	})
}

// EditSeason godoc
// @Summary EditSeason
// @Security ApiKeyAuth
// @Tags admin-movie-season-controller
// @ID edit-season
// @Accept  json
// @Produce  json
// @Param id path integer true "id"
// @Param seasonid path integer true "seasonid"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /admin/movie/{id}/season/{seasonid}/edit [get]
func EditSeason(c *gin.Context) {
	movieID := c.Param("id")
	seasonID := c.Param("seasonid")

	var season models.Season
	result := initializers.DB.Preload("Videos").
		Where("movie_id = ? and id = ?", movieID, seasonID).First(&season)

	if err := result.Error; err != nil {
		NewErrorResponse(c, http.StatusNotFound, "season not found")
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
