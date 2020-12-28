package model

import "fmt"

// Location indicates coordinates.
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (location Location) String() string {
	return fmt.Sprintf("%v, %v ", location.Latitude, location.Longitude)
}
