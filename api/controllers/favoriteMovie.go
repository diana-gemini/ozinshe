package controllers

import (
	"net/http"
	"ozinshe/db/initializers"
	"ozinshe/internal/helpers"
	"ozinshe/internal/models"
	"ozinshe/internal/validations"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FavoriteResponse struct {
	UserID  uint
	MovieID int
}

func AddMovieToFavorite(c *gin.Context) {
	authUser := helpers.GetAuthUser(c)
	movieIDFromParam := c.Param("id")

	movieID, err := strconv.Atoi(movieIDFromParam)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err,
		})
		return
	}

	if validations.IsUniqueTwoValue("favorites", "user_id", "movie_id", authUser.ID, movieID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"UniqueValue": "This favorite movie is already exist!",
		})
		return
	}

	favoriteMovie := models.Favorite{
		UserID:  authUser.ID,
		MovieID: movieID,
	}

	result := initializers.DB.Create(&favoriteMovie)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Cannot create favorite movie",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"UserID":        favoriteMovie.UserID,
		"FavoriteMovie": favoriteMovie.MovieID,
	})
}

func DeleteMovieFromFavorite(c *gin.Context) {
	authUser := helpers.GetAuthUser(c)
	movieID := c.Param("id")

	var favorite models.Favorite

	result := initializers.DB.Where("movie_id = ? AND user_id = ?", movieID, authUser.ID).First(&favorite)
	if err := result.Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err,
		})
		return
	}

	initializers.DB.Delete(&favorite)

	c.JSON(http.StatusOK, gin.H{
		"message": "The favorite movie has been deleted successfully",
	})
}

func GetAllFavoriteMovies(c *gin.Context) {
	user := helpers.GetAuthUser(c)
	var favoriteMovies []models.Favorite
	result := initializers.DB.Where("user_id = ?", user.ID).Find(&favoriteMovies)

	if err := result.Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"movie": "Record not found",
		})
		return
	}

	var response []FavoriteResponse

	for _, favorite := range favoriteMovies {
		response = append(response, FavoriteResponse{
			UserID:  favorite.UserID,
			MovieID: favorite.MovieID,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"favoriteMovies": response,
	})
}
