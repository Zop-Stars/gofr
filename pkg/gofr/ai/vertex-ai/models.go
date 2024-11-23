package vertex_ai

import (
	"strings"
)

// RequestPayload defines the structure of the JSON request to Vertex AI.
type RequestPayload struct {
	Contents          []Message         `json:"contents"`
	GenerationConfig  GenerationConfig  `json:"generationConfig"`
	SafetySettings    []SafetySetting   `json:"safetySettings"`
	SystemInstruction SystemInstruction `json:"systemInstruction"`
	Tools             []Tool            `json:"tools"`
}

// Message represents the content of a conversation.
type Message struct {
	Role  string `json:"role"`
	Parts []Part `json:"parts"`
}

type SystemInstruction struct {
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

type RetrievedContext struct {
	URI   string `json:"uri"`
	Title string `json:"title"`
	Text  string `json:"text"`
}

type GroundingChunk struct {
	RetrievedContext RetrievedContext `json:"retrievedContext"`
}

type GroundingMetadata struct {
	GroundingChunks []GroundingChunk `json:"groundingChunks"`
}

type Content struct {
	Role  string `json:"role"`
	Parts []struct {
		Text string `json:"text"`
	} `json:"parts"`
}

type Candidate struct {
	Content           Content           `json:"content"`
	FinishReason      string            `json:"finishReason,omitempty"`
	GroundingMetadata GroundingMetadata `json:"groundingMetadata,omitempty"`
}

type DataEntry struct {
	Candidates   []Candidate `json:"candidates"`
	ModelVersion string      `json:"modelVersion"`
}

func (d *DataEntry) ConcatenateParts() string {
	var builder strings.Builder

	for _, candidate := range d.Candidates {
		for _, part := range candidate.Content.Parts {
			builder.WriteString(part.Text)
		}
	}

	return builder.String()
}

func concatenateAllEntries(entries []DataEntry) string {
	var builder strings.Builder

	for _, entry := range entries {
		builder.WriteString(entry.ConcatenateParts())
	}

	return builder.String()
}
