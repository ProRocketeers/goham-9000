package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Repository struct {
	gorm.Model
	URL    string `json:"url"`
	Name   string `json:"name"`
	Branch string `json:"branch"`
	ImgUrl string `json:"img_url"`
	Status string `json:"status"`
}
