package ai

type VertexAI interface {
	SendMessage(prompt []map[string]string) (string, error)
	SendMessageUsingDatastore(prompt []map[string]string, datastore []string) (string, error)
	SendMessageUsingSystemInstruction(prompt []map[string]string, systemInstruction []string) (string, error)
	SendMessageUsingDatastoreAndSystemInstruction(prompt []map[string]string, datastore []string, systemInstruction []string) (string, error)

	UseLogger(logger any)
}
