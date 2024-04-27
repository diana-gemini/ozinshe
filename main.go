package main

import (
	"github.com/diana-gemini/ozinshe/api/router"
	"github.com/diana-gemini/ozinshe/config"
	"github.com/diana-gemini/ozinshe/db/initializers"

	_ "github.com/diana-gemini/ozinshe/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Ozinwe API
// @version 1.0
// @description API Server for Ozinwe

// @host localhost:3000
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func init() {
	config.LoadEnvVariables()
	initializers.ConnectDB()
}

func main() {
	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GetRoute(r)
	r.Run()
}
