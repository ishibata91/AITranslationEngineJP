package engine

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"aitranslationenginejp/internal/model"
	"aitranslationenginejp/internal/provider"
)

// SelectSupply は固有名の既訳の有無から供給源を選ぶこと。
// 既訳あり → 権威訳（Dest を埋め、書き込み宛先なし）、既訳なし → AI 訳（宛先は proper_noun）。
func TestSelectSupply(t *testing.T) {
	authoritative := map[string]string{
		"Riften":     "リフテン",
		"Dragonbane": "竜の災い",
	}

	// 既訳あり: 権威訳を返し、AI 訳しない（WriteTarget は空＝書き戻し不要）。
	hit := SelectSupply("Riften", authoritative)
	if hit.Kind != SupplyAuthoritative {
		t.Errorf("既訳ありの Kind = %d, want SupplyAuthoritative", hit.Kind)
	}
	if hit.Dest != "リフテン" {
		t.Errorf("既訳ありの Dest = %q, want リフテン", hit.Dest)
	}
	if hit.WriteTarget != "" {
		t.Errorf("既訳ありの WriteTarget = %q, want 空", hit.WriteTarget)
	}

	// 既訳なし: AI 訳対象で、宛先は proper_noun（master_term ではない）。
	miss := SelectSupply("Whiterun", authoritative)
	if miss.Kind != SupplyAITranslate {
		t.Errorf("既訳なしの Kind = %d, want SupplyAITranslate", miss.Kind)
	}
	if miss.Dest != "" {
		t.Errorf("既訳なしの Dest = %q, want 空", miss.Dest)
	}
	if miss.WriteTarget != properNounTable {
		t.Errorf("既訳なしの WriteTarget = %q, want %q", miss.WriteTarget, properNounTable)
	}
}

// 不変ルール: どの入力でも、AI 訳の書き込み宛先は proper_noun で、master_term を返さない（方針A）。
// 既訳あり・なしの両経路を網羅し、宛先が master_term になり得ないことを確かめる。
func TestSelectSupplyNeverWritesMasterTerm(t *testing.T) {
	authoritative := map[string]string{"Known": "既訳"}
	for _, source := range []string{"Known", "Unknown", "", "Another"} {
		sup := SelectSupply(source, authoritative)
		if sup.WriteTarget == masterTermTable {
			t.Errorf("source=%q で WriteTarget が master_term になった（不変ルール違反）", source)
		}
		if sup.Kind == SupplyAITranslate && sup.WriteTarget != properNounTable {
			t.Errorf("source=%q の AI 訳宛先 = %q, want %q", source, sup.WriteTarget, properNounTable)
		}
	}
}

// 空の既訳辞書では、すべて AI 訳（既訳ヒットなし）になること（境界）。
func TestSelectSupplyEmptyAuthoritative(t *testing.T) {
	sup := SelectSupply("Anything", map[string]string{})
	if sup.Kind != SupplyAITranslate || sup.WriteTarget != properNounTable {
		t.Errorf("空辞書での判定 = %+v, want AI 訳→proper_noun", sup)
	}
}

// deriveRunProperNouns の観測用 fixture。実行内で氏名の訳が確定した mod NPC を 1 人持ち、
// 派生した部分形の書き込み先（plugin・宛先テーブル）と、既出原語の除外を確かめる土台にする。
func storeWithConfirmedName(plugin string) *fakeStore {
	return &fakeStore{
		proper: []model.ProperNoun{
			{ID: 1, Plugin: plugin, Source: "Sorine Trueblade", Category: "NPC_", Dest: "ソリーヌ・トゥルーブレイド", Status: statusProvisional},
		},
		extracted: []model.ExtractedField{
			{Plugin: plugin, Rec: "NPC_", Field: "FULL", Source: "Sorine Trueblade"},
		},
	}
}

// R-2-4: 実行内で確定した氏名から作った部分形を、別の plugin を翻訳する実行へ持ち越さないこと。
// 派生行は対象 plugin の proper_noun（plugin スコープの非共有）へ書き、横断永続辞書 master_term へは書かない。
func TestDeriveRunProperNounsStaysInTargetPlugin(t *testing.T) {
	f := storeWithConfirmedName("Mod.esp")
	e := New(f, nil, nil, nil, nil)

	inserted, err := e.deriveRunProperNouns(context.Background(), "Mod.esp")
	if err != nil {
		t.Fatalf("部分形の派生: %v", err)
	}
	if inserted == 0 {
		t.Fatalf("部分形が 1 件も派生しなかった: %v", f.derivedPropers)
	}
	for _, pn := range f.derivedPropers {
		if pn.Plugin != "Mod.esp" {
			t.Errorf("派生した部分形 %q の plugin = %q, want Mod.esp", pn.Source, pn.Plugin)
		}
	}
	// 横断永続辞書へは 1 件も書かない（方針A の不変境界）。
	if len(f.insertedTerms) != 0 {
		t.Errorf("派生した部分形が master_term へ書かれた: %v", f.insertedTerms)
	}
	// 別 plugin を対象にした実行では、その plugin の確定固有名が無いので何も派生しない。
	other := storeWithConfirmedName("Mod.esp")
	if n, err := New(other, nil, nil, nil, nil).deriveRunProperNouns(context.Background(), "Other.esp"); err != nil || n != 0 {
		t.Errorf("別 plugin の実行で派生した: n=%d err=%v rows=%v", n, err, other.derivedPropers)
	}
}

// R-2-5: 横断辞書に既にある原語について、実行内で確定した氏名から別の訳を作らないこと。
// 既出原語（master_term ∪ その実行で確定済みの固有名）は派生に含めず、同じ原語へ 2 つの訳が立たないようにする。
func TestDeriveRunProperNounsSkipsExistingSources(t *testing.T) {
	f := storeWithConfirmedName("Mod.esp")
	// 横断辞書が同じ原語 Sorine を別の訳で既に持つ。
	f.terms = []model.MasterTerm{{Source: "Sorine", Dest: "ソリーン", Category: "NPC_"}}
	e := New(f, nil, nil, nil, nil)

	if _, err := e.deriveRunProperNouns(context.Background(), "Mod.esp"); err != nil {
		t.Fatalf("部分形の派生: %v", err)
	}
	for _, pn := range f.derivedPropers {
		if pn.Source == "Sorine" {
			t.Errorf("横断辞書に既にある原語 Sorine を派生した: dest=%q", pn.Dest)
		}
	}
	// 既出でない苗字側は従来どおり派生する（除外が広がりすぎていないこと）。
	found := false
	for _, pn := range f.derivedPropers {
		if pn.Source == "Trueblade" {
			found = true
		}
	}
	if !found {
		t.Errorf("既出でない原語 Trueblade が派生しなかった: %v", f.derivedPropers)
	}
}

func TestDeriveRunProperNounsDerivesHyphenatedNameParts(t *testing.T) {
	const plugin = "Mod.esp"
	store := &fakeStore{
		proper: []model.ProperNoun{{
			ID: 1, Plugin: plugin, Source: "Hoge Black-Briar", Category: recNPC,
			Dest: "ホゲ・ブラック・ブライア", Status: statusProvisional,
		}},
		extracted: []model.ExtractedField{{Plugin: plugin, Rec: recNPC, Field: fieldFull, Source: "Hoge Black-Briar"}},
	}
	if _, err := New(store, nil, nil, nil, nil).deriveRunProperNouns(context.Background(), plugin); err != nil {
		t.Fatalf("deriveRunProperNouns: %v", err)
	}
	want := map[string]string{"Hoge": "ホゲ", "Black-Briar": "ブラック・ブライア"}
	for _, row := range store.derivedPropers {
		if got, ok := want[row.Source]; ok && row.Dest == got && row.Plugin == plugin && row.Origin == model.OriginDerived {
			delete(want, row.Source)
		}
	}
	if len(want) != 0 {
		t.Errorf("ハイフンを含む姓の部分形 = %+v, missing=%v", store.derivedPropers, want)
	}
}

// 仕様: 固有名 1 件の翻訳が構造化出力の空で終わったとき、その固有名を未訳のまま残し、
// 残りの固有名と叙述文と台詞を訳し切ること。
// 実 LLM（7B）の空応答 1 件で実行全体が止まった不具合の再発防止（empty-translation-halts-run）。
func TestRunSkipsProperNounOnStructuredParseFailure(t *testing.T) {
	store := &fakeStore{
		proper: []model.ProperNoun{
			{ID: 1, Source: "Inigo", Category: "NPC_"},
			{ID: 2, Source: "Riften", Category: "CELL"},
		},
		untranslated: []model.Narration{{ID: 10, Source: "halls"}},
		lines:        []model.Line{{ID: 20, Source: "hello"}},
	}
	tr := &fakeTranslator{
		out: map[string]string{"Riften": "リフテン", "halls": "広間", "hello": "やあ"},
		// 固有名 Inigo の応答だけ translation が空（構造化出力の解析失敗）。
		errByUser: map[string]error{
			"Inigo": fmt.Errorf("%w: translation が空", provider.ErrStructuredParse),
		},
	}
	eng := New(store, tr, fakeLexicon{}, nil, nil)

	count, err := eng.Run(context.Background(), provider.Connection{Endpoint: "http://x"}, "model-x", "", nil)
	if err != nil {
		t.Fatalf("固有名の空応答で Run が停止した（未訳のまま残して続けるべき）: %v", err)
	}
	// 空応答の Inigo(ID=1) は書き戻さず、Riften(ID=2) だけを仮訳で確定すること。
	if len(store.properUpdates) != 1 || store.properUpdates[0] != (update{2, "リフテン", 3}) {
		t.Errorf("properUpdates = %v, want [{2 リフテン 3}]（空応答の固有名は未訳で残す）", store.properUpdates)
	}
	// 後続の叙述文・台詞フェーズまで進み、本文を訳し切ること。
	if len(store.updates) != 1 || store.updates[0] != (update{10, "広間", 3}) {
		t.Errorf("叙述文 updates = %v, want [{10 広間 3}]（固有名の失敗で本文が止まらない）", store.updates)
	}
	if len(store.lineUpdates) != 1 || store.lineUpdates[0] != (update{20, "やあ", 3}) {
		t.Errorf("台詞 lineUpdates = %v, want [{20 やあ 3}]（固有名の失敗で本文が止まらない）", store.lineUpdates)
	}
	// 未訳のまま残した固有名は進捗に数えないこと（固有名 1 + 叙述文 1 + 台詞 1 = 3）。
	if count != 3 {
		t.Errorf("count = %d, want 3（未訳で残した固有名は数えない）", count)
	}
}

// 仕様: 固有名 1 件の翻訳が応答エンベロープの読み取り失敗またはサーバ一時失敗で終わったとき、
// その固有名を未訳のまま残し、残り全件を訳し切ること。
func TestRunSkipsProperNounOnSkippableProviderFailures(t *testing.T) {
	store := &fakeStore{
		proper: []model.ProperNoun{
			{ID: 1, Source: "Inigo", Category: "NPC_"},
			{ID: 2, Source: "Lucien", Category: "NPC_"},
			{ID: 3, Source: "Riften", Category: "CELL"},
		},
		untranslated: []model.Narration{{ID: 10, Source: "halls"}},
	}
	tr := &fakeTranslator{
		out: map[string]string{"Riften": "リフテン", "halls": "広間"},
		errByUser: map[string]error{
			"Inigo":  fmt.Errorf("%w: choices が無い", provider.ErrResponseUnreadable),
			"Lucien": fmt.Errorf("%w: status 503", provider.ErrServerTransient),
		},
	}
	eng := New(store, tr, fakeLexicon{}, nil, nil)

	count, err := eng.Run(context.Background(), provider.Connection{Endpoint: "http://x"}, "model-x", "", nil)
	if err != nil {
		t.Fatalf("固有名の skippable 失敗で Run が停止した（未訳のまま残して続けるべき）: %v", err)
	}
	if len(store.properUpdates) != 1 || store.properUpdates[0] != (update{3, "リフテン", 3}) {
		t.Errorf("properUpdates = %v, want [{3 リフテン 3}]（skippable 失敗の固有名は未訳で残す）", store.properUpdates)
	}
	if len(store.updates) != 1 || store.updates[0] != (update{10, "広間", 3}) {
		t.Errorf("叙述文 updates = %v, want [{10 広間 3}]（固有名の失敗で本文が止まらない）", store.updates)
	}
	// 固有名 1 件 + 叙述文 1 件。未訳のまま残した 2 件は数えない。
	if count != 2 {
		t.Errorf("count = %d, want 2（未訳で残した固有名は数えない）", count)
	}
}

// 仕様: 固有名 1 件の翻訳が skippable な失敗のいずれにも当たらない失敗で終わったとき、
// 実行を止めて画面へ失敗を出すこと（engine は失敗を返し、後続フェーズへ進まない）。
func TestRunStopsOnFatalProperNounFailure(t *testing.T) {
	store := &fakeStore{
		proper: []model.ProperNoun{
			{ID: 1, Source: "Inigo", Category: "NPC_"},
			{ID: 2, Source: "Riften", Category: "CELL"},
		},
		untranslated: []model.Narration{{ID: 10, Source: "halls"}},
	}
	tr := &fakeTranslator{
		out: map[string]string{"Riften": "リフテン", "halls": "広間"},
		// 認証の失敗（設定起因の 4xx）。skippable 番兵でラップされない失敗。
		errByUser: map[string]error{"Inigo": errors.New("翻訳要求: status 401")},
	}
	eng := New(store, tr, fakeLexicon{}, nil, nil)

	_, err := eng.Run(context.Background(), provider.Connection{Endpoint: "http://x"}, "model-x", "", nil)
	if err == nil {
		t.Fatal("固有名の fatal 失敗で Run が止まらなかった")
	}
	// 失敗した固有名より後の固有名も、本文も処理しないこと。
	if len(store.properUpdates) != 0 {
		t.Errorf("properUpdates = %v, want 空（fatal 以降は処理しない）", store.properUpdates)
	}
	if len(store.updates) != 0 {
		t.Errorf("叙述文 updates = %v, want 空（fatal で本文フェーズへ進まない）", store.updates)
	}
}
