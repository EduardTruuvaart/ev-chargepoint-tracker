package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/EduardTruuvaart/ev-chargepoint-tracker/domain/model"
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

	//stationService := stations.NewStationService(apiKey)

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
		stationStr := fmt.Sprintf("%v. %v - %v - %.3f KM", index+1, element, element.Devices[0].Status, element.DistanceInKm)
		stationsStrArr = append(stationsStrArr, stationStr)
	}

	stringyfiedResults := strings.Join(stationsStrArr, "\n")
	return stringyfiedResults
}
