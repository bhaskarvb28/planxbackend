package models

type UserCredits struct {
	UserID  string `json:"user_id"`
	Credits int    `json:"credits"`
}