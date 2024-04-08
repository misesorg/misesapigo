package main

import (
	"log"
	"misesapigo/config"
	router "misesapigo/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		panic("Error loading .env file")
	}

	app := fiber.New()

	router.SetUpRoutes(app)
	router.SetUpApiRoutes(app)

	app.Use(cors.New())

	app.Use(func (c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNotFound)
	})

	log.Fatal(app.Listen("0.0.0.0:" + config.Config("PORT")))
}