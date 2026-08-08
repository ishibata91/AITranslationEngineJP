package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

type dictionaryChange struct {
	ID          int64  `db:"id" json:"id"`
	TargetTable string `db:"target_table" json:"target_table"`
	TargetID    int64  `db:"target_id" json:"target_id"`
	FieldName   string `db:"field_name" json:"field_name"`
	OldValue    string `db:"old_value" json:"old_value"`
	NewValue    string `db:"new_value" json:"new_value"`
	ChangedBy   string `db:"changed_by" json:"changed_by"`
	Reason      string `db:"reason" json:"reason"`
	CreatedAt   string `db:"created_at" json:"created_at"`
}

type changeHistory struct {
	Changes []dictionaryChange `json:"changes"`
}

func (s *store) history(ctx context.Context, targetTable string, targetID int64) (changeHistory, error) {
	targetTable = strings.TrimSpace(targetTable)
	allowed := map[string]bool{
		"dictionary_term":          true,
		"dictionary_sense":         true,
		"dictionary_occurrence":    true,
		"general_dictionary_match": true,
		"dictionary_review":        true,
	}
	if !allowed[targetTable] {
		return changeHistory{}, errors.New("変更履歴を取得できないtable")
	}
	out := changeHistory{Changes: []dictionaryChange{}}
	if err := s.db.SelectContext(ctx, &out.Changes, `
		SELECT id, target_table, target_id, field_name, old_value, new_value,
		       changed_by, reason, created_at
		FROM dictionary_change
		WHERE target_table = ? AND target_id = ?
		ORDER BY id`, targetTable, targetID); err != nil {
		return changeHistory{}, fmt.Errorf("%s %d の変更履歴取得: %w", targetTable, targetID, err)
	}
	return out, nil
}
