package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

type matchQueueFilter struct {
	Status            string
	InclusionDecision string
	ReviewStage       string
	AfterSenseID      int64
	Limit             int
}

type matchQueueItem struct {
	SenseID           int64          `db:"sense_id" json:"sense_id"`
	Revision          int64          `db:"revision" json:"revision"`
	Source            string         `db:"source" json:"source"`
	Dest              string         `db:"dest" json:"dest"`
	PartOfSpeech      string         `db:"part_of_speech" json:"part_of_speech"`
	Meaning           string         `db:"meaning" json:"meaning"`
	SkyrimCategories  string         `db:"skyrim_categories" json:"skyrim_categories"`
	InclusionDecision string         `db:"inclusion_decision" json:"inclusion_decision"`
	ReviewStage       string         `db:"review_stage" json:"review_stage"`
	Matches           []generalMatch `json:"matches"`
}

type matchQueueResult struct {
	Entries     []matchQueueItem `json:"entries"`
	NextSenseID int64            `json:"next_sense_id,omitempty"`
}

func (s *store) generalMatchQueue(ctx context.Context, filter matchQueueFilter) (matchQueueResult, error) {
	filter, err := normalizeMatchQueueFilter(filter)
	if err != nil {
		return matchQueueResult{}, err
	}
	out := matchQueueResult{Entries: []matchQueueItem{}}
	if err := s.db.SelectContext(ctx, &out.Entries, `
		SELECT s.id AS sense_id, s.revision, t.source, s.dest, s.part_of_speech, s.meaning,
		       COALESCE((
		           SELECT group_concat(category, ',')
		           FROM (
		               SELECT DISTINCT CASE
		                   WHEN o.skyrim_category <> '' THEN o.skyrim_category
		                   WHEN o.derivation_kind <> '' THEN 'derive:' || o.derivation_kind
		                   ELSE ''
		               END AS category
		               FROM dictionary_occurrence o
		               WHERE o.sense_id = s.id
		               ORDER BY category
		           )
		       ), '') AS skyrim_categories,
		       s.inclusion_decision, s.review_stage
		FROM dictionary_sense s
		JOIN dictionary_term t ON t.id = s.term_id
		WHERE s.id > ?
		  AND (? = '' OR s.inclusion_decision = ?)
		  AND (? = '' OR s.review_stage = ?)
		  AND EXISTS (
		      SELECT 1 FROM general_dictionary_match m
		      WHERE m.sense_id = s.id AND m.match_status = ?
		  )
		ORDER BY s.id
		LIMIT ?`, filter.AfterSenseID, filter.InclusionDecision,
		filter.InclusionDecision, filter.ReviewStage, filter.ReviewStage,
		filter.Status, filter.Limit); err != nil {
		return matchQueueResult{}, fmt.Errorf("一般辞書照合の確認対象取得: %w", err)
	}
	for i := range out.Entries {
		if err := s.loadQueueMatches(ctx, &out.Entries[i], filter.Status); err != nil {
			return matchQueueResult{}, err
		}
	}
	if len(out.Entries) > 0 {
		out.NextSenseID = out.Entries[len(out.Entries)-1].SenseID
	}
	return out, nil
}

func normalizeMatchQueueFilter(filter matchQueueFilter) (matchQueueFilter, error) {
	filter.Status = strings.TrimSpace(filter.Status)
	if filter.Status == "" {
		filter.Status = "same_mean_candidate"
	}
	if !isQueueStatus(filter.Status) {
		return matchQueueFilter{}, errors.New("一般辞書照合の確認対象にできない状態")
	}
	filter.InclusionDecision = strings.TrimSpace(filter.InclusionDecision)
	if filter.InclusionDecision != "" && filter.InclusionDecision != "undecided" &&
		filter.InclusionDecision != "include" && filter.InclusionDecision != "exclude" {
		return matchQueueFilter{}, errors.New("収録判断の絞り込み値が正しくない")
	}
	filter.ReviewStage = strings.TrimSpace(filter.ReviewStage)
	if filter.ReviewStage != "" && filter.ReviewStage != "unreviewed" &&
		filter.ReviewStage != "ai_reviewed" && filter.ReviewStage != "human_reviewed" {
		return matchQueueFilter{}, errors.New("レビュー段階の絞り込み値が正しくない")
	}
	if filter.Limit <= 0 {
		filter.Limit = 50
	}
	if filter.Limit > 200 {
		filter.Limit = 200
	}
	return filter, nil
}

func (s *store) loadQueueMatches(ctx context.Context, item *matchQueueItem, status string) error {
	item.Matches = []generalMatch{}
	if err := s.db.SelectContext(ctx, &item.Matches, `
		SELECT id, sense_id, dictionary_name, dictionary_version, external_sense_id,
		       part_of_speech, definition, japanese_lemma, match_status, reason
		FROM general_dictionary_match
		WHERE sense_id = ? AND match_status = ?
		ORDER BY id`, item.SenseID, status); err != nil {
		return fmt.Errorf("意味 %d の一般辞書照合候補取得: %w", item.SenseID, err)
	}
	return nil
}

func isQueueStatus(status string) bool {
	switch status {
	case "same_mean_candidate", "same_mean_and_translation", "different_meaning_or_translation":
		return true
	default:
		return false
	}
}
