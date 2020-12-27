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
	"github.com/EduardTruuvaart/ev-chargepoint-tracker/service"
)

func handle(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Print("Request body: ", request.Body)
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	apiKey := os.Getenv("APIKEY")
	stationService := service.NewStationService(apiKey)

	var update bot.Update
	json.Unmarshal([]byte(request.Body), &update)

	bot := service.NewTelegramBot(botToken)

	if update.Message.Location != nil {
		bot.Answer(update.Message.Chat.ID, "Here is all stations in 2 KM radius:")
		stations := stationService.Search(*update.Message.Location)
		stations = stationService.FulfillAllDetails(stations)
		stringyfiedResults := createStationsResponseString(stations)

		bot.Answer(update.Message.Chat.ID, stringyfiedResults)
		return createAPIResponse(200), nil
	}

	switch text := update.Message.Text; text {
	case "/start":
		bot.Answer(update.Message.Chat.ID, "Hello there! Just send me your location and I will find nearby stations!")
	case "/locate":
		bot.RequestLocation(update.Message.Chat.ID, "Please provide your location")
	case "/stop":
		bot.Answer(update.Message.Chat.ID, "Bye!")
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
		stationStr := fmt.Sprintf("%v. %v - <b><i>%v</i></b>", index+1, element, element.Status)
		stationsStrArr = append(stationsStrArr, stationStr)
	}

	stringyfiedResults := strings.Join(stationsStrArr, "\n")
	return stringyfiedResults
}

func main() {
	lambda.Start(handle)
}
