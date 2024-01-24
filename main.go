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
	"github.com/spf13/viper"
)

func main() {
	initDatabase()

	defer database.DBConn.Close()

	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	app := fiber.New()
	app.Use(logger.New())
	app.Get("/", helloKek)
	app.Get("/getgit", cloneRepo)
	app.Get("/pack", nixVersion)
	app.Post("build", nixBuild)
	app.Post("/newRepo", newRepo)
	app.Get("/uploadToReg", uploadToReg)
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

func uploadToReg(ctx *fiber.Ctx) error {
	kekel := lib.Tst("some_super_hash", "/nejaka/cesta")

	log.Debug("kekel")

	return ctx.SendString("Henlo uploaded to register " + kekel)
}
