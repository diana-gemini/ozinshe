package main

import (
	"fmt"
	"log"
	"os"

	"ozinshe/config"
	"ozinshe/db/initializers"
	"ozinshe/internal/models"

	"golang.org/x/crypto/bcrypt"
)

func init() {
	config.LoadEnvVariables()
	initializers.ConnectDB()
}

func main() {
	err := initializers.DB.Migrator().DropTable(models.User{}, models.Movie{},
		models.Favorite{}, models.Category{}, models.Type{}, models.AgeCategory{},
		models.Screenshot{}, models.Season{}, models.Video{}, "age")
	if err != nil {
		log.Fatal("Table dropping failed")
	}

	err = initializers.DB.AutoMigrate(models.User{}, models.AgeCategory{}, models.Category{}, models.Type{},
		models.Movie{}, models.Favorite{}, models.Screenshot{}, models.Season{}, models.Video{})

	if err != nil {
		log.Fatal("Migration failed")
	}

	CreateAdmin()
	// CreateCategories()
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

func CreateCategories() {
	categories := []models.Category{
		{CategoryName: "Horor"},
		{CategoryName: "Comedy"},
		{CategoryName: "Drama"},
	}

	for _, category := range categories {
		if err := initializers.DB.Create(&category).Error; err != nil {
			log.Fatalf("Failed to insert category: %v", err)
		}
	}
}
