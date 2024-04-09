package models

import (
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Movie struct {
	gorm.Model
	NameOfProject string                      `gorm:"unique;not null" json:"nameOfProject"`
	Category      string                      `gorm:"not null" json:"category"`
	TypeOfProject string                      `gorm:"not null" json:"typeOfProject"`
	AgeCategory   string                      `gorm:"not null" json:"ageCategory"`
	Year          string                      `gorm:"not null" json:"year"`
	Timing        string                      `gorm:"not null" json:"timing"`
	Keywords      string                      `gorm:"not null" json:"keywords"`
	Description   string                      `gorm:"not null" json:"description"`
	Director      string                      `gorm:"not null" json:"director"`
	Producer      string                      `gorm:"not null" json:"producer"`
	CountOfSeason datatypes.JSONSlice[Season] `gorm:"not null" json:"countOfSeason"`
	Cover         string                      `gorm:"not null" json:"cover"`
	Screenshots   pq.StringArray              `gorm:"type:text[];not null" json:"screenshots"`
	CountOfWatch  int                         `gorm:"primaryKey" json:"countOfWatch"`
}

type Season struct {
	Season       string         `gorm:"primaryKey" json:"season"`
	LinkOfSeries pq.StringArray `gorm:"primaryKey;type:text[];not null" json:"linkOfSeries"`
}
