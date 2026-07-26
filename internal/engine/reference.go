package engine

import (
	"context"
	"fmt"

	"aitranslationenginejp/internal/model"
)

// ReferenceStore は engine が既存訳（参照訳）を読むための中心データアクセス（使う分だけ宣言する）。
// reference_translation へ書くのは C# 抽出器（Data フォルダ全 plugin の走査）で、engine は読むだけ。
// ListReferenceTranslations は照合表と横断辞書の派生が読む全件、CountReferenceTranslations は観測ログ用の件数。
type ReferenceStore interface {
	ListReferenceTranslations(ctx context.Context) ([]model.ReferenceTranslation, error)
	CountReferenceTranslations(ctx context.Context) (int, error)
}

// referenceKey は既存訳の照合キー。reference_translation は対象横断で同一原文を再利用するため、
// form_id では絞らず (rec, field, source) で照合する。
type referenceKey struct {
	Rec    string
	Field  string
	Source string
}

// CountReferenceTranslations は既訳（参照訳）の件数を返す。翻訳前区間が供給の観測ログへ出す。
func (e *Engine) CountReferenceTranslations(ctx context.Context) (int, error) {
	n, err := e.store.CountReferenceTranslations(ctx)
	if err != nil {
		return 0, fmt.Errorf("既存訳の件数取得: %w", err)
	}
	return n, nil
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
