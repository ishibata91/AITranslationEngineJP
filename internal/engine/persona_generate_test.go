package engine

import (
	"context"
	"fmt"
	"testing"

	"aitranslationenginejp/internal/engine/tone"
	"aitranslationenginejp/internal/model"
)

// fakeGenStore は生成器の中心データアクセスの代替。line_analysis のキャッシュを保持し、保存呼び出しを記録する。
type fakeGenStore struct {
	sources    []model.SpeakerLineSource
	cache      map[string]model.LineAnalysis
	upsertedLA []model.LineAnalysis
	upsertedPC []model.PersonaCharacter
}

func (f *fakeGenStore) ListSpeakerLineSources(context.Context) ([]model.SpeakerLineSource, error) {
	return f.sources, nil
}

func (f *fakeGenStore) GetLineAnalyses(_ context.Context, hashes []string) (map[string]model.LineAnalysis, error) {
	out := make(map[string]model.LineAnalysis)
	for _, h := range hashes {
		if a, ok := f.cache[h]; ok {
			out[h] = a
		}
	}
	return out, nil
}

func (f *fakeGenStore) UpsertLineAnalysis(_ context.Context, a model.LineAnalysis) error {
	f.upsertedLA = append(f.upsertedLA, a)
	f.cache[a.SourceHash] = a
	return nil
}

func (f *fakeGenStore) UpsertPersonaCharacter(_ context.Context, pc model.PersonaCharacter) error {
	f.upsertedPC = append(f.upsertedPC, pc)
	return nil
}

// seedRude はキャッシュへ粗野（命令 1）な行解析を入れる。
func seedRude(cache map[string]model.LineAnalysis, src string) {
	cache[SourceHash(src)] = model.LineAnalysis{SourceHash: SourceHash(src), SentenceCount: 1, IsImperative: 1}
}

// seedPolite はキャッシュへ丁寧（丁寧定型 1）な行解析を入れる。
func seedPolite(cache map[string]model.LineAnalysis, src string) {
	cache[SourceHash(src)] = model.LineAnalysis{SourceHash: SourceHash(src), SentenceCount: 1, PoliteCount: 1}
}

// findPC は plugin・form_id で保存済みペルソナを探す。
func findPC(pcs []model.PersonaCharacter, plugin, formID string) (model.PersonaCharacter, bool) {
	for _, pc := range pcs {
		if pc.SpeakerPlugin == plugin && pc.SpeakerFormID == formID {
			return pc, true
		}
	}
	return model.PersonaCharacter{}, false
}

// Generate は話者ごとに台詞を束ねて基底口調を分類し保存する。キャッシュ命中時は prose を回さない。
func TestGenerateClassifiesPerSpeaker(t *testing.T) {
	cache := make(map[string]model.LineAnalysis)
	var sources []model.SpeakerLineSource
	// 話者 A: 粗野な命令 12 行（印 12 ≥ 10）。本文経路で尊大。
	for i := range 12 {
		src := fmt.Sprintf("Get out of here now %d.", i)
		seedRude(cache, src)
		sources = append(sources, model.SpeakerLineSource{SpeakerPlugin: "A.esp", SpeakerFormID: "F1", VoiceEDID: "MaleNord", Source: src})
	}
	// 話者 B: 丁寧な台詞 11 行。本文経路で丁寧。
	for i := range 11 {
		src := fmt.Sprintf("Please forgive me sir %d.", i)
		seedPolite(cache, src)
		sources = append(sources, model.SpeakerLineSource{SpeakerPlugin: "A.esp", SpeakerFormID: "F2", VoiceEDID: "FemaleNord", Source: src})
	}
	store := &fakeGenStore{sources: sources, cache: cache}

	n, err := NewPersonaGenerator(store, fakeLexicon{}).Generate(context.Background())
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if n != 2 {
		t.Fatalf("生成話者数 = %d, want 2", n)
	}
	if len(store.upsertedLA) != 0 {
		t.Errorf("キャッシュ命中なのに行解析を %d 件再計算した（prose を回した）", len(store.upsertedLA))
	}

	a, ok := findPC(store.upsertedPC, "A.esp", "F1")
	if !ok {
		t.Fatal("話者 A のペルソナが保存されていない")
	}
	if a.AttitudeBand != tone.AttitudeArrogant || a.DecisionPath != tone.PathDialogue || a.Marked != 12 {
		t.Errorf("話者 A = %+v, want 尊大・本文・印12", a)
	}
	b, ok := findPC(store.upsertedPC, "A.esp", "F2")
	if !ok {
		t.Fatal("話者 B のペルソナが保存されていない")
	}
	if b.AttitudeBand != tone.AttitudePolite || b.DecisionPath != tone.PathDialogue {
		t.Errorf("話者 B = %+v, want 丁寧・本文", b)
	}
}

// Generate はキャッシュ未命中の本文だけ prose で解析し、line_analysis へ保存する。
func TestGenerateComputesMissingAndCaches(t *testing.T) {
	src := "I am calm and steady here."
	store := &fakeGenStore{
		sources: []model.SpeakerLineSource{
			{SpeakerPlugin: "A.esp", SpeakerFormID: "F1", VoiceEDID: "MaleNord", Source: src},
		},
		cache: make(map[string]model.LineAnalysis),
	}

	n, err := NewPersonaGenerator(store, fakeLexicon{}).Generate(context.Background())
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if n != 1 {
		t.Fatalf("生成話者数 = %d, want 1", n)
	}
	if len(store.upsertedLA) != 1 || store.upsertedLA[0].SourceHash != SourceHash(src) {
		t.Fatalf("未命中本文の行解析キャッシュ保存 = %+v, want 1 件（本文ハッシュ一致）", store.upsertedLA)
	}
	// 印不足（マーカー無し）で voice 中立のため、voice 経路の中立になる。
	pc, ok := findPC(store.upsertedPC, "A.esp", "F1")
	if !ok {
		t.Fatal("ペルソナが保存されていない")
	}
	if pc.DecisionPath != tone.PathVoice || pc.AttitudeBand != tone.AttitudeNeutral {
		t.Errorf("ペルソナ = %+v, want voice・中立", pc)
	}
}
