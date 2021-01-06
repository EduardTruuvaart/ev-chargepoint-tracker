package model

import (
	"fmt"
	"time"
)

type Status struct {
	Description string
	Date        time.Time
	DateTitle   string
}

func (status Status) String() string {
	return fmt.Sprintf("Description: %v - Date: %v - DateTitle: %v", status.Description, status.DateTitle, status.Date)
}
