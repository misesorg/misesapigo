package routes

import (
	c "misesapigo/controllers"
	"misesapigo/middleware"

	"github.com/gofiber/fiber/v2"
)

func challengeRouter(api fiber.Router) {
	challenge := api.Group("/challenge")
	challenge.Get("/", middleware.AuthMiddleware, c.GetChallenge)
}