package vertexai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
	"os"

	"golang.org/x/oauth2/google"
)

// RequestPayload defines the structure of the JSON request to Vertex AI.
type RequestPayload struct {
	Contents         []Message        `json:"contents"`
	GenerationConfig GenerationConfig `json:"generationConfig"`
	SafetySettings   []SafetySetting  `json:"safetySettings"`
	Tools            []Tool           `json:"tools"`
}

// Message represents the content of a conversation.
type Message struct {
	Role  string `json:"role"`
	Parts []Part `json:"parts"`
}

// Part represents the user input message.
type Part struct {
	Text string `json:"text"`
}

// GenerationConfig defines generation parameters for the AI model.
type GenerationConfig struct {
	Temperature     float64 `json:"temperature"`
	MaxOutputTokens int     `json:"maxOutputTokens"`
	TopP            float64 `json:"topP"`
}

// SafetySetting defines safety category thresholds.
type SafetySetting struct {
	Category  string `json:"category"`
	Threshold string `json:"threshold"`
}

// Tool defines a retrieval tool for the AI model.
type Tool struct {
	Retrieval Retrieval `json:"retrieval"`
}

// Retrieval specifies the datastore to be used for search.
type Retrieval struct {
	VertexAiSearch VertexAiSearch `json:"vertexAiSearch"`
}

// VertexAiSearch holds the datastore identifier.
type VertexAiSearch struct {
	Datastore string `json:"datastore"`
}

// VertexAIClient handles interactions with Vertex AI.
type VertexAIClient struct {
	ProjectID   string
	LocationID  string
	APIEndpoint string
	ModelID     string
	httpClient  *http.Client
}

// NewVertexAIClient creates a new client for Vertex AI.
func NewVertexAIClientWithKey(projectID, locationID, apiEndpoint, modelID, credentialsFilePath string) (*VertexAIClient, error) {
	ctx := context.Background()

	// Load the service account credentials
	credentials, err := google.CredentialsFromJSON(ctx, []byte(credentialsFilePath), "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		return nil, fmt.Errorf("failed to load service account credentials: %w", err)
	}

	httpClient := oauth2.NewClient(ctx, credentials.TokenSource)
	return &VertexAIClient{
		ProjectID:   projectID,
		LocationID:  locationID,
		APIEndpoint: apiEndpoint,
		ModelID:     modelID,
		httpClient:  httpClient,
	}, nil
}

// GenerateContent sends a request to the Vertex AI endpoint and returns the response.
func (c *VertexAIClient) GenerateContent(payload RequestPayload) (string, error) {
	url := fmt.Sprintf(
		"https://%s/v1/projects/%s/locations/%s/publishers/google/models/%s:streamGenerateContent",
		c.APIEndpoint, c.ProjectID, c.LocationID, c.ModelID,
	)

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("ACCESS_TOKEN"))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("error response from Vertex AI: %s", string(body))
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(responseBody), nil
}
