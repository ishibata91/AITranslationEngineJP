package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"aitranslationenginejp/internal/model"
)

// StartBatchProgression は対象 plugin の batch 進行を開始（または reset）し、進行本体の id を返す。
// 同一 plugin に進行が既にある場合は同じ行を使い回し、provider と model を更新し、stage を固有名段へ戻し、
// 外部 batch ID を空へ戻す。前回の送信行対応（batch_request）も消して、まっさらな進行にする。
// 進行は plugin と 1 対 1（UNIQUE(plugin)）。1 トランザクションで reset と子の掃除を原子的に行う。
func (s *Store) StartBatchProgression(ctx context.Context, plugin, provider, model string) (int64, error) {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("batch 進行開始のトランザクション: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck // commit 済みなら no-op。失敗時の後始末用。

	if _, err := tx.ExecContext(ctx,
		`INSERT INTO batch_translation (plugin, provider, model, stage, proper_batch_id, body_batch_id)
		 VALUES (?, ?, ?, ?, '', '')
		 ON CONFLICT(plugin) DO UPDATE SET
		   provider = excluded.provider, model = excluded.model, stage = excluded.stage,
		   proper_batch_id = '', body_batch_id = ''`,
		plugin, provider, model, batchStageProperNoun); err != nil {
		return 0, fmt.Errorf("batch 進行の upsert: %w", err)
	}
	var id int64
	if err := tx.GetContext(ctx, &id, `SELECT id FROM batch_translation WHERE plugin = ?`, plugin); err != nil {
		return 0, fmt.Errorf("batch 進行 id の取得: %w", err)
	}
	// reset 時に前回の送信行対応を掃除する（同じ id を使い回すため）。
	if _, err := tx.ExecContext(ctx, `DELETE FROM batch_request WHERE batch_id = ?`, id); err != nil {
		return 0, fmt.Errorf("旧 batch_request の掃除: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("batch 進行開始のコミット: %w", err)
	}
	return id, nil
}

const batchStageProperNoun = model.BatchStageProperNoun

// RecordBatchExternalID は進行本体へ現在処理している外部 batch ID を記録する。
// stage が固有名段なら proper_batch_id、本文段なら body_batch_id を更新する。
func (s *Store) RecordBatchExternalID(ctx context.Context, batchID int64, stage, externalID string) error {
	col, err := batchExternalIDColumn(stage)
	if err != nil {
		return err
	}
	if _, err := s.db.ExecContext(ctx,
		//nolint:gosec // col は batchExternalIDColumn が固定文字列から選んだ列名で、外部入力ではない。
		fmt.Sprintf(`UPDATE batch_translation SET %s = ? WHERE id = ?`, col),
		externalID, batchID); err != nil {
		return fmt.Errorf("外部 batch ID の記録: %w", err)
	}
	return nil
}

// batchExternalIDColumn は進行段から外部 batch ID の列名を選ぶ。固有名段・本文段のみ受ける。
func batchExternalIDColumn(stage string) (string, error) {
	switch stage {
	case model.BatchStageProperNoun:
		return "proper_batch_id", nil
	case model.BatchStageBody:
		return "body_batch_id", nil
	default:
		return "", fmt.Errorf("外部 batch ID を持たない進行段: %q", stage)
	}
}

// AdvanceBatchStage は進行本体の stage を進める（固有名 → 本文 → 完了）。
func (s *Store) AdvanceBatchStage(ctx context.Context, batchID int64, stage string) error {
	if _, err := s.db.ExecContext(ctx,
		`UPDATE batch_translation SET stage = ? WHERE id = ?`, stage, batchID); err != nil {
		return fmt.Errorf("進行段の更新: %w", err)
	}
	return nil
}

// GetBatchProgression は対象 plugin の batch 進行を返す。無ければ ok=false。
func (s *Store) GetBatchProgression(ctx context.Context, plugin string) (model.BatchTranslation, bool, error) {
	var row model.BatchTranslation
	err := s.db.GetContext(ctx, &row,
		`SELECT id, plugin, provider, model, stage, proper_batch_id, body_batch_id, created_at
		 FROM batch_translation WHERE plugin = ?`, plugin)
	if errors.Is(err, sql.ErrNoRows) {
		return model.BatchTranslation{}, false, nil
	}
	if err != nil {
		return model.BatchTranslation{}, false, fmt.Errorf("batch 進行の取得: %w", err)
	}
	return row, true, nil
}

// ListActiveBatchProgressions は反映待ちの進行（stage が完了でない）を新しい順で返す。
// 起動時・反映操作で、対象 plugin をまたいで反映すべき進行を洗い出すのに使う。
func (s *Store) ListActiveBatchProgressions(ctx context.Context) ([]model.BatchTranslation, error) {
	var rows []model.BatchTranslation
	if err := s.db.SelectContext(ctx, &rows,
		`SELECT id, plugin, provider, model, stage, proper_batch_id, body_batch_id, created_at
		 FROM batch_translation WHERE stage != ? ORDER BY created_at DESC, plugin`,
		model.BatchStageDone); err != nil {
		return nil, fmt.Errorf("反映待ち batch 進行の一覧: %w", err)
	}
	return rows, nil
}

// InsertBatchRequests は送信行対応（custom_id ↔ 行）を一括で書く。追加できた件数を返す。
// 送信のたびに、その段のリクエスト群を記録する。UNIQUE(batch_id, custom_id) で二重を弾く。
func (s *Store) InsertBatchRequests(ctx context.Context, rows []model.BatchRequest) (int, error) {
	return batchInsert(ctx, s, rows, func(tx *sqlx.Tx, r model.BatchRequest) (int64, error) {
		res, err := tx.ExecContext(ctx,
			`INSERT OR IGNORE INTO batch_request (batch_id, external_batch_id, custom_id, kind, row_id, references_json, prompt_hash, send_state)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			r.BatchID, r.ExternalBatchID, r.CustomID, r.Kind, r.RowID, r.ReferencesJSON, r.PromptHash, r.SendState)
		if err != nil {
			return 0, fmt.Errorf("batch_request の書き込み: %w", err)
		}
		n, _ := res.RowsAffected()
		return n, nil
	})
}

// ListBatchRequests は外部 batch ID に載せた送信行対応を返す。結果反映時に custom_id → 行を引くのに使う。
func (s *Store) ListBatchRequests(ctx context.Context, externalBatchID string) ([]model.BatchRequest, error) {
	var rows []model.BatchRequest
	if err := s.db.SelectContext(ctx, &rows,
		`SELECT id, batch_id, external_batch_id, custom_id, kind, row_id, references_json, prompt_hash, send_state
		 FROM batch_request WHERE external_batch_id = ? ORDER BY id`, externalBatchID); err != nil {
		return nil, fmt.Errorf("送信行対応の取得: %w", err)
	}
	return rows, nil
}

// ListBatchRequestsByStage は同じ進行段ですでに外部 batch へ送った要求を返す。
// 固有名段は kind='p'、本文段は kind IN ('n','l') で区別する。
func (s *Store) ListBatchRequestsByStage(ctx context.Context, batchID int64, stage string) ([]model.BatchRequest, error) {
	var query string
	switch stage {
	case model.BatchStageProperNoun:
		query = `SELECT id, batch_id, external_batch_id, custom_id, kind, row_id, references_json, prompt_hash, send_state
		 FROM batch_request WHERE batch_id = ? AND kind = ? ORDER BY id`
	case model.BatchStageBody:
		query = `SELECT id, batch_id, external_batch_id, custom_id, kind, row_id, references_json, prompt_hash, send_state
		 FROM batch_request WHERE batch_id = ? AND kind IN (?, ?) ORDER BY id`
	default:
		return nil, fmt.Errorf("送信要求を持たない進行段: %q", stage)
	}
	var rows []model.BatchRequest
	var err error
	if stage == model.BatchStageProperNoun {
		err = s.db.SelectContext(ctx, &rows, query, batchID, model.BatchKindProper)
	} else {
		err = s.db.SelectContext(ctx, &rows, query, batchID, model.BatchKindNarration, model.BatchKindLine)
	}
	if err != nil {
		return nil, fmt.Errorf("進行段の送信行対応の取得: %w", err)
	}
	return rows, nil
}

// UpdateBatchRequestSendState は送信前仮状態を送信済みまたは送信失敗へまとめて更新する。
func (s *Store) UpdateBatchRequestSendState(ctx context.Context, batchID int64, customIDs []string, state, externalID string) error {
	if len(customIDs) == 0 {
		return nil
	}
	query, args, err := sqlx.In(`UPDATE batch_request SET send_state = ?, external_batch_id = ? WHERE batch_id = ? AND custom_id IN (?)`, state, externalID, batchID, customIDs)
	if err != nil {
		return fmt.Errorf("batch request送信状態のSQL生成: %w", err)
	}
	query = s.db.Rebind(query)
	if _, err := s.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("batch request送信状態の更新: %w", err)
	}
	return nil
}
