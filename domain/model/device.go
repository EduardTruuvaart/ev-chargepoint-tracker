package model

import (
	"fmt"
	"time"
)

type Device struct {
	ID                string
	Status            string
	lastHistoryStatus string
	StatusHistory     []Status
	Connectors        []Connector
}

func NewDevice(id string) *Device {
	return &Device{
		ID:            id,
		Connectors:    []Connector{},
		StatusHistory: []Status{},
	}
}

func (device Device) String() string {
	return fmt.Sprintf("%v - %v - StatusHistory: %v - Connectors: %v", device.ID, device.Status, device.StatusHistory, device.Connectors)
}

func (device Device) LastHistoryStatus() string {
	if len(device.StatusHistory) > 1 { // todo: add single item display
		duration := time.Since(device.StatusHistory[1].Date)
		return formatDuration(device.StatusHistory[1].Description, duration)
	} else if len(device.StatusHistory) == 1 {
		duration := time.Since(device.StatusHistory[0].Date)
		return formatDuration(device.StatusHistory[0].Description, duration)
	}

	return ""
}

func formatDuration(description string, duration time.Duration) string {
	const (
		Decisecond = 100 * time.Millisecond
		Day        = 24 * time.Hour
	)
	sign := time.Duration(1)
	if duration < 0 {
		sign = -1
		duration = -duration
	}
	duration += +Decisecond / 2
	d := sign * (duration / Day)
	duration = duration % Day
	h := duration / time.Hour
	duration = duration % time.Hour
	m := duration / time.Minute

	if d > 0 {
		return fmt.Sprintf("%v - %2d days %2d hours %2d minutes ago", description, d, h, m)
	} else if h > 0 {
		return fmt.Sprintf("%v - %2d hours %2d minutes ago", description, h, m)
	} else {
		return fmt.Sprintf("%v - %2d minutes ago", description, m)
	}
}
