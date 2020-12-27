package service

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type TelegramBot struct {
	BotToken string
	answer   func(chatID int, text string)
}

func NewTelegramBot(botToken string) *TelegramBot {
	return &TelegramBot{BotToken: botToken}
}

func (bot *TelegramBot) Answer(chatID int, text string) (string, error) {
	var jsonStr = []byte(fmt.Sprintf(`
{
		"chat_id": %d,
		"text": "%v",
		"parse_mode": "html"
}`, chatID, text))

	var telegramAPI string = "https://api.telegram.org/bot" + bot.BotToken + "/sendMessage"
	response, err := http.Post(telegramAPI, "application/json", bytes.NewBuffer(jsonStr))

	if err != nil {
		log.Printf("error when posting text to the chat: %s", err.Error())
		return "", err
	}
	defer response.Body.Close()

	var bodyBytes, errRead = ioutil.ReadAll(response.Body)
	if errRead != nil {
		log.Printf("error in parsing telegram answer %s", errRead.Error())
		return "", err
	}
	bodyString := string(bodyBytes)

	return bodyString, nil
}

func (bot *TelegramBot) RequestLocation(chatID int, text string) (string, error) {
	var jsonStr = []byte(fmt.Sprintf(`
{
		"chat_id": %d,
		"text": "%v",
		"reply_markup": {
		"keyboard": [
			[{
				"text": "Send Location",
				"request_location": true
			}]
		]
	}
}`, chatID, text))

	var telegramAPI string = "https://api.telegram.org/bot" + bot.BotToken + "/sendMessage"

	response, err := http.Post(telegramAPI, "application/json", bytes.NewBuffer(jsonStr))

	if err != nil {
		log.Printf("error when posting text to the chat: %s", err.Error())
		return "", err
	}
	defer response.Body.Close()

	var bodyBytes, errRead = ioutil.ReadAll(response.Body)
	if errRead != nil {
		log.Printf("error in parsing telegram answer %s", errRead.Error())
		return "", err
	}
	bodyString := string(bodyBytes)

	return bodyString, nil
}
