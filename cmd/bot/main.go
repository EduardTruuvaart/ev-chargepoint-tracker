package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/EduardTruuvaart/ev-chargepoint-tracker/domain/model"
	"github.com/EduardTruuvaart/ev-chargepoint-tracker/domain/model/bot"
	"github.com/EduardTruuvaart/ev-chargepoint-tracker/service/stations"
	"github.com/EduardTruuvaart/ev-chargepoint-tracker/service/telegram"
)

func handle(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Print("Request body: ", request.Body)
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	apiKey := os.Getenv("APIKEY")
	stationService := stations.NewStationService(apiKey)

	var update bot.Update
	json.Unmarshal([]byte(request.Body), &update)

	bot := telegram.NewTelegramBot(botToken)

	if update.CallbackQuery != nil {
		bot.AnswerCallbackQuery(update.CallbackQuery.ID, "")
		code := 200
		return createAPIResponse(code), nil
	}

	if update.Message.Location != nil {
		currentLocation := *update.Message.Location
		bot.AnswerWithRemoveKeyboard(update.Message.Chat.ID, "Here is all stations in 2 KM radius:")
		stations := stationService.Search(currentLocation)
		if len(stations) > 0 {
			stations = stationService.FulfillAllDetails(currentLocation, stations)

			stringyfiedResults := createStationsResponseString(stations)

			bot.Answer(update.Message.Chat.ID, stringyfiedResults)
			return createAPIResponse(200), nil
		}

		bot.Answer(update.Message.Chat.ID, "No stations found &#x1F614")
		return createAPIResponse(200), nil
	}

	switch text := strings.ToLower(update.Message.Text); text {
	case "/start":
		bot.RequestLocation(update.Message.Chat.ID, "Hello there! Just send me your location and I will find nearby stations!")
	case "/stop":
		bot.Answer(update.Message.Chat.ID, "Bye!")
	case "cancel":
		bot.Answer(update.Message.Chat.ID, "Oh well &#x1F644")
	default:
		log.Print("Unknown text")
	}

	code := 200
	return createAPIResponse(code), nil
}

func createAPIResponse(code int) events.APIGatewayProxyResponse {
	headers := map[string]string{"Access-Control-Allow-Origin": "*", "Access-Control-Allow-Headers": "Origin, X-Requested-With, Content-Type, Accept"}
	return events.APIGatewayProxyResponse{
		StatusCode:        code,
		Headers:           headers,
		MultiValueHeaders: http.Header{"Set-Cookie": {"Ding", "Ping"}},
		Body:              "Hello, World!",
		IsBase64Encoded:   false,
	}
}

func createStationsResponseString(stations []model.Station) string {
	stationsStrArr := []string{}
	for index, element := range stations {
		stationStr := fmt.Sprintf("%v. %v - <b><i>%v</i></b>", index+1, element, element.Devices[0].Status)
		stationsStrArr = append(stationsStrArr, stationStr)
	}

	stringyfiedResults := strings.Join(stationsStrArr, "\n")
	return stringyfiedResults
}

func main() {
	lambda.Start(handle)
}
