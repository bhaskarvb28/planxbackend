package services

import (
	"fmt"
	"planx/internal/models"
	"planx/internal/repository"
)

func GiveSignupCredits(userID string) error {
	return repository.InsertUserCredits(userID, 10)
}

func GetCredits(userID string) (*models.UserCredits, error) {
	result, err := repository.GetCredits(userID)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	return &result[0], nil
}

// 🔻 Deduct credits
func DeductCredits(userID string, amount int) error {
	if amount <= 0 {
		return fmt.Errorf("invalid amount")
	}
	
	credits, err := GetCredits(userID)
	if err != nil {
		return err
	}

	if credits == nil {
		return fmt.Errorf("user credits not found")
	}

	if credits.Credits < amount {
		return fmt.Errorf("not enough credits")
	}

	newCredits := credits.Credits - amount

	return repository.UpdateCredits(userID, newCredits)
}

// 🔺 Add credits
func AddCredits(userID string, amount int) error {
	
	if amount <= 0 {
		return fmt.Errorf("invalid amount")
	}

	current, err := GetCredits(userID)
	if err != nil {
		return err
	}

	if current == nil {
		// first time user
		return repository.InsertUserCredits(userID, amount)
	}

	newCredits := current.Credits + amount

	return repository.UpdateCredits(userID, newCredits)
}