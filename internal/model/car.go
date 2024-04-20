package model

type Car struct {
	RegistrationNumber string `db:"registration_number" json:"regNumber"`
	Mark               string `db:"mark" json:"mark"`
	Model              string `db:"model" json:"model"`
	Year               int    `db:"year" json:"year,omitempty"`
	OwnerName          string `db:"owner_name" json:"ownerName"`
	OwnerSurname       string `db:"owner_surname" json:"ownerSurname"`
}
