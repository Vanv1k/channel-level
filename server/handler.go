package server

import (
	"encoding/json"
	"log"

	"chat_channel_level/channel"

	"github.com/gin-gonic/gin"
)

type Segment struct {
	ID           string `json:"id"`
	TotalLength  int    `json:"total_length"`
	SegmentIndex int    `json:"segment_index"`
	Payload      string `json:"payload"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

// @Summary Передача данных через канальный уровень с учетом ошибок и потерь
// @Description Используется для передачи данных через канальный уровень с учетом возможности ошибок и потерь.
// @Tags Channel
// @Accept json
// @Produce json
// @Param data body Segment true "Данные для передачи"
// @Success 200 {object} Segment "Успешный ответ"
// @Success 204 "Пакет утерян"
// @Failure 400 {object} ErrorResponse "Некорректный запрос"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /channel-level [post]
func handleEncoding(c *gin.Context) {
	var segment Segment
	if err := c.ShouldBindJSON(&segment); err != nil {
		c.JSON(400, ErrorResponse{Message: "Неправильный формат ввода"})
		return
	}
	jsonData, err := json.Marshal(segment)
	if err != nil {
		c.JSON(500, ErrorResponse{Message: "Внутренняя ошибка сервера"})
		return
	}
	channelResponse, success, decodeError := channel.Channeltransmit(jsonData)
	if !success {
		c.Status(204)
		// c.JSON(204, ErrorResponse{Message: "Пакет утерян"})
		return
	}

	var responseData map[string]interface{}
	if err := json.Unmarshal(channelResponse, &responseData); err != nil {
		c.JSON(500, ErrorResponse{Message: "Ошибка декодирования ответа"})
		return
	}
	log.Println(decodeError)

	c.JSON(200, responseData)
}
