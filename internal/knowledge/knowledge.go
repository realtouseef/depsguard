package knowledge

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Entry struct {
	Summary     string `json:"summary"`
	ExplainedBy string `json:"explained_by"`
	ExpiresAt   string `json:"expires_at"`
}

func (e Entry) IsValid(now time.Time) bool {
	if e.Summary == "" || e.ExplainedBy == "" || e.ExpiresAt == "" {
		return false
	}
	timestamp, err := time.Parse(time.RFC3339, e.ExpiresAt)
	if err != nil {
		return false
	}
	return timestamp.After(now)
}

func Load(path string) (map[string]Entry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]Entry{}, nil
		}
		return nil, err
	}
	if len(data) == 0 {
		return map[string]Entry{}, nil
	}
	entries := map[string]Entry{}
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("failed to parse knowledge.json: %w", err)
	}
	return entries, nil
}

func Save(path string, entries map[string]Entry) error {
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
