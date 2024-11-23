package vertex_ai

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"net/http"
)

// VertexAIClient handles interactions with Vertex AI.
type VertexAIClient struct {
	ProjectID   string
	LocationID  string
	APIEndpoint string
	ModelID     string
	httpClient  *http.Client
	configs     *Configs
	logger      Logger
}

// NewVertexAIClientWithKey creates a new client for Vertex AI.
func NewVertexAIClientWithKey(configs *Configs) (*VertexAIClient, error) {
	ctx := context.Background()

	configs.setDefaults()

	err := configs.validate()
	if err != nil {
		return nil, err
	}

	// Load the service account credentials
	credentials, err := google.CredentialsFromJSON(ctx, []byte(configs.Credentials), "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		return nil, fmt.Errorf("failed to load service account credentials: %w", err)
	}

	httpClient := oauth2.NewClient(ctx, credentials.TokenSource)

	return &VertexAIClient{
		ProjectID:   configs.ProjectID,
		LocationID:  configs.LocationID,
		APIEndpoint: configs.APIEndpoint,
		ModelID:     configs.ModelID,
		httpClient:  httpClient,
		configs:     configs,
	}, nil
}

// GetResponse sends a request to the Vertex AI endpoint and returns the response.
func (c *VertexAIClient) GetResponse(prompt string, datastore ...string) (string, error) {
	url := fmt.Sprintf(
		"https://%s/v1/projects/%s/locations/%s/publishers/google/models/%s:streamGenerateContent",
		c.APIEndpoint, c.ProjectID, c.LocationID, c.ModelID,
	)

	payload := c.generateRequestPayload(prompt, datastore...)

	response, err := c.getResponseFromAPI(url, payload)
	if err != nil {
		return "", err
	}

	return concatenateAllEntries(response), nil
}

func (c *VertexAIClient) UseLogger(logger any) {
	if l, ok := logger.(Logger); ok {
		c.logger = l
	}
}
