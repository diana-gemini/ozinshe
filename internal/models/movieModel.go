package models

import (
	"gorm.io/gorm"
)

type Movie struct {
	gorm.Model
	NameOfProject string
	CategoryID    uint
	TypeID        uint
	AgeCategories []AgeCategory `gorm:"many2many:age;"`
	Screenshots   []Screenshot
	Seasons       []Season
	Year          string `gorm:"not null" json:"year"`
	Timing        string `gorm:"not null" json:"timing"`
	Keywords      string `gorm:"not null" json:"keywords"`
	Description   string `gorm:"not null" json:"description"`
	Director      string `gorm:"not null" json:"director"`
	Producer      string `gorm:"not null" json:"producer"`
	Cover         string `gorm:"not null" json:"cover"`
	CountOfWatch  int    `json:"countOfWatch"`
}

type Category struct {
	gorm.Model
	CategoryName string
	Movies       []Movie
}

type Type struct {
	gorm.Model
	TypeName string
	Movies   []Movie
}

type AgeCategory struct {
	gorm.Model
	AgeCategoryName string
	Movies          []Movie `gorm:"many2many:age;"`
}

type Screenshot struct {
	gorm.Model
	Link    string
	MovieID uint
}

type Season struct {
	gorm.Model
	Videos  []Video
	MovieID uint
}

type Video struct {
	gorm.Model
	Link     string
	SeasonID uint
}
