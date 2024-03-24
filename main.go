package main

import (
	"log"

	"chat_channel_level/server"
)

// @title Channel Level API
// @version 1.0
// @description API для передачи данных через канал с использованием кодирования Хэммингом (7,4)

// @host localhost:5000
// @BasePath /

func main() {
	log.Println("Server Started")
	server.Start()
	log.Println("Server Terminated")
}
