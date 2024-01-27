package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2/middleware/cors"
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
	// Initialize default config
	app.Use(cors.New())

	// Or extend your config for customization
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3001",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Use(logger.New())
	app.Get("/", root)
	app.Get("/version", nixVersion)
	app.Post("/build/:id", nixBuild)
	app.Post("/registry", uploadToDockerRegistry)
	app.Post("/clone/:id", cloneRepo)
	app.Post("/deploy/:id", deploy)
	app.Post("/repository", InitNewGitRepository)

	app.Get("/repository", GetGitRepositories)         // todo: not used
	app.Get("/repository/:id", GetGitRepository)       // todo: not used
	app.Put("/repository/:id", UpdateGitRepository)    // todo: not used
	app.Delete("/repository/:id", DeleteGitRepository) // todo: not used

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

func nixBuild(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	response, err := lib.NixpackBuildStep(id)
	if err != nil {
		return err
	}
	return ctx.JSON(response)
}

func GetGitRepositories(c *fiber.Ctx) error {
	repos, err := database.FetchAllRepos()
	if err != nil {
		c.Status(fiber.StatusInternalServerError).SendString("Error fetching repositories")
		return err
	}
	return c.JSON(repos)
}

func GetGitRepository(c *fiber.Ctx) error {
	id := c.Params("id")
	projectById, err := database.GetProjectById(id)
	if err != nil {
		return err
	}
	err = c.JSON(projectById)
	return err
}

func InitNewGitRepository(c *fiber.Ctx) error {
	repo, err := database.ParseRepoFromBody(c)
	if err != nil {
		return err
	}
	err = database.CreateRepo(repo)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).SendString("Error saving the repository")
		return err
	}
	return c.JSON(repo)
}

func uploadToDockerRegistry(ctx *fiber.Ctx) error {
	kekel := lib.Tst("some_super_hash", "/nejaka/cesta")

	log.Debug("kekel")

	return ctx.SendString("Henlo uploaded to register " + kekel)
}

func DeleteGitRepository(c *fiber.Ctx) error {
	id := c.Params("id")
	projectById, err := database.GetProjectById(id)
	if err != nil {
		return err
	}
	return database.DeleteRepoById(&projectById)
}

func UpdateGitRepository(c *fiber.Ctx) error {
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
	var ArgoGitRepostiroy = viper.Get("ARGO_GIT_REPOSITORY").(string)
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

	fmt.Println("Gimme obj by id")
	projectById, err := database.GetProjectById(payload.Id)
	if err != nil {
		return err
	}

	fmt.Println("Lets edit yaml")
	// Edit deploy files
	status, err := lib.DeployEditor(ArgoRepositoryFolderName, projectById)
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
