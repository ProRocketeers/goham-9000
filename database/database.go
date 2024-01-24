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

const (
	P_CREATED      = "CREATED"
	P_CLONED       = "CLONED"
	P_IMG_BUILT    = "IMG_BUILT"
	P_IMG_UPLOADED = "IMG_UPLOADED"
	P_DEPLOYED     = "DEPLOYED"
	P_CLEARED      = "CLEARED"
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
	log.Debug("Getting project by id")

	var repo model.Repository
	DBConn.Find(&repo, projectId)

	return repo, nil
}
func UpdateProjectStatus(projectId string, status string) (model.Repository, error) {
	log.Debug("Updating" + projectId + " status, to " + status)
	var repo model.Repository

	DBConn.Find(&repo, projectId).Update(&model.Repository{Status: status})

	return repo, nil
}
