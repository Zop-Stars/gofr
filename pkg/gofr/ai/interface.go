package ai

type VertexAI interface {
	GetResponse(prompt []map[string]string) (string, error)
	GetResponseUsingDatastore(prompt []map[string]string, datastore []string) (string, error)
	GetResponseUsingSystemInstruction(prompt []map[string]string, systemInstruction []string) (string, error)
	GetResponseUsingDatastoreAndSystemInstruction(prompt []map[string]string, datastore []string, systemInstruction []string) (string, error)

	UseLogger(logger any)
}
