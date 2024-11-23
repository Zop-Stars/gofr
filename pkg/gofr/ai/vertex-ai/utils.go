package vertex_ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (c *VertexAIClient) getResponseFromAPI(url string, payload *RequestPayload) ([]DataEntry, error) {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}

	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error response from Vertex AI: %s", string(responseBody))
	}

	entries := make([]DataEntry, 0)
	//decoder := json.NewDecoder(resp.Body)

	err = json.Unmarshal(responseBody, &entries)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	//for {
	//	var chunk DataEntry
	//
	//	if err := decoder.Decode(&chunk); err == io.EOF {
	//		break
	//	} else if err != nil {
	//		return nil, fmt.Errorf("error decoding response chunk: %w", err)
	//	}
	//	entries = append(entries, chunk)
	//}

	return entries, nil
}

func (c *VertexAIClient) generateRequestPayload(prompt string, datastores ...string) *RequestPayload {
	payload := &RequestPayload{
		Contents: []Message{
			{
				Role: "user",
				Parts: []Part{
					{Text: prompt},
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
	}

	if len(datastores) > 0 && len(datastores[0]) > 0 {
		payload.Tools = c.generateDatastoreForPayload(datastores)
	}

	if len(c.configs.SystemInstruction) != 0 {
		payload.SystemInstruction = c.generateSystemInstructionForPayload(strings.Split(c.configs.SystemInstruction, ","))

	}

	return payload
}

func (c *VertexAIClient) generateDatastoreForPayload(datastores []string) []Tool {
	tools := make([]Tool, 0)

	for _, ds := range datastores {
		tools = append(tools, Tool{
			Retrieval: Retrieval{
				VertexAiSearch: VertexAiSearch{
					Datastore: ds,
				},
			},
		})
	}

	return tools
}

func (c *VertexAIClient) generateSystemInstructionForPayload(systemInstructs []string) SystemInstruction {
	parts := make([]Part, 0)

	for _, si := range systemInstructs {
		parts = append(parts, Part{Text: si})
	}

	return SystemInstruction{Parts: parts}
}
