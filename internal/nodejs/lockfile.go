package nodejs

import (
	"encoding/json"
	"fmt"
	"os"
)

type Lockfile struct {
	Dependencies map[string]LockDependency `json:"dependencies"`
}

type LockDependency struct {
	Version string `json:"version"`
}

func LoadPackageLock(path string) (Lockfile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Lockfile{}, err
	}
	var lock Lockfile
	if err := json.Unmarshal(data, &lock); err != nil {
		return Lockfile{}, fmt.Errorf("failed to parse package-lock.json: %w", err)
	}
	return lock, nil
}
