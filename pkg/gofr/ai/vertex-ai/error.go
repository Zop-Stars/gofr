package vertex_ai

import (
	"fmt"
	"strings"
)

type MissingConfig struct {
	Configs []string
}

func (e MissingConfig) Error() string {
	if len(e.Configs) == 1 {
		return fmt.Sprintf("missing config %s", e.Configs[0])
	} else {
		return fmt.Sprintf("missing config %s", strings.Join(e.Configs, ", "))
	}
}
