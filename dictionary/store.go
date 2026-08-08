// Package main は事前作成辞書の生成 command と MCP server を提供する。
package main

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var dictionarySchema string

var errRevisionConflict = errors.New("辞書の意味は取得後に更新されている")

type store struct {
	db *sqlx.DB
}

type sense struct {
	ID                   int64              `db:"id" json:"id"`
	TermID               int64              `db:"term_id" json:"term_id"`
	Source               string             `db:"source" json:"source"`
	Dest                 string             `db:"dest" json:"dest"`
	PartOfSpeech         string             `db:"part_of_speech" json:"part_of_speech"`
	Meaning              string             `db:"meaning" json:"meaning"`
	ClassificationStatus string             `db:"classification_status" json:"classification_status"`
	GeneralMatchStatus   string             `db:"general_match_status" json:"general_match_status"`
	InclusionDecision    string             `db:"inclusion_decision" json:"inclusion_decision"`
	ReviewStage          string             `db:"review_stage" json:"review_stage"`
	Revision             int64              `db:"revision" json:"revision"`
	MatchKind            string             `json:"match_kind,omitempty"`
	Occurrences          []occurrence       `json:"occurrences,omitempty"`
	GeneralMatches       []generalMatch     `json:"general_matches,omitempty"`
	Reviews              []dictionaryReview `json:"reviews,omitempty"`
}

type occurrence struct {
	ID              int64  `db:"id" json:"id"`
	TermID          int64  `db:"term_id" json:"term_id"`
	SenseID         *int64 `db:"sense_id" json:"sense_id,omitempty"`
	ObservedDest    string `db:"observed_dest" json:"observed_dest"`
	SkyrimCategory  string `db:"skyrim_category" json:"skyrim_category"`
	OriginKind      string `db:"origin_kind" json:"origin_kind"`
	OriginReference string `db:"origin_reference" json:"origin_reference"`
	DerivationKind  string `db:"derivation_kind" json:"derivation_kind"`
}

type generalMatch struct {
	ID                int64  `db:"id" json:"id"`
	SenseID           int64  `db:"sense_id" json:"sense_id"`
	DictionaryName    string `db:"dictionary_name" json:"dictionary_name"`
	DictionaryVersion string `db:"dictionary_version" json:"dictionary_version"`
	ExternalSenseID   string `db:"external_sense_id" json:"external_sense_id"`
	PartOfSpeech      string `db:"part_of_speech" json:"part_of_speech"`
	Definition        string `db:"definition" json:"definition"`
	JapaneseLemma     string `db:"japanese_lemma" json:"japanese_lemma"`
	MatchStatus       string `db:"match_status" json:"match_status"`
	Reason            string `db:"reason" json:"reason"`
}

type dictionaryReview struct {
	ID                int64  `db:"id" json:"id"`
	SenseID           int64  `db:"sense_id" json:"sense_id"`
	ReviewerKind      string `db:"reviewer_kind" json:"reviewer_kind"`
	ReviewerReference string `db:"reviewer_reference" json:"reviewer_reference"`
	Decision          string `db:"decision" json:"decision"`
	Reason            string `db:"reason" json:"reason"`
	CreatedAt         string `db:"created_at" json:"created_at"`
}

type senseUpdate struct {
	ID           int64
	Revision     int64
	Dest         string
	PartOfSpeech string
	Meaning      string
	ChangedBy    string
	Reason       string
}

func openStore(path string) (*store, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		return nil, fmt.Errorf("辞書DBのdirectory作成: %w", err)
	}
	db, err := sqlx.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("辞書DBを開く: %w", err)
	}
	db.SetMaxOpenConns(1)
	if _, err := db.ExecContext(context.Background(), `PRAGMA busy_timeout = 5000; PRAGMA foreign_keys = ON;`); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("辞書DBの接続設定: %w", err)
	}
	if _, err := db.ExecContext(context.Background(), dictionarySchema); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("辞書DBのschema適用: %w", err)
	}
	s := &store{db: db}
	if err := s.migrateLegacy(context.Background()); err != nil {
		_ = db.Close()
		return nil, err
	}
	return s, nil
}

func (s *store) close() error {
	if err := s.db.Close(); err != nil {
		return fmt.Errorf("辞書DBを閉じる: %w", err)
	}
	return nil
}

func (s *store) get(ctx context.Context, id int64) (sense, error) {
	var out sense
	if err := s.db.GetContext(ctx, &out, `
		SELECT s.id, s.term_id, t.source, s.dest, s.part_of_speech, s.meaning,
		       s.classification_status, s.general_match_status, s.inclusion_decision,
		       s.review_stage, s.revision
		FROM dictionary_sense s
		JOIN dictionary_term t ON t.id = s.term_id
		WHERE s.id = ?`, id); err != nil {
		return sense{}, fmt.Errorf("辞書の意味 %d の取得: %w", id, err)
	}
	if err := s.db.SelectContext(ctx, &out.Occurrences, `
		SELECT id, term_id, sense_id, observed_dest, skyrim_category,
		       origin_kind, origin_reference, derivation_kind
		FROM dictionary_occurrence WHERE sense_id = ? ORDER BY id`, id); err != nil {
		return sense{}, fmt.Errorf("辞書の意味 %d の使用箇所取得: %w", id, err)
	}
	if err := s.db.SelectContext(ctx, &out.GeneralMatches, `
		SELECT id, sense_id, dictionary_name, dictionary_version, external_sense_id,
		       part_of_speech, definition, japanese_lemma, match_status, reason
		FROM general_dictionary_match WHERE sense_id = ? ORDER BY id`, id); err != nil {
		return sense{}, fmt.Errorf("辞書の意味 %d の一般辞書照合取得: %w", id, err)
	}
	if err := s.db.SelectContext(ctx, &out.Reviews, `
		SELECT id, sense_id, reviewer_kind, reviewer_reference, decision, reason, created_at
		FROM dictionary_review WHERE sense_id = ? ORDER BY id`, id); err != nil {
		return sense{}, fmt.Errorf("辞書の意味 %d のレビュー取得: %w", id, err)
	}
	return out, nil
}

func (s *store) add(ctx context.Context, source, dest, partOfSpeech, meaning string) (sense, error) {
	source, dest = strings.TrimSpace(source), strings.TrimSpace(dest)
	partOfSpeech, meaning = strings.TrimSpace(partOfSpeech), strings.TrimSpace(meaning)
	if source == "" || dest == "" {
		return sense{}, errors.New("原語と訳語は空にできない")
	}
	if partOfSpeech == "" {
		partOfSpeech = "unknown"
	}
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return sense{}, fmt.Errorf("辞書の意味追加transaction開始: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck
	termID, err := ensureTerm(ctx, tx, source)
	if err != nil {
		return sense{}, err
	}
	res, err := tx.ExecContext(ctx, `
		INSERT INTO dictionary_sense
		    (term_id, dest, part_of_speech, meaning, classification_status)
		VALUES (?, ?, ?, ?, 'classified')`, termID, dest, partOfSpeech, meaning)
	if err != nil {
		return sense{}, fmt.Errorf("辞書の意味追加: %w", err)
	}
	senseID, err := res.LastInsertId()
	if err != nil {
		return sense{}, fmt.Errorf("追加した意味のid取得: %w", err)
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO dictionary_occurrence
		    (term_id, sense_id, observed_dest, origin_kind, origin_reference)
		VALUES (?, ?, ?, 'manual', ?)`, termID, senseID, dest, fmt.Sprintf("sense:%d", senseID)); err != nil {
		return sense{}, fmt.Errorf("手動追加の出どころ保存: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return sense{}, fmt.Errorf("辞書の意味追加commit: %w", err)
	}
	return s.get(ctx, senseID)
}

func (s *store) update(ctx context.Context, in senseUpdate) (sense, error) {
	in.Dest = strings.TrimSpace(in.Dest)
	in.PartOfSpeech = strings.TrimSpace(in.PartOfSpeech)
	in.Meaning = strings.TrimSpace(in.Meaning)
	in.ChangedBy = strings.TrimSpace(in.ChangedBy)
	in.Reason = strings.TrimSpace(in.Reason)
	if in.Dest == "" || in.PartOfSpeech == "" || in.ChangedBy == "" || in.Reason == "" {
		return sense{}, errors.New("訳語、品詞、変更者、変更理由は空にできない")
	}
	current, err := s.get(ctx, in.ID)
	if err != nil {
		return sense{}, err
	}
	if current.Revision != in.Revision {
		return sense{}, errRevisionConflict
	}
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return sense{}, fmt.Errorf("辞書の意味更新transaction開始: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck
	res, err := tx.ExecContext(ctx, `
		UPDATE dictionary_sense
		SET dest = ?, part_of_speech = ?, meaning = ?, classification_status = 'classified',
		    revision = revision + 1, updated_at = ?
		WHERE id = ? AND revision = ?`,
		in.Dest, in.PartOfSpeech, in.Meaning,
		time.Now().UTC().Format(time.RFC3339Nano), in.ID, in.Revision)
	if err != nil {
		return sense{}, fmt.Errorf("辞書の意味 %d の更新: %w", in.ID, err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return sense{}, fmt.Errorf("辞書の意味 %d の更新件数取得: %w", in.ID, err)
	}
	if n == 0 {
		return sense{}, errRevisionConflict
	}
	changes := []struct{ field, old, new string }{
		{"dest", current.Dest, in.Dest},
		{"part_of_speech", current.PartOfSpeech, in.PartOfSpeech},
		{"meaning", current.Meaning, in.Meaning},
	}
	for _, change := range changes {
		if change.old == change.new {
			continue
		}
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO dictionary_change
			    (target_table, target_id, field_name, old_value, new_value, changed_by, reason)
			VALUES ('dictionary_sense', ?, ?, ?, ?, ?, ?)`,
			in.ID, change.field, change.old, change.new, in.ChangedBy, in.Reason); err != nil {
			return sense{}, fmt.Errorf("辞書の意味 %d の変更履歴保存: %w", in.ID, err)
		}
	}
	if err := tx.Commit(); err != nil {
		return sense{}, fmt.Errorf("辞書の意味更新commit: %w", err)
	}
	return s.get(ctx, in.ID)
}

type status struct {
	Terms                  int            `json:"terms"`
	Senses                 int            `json:"senses"`
	Occurrences            int            `json:"occurrences"`
	AssignedOccurrences    int            `json:"assigned_occurrences"`
	GeneralMatches         int            `json:"general_matches"`
	Reviews                int            `json:"reviews"`
	ClassificationStatuses map[string]int `json:"classification_statuses"`
	GeneralMatchStatuses   map[string]int `json:"general_match_statuses"`
	InclusionDecisions     map[string]int `json:"inclusion_decisions"`
	ReviewStages           map[string]int `json:"review_stages"`
	ReviewDecisions        map[string]int `json:"review_decisions"`
	Origins                map[string]int `json:"origins"`
	TermFTSEntries         int            `json:"term_fts_entries"`
	SenseFTSEntries        int            `json:"sense_fts_entries"`
}

func (s *store) status(ctx context.Context) (status, error) {
	var out status
	counts := []struct {
		dest  *int
		query string
	}{
		{&out.Terms, `SELECT COUNT(*) FROM dictionary_term`},
		{&out.Senses, `SELECT COUNT(*) FROM dictionary_sense`},
		{&out.Occurrences, `SELECT COUNT(*) FROM dictionary_occurrence`},
		{&out.AssignedOccurrences, `SELECT COUNT(*) FROM dictionary_occurrence WHERE sense_id IS NOT NULL`},
		{&out.GeneralMatches, `SELECT COUNT(*) FROM general_dictionary_match`},
		{&out.Reviews, `SELECT COUNT(*) FROM dictionary_review`},
		{&out.TermFTSEntries, `SELECT COUNT(*) FROM dictionary_term_fts`},
		{&out.SenseFTSEntries, `SELECT COUNT(*) FROM dictionary_sense_fts`},
	}
	for _, count := range counts {
		if err := s.db.GetContext(ctx, count.dest, count.query); err != nil {
			return status{}, fmt.Errorf("辞書状態の件数取得: %w", err)
		}
	}
	var err error
	if out.ClassificationStatuses, err = s.countBy(ctx, "dictionary_sense", "classification_status"); err != nil {
		return status{}, err
	}
	if out.GeneralMatchStatuses, err = s.countBy(ctx, "dictionary_sense", "general_match_status"); err != nil {
		return status{}, err
	}
	if out.InclusionDecisions, err = s.countBy(ctx, "dictionary_sense", "inclusion_decision"); err != nil {
		return status{}, err
	}
	if out.ReviewStages, err = s.countBy(ctx, "dictionary_sense", "review_stage"); err != nil {
		return status{}, err
	}
	if out.ReviewDecisions, err = s.countBy(ctx, "dictionary_review", "decision"); err != nil {
		return status{}, err
	}
	if out.Origins, err = s.countBy(ctx, "dictionary_occurrence", "origin_kind"); err != nil {
		return status{}, err
	}
	return out, nil
}

func (s *store) countBy(ctx context.Context, table, column string) (map[string]int, error) {
	type row struct {
		Key   string `db:"key"`
		Count int    `db:"count"`
	}
	var rows []row
	query := fmt.Sprintf("SELECT %s AS key, COUNT(*) AS count FROM %s GROUP BY %s ORDER BY %s", column, table, column, column) //nolint:gosec // tableとcolumnは呼び出し元の固定値だけを使う。
	if err := s.db.SelectContext(ctx, &rows, query); err != nil {
		return nil, fmt.Errorf("%s.%s の件数取得: %w", table, column, err)
	}
	out := make(map[string]int, len(rows))
	for _, item := range rows {
		out[item.Key] = item.Count
	}
	return out, nil
}
