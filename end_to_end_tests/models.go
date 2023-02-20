package end_to_end_tests

type findNumberPositionResponse struct {
	Number int    `json:"number,omitempty"`
	Error  string `json:"error,omitempty"`
}

type getUsersResponse struct {
	Users []*User `json:"users,omitempty"`
	Error string  `json:"error,omitempty"`
}

type User struct {
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	Birthday    string `json:"birthday,omitempty"`
	Address     string `json:"address,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
}
