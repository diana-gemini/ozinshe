package models

import (
	"time"
)

type User struct {
	// gorm.Model
	ID          uint   `gorm:"unique;primary key" json:"userID"`
	Email       string `json:"email" gorm:"unique;not null"`
	Password    string `json:"-"`
	RoleID      uint   `json:"roleID"`
	Username    string `json:"username"`
	MobilePhone string `json:"mobilephone"`
	BirthDate   string `json:"birthdate"`
	Favorites   []Favorite
}

type Favorite struct {
	ID        uint `gorm:"primaryKey"`
	MovieID   int  `gorm:"foreignkey:MovieID" json:"movieID"`
	UserID    uint `gorm:"foreignkey:UserID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	User      User
}
