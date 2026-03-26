package social

import "testing"

func TestNewSquadron(t *testing.T) {
	s := NewSquadron("Alpha")
	if s == nil {
		t.Fatal("NewSquadron() returned nil")
	}
	if s.Name != "Alpha" {
		t.Errorf("Name = %q, want %q", s.Name, "Alpha")
	}
	if len(s.Members) != 0 {
		t.Errorf("new squadron should have 0 members, got %d", len(s.Members))
	}
}

func TestSquadronAddMember(t *testing.T) {
	s := NewSquadron("Bravo")

	s.AddMember("player1")
	if len(s.Members) != 1 {
		t.Fatalf("expected 1 member, got %d", len(s.Members))
	}
	if s.Members[0] != "player1" {
		t.Errorf("Members[0] = %q, want %q", s.Members[0], "player1")
	}

	s.AddMember("player2")
	s.AddMember("player3")

	if len(s.Members) != 3 {
		t.Fatalf("expected 3 members, got %d", len(s.Members))
	}

	expectedMembers := []string{"player1", "player2", "player3"}
	for i, expected := range expectedMembers {
		if s.Members[i] != expected {
			t.Errorf("Members[%d] = %q, want %q", i, s.Members[i], expected)
		}
	}
}

func TestSquadronStruct(t *testing.T) {
	s := &Squadron{
		Name:    "Charlie",
		Members: []string{"a", "b", "c"},
	}

	if s.Name != "Charlie" {
		t.Errorf("Name = %q, want %q", s.Name, "Charlie")
	}
	if len(s.Members) != 3 {
		t.Errorf("len(Members) = %d, want 3", len(s.Members))
	}
}

func TestNewLeaderboard(t *testing.T) {
	lb := NewLeaderboard()
	if lb == nil {
		t.Fatal("NewLeaderboard() returned nil")
	}
	if len(lb.Entries) != 0 {
		t.Errorf("new leaderboard should have 0 entries, got %d", len(lb.Entries))
	}
}

func TestLeaderboardSubmit(t *testing.T) {
	lb := NewLeaderboard()

	entry1 := LeaderboardEntry{
		PlayerID: "player1",
		Score:    10000,
		Seed:     12345,
		Genre:    "scifi",
	}
	lb.Submit(entry1)

	if len(lb.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(lb.Entries))
	}
	if lb.Entries[0].PlayerID != "player1" {
		t.Errorf("Entries[0].PlayerID = %q, want %q", lb.Entries[0].PlayerID, "player1")
	}
	if lb.Entries[0].Score != 10000 {
		t.Errorf("Entries[0].Score = %d, want 10000", lb.Entries[0].Score)
	}
}

func TestLeaderboardSubmitMultiple(t *testing.T) {
	lb := NewLeaderboard()

	entries := []LeaderboardEntry{
		{PlayerID: "player1", Score: 10000, Seed: 12345, Genre: "scifi"},
		{PlayerID: "player2", Score: 15000, Seed: 12345, Genre: "scifi"},
		{PlayerID: "player3", Score: 8000, Seed: 12345, Genre: "fantasy"},
	}

	for _, e := range entries {
		lb.Submit(e)
	}

	if len(lb.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(lb.Entries))
	}

	for i, expected := range entries {
		if lb.Entries[i].PlayerID != expected.PlayerID {
			t.Errorf("Entries[%d].PlayerID = %q, want %q", i, lb.Entries[i].PlayerID, expected.PlayerID)
		}
		if lb.Entries[i].Score != expected.Score {
			t.Errorf("Entries[%d].Score = %d, want %d", i, lb.Entries[i].Score, expected.Score)
		}
	}
}

func TestLeaderboardEntryStruct(t *testing.T) {
	entry := LeaderboardEntry{
		PlayerID: "test-player",
		Score:    99999,
		Seed:     54321,
		Genre:    "horror",
	}

	if entry.PlayerID != "test-player" {
		t.Errorf("PlayerID = %q, want %q", entry.PlayerID, "test-player")
	}
	if entry.Score != 99999 {
		t.Errorf("Score = %d, want 99999", entry.Score)
	}
	if entry.Seed != 54321 {
		t.Errorf("Seed = %d, want 54321", entry.Seed)
	}
	if entry.Genre != "horror" {
		t.Errorf("Genre = %q, want %q", entry.Genre, "horror")
	}
}

func TestLeaderboardStruct(t *testing.T) {
	lb := &Leaderboard{
		Entries: []LeaderboardEntry{
			{PlayerID: "p1", Score: 100},
			{PlayerID: "p2", Score: 200},
		},
	}

	if len(lb.Entries) != 2 {
		t.Errorf("len(Entries) = %d, want 2", len(lb.Entries))
	}
}

func TestNewSquadronDifferentNames(t *testing.T) {
	names := []string{"", "Alpha", "The Best Squad", "日本語", "emoji-🚀"}

	for _, name := range names {
		s := NewSquadron(name)
		if s.Name != name {
			t.Errorf("NewSquadron(%q).Name = %q", name, s.Name)
		}
	}
}

func TestLeaderboardSubmitDuplicatePlayer(t *testing.T) {
	lb := NewLeaderboard()

	// Same player, different scores (should both be added - not replacing)
	lb.Submit(LeaderboardEntry{PlayerID: "player1", Score: 100})
	lb.Submit(LeaderboardEntry{PlayerID: "player1", Score: 200})

	if len(lb.Entries) != 2 {
		t.Errorf("expected 2 entries (no deduplication), got %d", len(lb.Entries))
	}
}
