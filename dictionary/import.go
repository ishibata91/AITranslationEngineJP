package main

import (
	"context"
	"fmt"
	"net/url"
	"path/filepath"
	"strconv"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

type importResult struct {
	Read      int `json:"read"`
	Created   int `json:"created"`
	Updated   int `json:"updated"`
	Unchanged int `json:"unchanged"`
}

type masterTerm struct {
	ID       int64  `db:"id"`
	Source   string `db:"source"`
	Dest     string `db:"dest"`
	Category string `db:"category"`
}

type importOutcome int

const (
	importCreated importOutcome = iota
	importUpdated
	importUnchanged
)

func importMasterTerms(ctx context.Context, sourcePath string, destination *store) (importResult, error) {
	abs, err := filepath.Abs(sourcePath)
	if err != nil {
		return importResult{}, fmt.Errorf("中心DBのpath解決: %w", err)
	}
	dsn := (&url.URL{Scheme: "file", Path: abs}).String() + "?mode=ro"
	source, err := sqlx.Open("sqlite", dsn)
	if err != nil {
		return importResult{}, fmt.Errorf("中心DBを読み取り専用で開く: %w", err)
	}
	defer source.Close() //nolint:errcheck
	rows, err := source.QueryxContext(ctx,
		`SELECT id, source, dest, category FROM master_term ORDER BY id`)
	if err != nil {
		return importResult{}, fmt.Errorf("中心DBのmaster_term取得: %w", err)
	}
	defer rows.Close() //nolint:errcheck

	tx, err := destination.db.BeginTxx(ctx, nil)
	if err != nil {
		return importResult{}, fmt.Errorf("移植transaction開始: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck

	var result importResult
	for rows.Next() {
		var term masterTerm
		if err := rows.StructScan(&term); err != nil {
			return importResult{}, fmt.Errorf("master_termの読み取り: %w", err)
		}
		result.Read++
		outcome, importErr := importMasterTerm(ctx, tx, term)
		if importErr != nil {
			return importResult{}, importErr
		}
		switch outcome {
		case importCreated:
			result.Created++
		case importUpdated:
			result.Updated++
		case importUnchanged:
			result.Unchanged++
		}
	}
	if err := rows.Err(); err != nil {
		return importResult{}, fmt.Errorf("master_termの走査: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return importResult{}, fmt.Errorf("移植transactionのcommit: %w", err)
	}
	return result, nil
}

func importMasterTerm(ctx context.Context, tx *sqlx.Tx, term masterTerm) (importOutcome, error) {
	reference := strconv.FormatInt(term.ID, 10)
	category, derivation := splitCategory(term.Category)
	type existingOccurrence struct {
		ID                   int64  `db:"id"`
		TermID               int64  `db:"term_id"`
		SenseID              *int64 `db:"sense_id"`
		Source               string `db:"source"`
		ObservedDest         string `db:"observed_dest"`
		SkyrimCategory       string `db:"skyrim_category"`
		DerivationKind       string `db:"derivation_kind"`
		ClassificationStatus string `db:"classification_status"`
		InclusionDecision    string `db:"inclusion_decision"`
		ReviewStage          string `db:"review_stage"`
	}
	var existing existingOccurrence
	err := tx.GetContext(ctx, &existing, `
		SELECT o.id, o.term_id, o.sense_id, t.source, o.observed_dest,
		       o.skyrim_category, o.derivation_kind,
		       COALESCE(s.classification_status, '') AS classification_status,
		       COALESCE(s.inclusion_decision, '') AS inclusion_decision,
		       COALESCE(s.review_stage, '') AS review_stage
		FROM dictionary_occurrence o
		JOIN dictionary_term t ON t.id = o.term_id
		LEFT JOIN dictionary_sense s ON s.id = o.sense_id
		WHERE o.origin_kind = 'master_term' AND o.origin_reference = ?`, reference)
	if isNoRows(err) {
		return insertMasterTerm(ctx, tx, term, reference, category, derivation)
	}
	if err != nil {
		return importUnchanged, fmt.Errorf("master_term %d の既存確認: %w", term.ID, err)
	}
	if existing.Source == term.Source && existing.ObservedDest == term.Dest &&
		existing.SkyrimCategory == category && existing.DerivationKind == derivation {
		return importUnchanged, nil
	}
	termID, err := ensureTerm(ctx, tx, term.Source)
	if err != nil {
		return importUnchanged, err
	}
	senseID, err := replaceImportedSense(ctx, tx, existing.SenseID, termID, term.Dest,
		existing.ClassificationStatus, existing.InclusionDecision, existing.ReviewStage)
	if err != nil {
		return importUnchanged, fmt.Errorf("master_term %d の意味候補更新: %w", term.ID, err)
	}
	if _, err := tx.ExecContext(ctx, `
		UPDATE dictionary_occurrence
		SET term_id = ?, sense_id = ?, observed_dest = ?, skyrim_category = ?,
		    derivation_kind = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?`, termID, senseID, term.Dest, category, derivation, existing.ID); err != nil {
		return importUnchanged, fmt.Errorf("master_term %d の使用箇所更新: %w", term.ID, err)
	}
	return importUpdated, nil
}

func insertMasterTerm(ctx context.Context, tx *sqlx.Tx, term masterTerm, reference, category, derivation string) (importOutcome, error) {
	termID, err := ensureTerm(ctx, tx, term.Source)
	if err != nil {
		return importUnchanged, err
	}
	res, err := tx.ExecContext(ctx, `
		INSERT INTO dictionary_sense (term_id, dest) VALUES (?, ?)`, termID, term.Dest)
	if err != nil {
		return importUnchanged, fmt.Errorf("master_term %d の意味候補追加: %w", term.ID, err)
	}
	senseID, err := res.LastInsertId()
	if err != nil {
		return importUnchanged, fmt.Errorf("master_term %d の意味候補id取得: %w", term.ID, err)
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO dictionary_occurrence
		    (term_id, sense_id, observed_dest, skyrim_category,
		     origin_kind, origin_reference, derivation_kind)
		VALUES (?, ?, ?, ?, 'master_term', ?, ?)`,
		termID, senseID, term.Dest, category, reference, derivation); err != nil {
		return importUnchanged, fmt.Errorf("master_term %d の使用箇所追加: %w", term.ID, err)
	}
	return importCreated, nil
}

func replaceImportedSense(ctx context.Context, tx *sqlx.Tx, currentSenseID *int64, termID int64, dest, classificationStatus, inclusionDecision, reviewStage string) (int64, error) {
	if currentSenseID != nil && classificationStatus != "classified" && inclusionDecision == "undecided" && reviewStage == "unreviewed" {
		var occurrences int
		if err := tx.GetContext(ctx, &occurrences,
			`SELECT COUNT(*) FROM dictionary_occurrence WHERE sense_id = ?`, *currentSenseID); err != nil {
			return 0, fmt.Errorf("意味候補の使用箇所数取得: %w", err)
		}
		if occurrences == 1 {
			if _, err := tx.ExecContext(ctx, `DELETE FROM general_dictionary_match WHERE sense_id = ?`, *currentSenseID); err != nil {
				return 0, fmt.Errorf("古い一般辞書照合の削除: %w", err)
			}
			if _, err := tx.ExecContext(ctx, `
				UPDATE dictionary_sense
				SET term_id = ?, dest = ?, classification_status = 'unclassified',
				    general_match_status = 'unchecked', revision = revision + 1,
				    updated_at = CURRENT_TIMESTAMP
				WHERE id = ?`, termID, dest, *currentSenseID); err != nil {
				return 0, fmt.Errorf("意味候補の再利用: %w", err)
			}
			return *currentSenseID, nil
		}
	}
	res, err := tx.ExecContext(ctx, `INSERT INTO dictionary_sense (term_id, dest) VALUES (?, ?)`, termID, dest)
	if err != nil {
		return 0, fmt.Errorf("新しい意味候補の追加: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("新しい意味候補のid取得: %w", err)
	}
	return id, nil
}
