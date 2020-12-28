package model

import "fmt"

type Device struct {
	ID         string
	Status     string
	Connectors []Connector
}

func NewDevice(id string) *Device {
	return &Device{
		ID:         id,
		Connectors: []Connector{},
	}
}

func (device Device) String() string {
	return fmt.Sprintf("%v - %v - Connectors: %v", device.ID, device.Status, device.Connectors)
}
