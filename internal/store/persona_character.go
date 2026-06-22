package store

import (
	"context"
	"fmt"
	"strings"

	"aitranslationenginejp/internal/model"
)

// hashChunkSize は IN 句に渡す host parameter の上限を避けるための分割サイズ（source_hash 群用）。
const hashChunkSize = 900

// stringPlaceholders は文字列群を IN 句用の "?,?,..." と引数列へ変換する。
func stringPlaceholders(vs []string) (string, []any) {
	marks := strings.TrimSuffix(strings.Repeat("?,", len(vs)), ",")
	args := make([]any, len(vs))
	for i, v := range vs {
		args[i] = v
	}
	return marks, args
}

// chunkString は文字列群を size ごとの塊へ分割する。IN 句の host parameter 上限回避に使う。
func chunkString(vs []string, size int) [][]string {
	var chunks [][]string
	for start := 0; start < len(vs); start += size {
		end := min(start+size, len(vs))
		chunks = append(chunks, vs[start:end])
	}
	return chunks
}

// lineAnalysisColumns は line_analysis の SELECT 列。model.LineAnalysis の db タグと対応する。
const lineAnalysisColumns = `id, source_hash, sentence_count, polite_count, insult_count, is_imperative, exclaim_count, elong_count, emotion_count`

// GetLineAnalyses は本文ハッシュ群に対応する行解析キャッシュを一括で読み、source_hash→特徴 の map を返す。
// キャッシュに無いハッシュは map に現れない（呼び出し側が prose で計算して UpsertLineAnalysis する）。
func (s *Store) GetLineAnalyses(ctx context.Context, hashes []string) (map[string]model.LineAnalysis, error) {
	out := make(map[string]model.LineAnalysis, len(hashes))
	if len(hashes) == 0 {
		return out, nil
	}
	for _, chunk := range chunkString(hashes, hashChunkSize) {
		marks, args := stringPlaceholders(chunk)
		var rows []model.LineAnalysis
		query := `SELECT ` + lineAnalysisColumns + ` FROM line_analysis WHERE source_hash IN (` + marks + `)`
		if err := s.db.SelectContext(ctx, &rows, query, args...); err != nil {
			return nil, fmt.Errorf("行解析キャッシュの一括取得: %w", err)
		}
		for _, row := range rows {
			out[row.SourceHash] = row
		}
	}
	return out, nil
}

// UpsertLineAnalysis は行解析キャッシュを 1 件保存する。本文ハッシュが既にあれば何もしない
// （同一本文の特徴は決定的で再計算不要）。
func (s *Store) UpsertLineAnalysis(ctx context.Context, a model.LineAnalysis) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO line_analysis
			(source_hash, sentence_count, polite_count, insult_count, is_imperative, exclaim_count, elong_count, emotion_count)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(source_hash) DO NOTHING`,
		a.SourceHash, a.SentenceCount, a.PoliteCount, a.InsultCount, a.IsImperative, a.ExclaimCount, a.ElongCount, a.EmotionCount)
	if err != nil {
		return fmt.Errorf("行解析キャッシュの保存: %w", err)
	}
	return nil
}

// ListSpeakerLineSources は台詞を所有する話者（名前付き）の、話者識別・声型・種族と本文 1 行を全件返す。
// 1 話者は複数行を持つため、呼び出し側が話者ごとに束ねて基底口調を集計する。
// 極端に短い台詞（15 文字以下）は対人・感情の判定に値しないため除く（PoC と同基準）。
func (s *Store) ListSpeakerLineSources(ctx context.Context) ([]model.SpeakerLineSource, error) {
	var rows []model.SpeakerLineSource
	query := `
		SELECT s.plugin AS speaker_plugin,
		       s.form_id AS speaker_form_id,
		       COALESCE(v.edid, '') AS voice_edid,
		       COALESCE(r.edid, '') AS race_edid,
		       l.source AS source
		FROM line_speaker ls
		JOIN speaker s ON s.id = ls.speaker_id
		JOIN line l ON l.id = ls.line_id
		LEFT JOIN voice_type v ON v.id = s.voice_type_id
		LEFT JOIN race r ON r.id = s.race_id
		WHERE TRIM(s.edid) <> '' AND TRIM(l.source) <> '' AND LENGTH(l.source) > 15
		ORDER BY s.plugin, s.form_id, ls.line_id`
	if err := s.db.SelectContext(ctx, &rows, query); err != nil {
		return nil, fmt.Errorf("生成入力（話者の台詞群）の取得: %w", err)
	}
	return rows, nil
}

// UpsertPersonaCharacter は生成した基底口調を 1 件保存する。既存なら更新するが、
// hand_edited=1（手修正済み）の行は更新しない（再生成で上書きしない保護）。
func (s *Store) UpsertPersonaCharacter(ctx context.Context, pc model.PersonaCharacter) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO persona_character
			(speaker_plugin, speaker_form_id, attitude_band, emotion_band, marked, decision_path, hand_edited)
		VALUES (?, ?, ?, ?, ?, ?, 0)
		ON CONFLICT(speaker_plugin, speaker_form_id) DO UPDATE SET
			attitude_band = excluded.attitude_band,
			emotion_band  = excluded.emotion_band,
			marked        = excluded.marked,
			decision_path = excluded.decision_path
		WHERE persona_character.hand_edited = 0`,
		pc.SpeakerPlugin, pc.SpeakerFormID, pc.AttitudeBand, pc.EmotionBand, pc.Marked, pc.DecisionPath)
	if err != nil {
		return fmt.Errorf("生成ペルソナの保存: %w", err)
	}
	return nil
}

// LoadLinePersonas は台詞 id 群に紐づく話者の生成済み基底口調を一括で読み、lineID→注入入力 の map を返す。
// persona_character を持つ話者の台詞だけが現れる。複数話者を指す純汎用台詞は id 昇順の先頭話者を採る。
func (s *Store) LoadLinePersonas(ctx context.Context, lineIDs []int64) (map[int64]model.LinePersonaInput, error) {
	out := make(map[int64]model.LinePersonaInput, len(lineIDs))
	if len(lineIDs) == 0 {
		return out, nil
	}
	for _, chunk := range chunkInt64(lineIDs, speakerChunkSize) {
		marks, args := placeholders(chunk)
		var rows []model.LinePersonaInput
		query := `
			SELECT ls.line_id AS line_id,
			       pc.attitude_band AS attitude_band,
			       pc.emotion_band AS emotion_band,
			       pc.marked AS marked,
			       pc.decision_path AS decision_path,
			       COALESCE(r.edid, '') AS race_edid
			FROM line_speaker ls
			JOIN speaker s ON s.id = ls.speaker_id
			JOIN persona_character pc ON pc.speaker_plugin = s.plugin AND pc.speaker_form_id = s.form_id
			LEFT JOIN race r ON r.id = s.race_id
			WHERE ls.line_id IN (` + marks + `)
			ORDER BY ls.line_id, s.id`
		if err := s.db.SelectContext(ctx, &rows, query, args...); err != nil {
			return nil, fmt.Errorf("注入入力（生成ペルソナ）の一括取得: %w", err)
		}
		for _, row := range rows {
			if _, seen := out[row.LineID]; seen {
				continue // 先頭話者（id 昇順）だけ採る。
			}
			out[row.LineID] = row
		}
	}
	return out, nil
}
