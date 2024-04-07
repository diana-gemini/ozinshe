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
}
