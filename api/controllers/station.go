package controllers

import (
	"net/http"
	"regexp"
)

type stationController struct {
	stationIDPattern *regexp.Regexp
}

func (s *stationController) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	//var station *model.Station = model.NewStation("123")
	rw.Write([]byte("123"))
}

func newStationController() *stationController {
	return &stationController{
		stationIDPattern: regexp.MustCompile(`^/stations/(\d+)/?`),
	}
}
