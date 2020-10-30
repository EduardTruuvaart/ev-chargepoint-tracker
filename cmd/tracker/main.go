package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/EduardTruuvaart/ev-chargepoint-tracker/service"
)

func main() {
	var stationID = getStationID(os.Args)
	apiKey := os.Getenv("APIKEY")
	if len(apiKey) == 0 {
		fmt.Println("APIKEY env var is not set!")
		return
	}

	stationService := new(service.StationService)
	station := stationService.GetStatus(stationID, apiKey)
	fmt.Println(station)
}

func getStationID(args []string) int64 {
	if len(args) == 2 {
		var stationID, _ = strconv.ParseInt(args[1], 10, 64)
		return stationID
	}

	return 806
}
