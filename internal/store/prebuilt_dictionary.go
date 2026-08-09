package store

import (
	"context"
	"fmt"

	"aitranslationenginejp/internal/model"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

// PrebuiltDictionary は事前作成済み翻訳辞書への読み取り専用接続である。
type PrebuiltDictionary struct {
	db *sqlx.DB
}

// OpenPrebuiltDictionary は指定pathの辞書DBを読み取り専用で開き、起動時にschema_versionを読めることを確認する。
func OpenPrebuiltDictionary(path string) (*PrebuiltDictionary, error) {
	conn, err := sqlx.Connect("sqlite", "file:"+path+"?mode=ro")
	if err != nil {
		return nil, fmt.Errorf("事前作成済み辞書DBを開けない (%s): %w", path, err)
	}
	var version int
	if err := conn.Get(&version, `PRAGMA schema_version`); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("事前作成済み辞書DBのschema_versionを読めない (%s): %w", path, err)
	}
	return &PrebuiltDictionary{db: conn}, nil
}

// Close は読み取り専用接続を閉じる。
func (d *PrebuiltDictionary) Close() error {
	if err := d.db.Close(); err != nil {
		return fmt.Errorf("事前作成済み辞書DBのクローズ: %w", err)
	}
	return nil
}

// ValidatePrebuiltDictionary は翻訳開始前に必要なschemaと全参考語の読取りを検証する。
func (d *PrebuiltDictionary) ValidatePrebuiltDictionary(ctx context.Context) error {
	var count int
	if err := d.db.GetContext(ctx, &count, `SELECT COUNT(*) FROM dictionary_term`); err != nil {
		return fmt.Errorf("dictionary_termの読取り: %w", err)
	}
	if err := d.db.GetContext(ctx, &count, `SELECT COUNT(*) FROM dictionary_sense`); err != nil {
		return fmt.Errorf("dictionary_senseの読取り: %w", err)
	}
	if err := d.db.GetContext(ctx, &count, `SELECT COUNT(*) FROM dictionary_occurrence`); err != nil {
		return fmt.Errorf("dictionary_occurrenceの読取り: %w", err)
	}
	_, err := d.References(ctx)
	if err != nil {
		return fmt.Errorf("事前作成済み辞書の全参考語読取り: %w", err)
	}
	return nil
}

// References は収録判断に依らず全ての意味単位の候補を返す。
// 同一の原語・訳語・品詞・意味・カテゴリだけを重複除去する。
func (d *PrebuiltDictionary) References(ctx context.Context) ([]model.PrebuiltDictionaryReference, error) {
	var rows []model.PrebuiltDictionaryReference
	err := d.db.SelectContext(ctx, &rows, `
		SELECT DISTINCT t.source, s.dest, s.part_of_speech, s.meaning, COALESCE(o.skyrim_category, '') AS skyrim_category
		FROM dictionary_term t
		JOIN dictionary_sense s ON s.term_id = t.id
		LEFT JOIN dictionary_occurrence o ON o.sense_id = s.id
		ORDER BY t.source, s.dest, s.part_of_speech, s.meaning, o.skyrim_category`)
	if err != nil {
		return nil, fmt.Errorf("辞書参考語の取得: %w", err)
	}
	return rows, nil
}
