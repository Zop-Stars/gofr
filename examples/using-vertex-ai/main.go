package main

import (
	"encoding/json"
	"gofr.dev/pkg/gofr"
)

func main() {
	// Create a new application
	a := gofr.New()

	//HTTP service with default health check endpoint
	a.AddHTTPService("anotherService", "http://localhost:9000")

	creds := a.Config.Get("SVC_ACC_CREDS")

	vertexAIClient, err := NewVertexAIClientWithKey("endless-fire-437206-j7", "us-central1",
		"us-central1-aiplatform.googleapis.com", "gemini-1.5-pro-002", creds)
	if err != nil {
		a.Logger().Fatal(err)
	}

	a.POST("/chat", func(c *gofr.Context) (interface{}, error) {
		var prompt struct {
			Prompt string `json:"prompt"`
		}

		err = c.Bind(&prompt)
		if err != nil {
			return nil, err
		}

		payload := RequestPayload{
			Contents: []Message{
				{
					Role: "user",
					Parts: []Part{
						{Text: prompt.Prompt},
					},
				},
			},
			GenerationConfig: GenerationConfig{
				Temperature:     1.0,
				MaxOutputTokens: 8192,
				TopP:            0.95,
			},
			SafetySettings: []SafetySetting{
				{Category: "HARM_CATEGORY_HATE_SPEECH", Threshold: "OFF"},
				{Category: "HARM_CATEGORY_DANGEROUS_CONTENT", Threshold: "OFF"},
				{Category: "HARM_CATEGORY_SEXUALLY_EXPLICIT", Threshold: "OFF"},
				{Category: "HARM_CATEGORY_HARASSMENT", Threshold: "OFF"},
			},
			Tools: []Tool{
				{
					Retrieval: Retrieval{
						VertexAiSearch: VertexAiSearch{
							Datastore: "projects/endless-fire-437206-j7/locations/global/collections/default_collection/dataStores/gofr-datastore_1732298621027",
						},
					},
				},
			},
		}

		respString, err := vertexAIClient.GenerateContent(payload)
		if err != nil {
			return nil, err
		}

		//return response.Raw{Data: struct {
		//	Response string `json:"response"`
		//}{resp}}, nil

		resp := make([]DataEntry, 0)

		err = json.Unmarshal([]byte(respString), &resp)
		if err != nil {
			return nil, err
		}

		return ConcatenateAllEntries(resp), nil
	})

	// Run the application
	a.Run()
}
