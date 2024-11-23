package ai

import "context"

type VertexAI interface {
	SendMessage(ctx context.Context, prompt []map[string]string) (string, error)
	SendMessageUsingDatastore(ctx context.Context, prompt []map[string]string, datastore []string) (string, error)
	SendMessageUsingSystemInstruction(ctx context.Context, prompt []map[string]string, systemInstruction []string) (string, error)
	SendMessageUsingDatastoreAndSystemInstruction(ctx context.Context, prompt []map[string]string, datastore []string, systemInstruction []string) (string, error)

	UseLogger(logger any)
}
