package store

import (
	"context"
	"fmt"

	"aitranslationenginejp/internal/model"
)

// GetPromptTemplate は単一行（id=1）のプロンプトテンプレートと話者なし台詞の口調設定を返す。
// base 指示は prompt_template（0004）、汎用・PC の口調と PC 性別は tone_default（0007）から、
// 単一行どうしの結合で 1 度に読む。engine（翻訳・実プロンプト再構成）と編集画面の初期表示が読む。
// 口調の雛形は directive テーブルの「口調」が唯一の正本で、prompt_template.persona_template は読まない。
func (s *Store) GetPromptTemplate(ctx context.Context) (model.PromptTemplate, error) {
	var t model.PromptTemplate
	if err := s.db.GetContext(ctx, &t,
		`SELECT p.base_directive,
		        d.generic_tone_text, d.pc_tone_text, d.pc_sex
		 FROM prompt_template p, tone_default d
		 WHERE p.id = 1 AND d.id = 1`); err != nil {
		return model.PromptTemplate{}, fmt.Errorf("prompt_template の取得: %w", err)
	}
	return t, nil
}

// SavePromptTemplate は base 指示（prompt_template）と、話者なし台詞の口調設定（tone_default）を
// 単一トランザクションで UPSERT する。編集画面の保存で使う。保存後の翻訳と実プロンプト再構成へ反映する。
// 口調の雛形は directive テーブルの「口調」が持つため、persona_template は書き換えない。
func (s *Store) SavePromptTemplate(ctx context.Context, t model.PromptTemplate) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("prompt_template 保存のトランザクション開始: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck // commit 済みなら no-op。失敗時の後始末用。
	// persona_template は NOT NULL のため INSERT 時だけ空文字を置く。列は残すが値は使わない
	// （口調の雛形は directive「口調」が持つ）。migration の seed が id=1 を作るので実際は UPDATE 側を通る。
	if _, err := tx.ExecContext(ctx,
		`INSERT INTO prompt_template (id, base_directive, persona_template) VALUES (1, ?, '')
		 ON CONFLICT(id) DO UPDATE SET
		   base_directive = excluded.base_directive`,
		t.BaseDirective); err != nil {
		return fmt.Errorf("prompt_template の保存: %w", err)
	}
	if _, err := tx.ExecContext(ctx,
		`INSERT INTO tone_default (id, generic_tone_text, pc_tone_text, pc_sex) VALUES (1, ?, ?, ?)
		 ON CONFLICT(id) DO UPDATE SET
		   generic_tone_text = excluded.generic_tone_text,
		   pc_tone_text = excluded.pc_tone_text,
		   pc_sex = excluded.pc_sex`,
		t.GenericToneText, t.PcToneText, t.PcSex); err != nil {
		return fmt.Errorf("tone_default の保存: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("prompt_template 保存のコミット: %w", err)
	}
	return nil
}
