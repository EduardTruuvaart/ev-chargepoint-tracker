package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/EduardTruuvaart/ev-chargepoint-tracker/domain/model"
	"github.com/EduardTruuvaart/ev-chargepoint-tracker/repository"
	"github.com/EduardTruuvaart/ev-chargepoint-tracker/service"
)

func main() {
	stationID, err := getStationID(os.Args)

	if err != nil {
		fmt.Println("Error occured: ", err)
		return
	}

	apiKey := os.Getenv("APIKEY")
	if len(apiKey) == 0 {
		fmt.Println("APIKEY env var is not set!")
		return
	}

	stationService := new(service.StationService)
	stationRepository := new(repository.StationRepository)

	savedStationStatus, err := stationRepository.FindByID(stationID)
	if err != nil {
		fmt.Println("Repository error occured: ", err)
		return
	}

	var currentStationStatus *model.Station = stationService.GetStatus(stationID, apiKey)
	fmt.Println(savedStationStatus)
	fmt.Println(currentStationStatus)
}

func getStationID(args []string) (string, error) {
	if len(args) == 2 {
		var stationID = args[1]
		return stationID, nil
	}

	stationID := os.Getenv("STATIONID")
	if len(stationID) > 0 {
		return stationID, nil
	}

	return "", errors.New("StationID is not provided")
}
