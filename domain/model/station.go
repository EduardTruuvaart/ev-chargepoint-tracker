package model

import "fmt"

type Station struct {
	ID               string
	Status           string
	Name             string
	NetworkName      string
	FormattedAddress string
	PostCode         string
}

func NewStation(id string) *Station {
	return &Station{ID: id}
}

func (station Station) String() string {
	return fmt.Sprintf("%v - %v - %v - %v", station.ID, station.Name, station.FormattedAddress, station.NetworkName)
}
