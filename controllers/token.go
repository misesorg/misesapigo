package controllers

import (
	"bytes"
	"encoding/json"
	"misesapigo/config"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func RefreshToken(c *fiber.Ctx) error {
	type refreshRequest struct {
		GrantType    string `json:"grant_type"`
		RefreshToken string `json:"refresh_token"`
	}
	
	type refreshResponse struct {
		AccessToken  string `json:"access_token"`
		ExpiresIn    string `json:"expires_in"`
		TokenType    string `json:"token_type"`
		RefreshToken string `json:"refresh_token"`
		IDToken      string `json:"id_token"`
		UserID       string `json:"user_id"`
		ProjectID    string `json:"project_id"`
	}

	payload, err := json.Marshal(refreshRequest{
		GrantType:    "refresh_token",
		RefreshToken: config.Config("REFRESHTOKEN"),
	})

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	requestUrl := "https://securetoken.googleapis.com/v1/token?key=" + config.Config("USERTOKEN")

	resp, err := http.Post(requestUrl, "application/json", bytes.NewReader(payload))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	defer resp.Body.Close()

	var js refreshResponse

	if err := json.NewDecoder(resp.Body).Decode(&js); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"AccessToken": js.AccessToken,
		"ExpiresIn": js.ExpiresIn,
		"TokenType": js.TokenType,
		"RefreshToken": js.RefreshToken,
		"IDToken": js.IDToken,
		"UserID": js.UserID,
		"ProjectID": js.ProjectID,
	})
}
