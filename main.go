package main

import (
	"fiber_proketo/database"
	"fiber_proketo/lib"
	"fiber_proketo/model"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func main() {
	initDatabase()

	defer database.DBConn.Close()
	app := fiber.New()
	app.Use(logger.New())
	app.Get("/", helloKek)
	app.Get("/getgit", cloneRepo)
	app.Get("/pack", nixVersion)
	app.Post("/newRepo", newRepo)
	app.Listen(":3000")
}

func helloKek(ctx *fiber.Ctx) error {
	log.Info("Hello Kek!")
	return ctx.SendString("Hello Kek!")
}
func cloneRepo(ctx *fiber.Ctx) error {
	path, err := lib.CloneRepository("https://github.com/Fenny/fiber-hello-world", "cool_bro")
	if err != nil {
		return err
	}
	return ctx.SendString("Repository cloned at " + path)

}
func nixVersion(ctx *fiber.Ctx) error {
	return ctx.SendString(lib.NixpackVersion())

}

func initDatabase() {
	var err error
	database.DBConn, err = gorm.Open("sqlite3", "repos.db")
	if err != nil {
		panic("Failed to connect to database")
	}
	fmt.Println("DB Connected")

	database.DBConn.AutoMigrate(&model.Repository{})
	fmt.Println("DB Migrated")

}

func newRepo(c *fiber.Ctx) error {
	db := database.DBConn
	repo := new(model.Repository)
	if err := c.BodyParser(repo); err != nil {
		c.Status(503).SendString("Error creating repo")
		return err
	}

	db.Create(&repo)
	err := c.JSON(repo)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
