package store

import (
	"context"
	"fmt"

	"aitranslationenginejp/internal/model"
)

// directiveColumns は directive の SELECT 列。model.Directive の db タグと対応する。
const directiveColumns = `key, instruction, variables`

// recordTypeColumns は record_type_master の SELECT 列。model.RecordType の db タグと対応する。
const recordTypeColumns = `rec, field, box, directive, logical_name`

// ListDirectives は指示文の正本を全件返す（編集画面のレコード別タブと、本文翻訳の口調指示に使う）。
// 件数は固定（migration の seed 数）で、ページングせず一括で読む。
func (s *Store) ListDirectives(ctx context.Context) ([]model.Directive, error) {
	var rows []model.Directive
	if err := s.db.SelectContext(ctx, &rows,
		`SELECT `+directiveColumns+` FROM directive ORDER BY key`); err != nil {
		return nil, fmt.Errorf("directive の取得: %w", err)
	}
	return rows, nil
}

// GetDirectiveInstruction は指定 key の指示文（instruction）を返す。本文フェーズが口調 directive を引くのに使う。
func (s *Store) GetDirectiveInstruction(ctx context.Context, key string) (string, error) {
	var instruction string
	if err := s.db.GetContext(ctx, &instruction,
		`SELECT instruction FROM directive WHERE key = ?`, key); err != nil {
		return "", fmt.Errorf("directive %q の取得: %w", key, err)
	}
	return instruction, nil
}

// SaveDirectiveInstruction は指定 key の指示文を更新する。編集画面のレコード別タブの保存で使う。
// key は seed 済みの固定集合のため新規行は作らず、既存行の instruction だけ書き換える。
func (s *Store) SaveDirectiveInstruction(ctx context.Context, key, instruction string) error {
	res, err := s.db.ExecContext(ctx,
		`UPDATE directive SET instruction = ? WHERE key = ?`, instruction, key)
	if err != nil {
		return fmt.Errorf("directive %q の保存: %w", key, err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("directive %q が存在しない", key)
	}
	return nil
}

// ListRecordTypeMaster は REC:FIELD → box・directive の割り当てを全件返す。
// 取込段の振り分けと、編集画面の対象一覧（directive ごとの対象 REC:FIELD）に使う。
func (s *Store) ListRecordTypeMaster(ctx context.Context) ([]model.RecordType, error) {
	var rows []model.RecordType
	if err := s.db.SelectContext(ctx, &rows,
		`SELECT `+recordTypeColumns+` FROM record_type_master ORDER BY box, rec, field`); err != nil {
		return nil, fmt.Errorf("record_type_master の取得: %w", err)
	}
	return rows, nil
}
