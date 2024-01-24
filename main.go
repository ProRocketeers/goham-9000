package main

import (
	"fiber_proketo/database"
	"fiber_proketo/lib"
	"fiber_proketo/model"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func main() {
	database.InitDatabase()
	app := fiber.New()
	app.Use(logger.New())
	app.Get("/", helloKek)
	app.Get("/getgit", cloneRepo)
	app.Get("/pack", nixVersion)

	app.Get("/repos", GetRepos)
	app.Get("/repo/:id", GetRepo)
	app.Post("/repo", NewRepo)
	app.Delete("/repo/:id", DeleteRepo)

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

func GetRepos(c *fiber.Ctx) error {
	db := database.DBConn

	var books []model.Repository
	db.Find(&books)
	c.JSON(books)
	return nil
}

func GetRepo(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn
	var book model.Repository

	db.Find(&book, id)
	c.JSON(book)

	return nil
}

func NewRepo(c *fiber.Ctx) error {
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

func DeleteRepo(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn
	var repo model.Repository

	db.First(&repo, id)

	if repo.URL == "" {
		c.Status(500).SendString("No book found with given ID")
		return nil
	}

	db.Delete(&repo)
	c.SendString("Book successfully deleted!")

	return nil
}
