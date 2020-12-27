package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/EduardTruuvaart/ev-chargepoint-tracker/domain/model"
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

	stationService := service.NewStationService(apiKey)

	location := model.Location{Latitude: 51.394284, Longitude: -0.304267}
	stations := stationService.Search(location)
	stations = stationService.FulfillAllDetails(stations)
	stringyfiedResults := createStationsResponseString(stations)
	fmt.Print(stringyfiedResults)

	// var stationRepository repository.StationRepository

	// savedStationStatus, err := stationRepository.FindByID(stationID)
	// if err != nil {
	// 	fmt.Println("Repository error occured: ", err)
	// 	return
	// }

	// var currentStationStatus *model.Station = stationService.GetStatus(stationID)

	// if savedStationStatus == nil {
	// 	stationRepository.Save(currentStationStatus)
	// 	notifyStatusChanged(currentStationStatus.Status)
	// 	return
	// }

	// if currentStationStatus.Status != savedStationStatus.Status {
	// 	stationRepository.Save(currentStationStatus)
	// 	notifyStatusChanged(currentStationStatus.Status)
	// 	return
	// }

	// fmt.Println("Status unchanged: ", currentStationStatus.Status)
}

func notifyStatusChanged(newStatus string) {
	fmt.Println("New status: ", newStatus)
}

func createStationsResponseString(stations []model.Station) string {
	stationsStrArr := []string{}
	for index, element := range stations {
		stationStr := fmt.Sprintf("%v. %v", index+1, element)
		stationsStrArr = append(stationsStrArr, stationStr)
	}

	stringyfiedResults := strings.Join(stationsStrArr, "\n")
	return stringyfiedResults
}
