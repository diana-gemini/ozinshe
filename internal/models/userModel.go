package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email       string     `json:"email" gorm:"unique;not null"`
	Password    string     `json:"-"`
	RoleID      uint       `json:"roleID"`
	Username    string     `json:"username"`
	MobilePhone string     `json:"mobilephone"`
	BirthDate   string     `json:"birthdate"`
	Favorites   []Favorite `json:"favorites"`
}

type Favorite struct {
	gorm.Model
	MovieID uint
	UserID  uint
}
