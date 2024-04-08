package routes

import (
	c "misesapigo/controllers"

	"github.com/gofiber/fiber/v2"
)

func SetUpRoutes(app *fiber.App) {
	app.Get("/", c.HomePage)
}