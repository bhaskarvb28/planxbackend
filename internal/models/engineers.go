package models

type Engineer struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Phone          string `json:"phone"`
	Email          string `json:"email"`
	Specialization string `json:"specialization"`
	City           string `json:"city"`
	CreatedAt      string `json:"created_at"`
}