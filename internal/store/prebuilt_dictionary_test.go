package store

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

// 事前作成済み辞書readerは収録判断に関係なく全候補を読み、meaning は内部読取りに残す。
func TestPrebuiltDictionaryReferencesAreReadOnlyAndIncludeAllSenses(t *testing.T) {
	path := filepath.Join(t.TempDir(), "dictionary.sqlite3")
	db, err := sqlx.Connect("sqlite", path)
	if err != nil {
		t.Fatalf("辞書DBの作成: %v", err)
	}
	schema, err := os.ReadFile(filepath.Join("..", "..", "dictionary", "schema.sql"))
	if err != nil {
		t.Fatalf("schemaの読取り: %v", err)
	}
	if _, err := db.Exec(string(schema)); err != nil {
		t.Fatalf("schemaの適用: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO dictionary_term (id, source) VALUES (1, 'Riften')`); err != nil {
		t.Fatalf("termの追加: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO dictionary_sense (id, term_id, dest, part_of_speech, meaning, inclusion_decision, review_stage) VALUES
		(1, 1, 'リフテン', 'noun', '城塞を守る都市', 'exclude', 'unreviewed'),
		(2, 1, 'リフテン', 'noun', '盗賊ギルドの町', 'include', 'human_reviewed')`); err != nil {
		t.Fatalf("senseの追加: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO dictionary_occurrence (term_id, sense_id, skyrim_category, origin_kind, origin_reference) VALUES (1, 1, 'city', 'test', 'a')`); err != nil {
		t.Fatalf("occurrenceの追加: %v", err)
	}
	if err := db.Close(); err != nil {
		t.Fatalf("作成接続の終了: %v", err)
	}

	reader, err := OpenPrebuiltDictionary(path)
	if err != nil {
		t.Fatalf("OpenPrebuiltDictionary: %v", err)
	}
	defer func() { _ = reader.Close() }()
	ctx := context.Background()
	if err := reader.ValidatePrebuiltDictionary(ctx); err != nil {
		t.Fatalf("ValidatePrebuiltDictionary: %v", err)
	}
	refs, err := reader.References(ctx)
	if err != nil {
		t.Fatalf("References: %v", err)
	}
	if len(refs) != 2 {
		t.Fatalf("候補数 = %d, want 2: %+v", len(refs), refs)
	}
	if refs[0].Meaning == "" || refs[1].Meaning == "" {
		t.Errorf("meaning がreader内部から失われた: %+v", refs)
	}
	if _, err := reader.db.Exec(`INSERT INTO dictionary_term (source) VALUES ('Whiterun')`); err == nil {
		t.Error("読み取り専用readerが書き込みを許可した")
	}
}

func TestOpenPrebuiltDictionaryFailsWhenFileIsMissing(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.sqlite3")
	if _, err := OpenPrebuiltDictionary(path); err == nil {
		t.Fatal("存在しない辞書DBを開けた")
	}
}
