package engine

import (
	"context"
	"errors"
	"testing"

	"aitranslationenginejp/internal/model"
	"aitranslationenginejp/internal/provider"
)

// 同一 INFO（同じ plugin・form_id）の複数応答行は、トークン予算内なら 1 リクエスト（キー配列経路）にまとまること。
// 返ったキーの台詞が仮訳（status=3）で書き戻ること。
func TestRunBatchesSameInfoLines(t *testing.T) {
	store := &fakeStore{lines: []model.Line{
		{ID: 10, Source: "alpha", Plugin: "P", FormID: "F"},
		{ID: 11, Source: "beta", Plugin: "P", FormID: "F"},
	}}
	tr := &fakeTranslator{out: map[string]string{"alpha": "ア", "beta": "ベ"}}
	eng := New(store, tr, fakeLexicon{}, nil, nil)

	count, err := eng.Run(context.Background(), provider.Connection{}, "m", "", 1000, nil)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if tr.batchCalls != 1 {
		t.Fatalf("batchCalls = %d, want 1（同一 INFO はまとめる）", tr.batchCalls)
	}
	if len(tr.gotBatchItems[0]) != 2 {
		t.Errorf("chunk items = %d, want 2", len(tr.gotBatchItems[0]))
	}
	if count != 2 {
		t.Errorf("count = %d, want 2", count)
	}
	want := []update{{10, "ア", 3}, {11, "ベ", 3}}
	if len(store.lineUpdates) != 2 || store.lineUpdates[0] != want[0] || store.lineUpdates[1] != want[1] {
		t.Errorf("lineUpdates = %v, want %v", store.lineUpdates, want)
	}
	// system メッセージへ base 指示が載ること。
	if len(tr.gotBatchSystem) != 1 || tr.gotBatchSystem[0] != testBaseDirective {
		t.Errorf("batch system = %v, want [%q]", tr.gotBatchSystem, testBaseDirective)
	}
}

// バルク応答で一部キーが欠けた（壊れた）場合、そのキーの台詞は未訳のまま残り（書き戻さない）、
// 返ったキーの台詞だけ確定すること。バッチ再送はしない。
func TestRunBatchPartialSuccessLeavesMissingUntranslated(t *testing.T) {
	store := &fakeStore{lines: []model.Line{
		{ID: 10, Source: "alpha", Plugin: "P", FormID: "F"},
		{ID: 11, Source: "beta", Plugin: "P", FormID: "F"},
	}}
	tr := &fakeTranslator{
		out:       map[string]string{"alpha": "ア", "beta": "ベ"},
		errByUser: map[string]error{"beta": errors.New("壊れたキー")}, // beta はキーを落とす
	}
	eng := New(store, tr, fakeLexicon{}, nil, nil)

	count, err := eng.Run(context.Background(), provider.Connection{}, "m", "", 1000, nil)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if tr.batchCalls != 1 {
		t.Fatalf("batchCalls = %d, want 1", tr.batchCalls)
	}
	if count != 1 {
		t.Errorf("count = %d, want 1（欠けたキーは未訳のまま）", count)
	}
	want := []update{{10, "ア", 3}}
	if len(store.lineUpdates) != 1 || store.lineUpdates[0] != want[0] {
		t.Errorf("lineUpdates = %v, want %v（beta は書き戻さない）", store.lineUpdates, want)
	}
}

// 別 INFO（form_id が違う）の台詞は、予算に余裕があってもまとめず、1 行ずつ単一経路で翻訳すること。
func TestRunDoesNotBatchAcrossInfo(t *testing.T) {
	store := &fakeStore{lines: []model.Line{
		{ID: 10, Source: "a", Plugin: "P", FormID: "F1"},
		{ID: 11, Source: "b", Plugin: "P", FormID: "F2"},
	}}
	tr := &fakeTranslator{out: map[string]string{"a": "ア", "b": "ビ"}}
	eng := New(store, tr, fakeLexicon{}, nil, nil)

	count, err := eng.Run(context.Background(), provider.Connection{}, "m", "", 1000, nil)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if tr.batchCalls != 0 {
		t.Errorf("batchCalls = %d, want 0（別 INFO はまとめない）", tr.batchCalls)
	}
	if count != 2 {
		t.Errorf("count = %d, want 2", count)
	}
	// 単一経路を通るので gotPrompts に各行が載る。
	if _, ok := tr.gotPrompts["a"]; !ok {
		t.Errorf("単一経路で a が送られていない")
	}
	if _, ok := tr.gotPrompts["b"]; !ok {
		t.Errorf("単一経路で b が送られていない")
	}
}

// tokenBudget が 0 以下なら、同一 INFO でもバルクせず 1 行ずつ翻訳すること（従来動作）。
func TestRunNoBatchWhenBudgetZero(t *testing.T) {
	store := &fakeStore{lines: []model.Line{
		{ID: 10, Source: "alpha", Plugin: "P", FormID: "F"},
		{ID: 11, Source: "beta", Plugin: "P", FormID: "F"},
	}}
	tr := &fakeTranslator{out: map[string]string{"alpha": "ア", "beta": "ベ"}}
	eng := New(store, tr, fakeLexicon{}, nil, nil)

	count, err := eng.Run(context.Background(), provider.Connection{}, "m", "", 0, nil)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if tr.batchCalls != 0 {
		t.Errorf("batchCalls = %d, want 0（予算 0 はバルクしない）", tr.batchCalls)
	}
	if count != 2 {
		t.Errorf("count = %d, want 2", count)
	}
}

// Run に対象 plugin を渡すと、その plugin の未訳台詞だけを翻訳し、別 plugin の未訳は触らないこと。
// 抽出した 1 plugin の実行が他 plugin の未訳を巻き込まない絞り込みを守る（実画面で発覚した既存挙動の修正）。
func TestRunScopesToTargetPlugin(t *testing.T) {
	store := &fakeStore{lines: []model.Line{
		{ID: 10, Source: "alpha", Plugin: "A.esp", FormID: "F1"},
		{ID: 20, Source: "bravo", Plugin: "B.esp", FormID: "F2"},
	}}
	tr := &fakeTranslator{out: map[string]string{"alpha": "ア", "bravo": "ブ"}}
	eng := New(store, tr, fakeLexicon{}, nil, nil)

	count, err := eng.Run(context.Background(), provider.Connection{}, "m", "A.esp", 0, nil)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if count != 1 {
		t.Errorf("count = %d, want 1（A.esp の 1 行だけ）", count)
	}
	want := []update{{10, "ア", 3}}
	if len(store.lineUpdates) != 1 || store.lineUpdates[0] != want[0] {
		t.Errorf("lineUpdates = %v, want %v（B.esp は触らない）", store.lineUpdates, want)
	}
	if _, sent := tr.gotPrompts["bravo"]; sent {
		t.Errorf("B.esp の bravo を翻訳へ送った。対象 plugin 外は送らない")
	}
}
