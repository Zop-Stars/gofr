package ai

type VertexAI interface {
	GetResponse(prompt string) (string, error)
	GetResponseUsingDatastore(prompt string, datastore []string) (string, error)
	GetResponseUsingSystemInstruction(prompt string, systemInstruction []string) (string, error)
	GetResponseUsingDatastoreAndSystemInstruction(prompt string, datastore []string, systemInstruction []string) (string, error)

	UseLogger(logger any)
}
