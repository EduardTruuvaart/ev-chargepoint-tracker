package model

type Station struct {
	ID     float64
	Status string
}

func NewStation(id float64, status string) *Station {
	return &Station{ID: id, Status: status}
}
