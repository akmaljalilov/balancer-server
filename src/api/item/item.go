package item

import "github.com/gofiber/fiber/v2"

func Router(r fiber.Router) {
	r.Get("/", getItemById)
}

func getItemById(c *fiber.Ctx) error {
	id := c.Query("id")
	return c.SendString(id)
}
