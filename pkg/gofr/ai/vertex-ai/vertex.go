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
func (c *VertexAIClient) SendMessage(ctx context.Context, prompt []map[string]string) (string, error) {
	url := fmt.Sprintf(
		"https://%s/v1/projects/%s/locations/%s/publishers/google/models/%s:streamGenerateContent",
		c.APIEndpoint, c.ProjectID, c.LocationID, c.ModelID,
	)

	payload := c.generateRequestPayload(prompt, nil, nil)

	response, err := c.getResponseFromAPI(ctx, url, payload)
	if err != nil {
		return "", err
	}

	return concatenateAllEntries(response), nil
}

// SendMessageUsingDatastore sends a request to the Vertex AI endpoint and returns the response.
func (c *VertexAIClient) SendMessageUsingDatastore(ctx context.Context, prompt []map[string]string, datastore []string) (string, error) {
	url := fmt.Sprintf(
		"https://%s/v1/projects/%s/locations/%s/publishers/google/models/%s:streamGenerateContent",
		c.APIEndpoint, c.ProjectID, c.LocationID, c.ModelID,
	)

	payload := c.generateRequestPayload(prompt, datastore, nil)

	response, err := c.getResponseFromAPI(ctx, url, payload)
	if err != nil {
		return "", err
	}

	return concatenateAllEntries(response), nil
}

// SendMessageUsingSystemInstruction sends a request to the Vertex AI endpoint and returns the response.
func (c *VertexAIClient) SendMessageUsingSystemInstruction(ctx context.Context, prompt []map[string]string, systemInstruction []string) (string, error) {
	url := fmt.Sprintf(
		"https://%s/v1/projects/%s/locations/%s/publishers/google/models/%s:streamGenerateContent",
		c.APIEndpoint, c.ProjectID, c.LocationID, c.ModelID,
	)

	payload := c.generateRequestPayload(prompt, nil, systemInstruction)

	response, err := c.getResponseFromAPI(ctx, url, payload)
	if err != nil {
		return "", err
	}

	return concatenateAllEntries(response), nil
}

// SendMessageUsingDatastoreAndSystemInstruction sends a request to the Vertex AI endpoint and returns the response.
func (c *VertexAIClient) SendMessageUsingDatastoreAndSystemInstruction(ctx context.Context, prompt []map[string]string, datastore []string, systemInstruction []string) (string, error) {
	url := fmt.Sprintf(
		"https://%s/v1/projects/%s/locations/%s/publishers/google/models/%s:streamGenerateContent",
		c.APIEndpoint, c.ProjectID, c.LocationID, c.ModelID,
	)

	payload := c.generateRequestPayload(prompt, datastore, systemInstruction)

	response, err := c.getResponseFromAPI(ctx, url, payload)
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
