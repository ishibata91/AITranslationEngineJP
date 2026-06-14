package store

import (
	"context"
	"fmt"

	"aitranslationenginejp/internal/model"
)

// lineColumns は line の SELECT 列。model.Line の db タグと対応する。1 箇所に集約する。
const lineColumns = `id, source, dest, status, response_order, plugin, form_id, edid, rec, field, ordinal`

// ListUntranslatedLines は未訳（status=0）の台詞を id 昇順で返す。
func (s *Store) ListUntranslatedLines(ctx context.Context) ([]model.Line, error) {
	return s.queryLines(ctx, `SELECT `+lineColumns+` FROM line WHERE status = 0 ORDER BY id`)
}

// LinesAfter は id が afterID より大きい台詞を id 昇順で最大 limit 件返す（keyset ページング用）。
// afterID=0 で先頭から、最後に読んだ id を次回 afterID に渡して次ページを得る。
func (s *Store) LinesAfter(ctx context.Context, afterID int64, limit int) ([]model.Line, error) {
	return s.queryLines(ctx, `SELECT `+lineColumns+` FROM line WHERE id > ? ORDER BY id LIMIT ?`, afterID, limit)
}

// CountLines は台詞の総件数を返す（ページャの総件数表示用）。
func (s *Store) CountLines(ctx context.Context) (int, error) {
	return s.count(ctx, `SELECT COUNT(*) FROM line`)
}

func (s *Store) queryLines(ctx context.Context, query string, args ...any) ([]model.Line, error) {
	var rows []model.Line
	if err := s.db.SelectContext(ctx, &rows, query, args...); err != nil {
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

// LoadLineSpeakers は台詞 id 群に紐づく話者の事実上の同定情報（種族・声型・所属勢力の EditorID）を
// 一括で読み、lineID→識別子 の map を返す。話者が紐づかない台詞は map に現れない。
// 純汎用台詞は複数話者を指しうるが、本実装は最小として id 昇順の先頭話者を採る（per-line 実装と同基準）。
// 話者属性の取得を台詞数に依存させないため、line ごとの個別問い合わせをやめて IN 句で一括取得する。
// IN 句の host parameter 上限を避けるため lineIDs を chunk して問い合わせる。
func (s *Store) LoadLineSpeakers(ctx context.Context, lineIDs []int64) (map[int64]model.SpeakerIdentity, error) {
	out := make(map[int64]model.SpeakerIdentity, len(lineIDs))
	if len(lineIDs) == 0 {
		return out, nil
	}

	// lineID → 先頭話者 speakerID。所属勢力の一括取得に使う。
	lineSpeaker := make(map[int64]int64, len(lineIDs))
	var speakerIDs []int64
	for _, chunk := range chunkInt64(lineIDs, speakerChunkSize) {
		marks, args := placeholders(chunk)
		var rows []struct {
			LineID    int64  `db:"line_id"`
			SpeakerID int64  `db:"speaker_id"`
			RaceEDID  string `db:"race_edid"`
			VoiceEDID string `db:"voice_edid"`
		}
		query := `
			SELECT ls.line_id AS line_id,
			       s.id AS speaker_id,
			       COALESCE(r.edid, '') AS race_edid,
			       COALESCE(v.edid, '') AS voice_edid
			FROM line_speaker ls
			JOIN speaker s ON s.id = ls.speaker_id
			LEFT JOIN race r ON r.id = s.race_id
			LEFT JOIN voice_type v ON v.id = s.voice_type_id
			WHERE ls.line_id IN (` + marks + `)
			ORDER BY ls.line_id, s.id`
		if err := s.db.SelectContext(ctx, &rows, query, args...); err != nil {
			return nil, fmt.Errorf("話者識別子の一括取得: %w", err)
		}
		for _, row := range rows {
			if _, seen := out[row.LineID]; seen {
				continue // 先頭話者（id 昇順）だけ採る。
			}
			out[row.LineID] = model.SpeakerIdentity{RaceEDID: row.RaceEDID, VoiceEDID: row.VoiceEDID}
			lineSpeaker[row.LineID] = row.SpeakerID
			speakerIDs = append(speakerIDs, row.SpeakerID)
		}
	}
	if len(speakerIDs) == 0 {
		return out, nil
	}

	// speakerID → 所属勢力 EditorID 群を一括取得する。
	factions := make(map[int64][]string, len(speakerIDs))
	for _, chunk := range chunkInt64(speakerIDs, speakerChunkSize) {
		marks, args := placeholders(chunk)
		var rows []struct {
			SpeakerID int64  `db:"speaker_id"`
			EDID      string `db:"edid"`
		}
		query := `
			SELECT sf.speaker_id AS speaker_id, f.edid AS edid
			FROM speaker_faction sf
			JOIN faction f ON f.id = sf.faction_id
			WHERE sf.speaker_id IN (` + marks + `) AND f.edid <> ''
			ORDER BY sf.speaker_id, f.id`
		if err := s.db.SelectContext(ctx, &rows, query, args...); err != nil {
			return nil, fmt.Errorf("所属勢力の一括取得: %w", err)
		}
		for _, row := range rows {
			factions[row.SpeakerID] = append(factions[row.SpeakerID], row.EDID)
		}
	}

	// 取得した所属勢力を各 line の識別子へ付ける。
	for lineID, speakerID := range lineSpeaker {
		if fs := factions[speakerID]; len(fs) > 0 {
			identity := out[lineID]
			identity.FactionEDIDs = fs
			out[lineID] = identity
		}
	}
	return out, nil
}
