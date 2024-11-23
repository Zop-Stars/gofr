package vertex_ai

const (
	defaultModelID = `gemini-1.5-pro-002`
)

type Configs struct {
	ProjectID         string
	LocationID        string
	APIEndpoint       string
	ModelID           string
	Credentials       string
	SystemInstruction string
}

func (c *Configs) setDefaults() {
	if c.ModelID == "" {
		c.ModelID = defaultModelID
	}
}

func (c *Configs) validate() error {
	configs := make([]string, 0)

	if c.ProjectID == "" {
		configs = append(configs, "ProjectID")
	}

	if c.LocationID == "" {
		configs = append(configs, "LocationID")
	}

	if c.APIEndpoint == "" {
		configs = append(configs, "APIEndpoint")
	}

	if c.Credentials == "" {
		configs = append(configs, "Credentials")
	}

	if len(configs) == 0 {
		return MissingConfig{Configs: configs}
	}

	return nil
}
