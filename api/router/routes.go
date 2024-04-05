package router

import (
	"net/http"

	"ozinshe/api/controllers"
	"ozinshe/api/middleware"
	"github.com/gin-gonic/gin"
)

func GetRoute(r *gin.Engine) {
	r.LoadHTMLGlob("templates/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Main page",
		})
	})
	r.GET("/signup", func(c *gin.Context) {
		c.HTML(http.StatusOK, "signup.html", gin.H{
			"title": "Sign up page",
		})
	})
	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"title": "Log in page",
		})
	})
	r.GET("/logout", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Log in page",
		})
	})
	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.Use(middleware.RequireAuth)
	r.POST("/logout", controllers.Logout)
}
