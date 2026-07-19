package harness

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite" // SQLite driver（pure-Go）を "sqlite" 名で登録する。seed 用の独立接続で使う。
)

// SeedExtractor は本番 dotnet 抽出子（C#）の代わりに、合成 fixture を中心 DB へ直接 seed する fake 抽出子。
// api.Extractor を満たし、RunExtractAndTranslate の抽出段に差し込まれる。
// store の private 接続には触れず、同じ temp DB ファイルへ独立接続を開いて raw SQL で投入する
// （store は extracted_field の投入 API を公開しないため、抽出子が書く層を直接模す）。
type SeedExtractor struct {
	DBPath  string
	Fixture Fixture
}

// Extract は fixture を中心 DB へ seed する。pluginPath は本番の plugin 選択に対応するが、fake は内容を fixture で固定する。
func (s *SeedExtractor) Extract(_ context.Context, _ string) error {
	db, err := sqlx.Connect("sqlite", s.DBPath)
	if err != nil {
		return fmt.Errorf("seed 用 DB を開けない: %w", err)
	}
	defer db.Close() //nolint:errcheck // 読み書き後の後始末。失敗しても seed 自体の成否は下の戻り値で判断する。
	return seedFixture(db, s.Fixture)
}

// seedFixture は話者素材（race・voice_type・speaker・INFO 橋渡し）と extracted_field を 1 トランザクションで投入する。
// 話者素材は extracted_field より先に入れ、取込段の話者連関解決（LinkLineSpeakersFromStaging）が成り立つようにする。
func seedFixture(db *sqlx.DB, f Fixture) error {
	tx, err := db.Beginx()
	if err != nil {
		return fmt.Errorf("seed トランザクション開始: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck // commit 済みなら no-op。失敗時の後始末用。

	for _, sp := range f.Speakers {
		if err := seedSpeaker(tx, sp); err != nil {
			return err
		}
	}
	for _, ef := range f.ExtractedFields {
		if _, err := tx.Exec(
			`INSERT INTO extracted_field (plugin, form_id, edid, rec, field, ordinal, source)
			 VALUES (?, ?, ?, ?, ?, ?, ?)`,
			ef.Plugin, ef.FormID, ef.EDID, ef.Rec, ef.Field, ef.Ordinal, ef.Source); err != nil {
			return fmt.Errorf("extracted_field の seed (%s:%s %s): %w", ef.Rec, ef.Field, ef.FormID, err)
		}
	}
	for _, em := range f.Emotions {
		if _, err := tx.Exec(
			`INSERT INTO extracted_info_emotion (info_plugin, info_form_id, ordinal, emotion_type)
			 VALUES (?, ?, ?, ?)`,
			em.Plugin, em.InfoFormID, em.Ordinal, em.EmotionType); err != nil {
			return fmt.Errorf("extracted_info_emotion の seed (%s ord=%d): %w", em.InfoFormID, em.Ordinal, err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("seed コミット: %w", err)
	}
	return nil
}

// seedSpeaker は 1 話者の素材（種族・声型・話者・INFO 橋渡し）を投入する。
// race・voice_type を先に入れて id を得て speaker の FK に結び、INFO（plugin・form_id）→ speaker.id の staging を残す。
// race・voice_type は (plugin, form_id) が UNIQUE で、実データでは複数話者が同じ種族・声型を共有する。
// そのため INSERT OR IGNORE して既存なら id を引き直し、同一 form_id の重複投入を許す（共有を再現する）。
func seedSpeaker(tx *sqlx.Tx, sp SeedSpeaker) error {
	raceID, err := ensureRowID(tx,
		`INSERT OR IGNORE INTO race (nature, plugin, form_id, edid) VALUES (?, ?, ?, ?)`,
		[]any{sp.RaceNature, sp.Plugin, sp.RaceFormID, sp.RaceEDID},
		`SELECT id FROM race WHERE plugin = ? AND form_id = ?`,
		[]any{sp.Plugin, sp.RaceFormID})
	if err != nil {
		return fmt.Errorf("race の seed (%s): %w", sp.RaceEDID, err)
	}
	voiceID, err := ensureRowID(tx,
		`INSERT OR IGNORE INTO voice_type (voice_id, voice_kind, nature, plugin, form_id, edid) VALUES (?, ?, ?, ?, ?, ?)`,
		[]any{sp.VoiceEDID, sp.VoiceKind, sp.VoiceNature, sp.Plugin, sp.VoiceFormID, sp.VoiceEDID},
		`SELECT id FROM voice_type WHERE plugin = ? AND form_id = ?`,
		[]any{sp.Plugin, sp.VoiceFormID})
	if err != nil {
		return fmt.Errorf("voice_type の seed (%s): %w", sp.VoiceEDID, err)
	}
	speakerID, err := ensureRowID(tx,
		`INSERT OR IGNORE INTO speaker (speaker_kind, sex, race_id, voice_type_id, plugin, form_id, edid)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		[]any{"NPC", sp.Sex, raceID, voiceID, sp.Plugin, sp.FormID, sp.EDID},
		`SELECT id FROM speaker WHERE plugin = ? AND form_id = ?`,
		[]any{sp.Plugin, sp.FormID})
	if err != nil {
		return fmt.Errorf("speaker の seed (%s): %w", sp.EDID, err)
	}
	if _, err := tx.Exec(
		`INSERT OR IGNORE INTO extracted_info_speaker (info_plugin, info_form_id, speaker_id) VALUES (?, ?, ?)`,
		sp.Plugin, sp.InfoFormID, speakerID); err != nil {
		return fmt.Errorf("INFO 橋渡しの seed (%s): %w", sp.InfoFormID, err)
	}
	return nil
}

// ensureRowID は INSERT OR IGNORE で行を確保し、(plugin, form_id) などの UNIQUE キーで id を引き直して返す。
// 既存行があっても同一 id を返すため、共有素材（複数話者が同じ種族・声型を指す）を冪等に seed できる。
func ensureRowID(tx *sqlx.Tx, insertSQL string, insertArgs []any, selectSQL string, selectArgs []any) (int64, error) {
	if _, err := tx.Exec(insertSQL, insertArgs...); err != nil {
		return 0, fmt.Errorf("INSERT 実行: %w", err)
	}
	var id int64
	if err := tx.Get(&id, selectSQL, selectArgs...); err != nil {
		return 0, fmt.Errorf("id の引き直し: %w", err)
	}
	return id, nil
}
