package database

import (
	"fiber_proketo/model"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
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
