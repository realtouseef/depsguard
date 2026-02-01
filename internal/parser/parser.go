package parser

import (
	"encoding/json"
	"fmt"
	"os"
)

type PackageJSON struct {
	Dependencies         map[string]string `json:"dependencies"`
	DevDependencies      map[string]string `json:"devDependencies"`
	PeerDependencies     map[string]string `json:"peerDependencies"`
	OptionalDependencies map[string]string `json:"optionalDependencies"`
}

func LoadPackageJSON(path string) (PackageJSON, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return PackageJSON{}, err
	}
	var pkg PackageJSON
	if err := json.Unmarshal(data, &pkg); err != nil {
		return PackageJSON{}, fmt.Errorf("failed to parse package.json: %w", err)
	}
	return pkg, nil
}

func (p PackageJSON) AllDependencies() map[string]string {
	merged := make(map[string]string)
	for name, version := range p.Dependencies {
		merged[name] = version
	}
	for name, version := range p.DevDependencies {
		merged[name] = version
	}
	return merged
}
