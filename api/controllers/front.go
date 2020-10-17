package controllers

import "net/http"

func RegisterControllers() {
	var sc *stationController = newStationController()

	http.Handle("/stations", sc)
	http.Handle("/stations/", sc)
}
