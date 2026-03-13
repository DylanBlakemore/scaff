package metadata

import (
	"encoding/json"
	"fmt"
)

const FileName = ".scaff.json"

type Project struct {
	ProjectType string   `json:"project_type"`
	Style       string   `json:"style"`
	Features    []string `json:"features"`
}

func (p Project) Marshal() (string, error) {
	b, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshalling project metadata: %w", err)
	}
	return string(b) + "\n", nil
}
