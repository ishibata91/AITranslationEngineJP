package store

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"aitranslationenginejp/internal/model"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

func openTermDictionaryEditorForTest(t *testing.T) *TermDictionaryEditor {
	t.Helper()
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
	if err := db.Close(); err != nil {
		t.Fatalf("作成接続の終了: %v", err)
	}
	editor, err := OpenTermDictionaryEditor(path)
	if err != nil {
		t.Fatalf("OpenTermDictionaryEditor: %v", err)
	}
	t.Cleanup(func() { _ = editor.Close() })
	return editor
}

func stringPointer(value string) *string { return &value }

func TestTermDictionaryEditorUpdatesOnlyChangedColumnsAndKeepsSharedSource(t *testing.T) {
	editor := openTermDictionaryEditorForTest(t)
	ctx := context.Background()
	first, err := editor.Create(ctx, model.TermDictionaryCreate{Source: "Companion", Destination: "同胞団", PartOfSpeech: "noun"})
	if err != nil {
		t.Fatalf("最初の用語作成: %v", err)
	}
	second, err := editor.Create(ctx, model.TermDictionaryCreate{Source: "Companion", Destination: "仲間", PartOfSpeech: "noun"})
	if err != nil {
		t.Fatalf("同じ原語の用語作成: %v", err)
	}

	updated, err := editor.Patch(ctx, model.TermDictionaryPatch{ID: first.ID, Revision: first.Revision, Source: stringPointer("Fellowship")})
	if err != nil {
		t.Fatalf("原語の更新: %v", err)
	}
	if updated.Source != "Fellowship" || updated.Destination != "同胞団" || updated.PartOfSpeech != "noun" {
		t.Fatalf("更新結果 = %+v", updated)
	}
	shared, err := editor.List(ctx, model.TermDictionaryFilter{Source: "Companion"}, 1)
	if err != nil {
		t.Fatalf("共有原語の検索: %v", err)
	}
	if len(shared.Entries) != 1 || shared.Entries[0].ID != second.ID {
		t.Fatalf("共有原語の検索結果 = %+v", shared.Entries)
	}

	if _, err := editor.Patch(ctx, model.TermDictionaryPatch{ID: first.ID, Revision: first.Revision, Destination: stringPointer("同胞団改")}); !errors.Is(err, ErrTermDictionaryRevisionConflict) {
		t.Fatalf("古い内容での更新 error = %v", err)
	}
	var changes []string
	if err := editor.db.Select(&changes, `SELECT field_name FROM dictionary_change WHERE target_id=? ORDER BY id`, first.ID); err != nil {
		t.Fatalf("変更内容の取得: %v", err)
	}
	if len(changes) != 1 || changes[0] != "source" {
		t.Fatalf("変更内容 = %v", changes)
	}
	if _, err := editor.Create(ctx, model.TermDictionaryCreate{Source: "Invalid", Destination: "無効", PartOfSpeech: "invalid"}); err == nil {
		t.Fatal("受け付けない品詞を作成できた")
	}
	if err := editor.Delete(ctx, second.ID, second.Revision); err != nil {
		t.Fatalf("用語の削除: %v", err)
	}
	remaining, err := editor.List(ctx, model.TermDictionaryFilter{Source: "Companion"}, 1)
	if err != nil {
		t.Fatalf("削除後の検索: %v", err)
	}
	if remaining.TotalCount != 0 {
		t.Fatalf("削除後の件数 = %d", remaining.TotalCount)
	}
}

func TestTermDictionaryEditorListsFiftyEntriesAndFiltersCategories(t *testing.T) {
	editor := openTermDictionaryEditorForTest(t)
	ctx := context.Background()
	for index := 1; index <= termDictionaryPageSize+1; index++ {
		result, err := editor.db.ExecContext(ctx, `INSERT INTO dictionary_term (source) VALUES (?)`, fmt.Sprintf("Source %03d", index))
		if err != nil {
			t.Fatalf("原語の追加: %v", err)
		}
		termID, err := result.LastInsertId()
		if err != nil {
			t.Fatalf("原語の識別子取得: %v", err)
		}
		result, err = editor.db.ExecContext(ctx, `INSERT INTO dictionary_sense (term_id,dest,part_of_speech) VALUES (?,?,'noun')`, termID, fmt.Sprintf("訳語 %03d", index))
		if err != nil {
			t.Fatalf("訳語の追加: %v", err)
		}
		senseID, err := result.LastInsertId()
		if err != nil {
			t.Fatalf("訳語の識別子取得: %v", err)
		}
		if _, err := editor.db.ExecContext(ctx, `INSERT INTO dictionary_occurrence (term_id,sense_id,skyrim_category,origin_kind,origin_reference) VALUES (?,?,?,'test',?)`, termID, senseID, "DIAL", fmt.Sprintf("test:%d", index)); err != nil {
			t.Fatalf("Skyrimカテゴリの追加: %v", err)
		}
	}

	first, err := editor.List(ctx, model.TermDictionaryFilter{Category: "DIAL"}, 1)
	if err != nil {
		t.Fatalf("最初のページ取得: %v", err)
	}
	if first.TotalCount != termDictionaryPageSize+1 || len(first.Entries) != termDictionaryPageSize {
		t.Fatalf("最初のページ = %+v", first)
	}
	second, err := editor.List(ctx, model.TermDictionaryFilter{Category: "DIAL"}, 2)
	if err != nil {
		t.Fatalf("次のページ取得: %v", err)
	}
	if len(second.Entries) != 1 || len(second.Entries[0].Categories) != 1 || second.Entries[0].Categories[0] != "DIAL" {
		t.Fatalf("次のページ = %+v", second)
	}
}
