package router

import (
	"github.com/diana-gemini/ozinshe/api/controllers"
	"github.com/diana-gemini/ozinshe/api/middleware"

	"github.com/gin-gonic/gin"
)

func GetRoute(r *gin.Engine) {
	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.POST("/passwordrecover", controllers.PasswordRecover)
	r.POST("/resetpassword", controllers.ResetPassword)

	r.Use(middleware.RequireAuth)
	r.POST("/logout", controllers.Logout)
	r.GET("/editprofile", controllers.EditUserProfile)
	r.PUT("/updateprofile", controllers.UpdateUserProfile)
	r.POST("/changepassword", controllers.ChangePassword)
	r.GET("/home", controllers.Home)
	r.GET("/trends", controllers.GetTrends)
	r.GET("/newprojects", controllers.GetNewprojects)
	r.GET("/telehikaya", controllers.GetTelehikaya)
	r.GET("/horor", controllers.Horor)
	r.GET("/anime", controllers.Anime)
	r.GET("/search", controllers.Search)
	r.GET("/all", controllers.GetAllMovies)

	r.GET("/movie/:id", controllers.GetMovieByID)
	r.GET("/movie/:id/series/:seasonid/:seriesid", controllers.GetMovieSeriesByID)

	r.POST("/movie/:id/favorite", controllers.AddMovieToFavorite)
	r.DELETE("/movie/:id/favorite", controllers.DeleteMovieFromFavorite)
	r.GET("/movie/favorite", controllers.GetAllFavoriteMovies)

	admin := r.Group("/admin")
	admin.Use(middleware.IsAdmin())
	{
		admin.POST("/movie/create", controllers.CreateMovie)
		admin.GET("/movie/:id/edit", controllers.EditMovie)
		admin.PUT("/movie/:id/update", controllers.UpdateMovie)
		admin.DELETE("/movie/:id/delete", controllers.DeleteMovie)

		admin.POST("/movie/:id/season/create", controllers.CreateSeason)
		admin.PUT("/movie/:id/season/:seasonid/update", controllers.UpdateSeason)
		admin.DELETE("/movie/:id/season/:seasonid/delete", controllers.DeleteSeason)

		admin.POST("/category/create", controllers.CreateCategory)
		admin.PUT("/category/:id/update", controllers.UpdateCategory)
		admin.DELETE("/category/:id/delete", controllers.DeleteCategory)

		admin.POST("/type/create", controllers.CreateTypeOfProject)
		admin.PUT("/type/:id/update", controllers.UpdateTypeOfProject)
		admin.DELETE("/type/:id/delete", controllers.DeleteTypeOfProject)

		admin.POST("/agecategory/create", controllers.CreateAgeCategory)
		admin.PUT("/agecategory/:id/update", controllers.UpdateAgeCategory)
		admin.DELETE("/agecategory/:id/delete", controllers.DeleteAgeCategory)
	}
}
