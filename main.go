package main

import (
	"fiber_proketo/lib"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New()
	app.Use(logger.New())
	app.Get("/", helloKek)
	app.Get("/getgit", cloneRepo)
	app.Get("/pack", nixVersion)
	app.Listen(":3000")
}

func helloKek(ctx *fiber.Ctx) error {
	return ctx.SendString("Hello Kek!")

}
func cloneRepo(ctx *fiber.Ctx) error {
	err := lib.CloneRepository("Keeek")
	if err != nil {
		return err
	}
	return ctx.SendString("Hello Kek!")

}
func nixVersion(ctx *fiber.Ctx) error {
	return ctx.SendString(lib.NixpackVersion())

}
