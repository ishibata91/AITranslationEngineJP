package store

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"aitranslationenginejp/internal/model"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

const termDictionaryPageSize = 50

var ErrTermDictionaryRevisionConflict = errors.New("用語は取得後に更新されている")

type TermDictionaryEditor struct{ db *sqlx.DB }

type TermDictionaryFilter = model.TermDictionaryFilter
type TermDictionaryEntry = model.TermDictionaryEntry
type TermDictionaryPage = model.TermDictionaryPage
type TermDictionaryCreate = model.TermDictionaryCreate
type TermDictionaryPatch = model.TermDictionaryPatch

func OpenTermDictionaryEditor(path string) (*TermDictionaryEditor, error) {
	db, err := sqlx.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("用語辞書編集DBを開けない (%s): %w", path, err)
	}
	db.SetMaxOpenConns(1)
	if _, err := db.ExecContext(context.Background(), `PRAGMA busy_timeout = 5000; PRAGMA foreign_keys = ON;`); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("用語辞書編集DBの接続設定: %w", err)
	}
	var count int
	if err := db.Get(&count, `SELECT COUNT(*) FROM dictionary_sense`); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("用語辞書編集DBのschema検証: %w", err)
	}
	return &TermDictionaryEditor{db: db}, nil
}
func (s *TermDictionaryEditor) Close() error {
	if err := s.db.Close(); err != nil {
		return fmt.Errorf("用語辞書編集DBのクローズ: %w", err)
	}
	return nil
}

func (s *TermDictionaryEditor) List(ctx context.Context, filter TermDictionaryFilter, page int) (TermDictionaryPage, error) {
	if page < 1 {
		page = 1
	}
	filter = trimFilter(filter)
	where, args := termDictionaryWhere(filter)
	var total int
	if err := s.db.GetContext(ctx, &total, `SELECT COUNT(*) FROM dictionary_sense s JOIN dictionary_term t ON t.id=s.term_id WHERE `+where, args...); err != nil {
		return TermDictionaryPage{}, fmt.Errorf("用語辞書の件数取得: %w", err)
	}
	type row struct {
		ID           int64  `db:"id"`
		Source       string `db:"source"`
		Destination  string `db:"dest"`
		PartOfSpeech string `db:"part_of_speech"`
		Revision     int64  `db:"revision"`
		Categories   string `db:"categories"`
	}
	args = append(args, termDictionaryPageSize, (page-1)*termDictionaryPageSize)
	var rows []row
	query := `SELECT s.id,t.source,s.dest,s.part_of_speech,s.revision,COALESCE(GROUP_CONCAT(DISTINCT o.skyrim_category), '') AS categories FROM dictionary_sense s JOIN dictionary_term t ON t.id=s.term_id LEFT JOIN dictionary_occurrence o ON o.sense_id=s.id WHERE ` + where + ` GROUP BY s.id ORDER BY t.source COLLATE NOCASE,s.dest COLLATE NOCASE,s.id LIMIT ? OFFSET ?`
	if err := s.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return TermDictionaryPage{}, fmt.Errorf("用語辞書の一覧取得: %w", err)
	}
	entries := make([]TermDictionaryEntry, len(rows))
	for i, r := range rows {
		entries[i] = TermDictionaryEntry{ID: r.ID, Source: r.Source, Destination: r.Destination, PartOfSpeech: r.PartOfSpeech, Revision: r.Revision, Categories: splitCategories(r.Categories)}
	}
	return TermDictionaryPage{Entries: entries, TotalCount: total, PageNumber: page}, nil
}
func (s *TermDictionaryEditor) Create(ctx context.Context, in TermDictionaryCreate) (TermDictionaryEntry, error) {
	in.Source, in.Destination, in.PartOfSpeech = strings.TrimSpace(in.Source), strings.TrimSpace(in.Destination), strings.TrimSpace(in.PartOfSpeech)
	if in.Source == "" || in.Destination == "" || in.PartOfSpeech == "" {
		return TermDictionaryEntry{}, errors.New("原語、訳語、品詞は空にできない")
	}
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return TermDictionaryEntry{}, fmt.Errorf("用語作成transaction開始: %w", err)
	}
	defer tx.Rollback()
	termID, err := ensureEditorTerm(ctx, tx, in.Source)
	if err != nil {
		return TermDictionaryEntry{}, err
	}
	result, err := tx.ExecContext(ctx, `INSERT INTO dictionary_sense (term_id,dest,part_of_speech,classification_status) VALUES (?,?,?,'classified')`, termID, in.Destination, in.PartOfSpeech)
	if err != nil {
		return TermDictionaryEntry{}, fmt.Errorf("用語作成: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return TermDictionaryEntry{}, fmt.Errorf("作成した用語の識別子取得: %w", err)
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO dictionary_occurrence (term_id,sense_id,observed_dest,origin_kind,origin_reference) VALUES (?,?,?,'manual',?)`, termID, id, in.Destination, fmt.Sprintf("sense:%d", id)); err != nil {
		return TermDictionaryEntry{}, fmt.Errorf("作成用語の用例保存: %w", err)
	}
	if err = tx.Commit(); err != nil {
		return TermDictionaryEntry{}, fmt.Errorf("用語作成commit: %w", err)
	}
	return s.entry(ctx, id)
}
func (s *TermDictionaryEditor) Patch(ctx context.Context, in TermDictionaryPatch) (TermDictionaryEntry, error) {
	if in.Source == nil && in.Destination == nil && in.PartOfSpeech == nil {
		return TermDictionaryEntry{}, errors.New("変更項目がない")
	}
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return TermDictionaryEntry{}, fmt.Errorf("用語更新transaction開始: %w", err)
	}
	defer tx.Rollback()
	var current struct {
		TermID       int64  `db:"term_id"`
		Source       string `db:"source"`
		Destination  string `db:"dest"`
		PartOfSpeech string `db:"part_of_speech"`
		Revision     int64  `db:"revision"`
	}
	if err = tx.GetContext(ctx, &current, `SELECT s.term_id,t.source,s.dest,s.part_of_speech,s.revision FROM dictionary_sense s JOIN dictionary_term t ON t.id=s.term_id WHERE s.id=?`, in.ID); err != nil {
		return TermDictionaryEntry{}, fmt.Errorf("更新対象用語の取得: %w", err)
	}
	if current.Revision != in.Revision {
		return TermDictionaryEntry{}, ErrTermDictionaryRevisionConflict
	}
	changes := make([]struct{ field, old, new string }, 0, 3)
	if in.Source != nil {
		v := strings.TrimSpace(*in.Source)
		if v == "" {
			return TermDictionaryEntry{}, errors.New("原語は空にできない")
		}
		if v != current.Source {
			newTerm, err := ensureEditorTerm(ctx, tx, v)
			if err != nil {
				return TermDictionaryEntry{}, err
			}
			if _, err = tx.ExecContext(ctx, `UPDATE dictionary_sense SET term_id=? WHERE id=?`, newTerm, in.ID); err != nil {
				return TermDictionaryEntry{}, fmt.Errorf("用語の原語移動: %w", err)
			}
			if _, err = tx.ExecContext(ctx, `UPDATE dictionary_occurrence SET term_id=? WHERE sense_id=?`, newTerm, in.ID); err != nil {
				return TermDictionaryEntry{}, fmt.Errorf("用例の原語移動: %w", err)
			}
			if _, err = tx.ExecContext(ctx, `DELETE FROM dictionary_term WHERE id=? AND NOT EXISTS (SELECT 1 FROM dictionary_sense WHERE term_id=?)`, current.TermID, current.TermID); err != nil {
				return TermDictionaryEntry{}, fmt.Errorf("空の原語削除: %w", err)
			}
			changes = append(changes, struct{ field, old, new string }{"source", current.Source, v})
		}
	}
	sets := []string{}
	values := []any{}
	if in.Destination != nil {
		v := strings.TrimSpace(*in.Destination)
		if v == "" {
			return TermDictionaryEntry{}, errors.New("訳語は空にできない")
		}
		if v != current.Destination {
			sets = append(sets, "dest = ?")
			values = append(values, v)
			changes = append(changes, struct{ field, old, new string }{"dest", current.Destination, v})
		}
	}
	if in.PartOfSpeech != nil {
		v := strings.TrimSpace(*in.PartOfSpeech)
		if v == "" {
			return TermDictionaryEntry{}, errors.New("品詞は空にできない")
		}
		if v != current.PartOfSpeech {
			sets = append(sets, "part_of_speech = ?")
			values = append(values, v)
			changes = append(changes, struct{ field, old, new string }{"part_of_speech", current.PartOfSpeech, v})
		}
	}
	if len(changes) == 0 {
		if err = tx.Commit(); err != nil {
			return TermDictionaryEntry{}, fmt.Errorf("用語更新commit: %w", err)
		}
		return s.entry(ctx, in.ID)
	}
	sets = append(sets, "revision = revision + 1", "updated_at = ?")
	values = append(values, time.Now().UTC().Format(time.RFC3339Nano), in.ID, in.Revision)
	result, err := tx.ExecContext(ctx, `UPDATE dictionary_sense SET `+strings.Join(sets, ",")+` WHERE id=? AND revision=?`, values...)
	if err != nil {
		return TermDictionaryEntry{}, fmt.Errorf("用語更新: %w", err)
	}
	n, err := result.RowsAffected()
	if err != nil {
		return TermDictionaryEntry{}, fmt.Errorf("用語更新件数取得: %w", err)
	}
	if n == 0 {
		return TermDictionaryEntry{}, ErrTermDictionaryRevisionConflict
	}
	for _, c := range changes {
		if _, err = tx.ExecContext(ctx, `INSERT INTO dictionary_change (target_table,target_id,field_name,old_value,new_value,changed_by,reason) VALUES ('dictionary_sense',?,?,?,?,?,?)`, in.ID, c.field, c.old, c.new, "term_dictionary_ui", "行内編集"); err != nil {
			return TermDictionaryEntry{}, fmt.Errorf("用語変更履歴保存: %w", err)
		}
	}
	if err = tx.Commit(); err != nil {
		return TermDictionaryEntry{}, fmt.Errorf("用語更新commit: %w", err)
	}
	return s.entry(ctx, in.ID)
}
func (s *TermDictionaryEditor) Delete(ctx context.Context, id, revision int64) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("用語削除transaction開始: %w", err)
	}
	defer tx.Rollback()
	var termID int64
	if err = tx.GetContext(ctx, &termID, `SELECT term_id FROM dictionary_sense WHERE id=? AND revision=?`, id, revision); err != nil {
		return fmt.Errorf("削除対象用語の取得: %w", err)
	}
	r, err := tx.ExecContext(ctx, `DELETE FROM dictionary_sense WHERE id=? AND revision=?`, id, revision)
	if err != nil {
		return fmt.Errorf("用語削除: %w", err)
	}
	n, err := r.RowsAffected()
	if err != nil {
		return fmt.Errorf("用語削除件数取得: %w", err)
	}
	if n == 0 {
		return ErrTermDictionaryRevisionConflict
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM dictionary_term WHERE id=? AND NOT EXISTS (SELECT 1 FROM dictionary_sense WHERE term_id=?)`, termID, termID); err != nil {
		return fmt.Errorf("空の原語削除: %w", err)
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO dictionary_change (target_table,target_id,field_name,changed_by,reason) VALUES ('dictionary_sense',?,'deleted','term_dictionary_ui','行削除')`, id); err != nil {
		return fmt.Errorf("用語削除履歴保存: %w", err)
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("用語削除commit: %w", err)
	}
	return nil
}
func (s *TermDictionaryEditor) entry(ctx context.Context, id int64) (TermDictionaryEntry, error) {
	page, err := s.List(ctx, TermDictionaryFilter{}, 1)
	if err != nil {
		return TermDictionaryEntry{}, err
	}
	for _, e := range page.Entries {
		if e.ID == id {
			return e, nil
		}
	}
	var e TermDictionaryEntry
	if err = s.db.GetContext(ctx, &e, `SELECT s.id,t.source,s.dest,s.part_of_speech,s.revision FROM dictionary_sense s JOIN dictionary_term t ON t.id=s.term_id WHERE s.id=?`, id); err != nil {
		return TermDictionaryEntry{}, fmt.Errorf("用語取得: %w", err)
	}
	return e, nil
}
func ensureEditorTerm(ctx context.Context, tx *sqlx.Tx, source string) (int64, error) {
	if _, err := tx.ExecContext(ctx, `INSERT INTO dictionary_term (source) VALUES (?) ON CONFLICT(source) DO NOTHING`, source); err != nil {
		return 0, fmt.Errorf("原語保存: %w", err)
	}
	var id int64
	if err := tx.GetContext(ctx, &id, `SELECT id FROM dictionary_term WHERE source=?`, source); err != nil {
		return 0, fmt.Errorf("原語取得: %w", err)
	}
	return id, nil
}
func trimFilter(f TermDictionaryFilter) TermDictionaryFilter {
	f.Source = strings.TrimSpace(f.Source)
	f.Destination = strings.TrimSpace(f.Destination)
	f.PartOfSpeech = strings.TrimSpace(f.PartOfSpeech)
	f.Category = strings.TrimSpace(f.Category)
	return f
}
func termDictionaryWhere(f TermDictionaryFilter) (string, []any) {
	clauses := []string{"1=1"}
	args := []any{}
	for _, v := range []struct{ column, value string }{{"t.source", f.Source}, {"s.dest", f.Destination}, {"s.part_of_speech", f.PartOfSpeech}} {
		if v.value != "" {
			clauses = append(clauses, v.column+" LIKE ? ESCAPE '\\' COLLATE NOCASE")
			args = append(args, "%"+escapeTermDictionaryLike(v.value)+"%")
		}
	}
	if f.Category != "" {
		clauses = append(clauses, `EXISTS (SELECT 1 FROM dictionary_occurrence o2 WHERE o2.sense_id=s.id AND o2.skyrim_category LIKE ? ESCAPE '\' COLLATE NOCASE)`)
		args = append(args, "%"+escapeTermDictionaryLike(f.Category)+"%")
	}
	return strings.Join(clauses, " AND "), args
}
func escapeTermDictionaryLike(v string) string {
	return strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(v)
}
func splitCategories(v string) []string {
	if v == "" {
		return []string{}
	}
	return strings.Split(v, ",")
}
