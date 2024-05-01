package controllers

import (
	"net/http"
	"strconv"

	"github.com/diana-gemini/ozinshe/db/initializers"
	"github.com/diana-gemini/ozinshe/internal/helpers"
	"github.com/diana-gemini/ozinshe/internal/models"
	"github.com/diana-gemini/ozinshe/internal/validations"

	"github.com/gin-gonic/gin"
)

type FavoriteResponse struct {
	UserID  uint
	MovieID uint
}

// AddMovieToFavorite godoc
// @Summary AddMovieToFavorite
// @Security ApiKeyAuth
// @Tags movie-controller
// @ID add-movie-to-favorite
// @Accept  json
// @Produce  json
// @Param id path integer true "movieID"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /movie/{id}/favorite [post]
func AddMovieToFavorite(c *gin.Context) {
	authUser := helpers.GetAuthUser(c)

	movieID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusNotFound, "cannot convert movieID to int ")
		return
	}

	if !validations.IsExistValue("movies", "id", movieID) {
		NewErrorResponse(c, http.StatusBadRequest, "movie id not found")
		return
	}

	if validations.IsUniqueTwoValue("favorites", "user_id", "movie_id", authUser.ID, uint(movieID)) {
		NewErrorResponse(c, http.StatusBadRequest, "this favorite movie is already exist")
		return
	}

	favoriteMovie := models.Favorite{
		UserID:  authUser.ID,
		MovieID: uint(movieID),
	}

	result := initializers.DB.Create(&favoriteMovie)
	if result.Error != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "cannot create favorite movie")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"UserID":        favoriteMovie.UserID,
		"FavoriteMovie": favoriteMovie.MovieID,
	})
}

// DeleteMovieFromFavorite godoc
// @Summary DeleteMovieFromFavorite
// @Security ApiKeyAuth
// @Tags movie-controller
// @ID delete-movie-from-favorite
// @Accept  json
// @Produce  json
// @Param id path integer true "movieID"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /movie/{id}/favorite [delete]
func DeleteMovieFromFavorite(c *gin.Context) {
	authUser := helpers.GetAuthUser(c)
	movieID := c.Param("id")

	var favorite models.Favorite

	result := initializers.DB.Where("movie_id = ? AND user_id = ?", movieID, authUser.ID).First(&favorite)
	if err := result.Error; err != nil {
		NewErrorResponse(c, http.StatusNotFound, "user favorite movie is not found")
		return
	}

	initializers.DB.Unscoped().Delete(&favorite)

	c.JSON(http.StatusOK, gin.H{
		"message": "favorite movie delete successfully",
	})
}

// GetAllFavoriteMovies godoc
// @Summary GetAllFavoriteMovies
// @Security ApiKeyAuth
// @Tags movie-controller
// @ID get-all-favorite-movies
// @Accept  json
// @Produce  json
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /movie/favorite [get]
func GetAllFavoriteMovies(c *gin.Context) {
	user := helpers.GetAuthUser(c)
	var favoriteMovies []models.Favorite
	result := initializers.DB.Where("user_id = ?", user.ID).Find(&favoriteMovies)

	if err := result.Error; err != nil {
		NewErrorResponse(c, http.StatusNotFound, "favorite movie is not found")
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
