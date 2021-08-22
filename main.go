package main

import (
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"os"
	"velox/server/src/api/item"
	"velox/server/src/migrate"
	"velox/server/src/service/consistenHashing"
)

// @title Server API
// @version 1.0
// @description API Server

// @in header
// @host localhost:8081
// @BasePath /
func main() {

	app := fiber.New()
	SwaggerRoute(app)
	item.Router(app.Group("item"))
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString(os.Getenv("NAME"))
	})
	consistenHashing.Start()
	migrate.Start()
	err := app.Listen("0.0.0.0:8081")
	panic(err)

}
func SwaggerRoute(a fiber.Router) {
	// Create routes group.
	route := a.Group("/swagger")

	// Routes for GET method:
	route.Get("*", swagger.Handler) // get one user by ID
}
