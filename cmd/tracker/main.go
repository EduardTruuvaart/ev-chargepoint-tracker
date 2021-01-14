package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/EduardTruuvaart/ev-chargepoint-tracker/domain/model"
	"github.com/EduardTruuvaart/ev-chargepoint-tracker/service/stations"
)

func main() {
	apiKey := os.Getenv("APIKEY")
	if len(apiKey) == 0 {
		panic("APIKEY env var is not set!")
	}

	location := model.Location{Latitude: 51.494698, Longitude: -0.153487}
	stationService := stations.NewStationService(apiKey)

	start := time.Now()

	stationIDs := stationService.Search(location)
	fmt.Printf("Stations count: %v\n", len(stationIDs))
	stationService.GetStations(location, stationIDs)

	elapsed := time.Since(start)
	fmt.Println(elapsed)

	// var wg sync.WaitGroup
	// wg.Add(10)
	// for i := 0; i < 10; i++ {
	// 	go func(j int) {
	// 		notifyStatusChanged(fmt.Sprintf("Hello %v", j))
	// 		wg.Done()
	// 	}(i)
	// }
	// wg.Wait()

	//stationService := stations.NewStationService(apiKey)
	//devices := stationService.GetStatus("892")

	//fmt.Println(devices[0].LastHistoryStatus())

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
	time.Sleep(time.Millisecond * 500)
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
