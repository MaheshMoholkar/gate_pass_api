package types

import (
	"time"
)

type User struct {
	ID       int    `json:"id,omitempty"`
	MobileNo string `json:"mobile_no"`
	Password string `json:"password"`
}

type Visitor struct {
	ID          int       `json:"id,omitempty"`
	Name        string    `json:"name"`
	Purpose     string    `json:"purpose"`
	Date        time.Time `json:"date"`
	Address     string    `json:"address,omitempty"`
	VehicleNo   int       `json:"vehicle_no,omitempty"`
	MobileNo    string    `json:"mobile_no"`
	Image       string    `json:"image,omitempty"`
	Appointment string    `json:"appointment,omitempty"`
	In          time.Time `json:"in"`
	Out         time.Time `json:"out,omitempty"`
}
