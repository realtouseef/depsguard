package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Baseline struct {
	Dependencies map[string]string `json:"dependencies"`
	LastUpdated  string            `json:"last_updated"`
}

func DepsguardDir() string {
	return ".depsguard"
}

func EnsureDepsguardDir() error {
	if _, err := os.Stat(DepsguardDir()); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("depsguard not initialized; run depsguard init")
		}
		return err
	}
	return nil
}

func NewBaseline(deps map[string]string) Baseline {
	copy := make(map[string]string, len(deps))
	for name, version := range deps {
		copy[name] = version
	}
	return Baseline{
		Dependencies: copy,
		LastUpdated:  time.Now().Format(time.RFC3339),
	}
}

func LoadBaseline(path string) (Baseline, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Baseline{}, fmt.Errorf("baseline.json not found; run depsguard init")
		}
		return Baseline{}, err
	}
	var base Baseline
	if err := json.Unmarshal(data, &base); err != nil {
		return Baseline{}, fmt.Errorf("failed to parse baseline.json: %w; delete .depsguard/baseline.json and rerun depsguard init", err)
	}
	if base.Dependencies == nil {
		base.Dependencies = map[string]string{}
	}
	return base, nil
}

func SaveBaseline(path string, base Baseline) error {
	data, err := json.MarshalIndent(base, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func NewDependencies(oldDeps, newDeps map[string]string) []string {
	added := make([]string, 0)
	for name := range newDeps {
		if _, ok := oldDeps[name]; !ok {
			added = append(added, name)
		}
	}
	return added
}
