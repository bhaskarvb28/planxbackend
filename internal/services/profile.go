package services

import ( 
	"planx/internal/repository"
	"strings"
)

func GetMyProfile(userID string) (map[string]interface{}, error) {
	result, err := repository.GetProfileByID(userID)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	return result[0], nil
}

func CreateProfile(userID string, name string) (map[string]interface{}, string, error) {
	// 1. Check existing
	existing, err := repository.GetProfileByID(userID)
	if err != nil {
		return nil, "", err
	}
	if len(existing) > 0 {
		return existing[0], "exists", nil
	}

	// 2. Get role_id
	roleID, err := repository.GetRoleIDByName("user")
	if err != nil {
		return nil, "", err
	}
	if roleID == "" {
		return nil, "", err
	}

	// 3. Fallback if name empty (safety)
	if name == "" {
		name = "New User"
	}

	// 4. Insert
	profile, err := repository.InsertProfile(userID, name, roleID)
	if err != nil {
		// Check if duplicate key error
		if strings.Contains(err.Error(), "duplicate key") {
			existing, err := repository.GetProfileByID(userID)
			if err != nil {
				return nil, "", err
			}
			if len(existing) > 0 {
				return existing[0], "exists", nil
			}
		}

		return nil, "", err
	}

	_ = GiveSignupCredits(userID)

	return profile, "created", nil
}

func GetUserRole(userID string) (string, error) {
	return repository.GetUserRole(userID)
}