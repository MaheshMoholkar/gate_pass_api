package types

type Staff struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	MobileNo string `json:"mobile_no"`
	Image    string `json:"image"`
}

type StaffEntry struct {
	Purpose string `json:"purpose,omitempty"`
}
