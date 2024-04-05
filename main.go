package main

import (
	"fmt"
	"ozinshe/api/router"
	"ozinshe/config"
	"ozinshe/db/initializers"

	"github.com/gin-gonic/gin"
)

func init() {
	config.LoadEnvVariables()
	initializers.ConnectDB()
}

func main() {
	fmt.Println("Hello auth")
	r := gin.Default()
	router.GetRoute(r)

	r.Run()
}
