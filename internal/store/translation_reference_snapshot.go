package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"aitranslationenginejp/internal/model"
)

// UpsertTranslationReferenceSnapshot は本文送信時の候補とprompt hashを保存する。
func (s *Store) UpsertTranslationReferenceSnapshot(ctx context.Context, row model.TranslationReferenceSnapshot) error {
	references, err := json.Marshal(row.References)
	if err != nil {
		return fmt.Errorf("参考語snapshotのJSON化: %w", err)
	}
	_, err = s.db.ExecContext(ctx, `
		INSERT INTO translation_reference_snapshot (plugin, kind, row_id, references_json, prompt_hash)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(plugin, kind, row_id) DO UPDATE SET
		 references_json = excluded.references_json, prompt_hash = excluded.prompt_hash`,
		row.Plugin, row.Kind, row.RowID, string(references), row.PromptHash)
	if err != nil {
		return fmt.Errorf("参考語snapshotの保存: %w", err)
	}
	return nil
}

// GetTranslationReferenceSnapshot は本文結果の送信時候補を返す。
func (s *Store) GetTranslationReferenceSnapshot(ctx context.Context, plugin, kind string, rowID int64) (model.TranslationReferenceSnapshot, bool, error) {
	var raw struct {
		ReferencesJSON string `db:"references_json"`
		PromptHash     string `db:"prompt_hash"`
	}
	err := s.db.GetContext(ctx, &raw, `
		SELECT references_json, prompt_hash FROM translation_reference_snapshot
		WHERE plugin = ? AND kind = ? AND row_id = ?`, plugin, kind, rowID)
	if errors.Is(err, sql.ErrNoRows) {
		return model.TranslationReferenceSnapshot{}, false, nil
	}
	if err != nil {
		return model.TranslationReferenceSnapshot{}, false, fmt.Errorf("参考語snapshotの取得: %w", err)
	}
	var refs []model.TranslationReference
	if err := json.Unmarshal([]byte(raw.ReferencesJSON), &refs); err != nil {
		return model.TranslationReferenceSnapshot{}, false, fmt.Errorf("参考語snapshotのJSON解析: %w", err)
	}
	return model.TranslationReferenceSnapshot{Plugin: plugin, Kind: kind, RowID: rowID, References: refs, PromptHash: raw.PromptHash}, true, nil
}
