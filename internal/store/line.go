package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"aitranslationenginejp/internal/model"
)

// lineColumns は line の SELECT 列。model.Line の db タグと対応する。1 箇所に集約する。
const lineColumns = `id, source, dest, status, response_order, plugin, form_id, edid, rec, field, ordinal`

// ListUntranslatedLines は未訳（status=0）の台詞を id 昇順で返す。
func (s *Store) ListUntranslatedLines(ctx context.Context) ([]model.Line, error) {
	return s.queryLines(ctx, `SELECT `+lineColumns+` FROM line WHERE status = 0 ORDER BY id`)
}

// ListLines は全ての台詞を id 昇順で返す（画面の結果一覧表示用）。
func (s *Store) ListLines(ctx context.Context) ([]model.Line, error) {
	return s.queryLines(ctx, `SELECT `+lineColumns+` FROM line ORDER BY id`)
}

func (s *Store) queryLines(ctx context.Context, query string) ([]model.Line, error) {
	var rows []model.Line
	if err := s.db.SelectContext(ctx, &rows, query); err != nil {
		return nil, fmt.Errorf("line の取得: %w", err)
	}
	return rows, nil
}

// UpdateLineDest は訳文と訳状態を書き戻す。
func (s *Store) UpdateLineDest(ctx context.Context, id int64, dest string, status int) error {
	if _, err := s.db.ExecContext(ctx,
		`UPDATE line SET dest = ?, status = ? WHERE id = ?`, dest, status, id); err != nil {
		return fmt.Errorf("line dest の更新: %w", err)
	}
	return nil
}

// LoadLineSpeaker は台詞に紐づく話者の事実上の同定情報（種族・声型・所属勢力の EditorID）を
// speaker・race・voice_type・faction から join して読む。純汎用台詞は複数話者を指しうるが、
// 本実装は最小として id 昇順の先頭話者を採る。話者が紐づかない場合は found=false を返す。
func (s *Store) LoadLineSpeaker(ctx context.Context, lineID int64) (model.SpeakerIdentity, bool, error) {
	var row struct {
		SpeakerID int64  `db:"speaker_id"`
		RaceEDID  string `db:"race_edid"`
		VoiceEDID string `db:"voice_edid"`
	}
	err := s.db.GetContext(ctx, &row, `
		SELECT s.id AS speaker_id,
		       COALESCE(r.edid, '') AS race_edid,
		       COALESCE(v.edid, '') AS voice_edid
		FROM line_speaker ls
		JOIN speaker s ON s.id = ls.speaker_id
		LEFT JOIN race r ON r.id = s.race_id
		LEFT JOIN voice_type v ON v.id = s.voice_type_id
		WHERE ls.line_id = ?
		ORDER BY s.id
		LIMIT 1`, lineID)
	if errors.Is(err, sql.ErrNoRows) {
		return model.SpeakerIdentity{}, false, nil
	}
	if err != nil {
		return model.SpeakerIdentity{}, false, fmt.Errorf("話者識別子の取得: %w", err)
	}

	var factionEDIDs []string
	if err := s.db.SelectContext(ctx, &factionEDIDs, `
		SELECT f.edid
		FROM speaker_faction sf
		JOIN faction f ON f.id = sf.faction_id
		WHERE sf.speaker_id = ? AND f.edid <> ''
		ORDER BY f.id`, row.SpeakerID); err != nil {
		return model.SpeakerIdentity{}, false, fmt.Errorf("所属勢力の取得: %w", err)
	}

	return model.SpeakerIdentity{
		RaceEDID:     row.RaceEDID,
		VoiceEDID:    row.VoiceEDID,
		FactionEDIDs: factionEDIDs,
	}, true, nil
}
