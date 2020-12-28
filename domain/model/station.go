package model

import (
	"fmt"
)

type Station struct {
	ID               string
	Name             string
	NetworkName      string
	FormattedAddress string
	PostCode         string
	Location         Location
	Devices          []Device
}

func NewStation(id string) *Station {
	return &Station{
		ID:      id,
		Devices: []Device{},
	}
}

func (station Station) String() string {
	return fmt.Sprintf("%v - %v - %v - %v", station.ID, station.Name, station.FormattedAddress, station.NetworkName)
}
