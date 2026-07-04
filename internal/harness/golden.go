package harness

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

// Capture は 1 回の RunExtractAndTranslate の観測結果。送信プロンプト列（順序つき）と DB 最終状態を持つ。
// この 2 つが非劣化の判定対象で、リファクタ前後で Serialize 出力が一致すれば決定的振る舞いは保たれている。
type Capture struct {
	Model           string
	TranslatedCount int // RunExtractAndTranslate が返した翻訳件数（取込振り分けが壊れて件数が変われば検出できる）
	Prompts         []RecordedPrompt
	DBState         string
}

// Serialize は Capture を golden 比較用の決定的な 1 文字列へ描く。
// 前半に送信プロンプト列（送信順・モデル・system・user）、後半に DB 最終状態を置く。
func (c Capture) Serialize() string {
	var b strings.Builder
	fmt.Fprintf(&b, "== SUMMARY ==\ntranslated_count=%d\n\n", c.TranslatedCount)
	b.WriteString("== PROMPTS (送信順) ==\n")
	for i, p := range c.Prompts {
		fmt.Fprintf(&b, "[%03d] model=%s\n", i, p.Model)
		b.WriteString("--- system ---\n")
		b.WriteString(p.System)
		b.WriteString("\n--- user ---\n")
		b.WriteString(p.User)
		b.WriteString("\n\n")
	}
	b.WriteString("== DB FINAL STATE ==\n")
	b.WriteString(c.DBState)
	return b.String()
}

// dbStateTable は captureDBState がダンプする 1 テーブルの仕様。query が ORDER BY で順序を固定し、format が 1 行を文字列化する。
type dbStateTable struct {
	table     string
	orderNote string
	query     string
	format    func(*sqlx.Rows) (string, error)
}

// dbStateTables は captureDBState がダンプするテーブルの仕様列。golden の出力順はこの並びで固定する。
// 取込・固有名確定・本文書き戻し・話者連関・口調生成・固有名派生の結果が現れるテーブルを並べる。
// 各テーブルの ORDER BY は自動採番 id でなく自然キー（UNIQUE 制約の列）で固定する。取込の INSERT 順が変わっても
// golden が揺れないようにし、id 採番順の差を内容劣化と取り違えないため。訳順の劣化は dest（連番訳）の値で別途観測される。
var dbStateTables = []dbStateTable{
	{"proper_noun", "category, source",
		`SELECT source, category, dest, status FROM proper_noun ORDER BY category, source`,
		formatProperNounRow},
	{"narration", "plugin, form_id, rec, field, ordinal",
		`SELECT rec, field, source, dest, status, style FROM narration ORDER BY plugin, form_id, rec, field, ordinal`,
		formatNarrationRow},
	// line.source は原文を保持する（機械置換は翻訳時の一過性で DB へは書かない）。本文への辞書適用の結果は
	// PROMPTS の user 行（置換済み原文）が正本の観測点で、dict.Apply の劣化はそこで検出する。
	{"line", "plugin, form_id, rec, field, ordinal",
		`SELECT rec, field, source, dest, status FROM line ORDER BY plugin, form_id, rec, field, ordinal`,
		formatLineRow},
	// line_speaker は raw な数値 id でなく、台詞・話者の自然キー（form_id）と話者 EDID で捕獲する。
	// 話者連関の解決（LinkLineSpeakersFromStaging）が別の話者へ誤接続した場合に、id の見かけでなく実体で検出する。
	// ORDER BY は form_id の文字列辞書順。値の総順序が一意なため golden は決定的に安定する（数値順との差は表示順だけ）。
	{"line_speaker", "line.form_id, speaker.form_id",
		`SELECT l.form_id AS line_form, s.edid AS speaker_edid, s.form_id AS speaker_form
		 FROM line_speaker ls
		 JOIN line l ON l.id = ls.line_id
		 JOIN speaker s ON s.id = ls.speaker_id
		 ORDER BY l.form_id, s.form_id`,
		formatLineSpeakerRow},
	// 言及テーブル（e4/e5）は注入の忠実な記録。stoplist の供給源選別が注入と言及の両方へ同じに
	// 効くこと（stoplist 語は言及されず、stoplist 外は従来どおり残ること）を golden で凍結する。
	// 行は raw な数値 id でなく、言及元（narration / line の自然キー）と相手（語彙の原語と供給源種別）で捕獲する。
	{"narration_mention", "narration の自然キー, 相手の原語",
		`SELECT n.rec, n.field, n.form_id, n.ordinal,
		        COALESCE(pn.source, mt.source) AS term_source,
		        CASE WHEN nm.proper_noun_id IS NOT NULL THEN 'proper_noun' ELSE 'master_term' END AS kind
		 FROM narration_mention nm
		 JOIN narration n ON n.id = nm.narration_id
		 LEFT JOIN proper_noun pn ON pn.id = nm.proper_noun_id
		 LEFT JOIN master_term mt ON mt.id = nm.master_term_id
		 ORDER BY n.plugin, n.form_id, n.rec, n.field, n.ordinal, term_source`,
		formatMentionRow},
	{"line_mention", "line の自然キー, 相手の原語",
		`SELECT l.rec, l.field, l.form_id, l.ordinal,
		        COALESCE(pn.source, mt.source) AS term_source,
		        CASE WHEN lm.proper_noun_id IS NOT NULL THEN 'proper_noun' ELSE 'master_term' END AS kind
		 FROM line_mention lm
		 JOIN line l ON l.id = lm.line_id
		 LEFT JOIN proper_noun pn ON pn.id = lm.proper_noun_id
		 LEFT JOIN master_term mt ON mt.id = lm.master_term_id
		 ORDER BY l.plugin, l.form_id, l.rec, l.field, l.ordinal, term_source`,
		formatMentionRow},
	{"persona_character", "speaker_plugin, speaker_form_id",
		`SELECT speaker_plugin, speaker_form_id, attitude_band, emotion_band, marked, decision_path FROM persona_character ORDER BY speaker_plugin, speaker_form_id`,
		formatPersonaCharacterRow},
	{"master_term", "category, source",
		`SELECT source, dest, category FROM master_term ORDER BY category, source`,
		formatMasterTermRow},
	// line_analysis は口調生成が書く行特徴キャッシュ。リファクタで特徴量計算が変わると訳の口調が変わるため、
	// 特徴量の劣化を直接観測できるよう捕獲する。id は map 反復順で揺れるため source_hash 順で安定化する。
	{"line_analysis", "source_hash",
		`SELECT source_hash, sentence_count, polite_count, insult_count, is_imperative, exclaim_count, elong_count, emotion_count
		 FROM line_analysis ORDER BY source_hash`,
		formatLineAnalysisRow},
}

// captureDBState は中心 DB の最終状態を決定的な文字列へ描く。dbStateTables の仕様順に各テーブルをダンプする。
func captureDBState(db *sqlx.DB) (string, error) {
	var b strings.Builder
	for _, t := range dbStateTables {
		if err := dumpRows(&b, db, t.table, t.orderNote, t.query, t.format); err != nil {
			return "", err
		}
	}
	return b.String(), nil
}

// formatProperNounRow は proper_noun の 1 行を描く。
func formatProperNounRow(rows *sqlx.Rows) (string, error) {
	var status int
	var source, category, dest string
	if err := rows.Scan(&source, &category, &dest, &status); err != nil {
		return "", fmt.Errorf("行の読み取り: %w", err)
	}
	return fmt.Sprintf("source=%q category=%q dest=%q status=%d", source, category, dest, status), nil
}

// formatNarrationRow は narration の 1 行を描く。
func formatNarrationRow(rows *sqlx.Rows) (string, error) {
	var status int
	var rec, field, source, dest, style string
	if err := rows.Scan(&rec, &field, &source, &dest, &status, &style); err != nil {
		return "", fmt.Errorf("行の読み取り: %w", err)
	}
	return fmt.Sprintf("%s:%s source=%q dest=%q status=%d style=%q", rec, field, source, dest, status, style), nil
}

// formatLineRow は line の 1 行を描く。
func formatLineRow(rows *sqlx.Rows) (string, error) {
	var status int
	var rec, field, source, dest string
	if err := rows.Scan(&rec, &field, &source, &dest, &status); err != nil {
		return "", fmt.Errorf("行の読み取り: %w", err)
	}
	return fmt.Sprintf("%s:%s source=%q dest=%q status=%d", rec, field, source, dest, status), nil
}

// formatLineSpeakerRow は line_speaker の 1 行を描く。
func formatLineSpeakerRow(rows *sqlx.Rows) (string, error) {
	var lineForm, speakerEDID, speakerForm string
	if err := rows.Scan(&lineForm, &speakerEDID, &speakerForm); err != nil {
		return "", fmt.Errorf("行の読み取り: %w", err)
	}
	return fmt.Sprintf("line=%s speaker_edid=%q speaker_form=%s", lineForm, speakerEDID, speakerForm), nil
}

// formatMentionRow は narration_mention / line_mention の 1 行を描く（言及元の自然キーと相手の語彙）。
func formatMentionRow(rows *sqlx.Rows) (string, error) {
	var ordinal int
	var rec, field, formID, termSource, kind string
	if err := rows.Scan(&rec, &field, &formID, &ordinal, &termSource, &kind); err != nil {
		return "", fmt.Errorf("行の読み取り: %w", err)
	}
	return fmt.Sprintf("%s:%s form=%s ord=%d term=%q kind=%s", rec, field, formID, ordinal, termSource, kind), nil
}

// formatPersonaCharacterRow は persona_character の 1 行を描く。
func formatPersonaCharacterRow(rows *sqlx.Rows) (string, error) {
	var plugin, formID, path string
	var att, emo, marked int
	if err := rows.Scan(&plugin, &formID, &att, &emo, &marked, &path); err != nil {
		return "", fmt.Errorf("行の読み取り: %w", err)
	}
	return fmt.Sprintf("%s/%s attitude=%d emotion=%d marked=%d path=%q", plugin, formID, att, emo, marked, path), nil
}

// formatMasterTermRow は master_term の 1 行を描く。
func formatMasterTermRow(rows *sqlx.Rows) (string, error) {
	var source, dest, category string
	if err := rows.Scan(&source, &dest, &category); err != nil {
		return "", fmt.Errorf("行の読み取り: %w", err)
	}
	return fmt.Sprintf("source=%q dest=%q category=%q", source, dest, category), nil
}

// formatLineAnalysisRow は line_analysis の 1 行を描く。
func formatLineAnalysisRow(rows *sqlx.Rows) (string, error) {
	var hash string
	var sent, polite, insult, imper, exclaim, elong, emotion int
	if err := rows.Scan(&hash, &sent, &polite, &insult, &imper, &exclaim, &elong, &emotion); err != nil {
		return "", fmt.Errorf("行の読み取り: %w", err)
	}
	return fmt.Sprintf("hash=%s sent=%d polite=%d insult=%d imper=%d exclaim=%d elong=%d emotion=%d",
		hash, sent, polite, insult, imper, exclaim, elong, emotion), nil
}

// dumpRows は 1 テーブルを決定的な行テキストへ描く。query は ORDER BY で順序を固定し、format が 1 行を文字列化する。
func dumpRows(b *strings.Builder, db *sqlx.DB, table, orderNote, query string, format func(*sqlx.Rows) (string, error)) error {
	rows, err := db.Queryx(query)
	if err != nil {
		return fmt.Errorf("%s の取得: %w", table, err)
	}
	defer func() { _ = rows.Close() }()
	fmt.Fprintf(b, "-- %s (order by %s) --\n", table, orderNote)
	for rows.Next() {
		line, err := format(rows)
		if err != nil {
			return fmt.Errorf("%s の行整形: %w", table, err)
		}
		b.WriteString(line)
		b.WriteString("\n")
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("%s の走査: %w", table, err)
	}
	b.WriteString("\n")
	return nil
}
