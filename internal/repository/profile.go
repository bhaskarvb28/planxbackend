package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func GetProfileByID(userID string) ([]map[string]interface{}, error) {
	supabaseURL := os.Getenv("SUPABASE_URL")
	serviceKey := os.Getenv("SUPABASE_SERVICE_KEY")

	url := supabaseURL + "/rest/v1/profiles?id=eq." + userID + "&select=*"

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("apikey", serviceKey)
	req.Header.Set("Authorization", "Bearer "+serviceKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	return result, nil
}

func GetRoleIDByName(role string) (string, error) {
	supabaseURL := os.Getenv("SUPABASE_URL")
	serviceKey := os.Getenv("SUPABASE_SERVICE_KEY")

	url := supabaseURL + "/rest/v1/roles?name=eq." + role + "&select=id"

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("apikey", serviceKey)
	req.Header.Set("Authorization", "Bearer "+serviceKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var roles []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&roles)

	if len(roles) == 0 {
		return "", nil
	}

	id, _ := roles[0]["id"].(string)
	return id, nil
}

func InsertProfile(userID, name, roleID string) (map[string]interface{}, error) {
	supabaseURL := os.Getenv("SUPABASE_URL")
	serviceKey := os.Getenv("SUPABASE_SERVICE_KEY")

	payload := map[string]interface{}{
		"id":      userID,
		"name":    name,
		"role_id": roleID,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		"POST",
		supabaseURL+"/rest/v1/profiles",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("apikey", serviceKey)
	req.Header.Set("Authorization", "Bearer "+serviceKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "return=representation")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 🔥 Handle non-200 properly
	if resp.StatusCode >= 300 {
		var errResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errResp)
		return nil, fmt.Errorf("supabase error: %v", errResp)
	}

	var result []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	return result[0], nil
}

func GetUserRole(userID string) (string, error) {
	url := os.Getenv("SUPABASE_URL") +
		"/rest/v1/profiles?id=eq." + userID + "&select=role:roles(name)"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	serviceKey := os.Getenv("SUPABASE_SERVICE_KEY")
	req.Header.Set("apikey", serviceKey)
	req.Header.Set("Authorization", "Bearer "+serviceKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result []struct {
		Role struct {
			Name string `json:"name"`
		} `json:"role"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result) == 0 {
		return "", fmt.Errorf("role not found")
	}

	return result[0].Role.Name, nil
}