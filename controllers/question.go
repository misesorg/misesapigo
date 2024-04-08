package controllers

import (
	"encoding/json"
	"net/http"
	"misesapigo/config"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/gofiber/fiber/v2"
)

func GetQuestion(c *fiber.Ctx) error {
	type responseRequest struct {
		Id          string   `json:"_id"`
		Options     []string `json:"options"`
		Tags        []string `json:"tags"`
		TestName    string   `json:"testName"`
		Content     string   `json:"content"`
		Markdown    string   `json:"markdown"`
		Level       int      `json:"level"`
		RightAnswer int      `json:"rightAnswer"`
		SubjectName string   `json:"subjectName"`
		CreatedAt   string   `json:"createdAt"`
		UpdatedAt   string   `json:"updatedAt"`
		IsFormatted bool     `json:"IsFormatted"`
	}

	categoryTranslator := map[string]string{
		"random":     "Aleatório",
		"biologia":   "Biologia",
		"quimica":    "Química",
		"fisica":     "Física",
		"matematica": "Matemática",
		"geografia":  "Geografia",
		"historia":   "História",
		"filosofia":  "Filosofia",
		"sociologia": "Sociologia",
		"portugues":  "Português",
		"literatura": "Literatura",
		"ingles":     "Inglês",
		"espanhol":   "Espanhol",
		"artes":      "Artes",
	}

	levelTranslator := map[int]string{
		0: "Fácil",
		1: "Médio",
		2: "Difícil",
	}

	payloadBody := struct {
		Category     string `json:"category"`
		Level        string `json:"level"`
		RefreshToken string `json:"refresh_token"`
	}{}

	if err := c.BodyParser(&payloadBody); err != nil {
		return err
	}

	requestUrl := config.Config("REQUEST_API_URL") + "api/exercise/category/" + payloadBody.Category + "/" + payloadBody.Level

	client := &http.Client{}
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	req.Header.Add("accept", "application/json, text/plain, */*")
	req.Header.Add("if-none-match", `W/"aa1-a07BLy2DCqtpwRBVGptJXSScZDs"`)
	req.Header.Add("authtoken", payloadBody.RefreshToken)
	req.Header.Add("user-agent", "RevisApp/3 CFNetwork/1402.0.8 Darwin/22.2.0")
	req.Header.Add("accept-language", "pt-BR,pt;q=0.9")
	req.Header.Add("accept-encoding", "gzip, deflate, br")

	resp, err := client.Do(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	defer resp.Body.Close()

	var js responseRequest

	if err := json.NewDecoder(resp.Body).Decode(&js); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	converter := md.NewConverter("", true, nil)
	markdown, err := converter.ConvertString(js.Content)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	if js.Content == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Token is not available or expired.",
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"Id":               js.Id,
		"Options":          js.Options,
		"Tags":             js.Tags,
		"TestName":         js.TestName,
		"Content":          js.Content,
		"Markdown":         markdown,
		"Level":            js.Level,
		"FormattedLevel":   levelTranslator[js.Level],
		"RightAnswer":      js.RightAnswer,
		"SubjectName":      js.SubjectName,
		"FormattedSubject": categoryTranslator[js.SubjectName],
		"CreatedAt":        js.CreatedAt,
		"UpdatedAt":        js.UpdatedAt,
		"IsFormatted":      js.IsFormatted,
	})
}
