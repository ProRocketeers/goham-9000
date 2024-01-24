package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Repository struct {
	gorm.Model
	GitUrl       string `json:"GitUrl"`
	Name         string `json:"Name"`
	Branch       string `json:"Branch"`
	DockerImgUrl string `json:"DockerImgUrl"`
	Status       string `json:"Status"`
	LocalImgPath string `json:"LocalImgPath"`
}
