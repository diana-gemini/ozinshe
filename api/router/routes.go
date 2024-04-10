package router

import (
	"ozinshe/api/controllers"
	"ozinshe/api/middleware"

	"github.com/gin-gonic/gin"
)

func GetRoute(r *gin.Engine) {
	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.Use(middleware.RequireAuth)
	r.POST("/logout", controllers.Logout)
	r.GET("/all", controllers.GetAllMovies)
	r.GET("/movie/:id", controllers.GetMovieByID)
	r.GET("/movie/:id/series/:seasonid/:seriesid", controllers.GetMovieSeriesByID)
	r.GET("/trends", controllers.GetTrends)
	r.GET("/newprojects", controllers.GetNewprojects)
	r.GET("/telehikaya", controllers.GetTelehikaya)
	r.GET("/search", controllers.Search)

	adminGroup := r.Group("/movie")
	adminGroup.Use(middleware.IsAdmin())
	{
		adminGroup.POST("/create", controllers.CreateMovie)
		adminGroup.GET("/:id/edit", controllers.EditMovie)
		adminGroup.PUT("/:id/update", controllers.UpdateMovie)
		adminGroup.DELETE("/:id/delete", controllers.DeleteMovie)
	}

}
