package services

import (
	"planx/internal/models"
	"planx/internal/repository"

	"fmt"
)

func CreateVendorApplication(userID, shopName, gstin string) error {
	return repository.InsertVendorApplication(userID, shopName, gstin)
}

func ListVendorApplications(status string) ([]models.VendorApplication, error) {
	return repository.GetVendorApplications(status)
}

func ListUserVendorApplications(userID, status string) ([]models.VendorApplication, error) {
	return repository.GetUserVendorApplications(userID, status)
}

func ApproveVendorApplication(appID string) error {
	// 1. Get application
	app, err := repository.GetVendorApplicationByID(appID)
	if err != nil {
		return err
	}

	if app == nil {
		return fmt.Errorf("application not found")
	}

	// 2. Validate status
	if app.Status != "pending" {
		return fmt.Errorf("application already processed")
	}

	// 3. Get vendor role ID
	roleID, err := repository.GetRoleIDByName("vendor")
	if err != nil {
		return err
	}

	if roleID == "" {
		return fmt.Errorf("vendor role not found")
	}

	// 4. Update user role
	err = repository.UpdateUserRole(app.UserID, roleID)
	if err != nil {
		return err
	}

	// 5. Update application status
	err = repository.UpdateApplicationStatus(appID, "approved")
	if err != nil {
		return err
	}

	return nil
}

func RejectVendorApplication(appID string) error {
	app, err := repository.GetVendorApplicationByID(appID)
	if err != nil {
		return err
	}

	if app == nil {
		return fmt.Errorf("application not found")
	}

	if app.Status != "pending" {
		return fmt.Errorf("application already processed")
	}

	return repository.UpdateApplicationStatus(appID, "rejected")
}



