package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `json:"email" gorm:"unique;not null"`
	Password string `json:"-"`
}

type WebPage struct {
	IsLoggedin bool
}

type ErrText struct {
	Email string
	Pass1 string
	Pass2 string
}