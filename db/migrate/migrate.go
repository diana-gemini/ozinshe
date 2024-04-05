package main

import (
	"log"

	"ozinshe/config"
	"ozinshe/db/initializers"
	"ozinshe/internal/models"
)

func init() {
	config.LoadEnvVariables()
	initializers.ConnectDB()
}

func main() {
	err := initializers.DB.Migrator().DropTable(models.User{})
	if err != nil {
		log.Fatal("Table dropping failed")
	}

	err = initializers.DB.AutoMigrate(models.User{})

	if err != nil {
		log.Fatal("Migration failed")
	}
}
