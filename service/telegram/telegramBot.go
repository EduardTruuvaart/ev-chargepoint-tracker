package telegram

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/EduardTruuvaart/ev-chargepoint-tracker/domain/model"
)

type TelegramBot struct {
	BotToken                    string
	answer                      func(chatID int, text string) (string, error)
	answerCallbackQuery         func(callbackQueryId string, text string) (string, error)
	requestLocation             func(chatID int, text string) (string, error)
	answerWithRemoveKeyboard    func(chatID int, text string) (string, error)
	sendStationSelectionButtons func(chatID int, text string, stations []model.Station) (string, error)
	sendStationDetails          func(chatID int, station *model.Station) (string, error)
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

func (bot *TelegramBot) SendStationSelectionButtons(chatID int, text string, stations []*model.Station) (string, error) {
	var jsonStr = fmt.Sprintf(`
{
		"chat_id": %d,
		"text": "%v",
		"reply_markup": {
			%v
		}
}`, chatID, text, renderInlineStationButtons(stations))

	var jsonStrBytes = []byte(jsonStr)

	var telegramAPI string = "https://api.telegram.org/bot" + bot.BotToken + "/sendMessage"

	response, err := http.Post(telegramAPI, "application/json", bytes.NewBuffer(jsonStrBytes))

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

func (bot *TelegramBot) SendStationDetails(chatID int, station *model.Station) (string, error) {
	message := renderStationDetailsText(station)
	return bot.Answer(chatID, message)
}

func renderStationDetailsText(station *model.Station) string {
	messageText := fmt.Sprintf(`<b>ID:</b> %v
<b>Name:</b> %v
<b>Network Name:</b> %v
<b>Address:</b> %v
	<b><i>Devices:</i></b>
		%v`, station.ID, station.Name, station.NetworkName, station.FormattedAddress, renderStationDevicesText(station.Devices))

	return messageText
}

func renderStationDevicesText(devices []*model.Device) string {
	var sb strings.Builder
	for index, element := range devices {
		sb.WriteString(fmt.Sprintf(`
		<b>ID:</b> %v
		<b>Status:</b> <i>%v</i>
		<b>History:</b> <i>%v</i>
			<b><i>Connectors:</i></b>
				%v`, element.ID, element.Status, element.LastHistoryStatus(), renderConnectorsText(element.Connectors)))

		if index != len(devices)-1 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

func renderConnectorsText(connectors []model.Connector) string {
	var sb strings.Builder
	for index, element := range connectors {
		sb.WriteString(fmt.Sprintf(`
				<b>Name:</b> %v
				<b>Speed:</b> %v
				<b>Status:</b> %v`, element.Name, element.Speed, element.Status))

		if index != len(connectors)-1 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

func renderInlineStationButtons(stations []*model.Station) string {
	var sb strings.Builder
	sb.WriteString(`
	"inline_keyboard": [
`)
	for index, element := range stations {
		stationStr := fmt.Sprintf("%v - %v - %.1f KM", element.Name, element.Devices[0].Status, element.DistanceInKm)
		sb.WriteString(fmt.Sprintf(`[{
                "text": "%v",
                "callback_data": "/details %v"
						}]`, stationStr, element.ID))
		if index != len(stations)-1 {
			sb.WriteString(",")
		}
	}
	sb.WriteString("]")

	return sb.String()
}
