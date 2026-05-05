package services

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// floorPlanPrompt builds a highly detailed, structured prompt for Imagen
// to produce clean 2D architectural floor plans.
//
// Imagen responds better to dense, comma-separated descriptive tags
// rather than natural language sentences. The negative guidance is embedded
// via explicit exclusion phrases since Imagen v3 does not support a separate
// negativePrompt field through this endpoint.
func floorPlanPrompt(userRequest string) string {
	return fmt.Sprintf(`Architectural 2D floor plan blueprint, strict overhead top-down view, %s.
Style: technical drawing, crisp black vector lines on solid white background, monochrome.
Walls: thick solid black lines, uniform stroke weight.
Doors: thin quarter-circle arc swing indicators.
Windows: parallel double lines across wall segments.
Rooms: clearly labeled with large readable sans-serif text (e.g. BEDROOM, LIVING ROOM, KITCHEN, BATHROOM).
Scale bar in bottom corner. Compass rose indicating north.
Clean, minimal, professional architectural drafting standard.
NO 3D perspective. NO isometric view. NO shading or shadows. NO furniture photographs.
NO exterior landscaping. NO color fills. NO decorative elements. NO human figures.`,
		userRequest,
	)
}

type imagenRequest struct {
	Instances  []imagenInstance  `json:"instances"`
	Parameters imagenParameters  `json:"parameters"`
}

type imagenInstance struct {
	Prompt string `json:"prompt"`
}

type imagenParameters struct {
	SampleCount     int    `json:"sampleCount"`
	AspectRatio     string `json:"aspectRatio"`
	SafetyFilterLevel string `json:"safetFilterLevel,omitempty"`
}

type imagenResponse struct {
	Predictions []struct {
		BytesBase64Encoded string `json:"bytesBase64Encoded"`
		MimeType           string `json:"mimeType"`
	} `json:"predictions"`
}

func GenerateFloorPlanImage(userPrompt string) ([]byte, error) {
	projectID := os.Getenv("GCP_PROJECT_ID")
	if projectID == "" {
		return nil, fmt.Errorf("missing GCP_PROJECT_ID env var")
	}

	accessToken := os.Getenv("GCP_ACCESS_TOKEN")
	if accessToken == "" {
		return nil, fmt.Errorf("missing GCP_ACCESS_TOKEN env var")
	}

	location := "us-central1"
	model := "imagen-3.0-generate-001"

	url := fmt.Sprintf(
		"https://%s-aiplatform.googleapis.com/v1/projects/%s/locations/%s/publishers/google/models/%s:predict",
		location, projectID, location, model,
	)

	reqBody := imagenRequest{
		Instances: []imagenInstance{
			{Prompt: floorPlanPrompt(userPrompt)},
		},
		Parameters: imagenParameters{
			SampleCount: 1,
			AspectRatio: "1:1",
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		errBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("imagen error %d: %s", resp.StatusCode, errBody)
	}

	var result imagenResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}

	if len(result.Predictions) == 0 {
		return nil, fmt.Errorf("no predictions returned")
	}

	imgBytes, err := base64.StdEncoding.DecodeString(result.Predictions[0].BytesBase64Encoded)
	if err != nil {
		return nil, fmt.Errorf("base64 decode error: %w", err)
	}

	return imgBytes, nil
}

func SaveImage(data []byte, filename string) error {
	return os.WriteFile(filename, data, 0644)
}