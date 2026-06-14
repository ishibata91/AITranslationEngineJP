package engine

import (
	"context"
	"errors"
	"strings"
	"testing"

	"aitranslationenginejp/internal/model"
	"aitranslationenginejp/internal/provider"
)

type fakeStore struct {
	untranslated []model.Narration
	lines        []model.Line
	speakers     map[int64]model.SpeakerIdentity // lineID → 識別子（無ければ話者なし）
	updates      []update
	lineUpdates  []update
}

type update struct {
	id     int64
	dest   string
	status int
}

func (f *fakeStore) ListUntranslatedNarrations(_ context.Context) ([]model.Narration, error) {
	return f.untranslated, nil
}

func (f *fakeStore) UpdateNarrationDest(_ context.Context, id int64, dest string, status int) error {
	f.updates = append(f.updates, update{id: id, dest: dest, status: status})
	return nil
}

func (f *fakeStore) ListUntranslatedLines(_ context.Context) ([]model.Line, error) {
	return f.lines, nil
}

func (f *fakeStore) LoadLineSpeaker(_ context.Context, lineID int64) (model.SpeakerIdentity, bool, error) {
	id, ok := f.speakers[lineID]
	return id, ok, nil
}

func (f *fakeStore) UpdateLineDest(_ context.Context, id int64, dest string, status int) error {
	f.lineUpdates = append(f.lineUpdates, update{id: id, dest: dest, status: status})
	return nil
}

type fakeTranslator struct {
	out           map[string]string
	gotModel      string
	gotDirectives map[string]string // source → 注入された directive
	err           error
}

func (f *fakeTranslator) Translate(_ context.Context, _ provider.Connection, model, source, directive string) (string, error) {
	f.gotModel = model
	if f.gotDirectives == nil {
		f.gotDirectives = map[string]string{}
	}
	f.gotDirectives[source] = directive
	if f.err != nil {
		return "", f.err
	}
	return f.out[source], nil
}

func (f *fakeTranslator) ListModels(_ context.Context, _ provider.Connection) ([]string, error) {
	return nil, nil
}

// 未訳の叙述文を順に翻訳し、訳文を status=3（仮訳）で書き戻すこと。叙述文はペルソナ指示を付けない。
func TestRunTranslatesUntranslatedAsProvisional(t *testing.T) {
	store := &fakeStore{untranslated: []model.Narration{
		{ID: 1, Source: "halls"},
		{ID: 2, Source: "cairn"},
	}}
	tr := &fakeTranslator{out: map[string]string{"halls": "広間", "cairn": "ケルン"}}
	eng := New(store, tr)

	count, err := eng.Run(context.Background(), provider.Connection{Endpoint: "http://x"}, "model-x", nil)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if count != 2 {
		t.Errorf("count = %d, want 2", count)
	}
	if tr.gotModel != "model-x" {
		t.Errorf("model passed = %q, want model-x", tr.gotModel)
	}
	if d := tr.gotDirectives["halls"]; d != "" {
		t.Errorf("叙述文に directive が付いた: %q", d)
	}
	want := []update{{1, "広間", 3}, {2, "ケルン", 3}}
	if len(store.updates) != len(want) {
		t.Fatalf("updates = %v, want %v", store.updates, want)
	}
	for i := range want {
		if store.updates[i] != want[i] {
			t.Errorf("update[%d] = %v, want %v", i, store.updates[i], want[i])
		}
	}
}

// provider がエラーを返したら、そのエラーを返すこと（壊れた訳を書き戻さない）。
func TestRunReturnsProviderError(t *testing.T) {
	store := &fakeStore{untranslated: []model.Narration{{ID: 1, Source: "x"}}}
	tr := &fakeTranslator{err: errors.New("connection refused")}
	eng := New(store, tr)

	_, err := eng.Run(context.Background(), provider.Connection{}, "m", nil)
	if err == nil {
		t.Fatal("want error, got nil")
	}
	if len(store.updates) != 0 {
		t.Errorf("should not write dest on provider error, got %v", store.updates)
	}
}

// 話者を持つ台詞は、話者属性から組んだ口調指示が翻訳プロンプトへ注入されること。
func TestRunTranslatesLinesWithPersonaDirective(t *testing.T) {
	store := &fakeStore{
		lines: []model.Line{{ID: 10, Source: "mother?"}},
		speakers: map[int64]model.SpeakerIdentity{
			10: {RaceEDID: "NordRace", VoiceEDID: "MaleChild"},
		},
	}
	tr := &fakeTranslator{out: map[string]string{"mother?": "母さん？"}}
	eng := New(store, tr)

	count, err := eng.Run(context.Background(), provider.Connection{}, "m", nil)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if count != 1 {
		t.Errorf("count = %d, want 1", count)
	}
	directive := tr.gotDirectives["mother?"]
	if !strings.Contains(directive, "幼い少年の声") {
		t.Errorf("声質が directive に無い: %q", directive)
	}
	if !strings.Contains(directive, "ノルド") {
		t.Errorf("種族の気質が directive に無い: %q", directive)
	}
	want := []update{{10, "母さん？", 3}}
	if len(store.lineUpdates) != 1 || store.lineUpdates[0] != want[0] {
		t.Errorf("lineUpdates = %v, want %v", store.lineUpdates, want)
	}
}

// 話者を解決できない台詞は、口調指示を注入せず空 directive で訳すこと。
func TestRunTranslatesLineWithoutSpeaker(t *testing.T) {
	store := &fakeStore{lines: []model.Line{{ID: 20, Source: "door"}}}
	tr := &fakeTranslator{out: map[string]string{"door": "扉"}}
	eng := New(store, tr)

	if _, err := eng.Run(context.Background(), provider.Connection{}, "m", nil); err != nil {
		t.Fatalf("Run: %v", err)
	}
	if d := tr.gotDirectives["door"]; d != "" {
		t.Errorf("話者なしの台詞に directive が付いた: %q", d)
	}
}

// 進捗 callback が、叙述文と台詞の合計を total に、処理ごとに done を増やして通知すること。
func TestRunReportsProgress(t *testing.T) {
	store := &fakeStore{
		untranslated: []model.Narration{{ID: 1, Source: "a"}},
		lines:        []model.Line{{ID: 10, Source: "b"}, {ID: 11, Source: "c"}},
	}
	tr := &fakeTranslator{out: map[string]string{"a": "あ", "b": "び", "c": "し"}}
	eng := New(store, tr)

	var seen [][2]int
	_, err := eng.Run(context.Background(), provider.Connection{}, "m", func(done, total int) {
		seen = append(seen, [2]int{done, total})
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	// 初回 0/3 と、3 件処理ごとの 1/3,2/3,3/3。
	want := [][2]int{{0, 3}, {1, 3}, {2, 3}, {3, 3}}
	if len(seen) != len(want) {
		t.Fatalf("progress = %v, want %v", seen, want)
	}
	for i := range want {
		if seen[i] != want[i] {
			t.Errorf("progress[%d] = %v, want %v", i, seen[i], want[i])
		}
	}
}

// LineDirective は解決できた話者の口調指示と短い要約を返し、話者なしは両方空を返すこと。
func TestLineDirective(t *testing.T) {
	store := &fakeStore{
		speakers: map[int64]model.SpeakerIdentity{
			10: {VoiceEDID: "FemaleOldGrumpy"},
		},
	}
	eng := New(store, &fakeTranslator{})

	directive, label, err := eng.LineDirective(context.Background(), 10)
	if err != nil {
		t.Fatalf("LineDirective: %v", err)
	}
	if !strings.Contains(directive, "気難しい老女の声") {
		t.Errorf("directive = %q", directive)
	}
	if label != "声質: 気難しい老女の声" {
		t.Errorf("label = %q", label)
	}

	d2, l2, err := eng.LineDirective(context.Background(), 99) // 話者なし
	if err != nil {
		t.Fatalf("LineDirective(no speaker): %v", err)
	}
	if d2 != "" || l2 != "" {
		t.Errorf("話者なしで directive=%q label=%q、空を期待", d2, l2)
	}
}
