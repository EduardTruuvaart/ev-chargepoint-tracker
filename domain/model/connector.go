package model

import "fmt"

type Connector struct {
	ID     string
	Name   string
	Speed  string
	Status string
}

func NewConnector(id string) *Connector {
	return &Connector{ID: id}
}

func (connector Connector) String() string {
	return fmt.Sprintf("%v - %v - %v - %v", connector.ID, connector.Name, connector.Speed, connector.Status)
}
