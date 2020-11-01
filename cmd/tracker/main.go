package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

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

func getStationID(args []string) (int64, error) {
	if len(args) == 2 {
		var stationID, _ = strconv.ParseInt(args[1], 10, 64)
		return stationID, nil
	}

	stationID := os.Getenv("STATIONID")
	if len(stationID) > 0 {
		stationIDNum, _ := strconv.ParseInt(stationID, 10, 64)
		return stationIDNum, nil
	}

	return 0, errors.New("StationID is not provided")
}
