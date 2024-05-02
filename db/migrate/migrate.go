package main

import (
	"fmt"
	"log"
	"os"

	"github.com/diana-gemini/ozinshe/config"
	"github.com/diana-gemini/ozinshe/db/initializers"
	"github.com/diana-gemini/ozinshe/internal/models"

	"golang.org/x/crypto/bcrypt"
)

func init() {
	config.LoadEnvVariables()
	initializers.ConnectDB()
}

func main() {
	//follow code delete table from DB

	// err := initializers.DB.Migrator().DropTable(models.User{}, models.Movie{},
	// 	models.Favorite{}, models.Category{}, models.Season{}, models.Type{}, models.AgeCategory{},
	// 	models.Video{}, models.Screenshot{}, "movie_category")

	// if err != nil {
	// 	log.Fatal("Table dropping failed")
	// }

	err := initializers.DB.AutoMigrate(models.User{}, models.AgeCategory{}, models.Category{}, models.Season{},
		models.Type{}, models.Movie{}, models.Screenshot{}, models.Favorite{}, models.Video{})

	if err != nil {
		log.Fatal("Migration failed")
	}

	CreateAdmin()
}

func CreateAdmin() {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(os.Getenv("ADMIN_PASSWORD")), 10)
	if err != nil {
		fmt.Println("Failed to hash admin assword")
		return
	}
	admin := models.User{
		Email:    os.Getenv("ADMIN_EMAIL"),
		Password: string(hashPassword),
		RoleID:   1,
	}

	result := initializers.DB.Create(&admin)

	if result.Error != nil {
		fmt.Println("Internal server error")
		return
	}
}
