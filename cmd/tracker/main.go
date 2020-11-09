package main

import (
	"fmt"
	"os"

	"github.com/EduardTruuvaart/ev-chargepoint-tracker/domain/model"
	"github.com/EduardTruuvaart/ev-chargepoint-tracker/repository"
	"github.com/EduardTruuvaart/ev-chargepoint-tracker/service"
)

func main() {

	stationID := os.Getenv("STATIONID")
	if len(stationID) == 0 {
		panic("StationID is not provided")
	}

	apiKey := os.Getenv("APIKEY")
	if len(apiKey) == 0 {
		panic("APIKEY env var is not set!")
	}

	var stationService service.StationService
	var stationRepository repository.StationRepository

	savedStationStatus, err := stationRepository.FindByID(stationID)
	if err != nil {
		fmt.Println("Repository error occured: ", err)
		return
	}

	var currentStationStatus *model.Station = stationService.GetStatus(stationID, apiKey)

	if savedStationStatus == nil {
		stationRepository.Save(currentStationStatus)
		notifyStatusChanged(currentStationStatus.Status)
		return
	}

	if currentStationStatus.Status != savedStationStatus.Status {
		stationRepository.Save(currentStationStatus)
		notifyStatusChanged(currentStationStatus.Status)
		return
	}

	fmt.Println("Status unchanged: ", currentStationStatus.Status)
}

func notifyStatusChanged(newStatus string) {
	fmt.Println("New status: ", newStatus)
}
