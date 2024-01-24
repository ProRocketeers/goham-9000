package main

import (
	"fmt"
	"goham-9000/database"
	"goham-9000/lib"
	"goham-9000/model"

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
	app.Get("/", root)
	app.Get("/version", nixVersion)
	app.Post("build", nixBuild)
	app.Get("/uploadToReg", uploadToReg)
	app.Post("/clone/:id", cloneRepo)
	app.Post("/deploy", deploy)

	app.Get("/repos", GetRepos)
	app.Get("/repo/:id", GetRepo)
	app.Post("/repo", NewRepo)
	app.Put("/repo/:id", UpdateRepo)
	app.Delete("/repo/:id", DeleteRepo)

	app.Listen(":" + viper.Get("PORT").(string))
}

func root(ctx *fiber.Ctx) error {
	return ctx.SendString("Hello from goham 9000")
}
func cloneRepo(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	_, err := database.GetProjectById(id)
	log.Debug("found project by id ", id)
	projectPath, err := lib.CloneRepoStep(id)
	if err != nil {
		return err
	}
	return ctx.SendString("Repository cloned at " + projectPath)

}
func nixVersion(ctx *fiber.Ctx) error {
	return ctx.SendString(lib.NixpackVersion())
}

func nixBuild(c *fiber.Ctx) error {
	payload := struct {
		Id string `json:"id"`
	}{}

	if err := c.BodyParser(&payload); err != nil {
		return err
	}

	var response = lib.NixpackBuild(payload.Id)
	return c.JSON(response)
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
	projectById, err := database.GetProjectById(id)
	if err != nil {
		return err
	}
	err = c.JSON(projectById)
	return err
}

func NewRepo(c *fiber.Ctx) error {
	db := database.DBConn
	repo := new(model.Repository)
	if err := c.BodyParser(repo); err != nil {
		c.Status(503).SendString("Error creating repo")
		return err
	}

	repo.Status = database.P_CREATED

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

func isValidStatus(status string) bool {
	switch status {
	case "CREATED", "BUILT", "CLONED", "UPLOADING", "DEPLOYED", "CLEARED":
		return true
	default:
		return false
	}
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

func deploy(ctx *fiber.Ctx) error {
	const ArgoGitRepostiroy = "git@github.com:ProRocketeers/goham-argo-repo.git" // todo: take me from env/app_cofnig
	const ArgoRepositoryFolderName = "argo-repo"
	payload := struct {
		Id string `json:"path"`
	}{}

	fmt.Println("Lets clone argo repo")
	// clone posledni verze repa
	path, err := lib.CloneRepository(ArgoGitRepostiroy, ArgoRepositoryFolderName)
	if err != nil {
		return err
	}

	fmt.Println("Lets edit yaml")
	// Edit deploy files
	status, err := lib.DeployEditor(ArgoRepositoryFolderName, payload.Id)
	if err != nil {
		return err
	}

	fmt.Println("Lets Commit and Push")
	err = lib.CommitAndPush(path)
	if err != nil {
		return err
	}

	return ctx.SendString(status)
}
