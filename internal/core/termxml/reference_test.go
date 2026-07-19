package termxml

import (
	"reflect"
	"testing"
)

// ParseReferences は REC を rec / field に分け、原文と既訳の揃った <String> を参照訳として採る。
func TestParseReferencesCollectsRecordLevelEntries(t *testing.T) {
	data := []byte(`<?xml version="1.0" encoding="utf-8"?>
<SSTXMLRessources>
  <Content>
    <String><EDID>TestSword</EDID><REC>WEAP:DESC</REC><Source>A fine blade.</Source><Dest>見事な刃。</Dest></String>
    <String><EDID>Greeting</EDID><REC>INFO:NAM1</REC><Source>Hello there.</Source><Dest>こんにちは。</Dest></String>
  </Content>
</SSTXMLRessources>`)
	got, err := ParseReferences(data)
	if err != nil {
		t.Fatalf("解析に失敗: %v", err)
	}
	want := []ReferenceEntry{
		{Rec: "WEAP", Field: "DESC", Source: "A fine blade.", Dest: "見事な刃。"},
		{Rec: "INFO", Field: "NAM1", Source: "Hello there.", Dest: "こんにちは。"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("参照訳が期待と違う:\n got=%#v\n want=%#v", got, want)
	}
}

// ParseReferences は既訳（Dest）が空の <String> を採らない（未訳の参照は置換に使えない）。
func TestParseReferencesSkipsEmptyDest(t *testing.T) {
	data := []byte(`<SSTXMLRessources><Content>
    <String><REC>WEAP:DESC</REC><Source>Untranslated blade.</Source><Dest></Dest></String>
    <String><REC>WEAP:DESC</REC><Source>Translated blade.</Source><Dest>訳済みの刃。</Dest></String>
  </Content></SSTXMLRessources>`)
	got, err := ParseReferences(data)
	if err != nil {
		t.Fatalf("解析に失敗: %v", err)
	}
	// 既訳ありの 1 件だけ採る。
	if len(got) != 1 || got[0].Source != "Translated blade." {
		t.Fatalf("既訳なしを除外できていない: %#v", got)
	}
}

// ParseReferences は REC が rec:field 形式でない <String> を採らない（照合できないため）。
func TestParseReferencesSkipsRecWithoutField(t *testing.T) {
	data := []byte(`<SSTXMLRessources><Content>
    <String><REC>WEAP</REC><Source>No field.</Source><Dest>フィールド無し。</Dest></String>
  </Content></SSTXMLRessources>`)
	got, err := ParseReferences(data)
	if err != nil {
		t.Fatalf("解析に失敗: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("field を持たない REC を採ってしまった: %#v", got)
	}
}

// ParseReferences は XML エンティティを復号し、原文の内部空白をそのまま保つ（厳密一致照合のため）。
func TestParseReferencesDecodesEntitiesAndPreservesSource(t *testing.T) {
	data := []byte(`<SSTXMLRessources><Content>
    <String><REC>BOOK:DESC</REC><Source>Tom &amp;  Jerry</Source><Dest>トムと ジェリー</Dest></String>
  </Content></SSTXMLRessources>`)
	got, err := ParseReferences(data)
	if err != nil {
		t.Fatalf("解析に失敗: %v", err)
	}
	if len(got) != 1 || got[0].Source != "Tom &  Jerry" {
		t.Fatalf("エンティティ復号または内部空白保持が期待と違う: %#v", got)
	}
}

// ReferencesFromFiles は複数 XML の参照訳をファイル入力順に連結する。
func TestReferencesFromFilesConcatenatesInOrder(t *testing.T) {
	files := []XMLFile{
		{Name: "A.xml", Data: []byte(`<SSTXMLRessources><Content><String><REC>WEAP:FULL</REC><Source>Iron Sword</Source><Dest>鉄の剣</Dest></String></Content></SSTXMLRessources>`)},
		{Name: "B.xml", Data: []byte(`<SSTXMLRessources><Content><String><REC>ARMO:FULL</REC><Source>Iron Shield</Source><Dest>鉄の盾</Dest></String></Content></SSTXMLRessources>`)},
	}
	got, err := ReferencesFromFiles(files)
	if err != nil {
		t.Fatalf("解析に失敗: %v", err)
	}
	want := []ReferenceEntry{
		{Rec: "WEAP", Field: "FULL", Source: "Iron Sword", Dest: "鉄の剣"},
		{Rec: "ARMO", Field: "FULL", Source: "Iron Shield", Dest: "鉄の盾"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("連結順が期待と違う:\n got=%#v\n want=%#v", got, want)
	}
}
