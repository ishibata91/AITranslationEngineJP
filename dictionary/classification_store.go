package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

func (s *store) assignOccurrence(ctx context.Context, occurrenceID, senseID int64, changedBy, reason string) (sense, error) {
	changedBy, reason = strings.TrimSpace(changedBy), strings.TrimSpace(reason)
	if changedBy == "" || reason == "" {
		return sense{}, errors.New("変更者と変更理由は空にできない")
	}
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return sense{}, fmt.Errorf("使用箇所割り当てtransaction開始: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck
	type assignment struct {
		TermID  int64  `db:"term_id"`
		SenseID *int64 `db:"sense_id"`
	}
	var current assignment
	if err := tx.GetContext(ctx, &current, `
		SELECT term_id, sense_id FROM dictionary_occurrence WHERE id = ?`, occurrenceID); err != nil {
		return sense{}, fmt.Errorf("使用箇所 %d の取得: %w", occurrenceID, err)
	}
	var targetTermID int64
	if err := tx.GetContext(ctx, &targetTermID, `SELECT term_id FROM dictionary_sense WHERE id = ?`, senseID); err != nil {
		return sense{}, fmt.Errorf("割り当て先の意味 %d の取得: %w", senseID, err)
	}
	if current.TermID != targetTermID {
		return sense{}, errors.New("使用箇所と意味の英語表記が一致しない")
	}
	if _, err := tx.ExecContext(ctx, `
		UPDATE dictionary_occurrence SET sense_id = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
		senseID, occurrenceID); err != nil {
		return sense{}, fmt.Errorf("使用箇所 %d の割り当て: %w", occurrenceID, err)
	}
	oldValue := ""
	if current.SenseID != nil {
		oldValue = strconv.FormatInt(*current.SenseID, 10)
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO dictionary_change
		    (target_table, target_id, field_name, old_value, new_value, changed_by, reason)
		VALUES ('dictionary_occurrence', ?, 'sense_id', ?, ?, ?, ?)`,
		occurrenceID, oldValue, strconv.FormatInt(senseID, 10), changedBy, reason); err != nil {
		return sense{}, fmt.Errorf("使用箇所 %d の変更履歴保存: %w", occurrenceID, err)
	}
	if err := tx.Commit(); err != nil {
		return sense{}, fmt.Errorf("使用箇所割り当てcommit: %w", err)
	}
	return s.get(ctx, senseID)
}

func (s *store) updateGeneralMatch(ctx context.Context, matchID int64, status, reason, changedBy string) (sense, error) {
	status, reason, changedBy = strings.TrimSpace(status), strings.TrimSpace(reason), strings.TrimSpace(changedBy)
	if reason == "" || changedBy == "" {
		return sense{}, errors.New("判定理由と変更者は空にできない")
	}
	if status != "same_mean_and_translation" && status != "different_meaning_or_translation" {
		return sense{}, errors.New("人間またはAIが確定できる状態を指定する")
	}
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return sense{}, fmt.Errorf("一般辞書照合更新transaction開始: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck
	type currentMatch struct {
		SenseID     int64  `db:"sense_id"`
		MatchStatus string `db:"match_status"`
	}
	var current currentMatch
	if err := tx.GetContext(ctx, &current, `
		SELECT sense_id, match_status FROM general_dictionary_match WHERE id = ?`, matchID); err != nil {
		return sense{}, fmt.Errorf("一般辞書照合 %d の取得: %w", matchID, err)
	}
	if _, err := tx.ExecContext(ctx, `
		UPDATE general_dictionary_match
		SET match_status = ?, reason = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?`, status, reason, matchID); err != nil {
		return sense{}, fmt.Errorf("一般辞書照合 %d の更新: %w", matchID, err)
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO dictionary_change
		    (target_table, target_id, field_name, old_value, new_value, changed_by, reason)
		VALUES ('general_dictionary_match', ?, 'match_status', ?, ?, ?, ?)`,
		matchID, current.MatchStatus, status, changedBy, reason); err != nil {
		return sense{}, fmt.Errorf("一般辞書照合 %d の変更履歴保存: %w", matchID, err)
	}
	if err := refreshSenseGeneralStatus(ctx, tx, current.SenseID); err != nil {
		return sense{}, err
	}
	if err := tx.Commit(); err != nil {
		return sense{}, fmt.Errorf("一般辞書照合更新commit: %w", err)
	}
	return s.get(ctx, current.SenseID)
}

func refreshSenseGeneralStatus(ctx context.Context, tx *sqlx.Tx, senseID int64) error {
	if _, err := tx.ExecContext(ctx, `
		UPDATE dictionary_sense
		SET general_match_status = CASE
		        WHEN EXISTS (
		            SELECT 1 FROM general_dictionary_match m
		            WHERE m.sense_id = dictionary_sense.id
		              AND m.match_status = 'same_mean_and_translation'
		        ) THEN 'same_mean_and_translation'
		        WHEN EXISTS (
		            SELECT 1 FROM general_dictionary_match m
		            WHERE m.sense_id = dictionary_sense.id
		              AND m.match_status = 'same_mean_candidate'
		        ) THEN 'same_mean_candidate'
		        WHEN EXISTS (
		            SELECT 1 FROM general_dictionary_match m
		            WHERE m.sense_id = dictionary_sense.id
		              AND m.match_status = 'same_surface_only'
		        ) THEN 'same_surface_only'
		        WHEN EXISTS (
		            SELECT 1 FROM general_dictionary_match m
		            WHERE m.sense_id = dictionary_sense.id
		              AND m.match_status = 'different_meaning_or_translation'
		        ) THEN 'different_meaning_or_translation'
		        ELSE 'no_english_headword'
		    END,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = ?`, senseID); err != nil {
		return fmt.Errorf("意味 %d の一般辞書状態更新: %w", senseID, err)
	}
	return nil
}
