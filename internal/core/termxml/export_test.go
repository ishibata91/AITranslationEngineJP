package termxml

import (
	"bytes"
	"strings"
	"testing"
)

func TestMarshalStrings(t *testing.T) {
	// 叙述文（未完了訳）と固有名（確定訳）を xTranslator の SSTXMLRessources 形式へ直列化することを確かめる。
	entries := []StringEntry{
		{EDID: "DB01", Rec: "QUST", Field: "NNAM", Source: "Kill Grelod", Dest: "グレロッドを倒せ", Partial: "1"},
		{EDID: "SwordEbony", Rec: "WEAP", Field: "FULL", Source: "Ebony Sword", Dest: "黒檀の剣", Partial: ""},
	}

	out, err := MarshalStrings(entries, "MyMod")
	if err != nil {
		t.Fatalf("MarshalStrings が失敗: %v", err)
	}

	// 先頭は UTF-8 BOM。
	if !bytes.HasPrefix(out, []byte{0xEF, 0xBB, 0xBF}) {
		t.Errorf("先頭に UTF-8 BOM が無い")
	}
	got := string(out)
	// XML 宣言（UTF-8・standalone）とルート SSTXMLRessources。
	if !strings.Contains(got, `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`) {
		t.Errorf("XML 宣言が無い:\n%s", got)
	}
	if !strings.Contains(got, "<SSTXMLRessources>") || !strings.Contains(got, "</SSTXMLRessources>") {
		t.Errorf("SSTXMLRessources ルートが無い:\n%s", got)
	}
	// Params ブロック（Addon・言語・Version）。
	for _, want := range []string{
		"<Addon>MyMod</Addon>", "<Source>english</Source>", "<Dest>japanese</Dest>", "<Version>2</Version>",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("Params の %q が無い:\n%s", want, got)
		}
	}
	// Content>String と、REC は "REC:FIELD" 結合、EDID・原文・訳文。
	for _, want := range []string{
		"<Content>", `<REC>QUST:NNAM</REC>`, `<REC>WEAP:FULL</REC>`,
		"<EDID>DB01</EDID>", "<Source>Kill Grelod</Source>", "<Dest>グレロッドを倒せ</Dest>",
		"<Dest>黒檀の剣</Dest>",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("期待要素 %q が無い:\n%s", want, got)
		}
	}
	// List=0 固定、sID は 1 始まりの連番 hex(6桁)。
	if !strings.Contains(got, `List="0"`) || !strings.Contains(got, `sID="000001"`) || !strings.Contains(got, `sID="000002"`) {
		t.Errorf("List/sID 属性が期待どおりでない:\n%s", got)
	}
	// Status 要素・FORMID 要素は持たない（この形式には存在しない）。
	if strings.Contains(got, "<Status>") || strings.Contains(got, "<FORMID>") {
		t.Errorf("Status/FORMID 要素が出ている（この形式には無いはず）:\n%s", got)
	}
}

func TestMarshalStringsPartial(t *testing.T) {
	// Partial は確定訳（空）で属性ごと省略、未完了訳("1")で属性が出ることを確かめる。
	out, err := MarshalStrings([]StringEntry{
		{EDID: "A", Rec: "WEAP", Field: "FULL", Source: "s1", Dest: "d1", Partial: ""},
		{EDID: "B", Rec: "WEAP", Field: "FULL", Source: "s2", Dest: "d2", Partial: "1"},
	}, "Mod")
	if err != nil {
		t.Fatalf("MarshalStrings が失敗: %v", err)
	}

	got := string(out)
	// 確定訳の String（sID=000001）は Partial 属性を持たない。
	first := got[strings.Index(got, `sID="000001"`):strings.Index(got, `sID="000002"`)]
	if strings.Contains(first, "Partial=") {
		t.Errorf("確定訳の String に Partial 属性が出ている:\n%s", first)
	}
	// 未完了訳は Partial="1" を持つ。
	if !strings.Contains(got, `Partial="1"`) {
		t.Errorf("未完了訳の Partial=\"1\" が無い:\n%s", got)
	}
}

func TestMarshalStringsRejectsInvalidPartial(t *testing.T) {
	// Partial が想定外（""・"1"・"2" 以外）ならエラーにすることを確かめる。
	if _, err := MarshalStrings([]StringEntry{{EDID: "X", Partial: "3"}}, "Mod"); err == nil {
		t.Errorf("不正な Partial でエラーにならない")
	}
}

func TestMarshalStringsEscapesSpecialChars(t *testing.T) {
	// 原文・訳文の XML 特殊文字（<、&）がエスケープされ、生の記号が本文へ漏れないことを確かめる。
	out, err := MarshalStrings([]StringEntry{
		{EDID: "E", Rec: "BOOK", Field: "DESC", Source: "a < b & c", Dest: "訳 <b>", Partial: "1"},
	}, "Mod")
	if err != nil {
		t.Fatalf("MarshalStrings が失敗: %v", err)
	}

	got := string(out)
	if !strings.Contains(got, "a &lt; b &amp; c") {
		t.Errorf("原文の特殊文字がエスケープされていない:\n%s", got)
	}
	if strings.Contains(got, "<Dest>訳 <b></Dest>") {
		t.Errorf("訳文の < が生のまま出ている:\n%s", got)
	}
}

func TestMarshalStringsEmpty(t *testing.T) {
	// 行が 1 つも無い場合も、宣言・ルート・Params を持つ XML を返すことを確かめる。
	out, err := MarshalStrings(nil, "Mod")
	if err != nil {
		t.Fatalf("MarshalStrings が失敗: %v", err)
	}

	got := string(out)
	if !strings.Contains(got, "<SSTXMLRessources>") || !strings.Contains(got, "<Addon>Mod</Addon>") {
		t.Errorf("空でもルートと Params が要る:\n%s", got)
	}
}
