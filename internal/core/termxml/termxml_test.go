package termxml

import (
	"reflect"
	"testing"

	"aitranslationenginejp/internal/core/termderive"
	"aitranslationenginejp/internal/core/termusage"
)

// xTranslator XML から NPC のフルネーム対・短名対と会話文を取り出すこと。
// 先頭 BOM の除去と、エンティティ（&lt;）の復号、対象 REC 以外の素通しを確かめる。
func TestParseTermXML(t *testing.T) {
	const fixture = "\uFEFF" + `<?xml version="1.0" encoding="UTF-8"?>
<SSTXMLRessources>
  <Content>
    <String><REC>NPC_:FULL</REC><Source>Grelod the Kind</Source><Dest>親切者のグレロッド</Dest></String>
    <String><REC>NPC_:SHRT</REC><Source>Maro</Source><Dest>マロ</Dest></String>
    <String><REC>INFO:NAM1</REC><Source>Bring it to &lt;Alias&gt;.</Source><Dest>x</Dest></String>
    <String><REC>WEAP:FULL</REC><Source>Iron Sword</Source><Dest>鉄の剣</Dest></String>
  </Content>
</SSTXMLRessources>`

	fulls, shrts, dialogues, err := ParseTermXML([]byte(fixture), true)
	if err != nil {
		t.Fatalf("ParseTermXML error: %v", err)
	}
	// NPC_:FULL だけを fulls に取り、base ゲーム印を付ける。
	if len(fulls) != 1 || fulls[0].Source != "Grelod the Kind" || !fulls[0].BaseGame {
		t.Errorf("fulls = %+v, want 1 件の Grelod the Kind（base ゲーム）", fulls)
	}
	// NPC_:SHRT だけを shrts に取る。
	if len(shrts) != 1 || shrts[0].Source != "Maro" || shrts[0].Dest != "マロ" {
		t.Errorf("shrts = %+v, want 1 件の Maro=>マロ", shrts)
	}
	// INFO:NAM1 の英語原文を会話文に取り、エンティティを復号する。
	if len(dialogues) != 1 || dialogues[0] != "Bring it to <Alias>." {
		t.Errorf("dialogues = %v, want 1 件の復号済み原文", dialogues)
	}
}

// ParseTermXML は壊れた XML（閉じない要素）をトークン走査の段階でエラーにする。
func TestParseTermXMLBroken(t *testing.T) {
	if _, _, _, err := ParseTermXML([]byte("<SSTXMLRessources><String>"), true); err == nil {
		t.Fatal("壊れた XML でエラーを期待")
	}
}

// IsBaseGame は base ゲーム接頭（Skyrim など）で true、第三者 mod 名で false を返す。
func TestIsBaseGame(t *testing.T) {
	for _, name := range []string{"Skyrim.xml", "Dawnguard_jp.xml", "Update.xml"} {
		if !IsBaseGame(name) {
			t.Errorf("IsBaseGame(%q) = false, want true", name)
		}
	}
	for _, name := range []string{"MyMod.xml", "Inigo.xml", ""} {
		if IsBaseGame(name) {
			t.Errorf("IsBaseGame(%q) = true, want false", name)
		}
	}
}

// DeriveTermsFromFiles は読み込み済み XML 群を ParseTermXML→BuildUsage→DeriveTerms の順に結線する。
// 期待値は同じ構成要素を手で合成したもので、結線（ファイル名の base ゲーム判定・用法集計・派生）を確かめる。
func TestDeriveTermsFromFiles(t *testing.T) {
	const xml = `<SSTXMLRessources><Content>
    <String><REC>NPC_:FULL</REC><Source>Grelod the Kind</Source><Dest>親切者のグレロッド</Dest></String>
    <String><REC>NPC_:SHRT</REC><Source>Maro</Source><Dest>マロ</Dest></String>
    <String><REC>INFO:NAM1</REC><Source>Maro greets Grelod warmly.</Source><Dest>x</Dest></String>
  </Content></SSTXMLRessources>`
	files := []XMLFile{{Name: "Skyrim_Test.xml", Data: []byte(xml)}}

	got, err := DeriveTermsFromFiles(files, map[string]bool{})
	if err != nil {
		t.Fatalf("DeriveTermsFromFiles error: %v", err)
	}
	fulls, shrts, dialogues, _ := ParseTermXML([]byte(xml), true)
	want := termderive.DeriveTerms(fulls, shrts, termusage.BuildUsage(dialogues),
		map[string]bool{}, termderive.DefaultDeriveConfig())
	if !reflect.DeepEqual(got, want) {
		t.Errorf("DeriveTermsFromFiles = %+v, want（手合成）= %+v", got, want)
	}
}
