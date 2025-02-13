package entities

type User struct {
	UUID      string `json:"uuid"`
	FirstName string `json:"name.first"`
	LastName  string `json:"name.last"`
	Email     string `json:"email"`
	City      string `json:"location.city"`
	Country   string `json:"location.country"`
	Gender    string `json:"gender"`
}
