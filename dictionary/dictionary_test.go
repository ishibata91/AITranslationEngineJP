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

func TestImportMasterTermsCreatesTermsSensesAndOccurrences(t *testing.T) {
	ctx := context.Background()
	sourcePath := filepath.Join(t.TempDir(), "central.sqlite3")
	source := openSourceDB(t, sourcePath)
	execTestSQL(t, source, `
		CREATE TABLE master_term (id INTEGER PRIMARY KEY, source TEXT NOT NULL, dest TEXT NOT NULL, category TEXT NOT NULL);
		CREATE TABLE proper_noun (id INTEGER PRIMARY KEY, plugin TEXT NOT NULL, source TEXT NOT NULL, dest TEXT NOT NULL, category TEXT NOT NULL);
		INSERT INTO master_term VALUES
			(1, 'Imperial', 'インペリアル', 'RACE'),
			(2, 'Imperial', 'インペリアル', 'NPC_'),
			(3, 'Wight', 'ワイト', 'derive:two');
		INSERT INTO proper_noun VALUES
			(1, 'inigo.esp', 'Inigo', 'イニゴ', 'NPC_');
	`)
	if err := source.Close(); err != nil {
		t.Fatal(err)
	}

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
	got, err := destination.search(ctx, searchFilter{Query: "Imperial", Limit: 50})
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Entries) != 2 {
		t.Fatalf("Imperial search = %+v", got.Entries)
	}
	gotStatus, err := destination.status(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if gotStatus.Terms != 2 || gotStatus.Senses != 3 || gotStatus.Occurrences != 3 ||
		gotStatus.AssignedOccurrences != 3 || gotStatus.Origins["master_term"] != 3 ||
		gotStatus.ClassificationStatuses["unclassified"] != 3 {
		t.Fatalf("status = %+v", gotStatus)
	}
	var derived occurrence
	if err := destination.db.GetContext(ctx, &derived, `
		SELECT id, term_id, sense_id, observed_dest, skyrim_category,
		       origin_kind, origin_reference, derivation_kind
		FROM dictionary_occurrence WHERE origin_reference = '3'`); err != nil {
		t.Fatal(err)
	}
	if derived.SkyrimCategory != "" || derived.DerivationKind != "two" {
		t.Fatalf("derived occurrence = %+v", derived)
	}
}

func TestOpenStoreMigratesLegacyEntriesIntoSeparatedTables(t *testing.T) {
	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "legacy.sqlite3")
	legacy := openSourceDB(t, path)
	execTestSQL(t, legacy, `
		CREATE TABLE dictionary_entry (
			id INTEGER PRIMARY KEY, source TEXT NOT NULL, dest TEXT NOT NULL,
			category TEXT NOT NULL, revision INTEGER NOT NULL
		);
		CREATE TABLE dictionary_entry_source (
			entry_id INTEGER NOT NULL, kind TEXT NOT NULL, reference TEXT NOT NULL
		);
		INSERT INTO dictionary_entry VALUES
			(7, 'Imperial', 'インペリアル', 'RACE', 3),
			(8, 'Ancient', 'エンシェント', 'derive:two', 1);
		INSERT INTO dictionary_entry_source VALUES
			(7, 'master_term', '7'),
			(8, 'master_term', '8');
	`)
	if err := legacy.Close(); err != nil {
		t.Fatal(err)
	}
	s, err := openStore(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if closeErr := s.close(); closeErr != nil {
			t.Error(closeErr)
		}
	})
	got, err := s.get(ctx, 7)
	if err != nil {
		t.Fatal(err)
	}
	if got.Source != "Imperial" || got.Revision != 3 || len(got.Occurrences) != 1 || got.Occurrences[0].SkyrimCategory != "RACE" {
		t.Fatalf("migrated sense = %+v", got)
	}
	derived, err := s.get(ctx, 8)
	if err != nil {
		t.Fatal(err)
	}
	if len(derived.Occurrences) != 1 || derived.Occurrences[0].SkyrimCategory != "" || derived.Occurrences[0].DerivationKind != "two" {
		t.Fatalf("migrated derivation = %+v", derived)
	}
}

func TestSearchAddUpdateAndHistory(t *testing.T) {
	ctx := context.Background()
	s := openTestStore(t)
	added, err := s.add(ctx, "Imperial", "インペリアル", "noun", "Skyrimの種族")
	if err != nil {
		t.Fatal(err)
	}
	if _, addErr := s.add(ctx, "Imperial Legion", "帝国軍", "noun", "Skyrimの軍事組織"); addErr != nil {
		t.Fatal(addErr)
	}
	tests := []struct {
		query string
		want  string
		kind  string
	}{
		{"imperial", "Imperial", "source_exact"},
		{"Imperial L", "Imperial Legion", "source_prefix"},
		{"perial Leg", "Imperial Legion", "substring"},
		{"軍事組織", "Imperial Legion", "substring"},
	}
	for _, test := range tests {
		t.Run(test.query, func(t *testing.T) {
			got, searchErr := s.search(ctx, searchFilter{Query: test.query, Limit: 50})
			if searchErr != nil {
				t.Fatal(searchErr)
			}
			if len(got.Entries) == 0 || got.Entries[0].Source != test.want || got.Entries[0].MatchKind != test.kind {
				t.Fatalf("search(%q) = %+v", test.query, got.Entries)
			}
		})
	}
	updated, err := s.update(ctx, senseUpdate{
		ID: added.ID, Revision: added.Revision, Dest: "インペリアル種族",
		PartOfSpeech: "noun", Meaning: "Skyrimのプレイヤー種族",
		ChangedBy: "test-human", Reason: "種族の意味を確認した",
	})
	if err != nil {
		t.Fatal(err)
	}
	if updated.Revision != added.Revision+1 || updated.Dest != "インペリアル種族" {
		t.Fatalf("updated = %+v", updated)
	}
	reviewed, err := s.addReview(ctx, reviewInput{
		SenseID: updated.ID, Revision: updated.Revision, ReviewerKind: "human",
		ReviewerReference: "test-human", Decision: "include", Reason: "Skyrim固有の種族名として収録する",
	})
	if err != nil {
		t.Fatal(err)
	}
	if reviewed.InclusionDecision != "include" || reviewed.ReviewStage != "human_reviewed" || len(reviewed.Reviews) != 1 {
		t.Fatalf("reviewed = %+v", reviewed)
	}
	if _, updateErr := s.update(ctx, senseUpdate{
		ID: added.ID, Revision: added.Revision, Dest: added.Dest,
		PartOfSpeech: "noun", Meaning: added.Meaning,
		ChangedBy: "test-human", Reason: "古い更新",
	}); !errors.Is(updateErr, errRevisionConflict) {
		t.Fatalf("stale update error = %v", updateErr)
	}
	history, err := s.history(ctx, "dictionary_sense", added.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(history.Changes) != 4 {
		t.Fatalf("history = %+v", history.Changes)
	}
}

func TestClassifyGeneralDictionarySavesCandidatesWithoutExcluding(t *testing.T) {
	ctx := context.Background()
	s := openTestStore(t)
	for _, item := range []struct{ source, dest string }{
		{"Fence", "柵"}, {"Fence", "盗品売人"}, {"Frost", "フロスト"},
		{"Imperial", "インペリアル"}, {"Solitude", "ソリチュード"},
	} {
		if _, err := s.add(ctx, item.source, item.dest, "unknown", ""); err != nil {
			t.Fatal(err)
		}
	}
	wordNetPath := createTestWordNet(t)
	first, err := s.classifyGeneralDictionary(ctx, wordNetPath)
	if err != nil {
		t.Fatal(err)
	}
	if first.Senses != 5 || first.Statuses["same_mean_candidate"] != 2 ||
		first.Statuses["same_surface_only"] != 2 || first.Statuses["no_english_headword"] != 1 {
		t.Fatalf("classification = %+v", first)
	}
	var frostMatches int
	if queryErr := s.db.GetContext(ctx, &frostMatches, `
		SELECT COUNT(*) FROM general_dictionary_match m
		JOIN dictionary_sense s ON s.id = m.sense_id
		JOIN dictionary_term t ON t.id = s.term_id
		WHERE t.source = 'Frost' AND m.match_status = 'same_mean_candidate'`); queryErr != nil {
		t.Fatal(queryErr)
	}
	if frostMatches != 2 {
		t.Fatalf("Frost candidates = %d, want 2", frostMatches)
	}
	var excluded int
	if queryErr := s.db.GetContext(ctx, &excluded, `
		SELECT COUNT(*) FROM dictionary_sense WHERE inclusion_decision = 'exclude'`); queryErr != nil {
		t.Fatal(queryErr)
	}
	if excluded != 0 {
		t.Fatalf("automatic exclusions = %d, want 0", excluded)
	}
	second, err := s.classifyGeneralDictionary(ctx, wordNetPath)
	if err != nil {
		t.Fatal(err)
	}
	if second.Matches != first.Matches {
		t.Fatalf("second matches = %d, first = %d", second.Matches, first.Matches)
	}
}

func TestGeneralMatchAndOccurrenceChangesKeepReasons(t *testing.T) {
	ctx := context.Background()
	s := openTestStore(t)
	fence, err := s.add(ctx, "Fence", "柵", "noun", "区域を囲う構造物")
	if err != nil {
		t.Fatal(err)
	}
	secondSense, err := s.add(ctx, "Fence", "盗品売人", "noun", "盗品を扱う人物")
	if err != nil {
		t.Fatal(err)
	}
	wordNetPath := createTestWordNet(t)
	if _, classifyErr := s.classifyGeneralDictionary(ctx, wordNetPath); classifyErr != nil {
		t.Fatal(classifyErr)
	}
	detail, err := s.get(ctx, fence.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(detail.GeneralMatches) == 0 || len(detail.Occurrences) == 0 {
		t.Fatalf("detail = %+v", detail)
	}
	if _, reviewErr := s.addReview(ctx, reviewInput{
		SenseID: detail.ID, Revision: detail.Revision, ReviewerKind: "ai",
		ReviewerReference: "test-ai", Decision: "exclude", Reason: "一般語として除外する",
	}); reviewErr == nil {
		t.Fatal("同じ意味と訳の確定前に除外できた")
	}
	confirmed, err := s.updateGeneralMatch(ctx, detail.GeneralMatches[0].ID,
		"same_mean_and_translation", "Skyrimでも囲いの意味と柵の訳を使う", "test-ai")
	if err != nil {
		t.Fatal(err)
	}
	if confirmed.GeneralMatchStatus != "same_mean_and_translation" || confirmed.InclusionDecision != "undecided" {
		t.Fatalf("confirmed = %+v", confirmed)
	}
	excluded, err := s.addReview(ctx, reviewInput{
		SenseID: confirmed.ID, Revision: confirmed.Revision, ReviewerKind: "ai",
		ReviewerReference: "test-ai", Decision: "exclude", Reason: "一般辞書と同じ意味と訳である",
	})
	if err != nil {
		t.Fatal(err)
	}
	if excluded.InclusionDecision != "exclude" || excluded.ReviewStage != "ai_reviewed" || len(excluded.Reviews) != 1 {
		t.Fatalf("excluded = %+v", excluded)
	}
	gotStatus, err := s.status(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if gotStatus.Reviews != 1 || gotStatus.ReviewDecisions["exclude"] != 1 {
		t.Fatalf("status after review = %+v", gotStatus)
	}
	assigned, err := s.assignOccurrence(ctx, detail.Occurrences[0].ID, secondSense.ID,
		"test-human", "使用箇所は盗品売人を表す")
	if err != nil {
		t.Fatal(err)
	}
	if assigned.ID != secondSense.ID || len(assigned.Occurrences) != 2 {
		t.Fatalf("assigned = %+v", assigned)
	}
	matchHistory, err := s.history(ctx, "general_dictionary_match", detail.GeneralMatches[0].ID)
	if err != nil || len(matchHistory.Changes) != 1 {
		t.Fatalf("match history = %+v, err = %v", matchHistory, err)
	}
	occurrenceHistory, err := s.history(ctx, "dictionary_occurrence", detail.Occurrences[0].ID)
	if err != nil || len(occurrenceHistory.Changes) != 1 {
		t.Fatalf("occurrence history = %+v, err = %v", occurrenceHistory, err)
	}
}

func TestMCPServerExposesR4Tools(t *testing.T) {
	ctx := context.Background()
	s := openTestStore(t)
	wordNetPath := createTestWordNet(t)
	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	serverSession, err := newMCPServer(s, wordNetPath).Connect(ctx, serverTransport, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if closeErr := serverSession.Close(); closeErr != nil {
			t.Error(closeErr)
		}
	})
	client := mcp.NewClient(&mcp.Implementation{Name: "dictionary-test", Version: "0.1.0"}, nil)
	clientSession, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if closeErr := clientSession.Close(); closeErr != nil {
			t.Error(closeErr)
		}
	})
	tools, err := clientSession.ListTools(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(tools.Tools) != 11 {
		t.Fatalf("tools = %d, want 11", len(tools.Tools))
	}
	added := callMCPTool[sense](ctx, t, clientSession, "dictionary_sense_add", map[string]any{
		"source": "Fence", "dest": "柵", "part_of_speech": "noun", "meaning": "区域を囲う構造物",
	})
	callMCPTool[classifyResult](ctx, t, clientSession, "dictionary_classify", map[string]any{})
	queue := callMCPTool[matchQueueResult](ctx, t, clientSession, "dictionary_general_match_queue", map[string]any{
		"status": "same_mean_candidate", "inclusion_decision": "undecided",
		"review_stage": "unreviewed", "limit": 1,
	})
	if len(queue.Entries) != 1 || queue.Entries[0].Source != "Fence" || queue.Entries[0].SkyrimCategories != "" {
		t.Fatalf("MCP match queue = %+v", queue)
	}
	searched := callMCPTool[searchResult](ctx, t, clientSession, "dictionary_search", map[string]any{
		"query": "Fence", "general_match_status": "same_mean_candidate",
	})
	if len(searched.Entries) != 1 {
		t.Fatalf("MCP search = %+v", searched)
	}
	detail := callMCPTool[sense](ctx, t, clientSession, "dictionary_get", map[string]any{"id": added.ID})
	updated := callMCPTool[sense](ctx, t, clientSession, "dictionary_sense_update", map[string]any{
		"id": detail.ID, "revision": detail.Revision, "dest": detail.Dest,
		"part_of_speech": "noun", "meaning": "柵の意味",
		"changed_by": "test", "reason": "説明を確認した",
	})
	confirmed := callMCPTool[sense](ctx, t, clientSession, "dictionary_general_match_update", map[string]any{
		"match_id": updated.GeneralMatches[0].ID, "status": "same_mean_and_translation",
		"changed_by": "test", "reason": "意味と訳が一致する",
	})
	callMCPTool[sense](ctx, t, clientSession, "dictionary_review_add", map[string]any{
		"sense_id": confirmed.ID, "revision": confirmed.Revision, "reviewer_kind": "ai",
		"reviewer_reference": "test", "decision": "exclude", "reason": "一般語として除外する",
	})
	second := callMCPTool[sense](ctx, t, clientSession, "dictionary_sense_add", map[string]any{
		"source": "Fence", "dest": "盗品売人", "part_of_speech": "noun", "meaning": "盗品を扱う人物",
	})
	callMCPTool[sense](ctx, t, clientSession, "dictionary_occurrence_assign", map[string]any{
		"occurrence_id": detail.Occurrences[0].ID, "sense_id": second.ID,
		"changed_by": "test", "reason": "別の意味へ割り当てる",
	})
	history := callMCPTool[changeHistory](ctx, t, clientSession, "dictionary_history", map[string]any{
		"target_table": "dictionary_sense", "target_id": detail.ID,
	})
	if len(history.Changes) == 0 {
		t.Fatal("MCP history is empty")
	}
	gotStatus := callMCPTool[status](ctx, t, clientSession, "dictionary_status", map[string]any{})
	if gotStatus.Terms != 1 || gotStatus.Senses != 2 {
		t.Fatalf("MCP status = %+v", gotStatus)
	}
}

func callMCPTool[T any](ctx context.Context, t *testing.T, session *mcp.ClientSession, name string, arguments map[string]any) T {
	t.Helper()
	result, err := session.CallTool(ctx, &mcp.CallToolParams{Name: name, Arguments: arguments})
	if err != nil {
		t.Fatal(err)
	}
	raw, err := json.Marshal(result.StructuredContent)
	if err != nil {
		t.Fatal(err)
	}
	var out T
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatal(err)
	}
	return out
}

func createTestWordNet(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "wordnet.sqlite3")
	db := openSourceDB(t, path)
	execTestSQL(t, db, `
		CREATE TABLE word (wordid INTEGER PRIMARY KEY, lang TEXT, lemma TEXT, pron TEXT, pos TEXT);
		CREATE TABLE sense (synset TEXT, wordid INTEGER, lang TEXT, rank TEXT, lexid INTEGER, freq INTEGER, src TEXT);
		CREATE TABLE synset_def (synset TEXT, lang TEXT, def TEXT, sid TEXT);
		INSERT INTO word (wordid, lang, lemma, pos) VALUES
			(1, 'eng', 'fence', 'n'), (2, 'jpn', '柵', 'n'),
			(3, 'eng', 'frost', 'n'), (4, 'jpn', 'フロスト', 'n'),
			(5, 'eng', 'imperial', 'a');
		INSERT INTO sense (synset, wordid, lang) VALUES
			('fence-barrier', 1, 'eng'), ('fence-barrier', 2, 'jpn'),
			('fence-dealer', 1, 'eng'),
			('frost-person', 3, 'eng'), ('frost-person', 4, 'jpn'),
			('frost-weather', 3, 'eng'), ('frost-weather', 4, 'jpn'),
			('imperial-adjective', 5, 'eng');
		INSERT INTO synset_def (synset, lang, def, sid) VALUES
			('fence-barrier', 'eng', 'a barrier enclosing an area', '1'),
			('fence-dealer', 'eng', 'a dealer in stolen property', '1'),
			('frost-person', 'eng', 'a person named Frost', '1'),
			('frost-weather', 'eng', 'ice crystals caused by cold', '1'),
			('imperial-adjective', 'eng', 'relating to an empire', '1');
	`)
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	return path
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
	t.Cleanup(func() {
		if closeErr := s.close(); closeErr != nil {
			t.Error(closeErr)
		}
	})
	return s
}

func execTestSQL(t *testing.T, db *sqlx.DB, query string) {
	t.Helper()
	if _, err := db.ExecContext(context.Background(), query); err != nil {
		t.Fatal(err)
	}
}
