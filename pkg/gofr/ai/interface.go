package ai

type VertexAI interface {
	GetResponse(prompt string, datastore ...string) (string, error)

	UseLogger(logger any)
}
