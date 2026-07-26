package engine

import (
	"context"
	"testing"

	"aitranslationenginejp/internal/model"
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
