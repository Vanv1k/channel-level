package server

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"chat_channel_level/channel"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Segment struct {
	ID           int    `json:"time"`
	TotalLength  int    `json:"total"`
	SegmentIndex int    `json:"number"`
	Payload      string `json:"payload"`
}

type ResponseMessage struct {
	Message string `json:"message"`
}

// @Summary Передача данных через канальный уровень с учетом ошибок и потерь
// @Description Используется для передачи данных через канальный уровень с учетом возможности ошибок и потерь.
// @Tags Channel
// @Accept json
// @Produce json
// @Param data body Segment true "Данные для передачи"
// @Success 200 {object} ResponseMessage "Успешный ответ"
// @Success 204 {object} ResponseMessage "Пакет утерян"
// @Failure 400 {object} ResponseMessage "Некорректный запрос"
// @Failure 500 {object} ResponseMessage "Внутренняя ошибка сервера"
// @Router /code [post]
func handleEncoding(c *gin.Context) {
	var segment Segment
	if err := c.ShouldBindJSON(&segment); err != nil {
		c.JSON(400, ResponseMessage{Message: "Неправильный формат ввода"})
		return
	}
	jsonData, err := json.Marshal(segment)
	if err != nil {
		c.JSON(500, ResponseMessage{Message: "Внутренняя ошибка сервера"})
		return
	}
	channelResponse, success, decodeError := channel.Channeltransmit(jsonData)
	if !success {
		c.Status(204)
		return
	}

	var responseData map[string]interface{}
	if err := json.Unmarshal(channelResponse, &responseData); err != nil {
		c.JSON(500, ResponseMessage{Message: "Ошибка декодирования ответа"})
		return
	}
	log.Println(decodeError)
	log.Println(responseData)
	responseData["error"] = decodeError
	log.Println(responseData)
	reqBody := &bytes.Buffer{}
	if err := json.NewEncoder(reqBody).Encode(responseData); err != nil {
		c.JSON(500, ResponseMessage{Message: "Ошибка при подготовке данных к отправке"})
	}

	if err := godotenv.Load(); err != nil {
		log.Fatal("Ошибка загрузки файла .env")
	}

	apiUrl := os.Getenv("API_URL")
	log.Println("URL API:", apiUrl)
	log.Println(reqBody)
	resp, err := http.Post(apiUrl, "application/json", reqBody)
	if err != nil {
		c.JSON(400, ResponseMessage{Message: "Ошибка при отправке сегмента на эндпоинт"})
	}

	defer resp.Body.Close()
	log.Println(resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		c.JSON(400, ResponseMessage{Message: "Ошибка: неверный код состояния ответа"})
		return
	}

	c.JSON(200, ResponseMessage{Message: "Сегмент успешно отправлен"})
}
