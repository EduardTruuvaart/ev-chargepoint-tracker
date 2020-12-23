package service

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type TelegramBot struct {
	BotToken string
}

func NewTelegramBot(botToken string) *TelegramBot {
	return &TelegramBot{BotToken: botToken}
}

func (bot *TelegramBot) Answer(chatID int, text string) (string, error) {
	var telegramAPI string = "https://api.telegram.org/bot" + bot.BotToken + "/sendMessage"
	response, err := http.PostForm(
		telegramAPI,
		url.Values{
			"chat_id": {strconv.Itoa(chatID)},
			"text":    {text},
		})

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
