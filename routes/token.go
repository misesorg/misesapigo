package routes

import (
	c "misesapigo/controllers"
	"misesapigo/middleware"

	"github.com/gofiber/fiber/v2"
)

func tokenRouter(api fiber.Router) {
	token := api.Group("/token")
	token.Get("/refresh", middleware.AuthMiddleware, c.RefreshToken)
}