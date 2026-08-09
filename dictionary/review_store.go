package main

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type reviewInput struct {
	SenseID           int64
	Revision          int64
	ReviewerKind      string
	ReviewerReference string
	Decision          string
	Reason            string
}

func (s *store) addReview(ctx context.Context, in reviewInput) (sense, error) {
	in, err := normalizeReviewInput(in)
	if err != nil {
		return sense{}, err
	}
	current, err := s.get(ctx, in.SenseID)
	if err != nil {
		return sense{}, err
	}
	if err := validateReviewTransition(current, in); err != nil {
		return sense{}, err
	}
	inclusionDecision, reviewStage := reviewSummary(in)
	if err := s.persistReview(ctx, current, in, inclusionDecision, reviewStage); err != nil {
		return sense{}, err
	}
	return s.get(ctx, in.SenseID)
}

func normalizeReviewInput(in reviewInput) (reviewInput, error) {
	in.ReviewerKind = strings.TrimSpace(in.ReviewerKind)
	in.ReviewerReference = strings.TrimSpace(in.ReviewerReference)
	in.Decision = strings.TrimSpace(in.Decision)
	in.Reason = strings.TrimSpace(in.Reason)
	if in.ReviewerReference == "" || in.Reason == "" {
		return reviewInput{}, errors.New("レビュー担当と判断理由は空にできない")
	}
	if in.ReviewerKind != "ai" && in.ReviewerKind != "human" {
		return reviewInput{}, errors.New("reviewer_kindはaiまたはhumanを指定する")
	}
	if in.Decision != "include" && in.Decision != "exclude" && in.Decision != "needs_human" {
		return reviewInput{}, errors.New("decisionはinclude、exclude、needs_humanのいずれかを指定する")
	}
	return in, nil
}

func validateReviewTransition(current sense, in reviewInput) error {
	if current.Revision != in.Revision {
		return errRevisionConflict
	}
	if current.ReviewStage == "human_reviewed" && in.ReviewerKind == "ai" {
		return errors.New("人間レビュー済みの判断をAIレビューで変更できない")
	}
	return nil
}

func reviewSummary(in reviewInput) (string, string) {
	inclusionDecision := in.Decision
	if in.Decision == "needs_human" {
		inclusionDecision = "undecided"
	}
	reviewStage := "ai_reviewed"
	if in.ReviewerKind == "human" {
		reviewStage = "human_reviewed"
	}
	return inclusionDecision, reviewStage
}

func (s *store) persistReview(ctx context.Context, current sense, in reviewInput, inclusionDecision, reviewStage string) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("レビュー保存transaction開始: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck
	if in.Decision == "exclude" {
		if validationErr := validateGeneralWordExclusion(ctx, tx, in.SenseID); validationErr != nil {
			return validationErr
		}
	}
	if _, execErr := tx.ExecContext(ctx, `
		INSERT INTO dictionary_review
		    (sense_id, reviewer_kind, reviewer_reference, decision, reason)
		VALUES (?, ?, ?, ?, ?)`, in.SenseID, in.ReviewerKind,
		in.ReviewerReference, in.Decision, in.Reason); execErr != nil {
		return fmt.Errorf("意味 %d のレビュー保存: %w", in.SenseID, execErr)
	}
	res, err := tx.ExecContext(ctx, `
		UPDATE dictionary_sense
		SET inclusion_decision = ?, review_stage = ?, classification_status = 'classified',
		    revision = revision + 1, updated_at = CURRENT_TIMESTAMP
		WHERE id = ? AND revision = ?`, inclusionDecision, reviewStage, in.SenseID, in.Revision)
	if err != nil {
		return fmt.Errorf("意味 %d のレビュー状態更新: %w", in.SenseID, err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("意味 %d のレビュー更新件数取得: %w", in.SenseID, err)
	}
	if n == 0 {
		return errRevisionConflict
	}
	for _, change := range []struct{ field, old, new string }{
		{"inclusion_decision", current.InclusionDecision, inclusionDecision},
		{"review_stage", current.ReviewStage, reviewStage},
	} {
		if change.old == change.new {
			continue
		}
		if _, execErr := tx.ExecContext(ctx, `
			INSERT INTO dictionary_change
			    (target_table, target_id, field_name, old_value, new_value, changed_by, reason)
			VALUES ('dictionary_sense', ?, ?, ?, ?, ?, ?)`, in.SenseID, change.field,
			change.old, change.new, in.ReviewerReference, in.Reason); execErr != nil {
			return fmt.Errorf("意味 %d のレビュー変更履歴保存: %w", in.SenseID, execErr)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("レビュー保存commit: %w", err)
	}
	return nil
}

func validateGeneralWordExclusion(ctx context.Context, tx *sqlx.Tx, senseID int64) error {
	var confirmed int
	if err := tx.GetContext(ctx, &confirmed, `
		SELECT COUNT(*)
		FROM general_dictionary_match
		WHERE sense_id = ?
		  AND match_status = 'same_mean_and_translation'
		  AND trim(dictionary_version) <> ''
		  AND trim(external_sense_id) <> ''
		  AND trim(reason) <> ''`, senseID); err != nil {
		return fmt.Errorf("意味 %d の一般辞書照合確認: %w", senseID, err)
	}
	if confirmed == 0 {
		return errors.New("同じ意味と訳を確認した一般辞書の意味がないため除外できない")
	}
	return nil
}
