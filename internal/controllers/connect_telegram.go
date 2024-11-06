package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const ChatID = "-4284177786"
const BotToken = "7462110864:AAHrgoOxRVPhz6B-U7ACY0XitlfSiZjKXD8"

// NotifiToTelegram sends a message to a specified Telegram chat
func NotifiToTelegram(deviceName, statusDevice, deviceType string) {
	fmt.Println("Sending notification to Telegram")

	message := fmt.Sprintf("Xin chào các quản trị viên hiện tại %s %s %s", deviceType, deviceName, statusDevice)
	data := map[string]string{
		"chat_id": ChatID,
		"text":    message,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Default().Println("Error marshaling JSON:", err)
		return
	}

	resp, err := http.Post(fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", BotToken), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Default().Println("Error:", err)
		return
	}
	defer resp.Body.Close()
	log.Default().Println("Notification sent successfully")
}
