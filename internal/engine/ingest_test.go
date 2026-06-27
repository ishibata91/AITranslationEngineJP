package engine

import (
	"testing"

	"aitranslationenginejp/internal/model"
)

// Dispatch は extracted_field を record_type_master の box で振り分けること。
// 固有名 → proper_noun（category＝rec）、叙述文・定型句 → narration（style＝directive キー）、
// 台詞 → line（response_order＝ordinal）。master に無い (rec, field) と未知の box は Skipped に数える。
func TestDispatch(t *testing.T) {
	master := map[RecordKey]model.RecordType{
		{Rec: "BOOK", Field: "DESC"}: {Rec: "BOOK", Field: "DESC", Box: boxNarration, Directive: "書物体"},
		{Rec: "ACTI", Field: "RNAM"}: {Rec: "ACTI", Field: "RNAM", Box: boxSetPhrase, Directive: "定型句"},
		{Rec: "WEAP", Field: "FULL"}: {Rec: "WEAP", Field: "FULL", Box: boxProperNoun, Directive: "固有名"},
		{Rec: "INFO", Field: "NAM1"}: {Rec: "INFO", Field: "NAM1", Box: boxLine, Directive: "口調"},
		{Rec: "BAD_", Field: "XXXX"}: {Rec: "BAD_", Field: "XXXX", Box: "未知の箱", Directive: "x"},
	}
	fields := []model.ExtractedField{
		{Plugin: "P.esm", FormID: "0x01", EDID: "Book1", Rec: "BOOK", Field: "DESC", Ordinal: 0, Source: "a book"},
		{Plugin: "P.esm", FormID: "0x02", EDID: "Door1", Rec: "ACTI", Field: "RNAM", Ordinal: 0, Source: "Open"},
		{Plugin: "P.esm", FormID: "0x03", EDID: "Sword1", Rec: "WEAP", Field: "FULL", Ordinal: 0, Source: "Dragonbane"},
		{Plugin: "P.esm", FormID: "0x04", EDID: "Info1", Rec: "INFO", Field: "NAM1", Ordinal: 2, Source: "Hello."},
		{Plugin: "P.esm", FormID: "0x05", EDID: "Header", Rec: "TES4", Field: "CNAM", Ordinal: 0, Source: "author"}, // master に無い → 読み込まない
		{Plugin: "P.esm", FormID: "0x06", EDID: "Bad", Rec: "BAD_", Field: "XXXX", Ordinal: 0, Source: "x"},         // 未知の box
	}

	d := Dispatch(fields, master)

	if len(d.ProperNouns) != 1 || d.ProperNouns[0].Source != "Dragonbane" || d.ProperNouns[0].Category != "WEAP" {
		t.Errorf("固有名 = %+v, want [Dragonbane/WEAP]", d.ProperNouns)
	}
	if len(d.Narrations) != 2 {
		t.Fatalf("叙述文・定型句 = %d 件, want 2", len(d.Narrations))
	}
	// 叙述文（書物体）の style に directive キーが入ること。
	if d.Narrations[0].Style != "書物体" || d.Narrations[0].Rec != "BOOK" {
		t.Errorf("叙述文[0] = %+v, want style=書物体 rec=BOOK", d.Narrations[0])
	}
	// 定型句も narration へ入り、style に定型句 directive が入ること。
	if d.Narrations[1].Style != "定型句" || d.Narrations[1].Rec != "ACTI" {
		t.Errorf("定型句 = %+v, want style=定型句 rec=ACTI", d.Narrations[1])
	}
	if len(d.Lines) != 1 || d.Lines[0].Source != "Hello." || d.Lines[0].ResponseOrder != 2 {
		t.Errorf("台詞 = %+v, want [Hello./order2]", d.Lines)
	}
	// master に無い TES4:CNAM と未知 box の 2 件は読み込まない。
	if d.Skipped != 2 {
		t.Errorf("Skipped = %d, want 2", d.Skipped)
	}
}

// Dispatch は空入力で空の結果を返すこと（境界）。
func TestDispatchEmpty(t *testing.T) {
	d := Dispatch(nil, map[RecordKey]model.RecordType{})
	if len(d.Narrations) != 0 || len(d.ProperNouns) != 0 || len(d.Lines) != 0 || d.Skipped != 0 {
		t.Errorf("空入力で空でない: %+v", d)
	}
}

// recordMasterMap は record_type_master の行群を (rec, field) の map に畳むこと。
func TestRecordMasterMap(t *testing.T) {
	rows := []model.RecordType{
		{Rec: "BOOK", Field: "DESC", Box: boxNarration, Directive: "書物体"},
		{Rec: "WEAP", Field: "FULL", Box: boxProperNoun, Directive: "固有名"},
	}
	m := recordMasterMap(rows)
	if len(m) != 2 {
		t.Fatalf("map size = %d, want 2", len(m))
	}
	if m[RecordKey{Rec: "WEAP", Field: "FULL"}].Box != boxProperNoun {
		t.Errorf("WEAP:FULL の box が引けない: %+v", m)
	}
}
