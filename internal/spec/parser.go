package spec

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

func Parse(data []byte) (*Profile, error) {
	var p Profile
	if err := yaml.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("yaml parse: %w", err)
	}
	if err := Validate(&p); err != nil {
		return nil, err
	}
	return &p, nil
}
