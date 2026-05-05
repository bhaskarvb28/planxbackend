package repository

import (
	"os"
	"encoding/json"
	"net/http"
	"bytes"
	"fmt"

	"planx/internal/models"
)

func InsertVendorApplication(userID, shopName, gstin string) error {
	url := os.Getenv("SUPABASE_URL") + "/rest/v1/vendor_applications"

	payload := map[string]interface{}{
		"user_id":   userID,
		"shop_name": shopName,
		"gstin":     gstin,
	}

	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))

	serviceKey := os.Getenv("SUPABASE_SERVICE_KEY")
	req.Header.Set("apikey", serviceKey)
	req.Header.Set("Authorization", "Bearer "+serviceKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("failed to insert application")
	}

	return nil
}

func GetVendorApplications(status string) ([]models.VendorApplication, error) {
	supabaseURL := os.Getenv("SUPABASE_URL")
	serviceKey := os.Getenv("SUPABASE_SERVICE_KEY")

	url := supabaseURL +  "/rest/v1/vendor_applications?select=*"

	// filter by status if provided
	if status != "" {
		url += "&status=eq." + status
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("apikey", serviceKey)
	req.Header.Set("Authorization", "Bearer "+serviceKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result []models.VendorApplication
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func GetUserVendorApplications(userID, status string) ([]models.VendorApplication, error) {
	url := os.Getenv("SUPABASE_URL") +
		"/rest/v1/vendor_applications?user_id=eq." + userID + "&select=*"

	if status != "" {
		url += "&status=eq." + status
	}

	req, _ := http.NewRequest("GET", url, nil)

	serviceKey := os.Getenv("SUPABASE_SERVICE_KEY")
	req.Header.Set("apikey", serviceKey)
	req.Header.Set("Authorization", "Bearer "+serviceKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result []models.VendorApplication
	json.NewDecoder(resp.Body).Decode(&result)

	return result, nil
}

func GetVendorApplicationByID(appID string) (*models.VendorApplication, error) {
	url := os.Getenv("SUPABASE_URL") +
		"/rest/v1/vendor_applications?id=eq." + appID + "&select=*"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	serviceKey := os.Getenv("SUPABASE_SERVICE_KEY")
	req.Header.Set("apikey", serviceKey)
	req.Header.Set("Authorization", "Bearer "+serviceKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result []models.VendorApplication
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	return &result[0], nil
}

func UpdateUserRole(userID string, roleID string) error {
	url := os.Getenv("SUPABASE_URL") +
		"/rest/v1/profiles?id=eq." + userID

	payload := map[string]interface{}{
		"role_id": roleID,
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	serviceKey := os.Getenv("SUPABASE_SERVICE_KEY")
	req.Header.Set("apikey", serviceKey)
	req.Header.Set("Authorization", "Bearer "+serviceKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("failed to update user role")
	}

	return nil
}

func UpdateApplicationStatus(appID string, status string) error {
	url := os.Getenv("SUPABASE_URL") +
		"/rest/v1/vendor_applications?id=eq." + appID

	payload := map[string]interface{}{
		"status": status,
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	serviceKey := os.Getenv("SUPABASE_SERVICE_KEY")
	req.Header.Set("apikey", serviceKey)
	req.Header.Set("Authorization", "Bearer "+serviceKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("failed to update application status")
	}

	return nil
}