package model

type Station struct {
	ID     string
	Status string
}

func NewStation(id string, status string) *Station {
	return &Station{ID: id, Status: status}
}
