// Package social provides squadron management, leaderboards, and federation identity.
//
// TODO(v5.0): This package is a stub. Full implementation planned for v5.0.
// See ROADMAP.md v5.0 milestone and GAPS.md "v5.0+ Features Are Stubs" section.
package social

// Squadron represents a named group of players.
type Squadron struct {
	Name    string
	Members []string
}

// NewSquadron creates a new squadron with the given name.
func NewSquadron(name string) *Squadron {
	return &Squadron{Name: name}
}

// AddMember adds a player to the squadron.
func (s *Squadron) AddMember(playerID string) {
	s.Members = append(s.Members, playerID)
}

// LeaderboardEntry represents a single high-score entry.
type LeaderboardEntry struct {
	PlayerID string
	Score    int64
	Seed     int64
	Genre    string
}

// Leaderboard holds ranked score entries.
type Leaderboard struct {
	Entries []LeaderboardEntry
}

// NewLeaderboard creates an empty leaderboard.
func NewLeaderboard() *Leaderboard {
	return &Leaderboard{}
}

// Submit adds a new entry to the leaderboard.
func (lb *Leaderboard) Submit(entry LeaderboardEntry) {
	lb.Entries = append(lb.Entries, entry)
}
