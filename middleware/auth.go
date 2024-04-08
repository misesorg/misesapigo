package middleware

import (
	"strings"
	"misesapigo/config"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(c *fiber.Ctx) error {
	getToken := c.GetReqHeaders()
	tokenString := getToken["Authorization"]
	
	if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Please provide a valid bearer token in the authorization header.",
		})
	}

	bearerToken := config.Config("APP_BEARER_TOKEN");

	if tokenString[7:] != bearerToken {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "The provided bearer token is invalid or has expired.",
		})
	}

	return c.Next()
}
