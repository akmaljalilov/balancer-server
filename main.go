package main

import (
	"github.com/gofiber/fiber/v2"
	"os"
	"velox/server/src/api/item"
)

func main() {

	app := fiber.New()
	item.Router(app.Group("item"))

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString(os.Getenv("NAME"))
	})
	err := app.Listen("0.0.0.0:8081")
	panic(err)

}
