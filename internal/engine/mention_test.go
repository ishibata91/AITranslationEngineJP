package engine

import (
	"context"
	"strings"
	"testing"

	"aitranslationenginejp/internal/core/dictionary"
	"aitranslationenginejp/internal/model"
)

// Ingest の言及段が、叙述文・台詞の本文中の既知固有名を検出して narration_mention / line_mention へ
// 投入し、説明対象（e3）の解決を呼ぶこと。語彙は master_term ∪ proper_noun で、検出結果は
// 供給源に応じた排他列（master_term_id / proper_noun_id）へ写ること。
func TestIngestRecordsMentions(t *testing.T) {
	store := &fakeStore{
		terms:  []model.MasterTerm{{ID: 10, Source: "Whiterun", Dest: "ホワイトラン"}},
		proper: []model.ProperNoun{{ID: 20, Source: "Dragonbane", Category: "WEAP"}},
		allNarrations: []model.Narration{
			{ID: 1, Source: "Dragonbane was forged in Whiterun."}, // 両供給源の言及
			{ID: 2, Source: "No names here."},                     // 言及なし
		},
		allLines: []model.Line{
			{ID: 5, Source: "Take Dragonbane and run."}, // proper_noun の言及
		},
	}
	eng := New(store, &fakeTranslator{}, fakeLexicon{}, nil, nil)

	counts, err := eng.Ingest(context.Background())
	if err != nil {
		t.Fatalf("Ingest: %v", err)
	}

	wantNarr := []model.NarrationMention{
		{NarrationID: 1, ProperNounID: 20},
		{NarrationID: 1, MasterTermID: 10},
	}
	if len(store.narrationMentions) != 2 ||
		store.narrationMentions[0] != wantNarr[0] || store.narrationMentions[1] != wantNarr[1] {
		t.Errorf("narration_mention = %+v, want %+v", store.narrationMentions, wantNarr)
	}
	wantLine := model.LineMention{LineID: 5, ProperNounID: 20}
	if len(store.lineMentions) != 1 || store.lineMentions[0] != wantLine {
		t.Errorf("line_mention = %+v, want [%+v]", store.lineMentions, wantLine)
	}
	if !store.describedCalled {
		t.Error("説明対象（e3）の解決が呼ばれていない")
	}
	if counts.Mentions.NarrationMentions != 2 || counts.Mentions.LineMentions != 1 {
		t.Errorf("Mentions counts = %+v, want narration=2 line=1", counts.Mentions)
	}
}

// 同一原語が master_term と proper_noun の両方にある場合、機械置換辞書と同じ先勝ち
// （master_term 優先）で言及の相手が決まること。注入された側と言及の相手を一致させる。
func TestIngestMentionPrefersMasterTermOnDuplicateSource(t *testing.T) {
	store := &fakeStore{
		terms:         []model.MasterTerm{{ID: 10, Source: "Dragonbane", Dest: "ドラゴンベイン"}},
		proper:        []model.ProperNoun{{ID: 20, Source: "Dragonbane", Category: "WEAP"}},
		allNarrations: []model.Narration{{ID: 1, Source: "Dragonbane hums."}},
	}
	eng := New(store, &fakeTranslator{}, fakeLexicon{}, nil, nil)

	if _, err := eng.Ingest(context.Background()); err != nil {
		t.Fatalf("Ingest: %v", err)
	}

	want := model.NarrationMention{NarrationID: 1, MasterTermID: 10}
	if len(store.narrationMentions) != 1 || store.narrationMentions[0] != want {
		t.Errorf("narration_mention = %+v, want [%+v]（master_term 先勝ち）", store.narrationMentions, want)
	}
}

// 一般語 stoplist が注入（機械置換辞書）と言及（検出語彙）の供給を同じ選別で除くこと。
// 供給源は master_term・proper_noun の両方を跨ぎ、stoplist 語（Yes・No）は置換も言及もされない。
// stoplist に無い固有名（Riften）は両方に残る。片側だけに効く除外が無いことをこの 1 つの store で固定する。
func TestStoplistFiltersInjectionAndMentionAlike(t *testing.T) {
	const narrationText = "Yes and No, said the guard of Riften."
	store := &fakeStore{
		terms: []model.MasterTerm{
			{ID: 10, Source: "Yes", Dest: "はい"},
			{ID: 11, Source: "Riften", Dest: "リフテン"},
		},
		proper:        []model.ProperNoun{{ID: 20, Source: "No", Category: "FACT", Dest: "いいえ"}},
		allNarrations: []model.Narration{{ID: 1, Source: narrationText}},
	}
	stop, err := dictionary.ParseStoplist(strings.NewReader("yes\nno\n"))
	if err != nil {
		t.Fatalf("ParseStoplist: %v", err)
	}
	eng := New(store, &fakeTranslator{}, fakeLexicon{}, nil, stop)

	// 言及側: Ingest の言及段が stoplist 外の Riften だけを記録する。
	if _, ingestErr := eng.Ingest(context.Background()); ingestErr != nil {
		t.Fatalf("Ingest: %v", ingestErr)
	}
	wantMention := model.NarrationMention{NarrationID: 1, MasterTermID: 11}
	if len(store.narrationMentions) != 1 || store.narrationMentions[0] != wantMention {
		t.Errorf("narration_mention = %+v, want [%+v]（stoplist 語は言及されない）", store.narrationMentions, wantMention)
	}

	// 注入側: LoadDictionary の辞書も stoplist 外の Riften だけを置換する。
	dict, err := eng.LoadDictionary(context.Background())
	if err != nil {
		t.Fatalf("LoadDictionary: %v", err)
	}
	got, used := dict.Apply(narrationText)
	if want := "Yes and No, said the guard of リフテン."; got != want {
		t.Errorf("Apply = %q, want %q（stoplist 語は置換されない）", got, want)
	}
	if len(used) != 1 || used[0].Source != "Riften" {
		t.Errorf("Apply used = %+v, want [Riften]", used)
	}
}
