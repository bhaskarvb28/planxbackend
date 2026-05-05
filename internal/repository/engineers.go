package repository

import (
	"planx/internal/models"
	"net/http"
	"encoding/json"
	"os"
	"fmt"
	"bytes"
)

func GetEngineers(limit, offset int, specialization, city string) ([]models.Engineer, error) {
	baseURL := os.Getenv("SUPABASE_URL") + "/rest/v1/engineers"

	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Set("select", "*")

	// filters
	if specialization != "" {
		q.Set("specialization", "ilike.*"+specialization+"*")
	}
	if city != "" {
		q.Set("city", "ilike.*"+city+"*")
	}

	// pagination
	q.Set("limit", fmt.Sprintf("%d", limit))
	q.Set("offset", fmt.Sprintf("%d", offset))

	req.URL.RawQuery = q.Encode()

	serviceKey := os.Getenv("SUPABASE_SERVICE_KEY")
	req.Header.Set("apikey", serviceKey)
	req.Header.Set("Authorization", "Bearer "+serviceKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result []models.Engineer
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func InsertEngineer(name, phone, email, specialization, city string) error {
	url := os.Getenv("SUPABASE_URL") + "/rest/v1/engineers"

	payload := map[string]interface{}{
		"name":           name,
		"phone":          phone,
		"email":          email,
		"specialization": specialization,
		"city":           city,
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
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
		return fmt.Errorf("failed to insert engineer")
	}

	return nil
}