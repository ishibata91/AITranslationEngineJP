package engine

import (
	"context"
	"fmt"
	"strings"

	"aitranslationenginejp/internal/model"
)

// ReferenceStore は engine が既存訳（参照訳）を読み書きするための中心データアクセス（使う分だけ宣言する）。
// InsertReferenceTranslations は取込（extracted_field の英日対 → 表）、ListReferenceTranslations は照合表の読み出し。
type ReferenceStore interface {
	InsertReferenceTranslations(ctx context.Context, refs []model.ReferenceTranslation) (int, error)
	ListReferenceTranslations(ctx context.Context) ([]model.ReferenceTranslation, error)
}

// referenceKey は既存訳の照合キー。reference_translation は対象横断で同一原文を再利用するため、
// form_id では絞らず (rec, field, source) で照合する。
type referenceKey struct {
	Rec    string
	Field  string
	Source string
}

// LoadReferenceTranslations は extracted_field の英日対（source・dest 非空行）から record 単位の既存訳を
// 取り込み、reference_translation へ追記する。抽出（extracted_field への書き込み）の後、翻訳 Run の前に呼ぶ。
// 追記した件数を返す。INSERT OR IGNORE のため二重実行でも増えない（冪等）。
func (e *Engine) LoadReferenceTranslations(ctx context.Context) (int, error) {
	fields, err := e.store.ListExtractedFields(ctx)
	if err != nil {
		return 0, fmt.Errorf("既存訳の供給元（extracted_field）の取得: %w", err)
	}
	refs := make([]model.ReferenceTranslation, 0, len(fields))
	for _, f := range fields {
		if strings.TrimSpace(f.Source) == "" || strings.TrimSpace(f.Dest) == "" {
			continue
		}
		refs = append(refs, model.ReferenceTranslation{Rec: f.Rec, Field: f.Field, Source: f.Source, Dest: f.Dest})
	}
	inserted, err := e.store.InsertReferenceTranslations(ctx, refs)
	if err != nil {
		return inserted, fmt.Errorf("既存訳の追記: %w", err)
	}
	return inserted, nil
}

// referenceIndex は reference_translation を (rec, field, source)→dest の照合表へ畳む。
// 翻訳ループが本文の既訳流用（完全一致置換）で 1 度だけ組んで使う。同一キーは先勝ち（先に読んだ既訳を残す）。
func (e *Engine) referenceIndex(ctx context.Context) (map[referenceKey]string, error) {
	refs, err := e.store.ListReferenceTranslations(ctx)
	if err != nil {
		return nil, fmt.Errorf("既存訳の取得: %w", err)
	}
	idx := make(map[referenceKey]string, len(refs))
	for _, r := range refs {
		k := referenceKey{Rec: r.Rec, Field: r.Field, Source: r.Source}
		if _, ok := idx[k]; !ok {
			idx[k] = r.Dest
		}
	}
	return idx, nil
}
