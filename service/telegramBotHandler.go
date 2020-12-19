package service

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/EduardTruuvaart/ev-chargepoint-tracker/domain/model/bot"
)

func parseTelegramRequest(r *http.Request) (*bot.Update, error) {
	var update bot.Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		log.Printf("could not decode incoming update %s", err.Error())
		return nil, err
	}
	return &update, nil
}

// HandleTelegramWebHook sends a message back to the chat with a punchline starting by the message provided by the user.
func HandleTelegramWebHook(w http.ResponseWriter, r *http.Request) {
	// Parse incoming request
	var update, err = parseTelegramRequest(r)
	if err != nil {
		log.Printf("error parsing update, %s", err.Error())
		return
	}

	// Sanitize input
	var sanitizedSeed = sanitize(update.Message.Text)

	// Call RapLyrics to get a punchline
	var lyric, errRapLyrics = getPunchline(sanitizedSeed)
	if errRapLyrics != nil {
		log.Printf("got error when calling RapLyrics API %s", errRapLyrics.Error())
		return
	}

	// Send the punchline back to Telegram
	var telegramResponseBody, errTelegram = sendTextToTelegramChat(update.Message.Chat.ID, lyric)
	if errTelegram != nil {
		log.Printf("got error %s from telegram, reponse body is %s", errTelegram.Error(), telegramResponseBody)
	} else {
		log.Printf("punchline %s successfuly distributed to chat id %d", lyric, update.Message.Chat.ID)
	}
}

func sanitize(text string) string {
	return text
}

func getPunchline(text string) (string, error) {
	return text, nil
}

// sendTextToTelegramChat sends a text message to the Telegram chat identified by its chat Id
func sendTextToTelegramChat(chatID int, text string) (string, error) {

	log.Printf("Sending %s to chat_id: %d", text, chatID)
	var telegramAPI string = "https://api.telegram.org/bot" + os.Getenv("TELEGRAM_BOT_TOKEN") + "/sendMessage"
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
	log.Printf("Body of Telegram Response: %s", bodyString)

	return bodyString, nil
}
