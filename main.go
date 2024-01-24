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
	"github.com/spf13/viper"
)

func main() {
	database.InitDatabase()

	defer database.DBConn.Close()

	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	database.InitDatabase()
	app := fiber.New()
	app.Use(logger.New())
	app.Get("/", helloKek)
	app.Get("/getgit", cloneRepo)
	app.Get("/version", nixVersion)
	app.Post("build", nixBuild)
	app.Get("/uploadToReg", uploadToReg)

	app.Get("/repos", GetRepos)
	app.Get("/repo/:id", GetRepo)
	app.Post("/repo", NewRepo)
	app.Put("/repo/:id", UpdateRepo)
	app.Delete("/repo/:id", DeleteRepo)

	app.Listen(":" + viper.Get("PORT").(string))
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

func nixBuild(ctx *fiber.Ctx) error {
	payload := struct {
		Path string `json:"path"`
	}{}

	if err := ctx.BodyParser(&payload); err != nil {
		return err
	}
	lib.NixpackBuild(payload.Path)
	return ctx.JSON(payload)
}

func GetRepos(c *fiber.Ctx) error {
	db := database.DBConn

	var books []model.Repository
	db.Find(&books)
	err := c.JSON(books)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func GetRepo(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn
	var book model.Repository

	db.Find(&book, id)
	err := c.JSON(book)

	if err != nil {
		fmt.Println(err)
		return err
	}
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

func uploadToReg(ctx *fiber.Ctx) error {
	kekel := lib.Tst("some_super_hash", "/nejaka/cesta")

	log.Debug("kekel")

	return ctx.SendString("Henlo uploaded to register " + kekel)
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

func UpdateRepo(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn

	var repo model.Repository
	result := db.First(&repo, id)
	if result.Error != nil {
		c.Status(404).SendString("No repository found with given ID")
		return result.Error
	}

	var updatedRepo model.Repository
	if err := c.BodyParser(&updatedRepo); err != nil {
		c.Status(503).SendString("Error parsing request")
		return err
	}
	db.Model(&repo).Updates(updatedRepo)

	err := c.JSON(repo)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
