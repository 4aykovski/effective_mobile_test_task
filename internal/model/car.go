package model

type Car struct {
	RegistrationNumber string `json:"registration_number"`
	Mark               string `json:"mark"`
	Model              string `json:"model"`
	Year               int    `json:"year,omitempty"`
	OwnerName          string `json:"owner_name"`
	OwnerSurname       string `json:"owner_surname"`
}
