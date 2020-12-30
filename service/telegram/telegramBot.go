package telegram

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type TelegramBot struct {
	BotToken                 string
	answer                   func(chatID int, text string) (string, error)
	answerCallbackQuery      func(callbackQueryId string, text string) (string, error)
	requestLocation          func(chatID int, text string) (string, error)
	answerWithRemoveKeyboard func(chatID int, text string) (string, error)
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

func (bot *TelegramBot) AnswerCallbackQuery(callbackQueryId string, text string) (string, error) {
	var jsonStr []byte
	if len(text) == 0 {
		jsonStr = []byte(fmt.Sprintf(`
{
		"callback_query_id": %v
}`, callbackQueryId))
	} else {
		jsonStr = []byte(fmt.Sprintf(`
{
		"callback_query_id": %v,
		"text": "%v",
		"parse_mode": "html"
}`, callbackQueryId, text))
	}

	var telegramAPI string = "https://api.telegram.org/bot" + bot.BotToken + "/answerCallbackQuery"
	response, err := http.Post(telegramAPI, "application/json", bytes.NewBuffer(jsonStr))

	if err != nil {
		log.Printf("error when posting callbackQuery to the chat: %s", err.Error())
		return "", err
	}
	defer response.Body.Close()

	var bodyBytes, errRead = ioutil.ReadAll(response.Body)
	if errRead != nil {
		log.Printf("error in parsing telegram callbackQuery %s", errRead.Error())
		return "", err
	}
	bodyString := string(bodyBytes)

	return bodyString, nil
}

func (bot *TelegramBot) AnswerWithRemoveKeyboard(chatID int, text string) (string, error) {
	var jsonStr = []byte(fmt.Sprintf(`
{
		"chat_id": %d,
		"text": "%v",
		"reply_markup": {
			"remove_keyboard": true
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

func (bot *TelegramBot) RequestLocation(chatID int, text string) (string, error) {
	var jsonStr = []byte(fmt.Sprintf(`
{
		"chat_id": %d,
		"text": "%v",
		"reply_markup": {
			"one_time_keyboard": true,
			"keyboard": [
				[{
					"text": "Send my location",
					"request_location": true
				}],
				[
					{"text": "Cancel"}
				]
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
