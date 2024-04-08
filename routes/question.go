package routes

import (
	c "misesapigo/controllers"
	"misesapigo/middleware"

	"github.com/gofiber/fiber/v2"
)

func questionRouter(api fiber.Router) {
	question := api.Group("/question")
	question.Get("/", middleware.AuthMiddleware, c.GetQuestion)
}