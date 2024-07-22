package types

import "time"

type Staff struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	MobileNo string `json:"mobile_no"`
	Image    string `json:"image"`
}

type StaffEntry struct {
	ID       int       `json:"id,omitempty"`
	Name     string    `json:"name"`
	Purpose  string    `json:"purpose,omitempty"`
	In       time.Time `json:"in,omitempty"`
	Out      time.Time `json:"out,omitempty"`
	MobileNo string    `json:"mobile_no,omitempty"`
	Image    string    `json:"image,omitempty"`
}
