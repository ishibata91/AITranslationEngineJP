package store

import (
	"context"
	"fmt"

	"aitranslationenginejp/internal/model"
)

// promptTemplateColumns は prompt_template の SELECT 列。model.PromptTemplate の db タグと対応する。
const promptTemplateColumns = `base_directive, persona_template`

// GetPromptTemplate は単一行（id=1）のプロンプトテンプレートを返す。
// テーブルは migration 0004 で既定値を seed 済みのため、通常は行が存在する。
// engine（翻訳実行・実プロンプト再構成）と編集画面の初期表示が読む。
func (s *Store) GetPromptTemplate(ctx context.Context) (model.PromptTemplate, error) {
	var t model.PromptTemplate
	if err := s.db.GetContext(ctx, &t,
		`SELECT `+promptTemplateColumns+` FROM prompt_template WHERE id = 1`); err != nil {
		return model.PromptTemplate{}, fmt.Errorf("prompt_template の取得: %w", err)
	}
	return t, nil
}

// SavePromptTemplate は単一行（id=1）へテンプレートを UPSERT する。編集画面の保存で使う。
// 保存後の翻訳と実プロンプト再構成は、この保存値を読んで反映する。
func (s *Store) SavePromptTemplate(ctx context.Context, t model.PromptTemplate) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO prompt_template (id, base_directive, persona_template) VALUES (1, ?, ?)
		 ON CONFLICT(id) DO UPDATE SET
		   base_directive = excluded.base_directive,
		   persona_template = excluded.persona_template`,
		t.BaseDirective, t.PersonaTemplate)
	if err != nil {
		return fmt.Errorf("prompt_template の保存: %w", err)
	}
	return nil
}
