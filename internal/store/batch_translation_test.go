package store

import (
	"context"
	"path/filepath"
	"strconv"
	"testing"

	"aitranslationenginejp/internal/model"
)

// itoa は int64 の id を count クエリへ差し込むための短い整数→文字列変換。
func itoa(id int64) string { return strconv.FormatInt(id, 10) }

// R-1-4: migration 前からある batch 進行と同じ形で provider を省略した行は xAI として読めること。
func Test既存Batch進行のProviderはXAIになる(t *testing.T) {
	s, err := Open(filepath.Join(t.TempDir(), "central.sqlite3"))
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer func() { _ = s.Close() }()
	ctx := context.Background()
	if _, err := s.db.ExecContext(ctx, `INSERT INTO batch_translation (plugin, model, stage) VALUES ('legacy.esp', 'grok-4', 'proper_noun')`); err != nil {
		t.Fatalf("旧形式相当の行を追加: %v", err)
	}
	prog, ok, err := s.GetBatchProgression(ctx, "legacy.esp")
	if err != nil || !ok || prog.Provider != "xai" {
		t.Errorf("既存進行 = %+v, ok=%v err=%v", prog, ok, err)
	}
}

// StartBatchProgression が対象 plugin の進行を開始し、固有名段で id を返すこと。
// 同一 plugin の再開始（reset）で行を使い回し、stage を固有名へ戻し、外部 ID を空へ戻し、
// 前回の送信行対応（batch_request）を消すこと。進行は plugin と 1 対 1。
func TestStartBatchProgressionResets(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "central.sqlite3")
	s, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer func() { _ = s.Close() }()
	ctx := context.Background()

	id, err := s.StartBatchProgression(ctx, "A.esp", "xai", "grok-4")
	if err != nil {
		t.Fatalf("1 回目の開始: %v", err)
	}
	// 進行を本文段へ進め、外部 ID を記録し、送信行対応を入れる。
	if err = s.RecordBatchExternalID(ctx, id, model.BatchStageProperNoun, "ext-p"); err != nil {
		t.Fatalf("固有名 外部 ID: %v", err)
	}
	if err = s.AdvanceBatchStage(ctx, id, model.BatchStageBody); err != nil {
		t.Fatalf("本文段へ: %v", err)
	}
	if _, err = s.InsertBatchRequests(ctx, []model.BatchRequest{
		{BatchID: id, ExternalBatchID: "ext-p", CustomID: "p:1", Kind: model.BatchKindProper, RowID: 1},
	}); err != nil {
		t.Fatalf("送信行対応の挿入: %v", err)
	}

	// 再開始（reset）。同じ plugin なので id は同一で、stage・外部 ID は初期化、batch_request は消える。
	id2, err := s.StartBatchProgression(ctx, "A.esp", "openai", "grok-4-fast")
	if err != nil {
		t.Fatalf("再開始: %v", err)
	}
	if id2 != id {
		t.Errorf("reset で id が変わった: %d -> %d", id, id2)
	}
	prog, ok, err := s.GetBatchProgression(ctx, "A.esp")
	if err != nil || !ok {
		t.Fatalf("GetBatchProgression: ok=%v err=%v", ok, err)
	}
	if prog.Stage != model.BatchStageProperNoun {
		t.Errorf("reset 後 stage = %q, want proper_noun", prog.Stage)
	}
	if prog.ProperBatchID != "" || prog.BodyBatchID != "" {
		t.Errorf("reset 後 外部 ID が残った: proper=%q body=%q", prog.ProperBatchID, prog.BodyBatchID)
	}
	if prog.Model != "grok-4-fast" {
		t.Errorf("reset 後 model = %q, want grok-4-fast", prog.Model)
	}
	if prog.Provider != "openai" {
		t.Errorf("reset 後 provider = %q, want openai", prog.Provider)
	}
	if n := countRows(t, dbPath, `SELECT COUNT(*) FROM batch_request WHERE batch_id = `+itoa(id)); n != 0 {
		t.Errorf("reset 後 batch_request が残った（%d 件）", n)
	}
	if n := countRows(t, dbPath, `SELECT COUNT(*) FROM batch_translation WHERE plugin = 'A.esp'`); n != 1 {
		t.Errorf("batch_translation 行数 = %d, want 1（plugin と 1 対 1）", n)
	}
}

// RecordBatchExternalID が段に応じた列（proper_batch_id / body_batch_id）へ外部 ID を書くこと。
func TestRecordBatchExternalIDByStage(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "central.sqlite3")
	s, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer func() { _ = s.Close() }()
	ctx := context.Background()

	id, _ := s.StartBatchProgression(ctx, "A.esp", "xai", "grok-4")
	if err = s.RecordBatchExternalID(ctx, id, model.BatchStageProperNoun, "ext-p"); err != nil {
		t.Fatalf("固有名 外部 ID: %v", err)
	}
	if err = s.RecordBatchExternalID(ctx, id, model.BatchStageBody, "ext-b"); err != nil {
		t.Fatalf("本文 外部 ID: %v", err)
	}
	prog, _, _ := s.GetBatchProgression(ctx, "A.esp")
	if prog.ProperBatchID != "ext-p" || prog.BodyBatchID != "ext-b" {
		t.Errorf("外部 ID = proper %q / body %q, want ext-p / ext-b", prog.ProperBatchID, prog.BodyBatchID)
	}
	// 完了段は外部 ID 列を持たないため拒否する。
	if err = s.RecordBatchExternalID(ctx, id, model.BatchStageDone, "x"); err == nil {
		t.Errorf("完了段への外部 ID 記録が拒否されなかった")
	}
}

// ListActiveBatchProgressions が反映待ち（完了でない）進行だけを返すこと。
func TestListActiveBatchProgressionsExcludesDone(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "central.sqlite3")
	s, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer func() { _ = s.Close() }()
	ctx := context.Background()

	idA, _ := s.StartBatchProgression(ctx, "A.esp", "xai", "grok-4")
	idB, _ := s.StartBatchProgression(ctx, "B.esp", "openai", "gpt-5.6-luna")
	// A は完了、B は本文段（反映待ち）。
	if err = s.AdvanceBatchStage(ctx, idA, model.BatchStageDone); err != nil {
		t.Fatalf("A 完了: %v", err)
	}
	if err = s.AdvanceBatchStage(ctx, idB, model.BatchStageBody); err != nil {
		t.Fatalf("B 本文段: %v", err)
	}
	active, err := s.ListActiveBatchProgressions(ctx)
	if err != nil {
		t.Fatalf("ListActiveBatchProgressions: %v", err)
	}
	if len(active) != 1 || active[0].Plugin != "B.esp" {
		t.Errorf("反映待ち = %+v, want B.esp のみ", active)
	}
}

// InsertBatchRequests / ListBatchRequests が外部 batch ID で送信行対応を出し入れできること。
func TestInsertAndListBatchRequests(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "central.sqlite3")
	s, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer func() { _ = s.Close() }()
	ctx := context.Background()

	id, _ := s.StartBatchProgression(ctx, "A.esp", "xai", "grok-4")
	rows := []model.BatchRequest{
		{BatchID: id, ExternalBatchID: "ext-b", CustomID: "n:10", Kind: model.BatchKindNarration, RowID: 10},
		{BatchID: id, ExternalBatchID: "ext-b", CustomID: "l:20", Kind: model.BatchKindLine, RowID: 20},
		{BatchID: id, ExternalBatchID: "ext-p", CustomID: "p:30", Kind: model.BatchKindProper, RowID: 30},
	}
	n, err := s.InsertBatchRequests(ctx, rows)
	if err != nil {
		t.Fatalf("InsertBatchRequests: %v", err)
	}
	if n != 3 {
		t.Errorf("挿入件数 = %d, want 3", n)
	}
	body, err := s.ListBatchRequests(ctx, "ext-b")
	if err != nil {
		t.Fatalf("ListBatchRequests: %v", err)
	}
	if len(body) != 2 {
		t.Fatalf("ext-b の送信行対応 = %d 件, want 2", len(body))
	}
	if body[0].CustomID != "n:10" || body[0].Kind != model.BatchKindNarration || body[0].RowID != 10 {
		t.Errorf("1 件目 = %+v", body[0])
	}
	if body[1].CustomID != "l:20" || body[1].RowID != 20 {
		t.Errorf("2 件目 = %+v", body[1])
	}
}

// DeleteTargetPlugin が対象 plugin の batch 進行と送信行対応を消し、別 plugin の分を残すこと。
// 手続き削除リストへ batch テーブルを子→親順で追記した効果を保証する。
func TestDeleteTargetPluginRemovesBatch(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "central.sqlite3")
	s, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer func() { _ = s.Close() }()
	ctx := context.Background()

	if err = s.UpsertTargetPlugin(ctx, "A.esp", "/data/A.esp"); err != nil {
		t.Fatalf("A 登録: %v", err)
	}
	if err = s.UpsertTargetPlugin(ctx, "B.esp", "/data/B.esp"); err != nil {
		t.Fatalf("B 登録: %v", err)
	}
	idA, _ := s.StartBatchProgression(ctx, "A.esp", "xai", "grok-4")
	idB, _ := s.StartBatchProgression(ctx, "B.esp", "openai", "gpt-5.6-luna")
	if _, err = s.InsertBatchRequests(ctx, []model.BatchRequest{
		{BatchID: idA, ExternalBatchID: "ext-a", CustomID: "n:1", Kind: model.BatchKindNarration, RowID: 1},
		{BatchID: idB, ExternalBatchID: "ext-b", CustomID: "n:2", Kind: model.BatchKindNarration, RowID: 2},
	}); err != nil {
		t.Fatalf("送信行対応の挿入: %v", err)
	}

	if err = s.DeleteTargetPlugin(ctx, "A.esp"); err != nil {
		t.Fatalf("DeleteTargetPlugin: %v", err)
	}

	if n := countRows(t, dbPath, `SELECT COUNT(*) FROM batch_translation WHERE plugin = 'A.esp'`); n != 0 {
		t.Errorf("A の batch_translation が残った（%d 件）", n)
	}
	if n := countRows(t, dbPath, `SELECT COUNT(*) FROM batch_request WHERE batch_id = `+itoa(idA)); n != 0 {
		t.Errorf("A の batch_request が残った（%d 件）", n)
	}
	if n := countRows(t, dbPath, `SELECT COUNT(*) FROM batch_translation WHERE plugin = 'B.esp'`); n != 1 {
		t.Errorf("B の batch_translation が消えた（%d 件、want 1）", n)
	}
	if n := countRows(t, dbPath, `SELECT COUNT(*) FROM batch_request WHERE batch_id = `+itoa(idB)); n != 1 {
		t.Errorf("B の batch_request が消えた（%d 件、want 1）", n)
	}
}
