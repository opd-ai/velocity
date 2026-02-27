// Package saveload provides save/load serialization and run-state snapshots.
package saveload

import (
	"encoding/json"
	"os"
)

// RunState represents the serializable state of a game run.
type RunState struct {
	Version    int    `json:"version"`
	Seed       int64  `json:"seed"`
	Genre      string `json:"genre"`
	Wave       int    `json:"wave"`
	Score      int64  `json:"score"`
	PlayerData []byte `json:"player_data,omitempty"`
}

// Save writes the run state to a file.
func Save(path string, state *RunState) error {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// Load reads a run state from a file.
func Load(path string) (*RunState, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var state RunState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}
	return &state, nil
}
