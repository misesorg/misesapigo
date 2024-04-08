package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"misesapigo/config"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/gofiber/fiber/v2"
)

func GetChallenge(c *fiber.Ctx) error {
	type ResponseRequest struct {
		Id       string `json:"_id"`
		Question struct {
			Id          string   `json:"_id"`
			Options     []string `json:"options"`
			Tags        []string `json:"tags"`
			TestName    string   `json:"testName"`
			Content     string   `json:"content"`
			Level       int      `json:"level"`
			RightAnswer int      `json:"rightAnswer"`
			SubjectName string   `json:"subjectName"`
			CreatedAt   string   `json:"createdAt"`
			UpdatedAt   string   `json:"updatedAt"`
			IsFormatted bool     `json:"isFormatted"`
		} `json:"question"`
		QuestionId   string   `json:"questionId"`
		RightAnswers int      `json:"rightAnswers"`
		WrongAnswers int      `json:"wrongAnswers"`
		Comments     []string `json:"comments"`
		CreatedAt    string   `json:"createdAt"`
		UpdatedAt    string   `json:"updatedAt"`
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
		RefreshToken string `json:"refresh_token"`
	}{}

	if err := c.BodyParser(&payloadBody); err != nil {
		return err
	}

	if payloadBody.RefreshToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "You must provide a refresh token",
		})
	}

	requestUrl := config.Config("REQUEST_API_URL") + "/api/daily-challenge"

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

	var js ResponseRequest

	if err := json.NewDecoder(resp.Body).Decode(&js); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	converter := md.NewConverter("", true, nil)
	markdown, err := converter.ConvertString(js.Question.Content)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	rightAnswers := js.RightAnswers
	wrongAnswers := js.WrongAnswers
	totalAnswers := rightAnswers + wrongAnswers
	averagePercent := float64(rightAnswers) / float64(totalAnswers) * 100
	roundedPercentage := strconv.FormatFloat(averagePercent, 'f', 2, 64)

	response := make(map[string]interface{})
	response["Id"] = js.Id
	response["Question"] = map[string]interface{}{
		"Id":               js.Question.Id,
		"Options":          js.Question.Options,
		"Tags":             js.Question.Tags,
		"TestName":         js.Question.TestName,
		"Content":          js.Question.Content,
		"Markdown":         markdown,
		"Level":            js.Question.Level,
		"FormattedLevel":   levelTranslator[js.Question.Level],
		"RightAnswer":      js.Question.RightAnswer,
		"SubjectName":      js.Question.SubjectName,
		"FormattedSubject": categoryTranslator[js.Question.SubjectName],
		"CreatedAt":        js.Question.CreatedAt,
		"UpdatedAt":        js.Question.UpdatedAt,
		"IsFormatted":      js.Question.IsFormatted,
	}
	response["QuestionId"] = js.QuestionId
	response["RightAnswers"] = js.RightAnswers
	response["WrongAnswers"] = js.WrongAnswers
	response["AveragePercent"] = roundedPercentage
	response["Comments"] = js.Comments
	response["CreatedAt"] = js.CreatedAt
	response["UpdatedAt"] = js.UpdatedAt

	if js.QuestionId == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Token is not available or expired.",
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(response)
}
