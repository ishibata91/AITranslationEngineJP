package main

import (
	"context"
	"fmt"
	"strings"
)

type searchFilter struct {
	Query              string
	Category           string
	GeneralMatchStatus string
	InclusionDecision  string
	ReviewStage        string
	Limit              int
}

type searchResult struct {
	Entries []sense `json:"entries"`
}

type searchCollector struct {
	store  *store
	filter searchFilter
	found  map[int64]sense
	order  []int64
}

func (s *store) search(ctx context.Context, filter searchFilter) (searchResult, error) {
	filter.Query = strings.TrimSpace(filter.Query)
	filter.Category = strings.TrimSpace(filter.Category)
	filter.GeneralMatchStatus = strings.TrimSpace(filter.GeneralMatchStatus)
	filter.InclusionDecision = strings.TrimSpace(filter.InclusionDecision)
	filter.ReviewStage = strings.TrimSpace(filter.ReviewStage)
	if filter.Limit <= 0 {
		filter.Limit = 50
	}
	if filter.Limit > 200 {
		filter.Limit = 200
	}
	collector := searchCollector{
		store:  s,
		filter: filter,
		found:  make(map[int64]sense),
		order:  make([]int64, 0, filter.Limit),
	}
	if filter.Query == "" {
		if err := collector.appendRows(ctx, "filtered", "1 = 1"); err != nil {
			return searchResult{}, fmt.Errorf("辞書の絞り込み: %w", err)
		}
		return collector.result(), nil
	}
	if err := collector.appendRows(ctx, "source_exact", `t.source = ? COLLATE NOCASE`, filter.Query); err != nil {
		return searchResult{}, fmt.Errorf("原語の完全一致検索: %w", err)
	}
	if err := collector.appendRows(ctx, "dest_exact", `s.dest = ? COLLATE NOCASE`, filter.Query); err != nil {
		return searchResult{}, fmt.Errorf("訳語の完全一致検索: %w", err)
	}
	prefix := escapeLike(filter.Query) + "%"
	upper := filter.Query + "\U0010FFFF"
	if err := collector.appendRows(ctx, "source_prefix", `
		t.source >= ? COLLATE NOCASE AND t.source < ? COLLATE NOCASE
		AND t.source LIKE ? ESCAPE '\' COLLATE NOCASE`, filter.Query, upper, prefix); err != nil {
		return searchResult{}, fmt.Errorf("原語の前方一致検索: %w", err)
	}
	if err := collector.appendRows(ctx, "dest_prefix", `
		s.dest >= ? COLLATE NOCASE AND s.dest < ? COLLATE NOCASE
		AND s.dest LIKE ? ESCAPE '\' COLLATE NOCASE`, filter.Query, upper, prefix); err != nil {
		return searchResult{}, fmt.Errorf("訳語の前方一致検索: %w", err)
	}
	if len([]rune(filter.Query)) >= 3 {
		phrase := `"` + strings.ReplaceAll(filter.Query, `"`, `""`) + `"`
		if err := collector.appendRows(ctx, "substring", `
			t.id IN (SELECT rowid FROM dictionary_term_fts WHERE dictionary_term_fts MATCH ?)`, phrase); err != nil {
			return searchResult{}, fmt.Errorf("原語の部分一致検索: %w", err)
		}
		if err := collector.appendRows(ctx, "substring", `
			s.id IN (SELECT rowid FROM dictionary_sense_fts WHERE dictionary_sense_fts MATCH ?)`, phrase); err != nil {
			return searchResult{}, fmt.Errorf("訳語と意味の部分一致検索: %w", err)
		}
	}
	return collector.result(), nil
}

func (c *searchCollector) appendRows(ctx context.Context, matchKind, where string, args ...any) error {
	if len(c.order) >= c.filter.Limit {
		return nil
	}
	query := `
		SELECT s.id, s.term_id, t.source, s.dest, s.part_of_speech, s.meaning,
		       s.classification_status, s.general_match_status, s.inclusion_decision,
		       s.review_stage, s.revision
		FROM dictionary_sense s
		JOIN dictionary_term t ON t.id = s.term_id
		WHERE ` + where
	query, args = c.appendFilters(query, args)
	query += ` ORDER BY length(t.source), t.source COLLATE NOCASE, s.dest, s.id LIMIT ?`
	args = append(args, c.filter.Limit)
	var rows []sense
	if err := c.store.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return fmt.Errorf("意味候補の検索: %w", err)
	}
	for _, item := range rows {
		if _, ok := c.found[item.ID]; ok || len(c.order) >= c.filter.Limit {
			continue
		}
		item.MatchKind = matchKind
		c.found[item.ID] = item
		c.order = append(c.order, item.ID)
	}
	return nil
}

func (c *searchCollector) appendFilters(query string, args []any) (string, []any) {
	if c.filter.Category != "" {
		query += ` AND EXISTS (
			SELECT 1 FROM dictionary_occurrence o
			WHERE o.sense_id = s.id AND o.skyrim_category = ?
		)`
		args = append(args, c.filter.Category)
	}
	if c.filter.GeneralMatchStatus != "" {
		query += ` AND s.general_match_status = ?`
		args = append(args, c.filter.GeneralMatchStatus)
	}
	if c.filter.InclusionDecision != "" {
		query += ` AND s.inclusion_decision = ?`
		args = append(args, c.filter.InclusionDecision)
	}
	if c.filter.ReviewStage != "" {
		query += ` AND s.review_stage = ?`
		args = append(args, c.filter.ReviewStage)
	}
	return query, args
}

func (c *searchCollector) result() searchResult {
	out := searchResult{Entries: make([]sense, 0, len(c.order))}
	for _, id := range c.order {
		out.Entries = append(out.Entries, c.found[id])
	}
	return out
}

func escapeLike(value string) string {
	replacer := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)
	return replacer.Replace(value)
}
