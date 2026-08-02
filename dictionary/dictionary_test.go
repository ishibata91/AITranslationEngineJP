package main

import (
	"context"
	"encoding/json"
	"errors"
	"path/filepath"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	_ "modernc.org/sqlite"
)

func TestImportMasterTermsKeepsRowsAndExcludesProperNoun(t *testing.T) {
	ctx := context.Background()
	sourcePath := filepath.Join(t.TempDir(), "central.sqlite3")
	source := openSourceDB(t, sourcePath)
	execTestSQL(t, source, `
		CREATE TABLE master_term (id INTEGER PRIMARY KEY, source TEXT NOT NULL, dest TEXT NOT NULL, category TEXT NOT NULL);
		CREATE TABLE proper_noun (id INTEGER PRIMARY KEY, plugin TEXT NOT NULL, source TEXT NOT NULL, dest TEXT NOT NULL, category TEXT NOT NULL);
		INSERT INTO master_term VALUES
			(1, 'Imperial', 'インペリアル', 'RACE'),
			(2, 'Imperial', 'インペリアル', 'NPC_'),
			(3, 'Solitude', 'ソリチュード', 'LCTN');
		INSERT INTO proper_noun VALUES
			(1, 'inigo.esp', 'Inigo', 'イニゴ', 'NPC_');
	`)
	source.Close()

	destination := openTestStore(t)
	first, err := importMasterTerms(ctx, sourcePath, destination)
	if err != nil {
		t.Fatal(err)
	}
	if first.Read != 3 || first.Created != 3 || first.Updated != 0 || first.Unchanged != 0 {
		t.Fatalf("first import = %+v", first)
	}

	second, err := importMasterTerms(ctx, sourcePath, destination)
	if err != nil {
		t.Fatal(err)
	}
	if second.Read != 3 || second.Created != 0 || second.Updated != 0 || second.Unchanged != 3 {
		t.Fatalf("second import = %+v", second)
	}

	got, err := destination.search(ctx, "Imperial", "", 50)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Entries) != 2 || got.Entries[0].Source != "Imperial" || got.Entries[1].Source != "Imperial" {
		t.Fatalf("Imperial search = %+v", got.Entries)
	}
	status, err := destination.status(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if status.Entries != 3 || status.FTSEntries != 3 || status.Origins["master_term"] != 3 || status.Origins["proper_noun"] != 0 {
		t.Fatalf("status = %+v", status)
	}
}

func TestSearchAddAndRevisionUpdate(t *testing.T) {
	ctx := context.Background()
	s := openTestStore(t)
	imperialRace, err := s.add(ctx, "Imperial", "インペリアル", "RACE")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.add(ctx, "Imperial Legion", "帝国軍", "FACT"); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		query string
		want  string
		kind  string
	}{
		{"imperial", "Imperial", "source_exact"},
		{"Imperial L", "Imperial Legion", "source_prefix"},
		{"perial Leg", "Imperial Legion", "substring"},
		{"ペリア", "Imperial", "substring"},
	}
	for _, tt := range tests {
		t.Run(tt.query, func(t *testing.T) {
			got, err := s.search(ctx, tt.query, "", 50)
			if err != nil {
				t.Fatal(err)
			}
			if len(got.Entries) == 0 || got.Entries[0].Source != tt.want || got.Entries[0].MatchKind != tt.kind {
				t.Fatalf("search(%q) = %+v", tt.query, got.Entries)
			}
		})
	}

	updated, err := s.update(ctx, imperialRace.ID, imperialRace.Revision, "Imperial", "インペリアル種族", "RACE")
	if err != nil {
		t.Fatal(err)
	}
	if updated.Revision != imperialRace.Revision+1 || updated.Dest != "インペリアル種族" {
		t.Fatalf("updated = %+v", updated)
	}
	if _, err := s.update(ctx, imperialRace.ID, imperialRace.Revision, "Imperial", "古い更新", "RACE"); !errors.Is(err, errRevisionConflict) {
		t.Fatalf("stale update error = %v", err)
	}
}

func TestMCPServerExposesAndCallsDictionaryTools(t *testing.T) {
	ctx := context.Background()
	s := openTestStore(t)
	if _, err := s.add(ctx, "Imperial", "インペリアル", "RACE"); err != nil {
		t.Fatal(err)
	}

	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	serverSession, err := newMCPServer(s).Connect(ctx, serverTransport, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer serverSession.Close()
	client := mcp.NewClient(&mcp.Implementation{Name: "dictionary-test", Version: "0.1.0"}, nil)
	clientSession, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer clientSession.Close()

	tools, err := clientSession.ListTools(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(tools.Tools) != 5 {
		t.Fatalf("tools = %d, want 5", len(tools.Tools))
	}
	result, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name: "dictionary_search",
		Arguments: map[string]any{
			"query": "Imperial",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	raw, err := json.Marshal(result.StructuredContent)
	if err != nil {
		t.Fatal(err)
	}
	var got searchResult
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatal(err)
	}
	if len(got.Entries) != 1 || got.Entries[0].Source != "Imperial" {
		t.Fatalf("MCP search = %+v", got.Entries)
	}
}

func openSourceDB(t *testing.T, path string) *sqlx.DB {
	t.Helper()
	db, err := sqlx.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func openTestStore(t *testing.T) *store {
	t.Helper()
	s, err := openStore(filepath.Join(t.TempDir(), "dictionary.sqlite3"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { s.close() })
	return s
}

func execTestSQL(t *testing.T, db *sqlx.DB, query string) {
	t.Helper()
	if _, err := db.Exec(query); err != nil {
		t.Fatal(err)
	}
}
