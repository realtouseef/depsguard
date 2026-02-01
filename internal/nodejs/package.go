package nodejs

import (
	"encoding/json"
	"fmt"
	"os"
)

type PackageJSON struct {
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

func LoadPackageJSON(path string) (PackageJSON, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return PackageJSON{}, fmt.Errorf("package.json not found; run depsguard from a Node.js project root")
		}
		return PackageJSON{}, err
	}
	var pkg PackageJSON
	if err := json.Unmarshal(data, &pkg); err != nil {
		return PackageJSON{}, fmt.Errorf("failed to parse package.json: %w", err)
	}
	return pkg, nil
}

func (p PackageJSON) GetAllDependencies() map[string]string {
	all := make(map[string]string)
	for name, version := range p.Dependencies {
		all[name] = version
	}
	for name, version := range p.DevDependencies {
		all[name] = version
	}
	return all
}
