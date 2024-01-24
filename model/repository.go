package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Repository struct {
	gorm.Model
	URL    string `json:"Url"`
	Name   string `json:"Name"`
	Branch string `json:"Branch"`
	ImgUrl string `json:"ImgUrl"`
	Status string `json:"Status"`
}
