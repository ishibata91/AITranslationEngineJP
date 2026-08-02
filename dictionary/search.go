package main

import (
	"context"
	"fmt"
	"strings"
)

type searchResult struct {
	Entries []entry `json:"entries"`
}

func (s *store) search(ctx context.Context, query, category string, limit int) (searchResult, error) {
	query, category = strings.TrimSpace(query), strings.TrimSpace(category)
	if query == "" {
		return searchResult{}, fmt.Errorf("検索文字列は空にできない")
	}
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}

	found := make(map[int64]entry)
	order := make([]int64, 0, limit)
	appendRows := func(matchKind, where string, args ...any) error {
		if len(order) >= limit {
			return nil
		}
		q := `SELECT id, source, dest, category, revision FROM dictionary_entry WHERE ` + where
		if category != "" {
			q += ` AND category = ?`
			args = append(args, category)
		}
		q += ` ORDER BY length(source), source COLLATE NOCASE, category, id LIMIT ?`
		args = append(args, limit)
		var rows []entry
		if err := s.db.SelectContext(ctx, &rows, q, args...); err != nil {
			return err
		}
		for _, e := range rows {
			if _, ok := found[e.ID]; ok || len(order) >= limit {
				continue
			}
			e.MatchKind = matchKind
			found[e.ID] = e
			order = append(order, e.ID)
		}
		return nil
	}

	if err := appendRows("source_exact", `source = ? COLLATE NOCASE`, query); err != nil {
		return searchResult{}, fmt.Errorf("原語の完全一致検索: %w", err)
	}
	if err := appendRows("dest_exact", `dest = ? COLLATE NOCASE`, query); err != nil {
		return searchResult{}, fmt.Errorf("訳語の完全一致検索: %w", err)
	}
	prefix := escapeLike(query) + "%"
	upper := query + "\U0010FFFF"
	if err := appendRows("source_prefix", `source >= ? COLLATE NOCASE AND source < ? COLLATE NOCASE AND source LIKE ? ESCAPE '\' COLLATE NOCASE`, query, upper, prefix); err != nil {
		return searchResult{}, fmt.Errorf("原語の前方一致検索: %w", err)
	}
	if err := appendRows("dest_prefix", `dest >= ? COLLATE NOCASE AND dest < ? COLLATE NOCASE AND dest LIKE ? ESCAPE '\' COLLATE NOCASE`, query, upper, prefix); err != nil {
		return searchResult{}, fmt.Errorf("訳語の前方一致検索: %w", err)
	}

	if len([]rune(query)) >= 3 && len(order) < limit {
		var ids []int64
		phrase := `"` + strings.ReplaceAll(query, `"`, `""`) + `"`
		if err := s.db.SelectContext(ctx, &ids,
			`SELECT rowid FROM dictionary_entry_fts WHERE dictionary_entry_fts MATCH ? LIMIT ?`, phrase, limit); err != nil {
			return searchResult{}, fmt.Errorf("部分一致検索: %w", err)
		}
		for _, id := range ids {
			if _, ok := found[id]; ok || len(order) >= limit {
				continue
			}
			e, err := s.get(ctx, id)
			if err != nil {
				return searchResult{}, err
			}
			if category != "" && e.Category != category {
				continue
			}
			e.MatchKind = "substring"
			found[id] = e
			order = append(order, id)
		}
	}

	out := searchResult{Entries: make([]entry, 0, len(order))}
	for _, id := range order {
		e := found[id]
		if len(e.Sources) == 0 {
			withSources, err := s.get(ctx, id)
			if err != nil {
				return searchResult{}, err
			}
			withSources.MatchKind = e.MatchKind
			e = withSources
		}
		out.Entries = append(out.Entries, e)
	}
	return out, nil
}

func escapeLike(s string) string {
	r := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)
	return r.Replace(s)
}
