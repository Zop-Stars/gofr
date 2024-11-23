package gofr

import "gofr.dev/pkg/gofr/ai"

func (a *App) UseVertexAI(vertexAIClient ai.VertexAI) {
	vertexAIClient.UseLogger(a.Logger())

	a.container.VertexAI = vertexAIClient
}
