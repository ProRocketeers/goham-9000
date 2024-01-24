package database

import (
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"goham-9000/model"
)

var (
	DBConn *gorm.DB
)

func InitDatabase() {
	var err error
	DBConn, err = gorm.Open("sqlite3", "repos.db")
	if err != nil {
		panic("Failed to connect to database")
	}
	fmt.Println("DB Connected")

	DBConn.AutoMigrate(&model.Repository{})
	fmt.Println("DB Migrated")

}

func GetProjectById(projectId string) (model.Repository, error) {
	var repo model.Repository
	log.Debug("Getting project by id")
	DBConn.Find(&repo, projectId)

	return repo, nil
}
