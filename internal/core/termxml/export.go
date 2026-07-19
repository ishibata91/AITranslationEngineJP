package termxml

import (
	"encoding/xml"
	"fmt"
)

// bomUTF8 は UTF-8 のバイトオーダーマーク（U+FEFF）。base 辞書の出力に合わせ、宣言の前に置く。
// ソースへ BOM 文字を直書きしないためバイト列で表す。
var bomUTF8 = []byte{0xEF, 0xBB, 0xBF}

// xmlHeader は xTranslator（SSTXMLRessources 形式）の宣言行。UTF-8（大文字）＋ standalone="yes"。
const xmlHeader = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n"

// StringEntry は xTranslator の SSTXMLRessources 形式の <String> 1 行分の書き出し値。
// 翻訳結果（叙述文・台詞・固有名）を xTranslator が読める形へ写した中間表現。
type StringEntry struct {
	EDID   string // レコードの EditorID
	Rec    string // 4 文字レコード種別（例 QUST・NPC_）。REC 要素へ "REC:FIELD" として結合する。
	Field  string // 4 文字フィールド（例 NNAM・NAM1・FULL）
	Source string // 原文
	Dest   string // 訳文
	// Partial は訳状態フラグ。""=確定訳（translated）、"1"=未完了訳（incompleteTrans）、"2"=ロック訳（lockedTrans）。
	// xTranslator は Status 要素を持たず、訳状態をこの属性で表す。
	Partial string
}

// xmlExportString は <String> 要素の marshal 用。属性 List・sID・Partial と、子要素 EDID・REC・Source・Dest を持つ。
// Partial は omitempty で、確定訳（空）のときに属性ごと省略する（xTranslator の「属性なし＝translated」に合わせる）。
type xmlExportString struct {
	List    string `xml:"List,attr"`
	SID     string `xml:"sID,attr"`
	Partial string `xml:"Partial,attr,omitempty"`
	EDID    string `xml:"EDID"`
	REC     string `xml:"REC"`
	Source  string `xml:"Source"`
	Dest    string `xml:"Dest"`
}

// xmlExportParams は Params ブロック。Addon は対象 plugin 名、Source/Dest は言語、Version は 2（現行形式）。
type xmlExportParams struct {
	Addon   string `xml:"Addon"`
	Source  string `xml:"Source"`
	Dest    string `xml:"Dest"`
	Version int    `xml:"Version"`
}

// xmlExportResources は SSTXMLRessources ルート。Params と Content>String 群を持つ。
type xmlExportResources struct {
	XMLName xml.Name          `xml:"SSTXMLRessources"`
	Params  xmlExportParams   `xml:"Params"`
	Strings []xmlExportString `xml:"Content>String"`
}

// MarshalStrings は StringEntry 群を xTranslator の SSTXMLRessources 形式 XML へ直列化する。
// addon は Params の Addon（対象 plugin 名）。Source=english・Dest=japanese・Version=2 を付ける。
// 各 <String> の sID は 1 始まりの連番 hex（6 桁大文字）、List は 0 固定（インポートの照合は EDID・REC・原文で行われる）。
// REC は "REC:FIELD" に結合する。Partial は StringEntry の値をそのまま属性へ写す（空なら省略）。
// XML 特殊文字のエスケープは encoding/xml に委ねる。ファイル入出力は持たない（呼び出し側が書き込む）。
func MarshalStrings(entries []StringEntry, addon string) ([]byte, error) {
	res := xmlExportResources{
		Params:  xmlExportParams{Addon: addon, Source: "english", Dest: "japanese", Version: 2},
		Strings: make([]xmlExportString, len(entries)),
	}
	for i, e := range entries {
		if e.Partial != "" && e.Partial != "1" && e.Partial != "2" {
			return nil, fmt.Errorf("不正な Partial 属性: EDID=%s Partial=%q", e.EDID, e.Partial)
		}
		res.Strings[i] = xmlExportString{
			List:    "0",
			SID:     fmt.Sprintf("%06X", i+1),
			Partial: e.Partial,
			EDID:    e.EDID,
			REC:     fmt.Sprintf("%s:%s", e.Rec, e.Field),
			Source:  e.Source,
			Dest:    e.Dest,
		}
	}
	body, err := xml.MarshalIndent(res, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("xTranslator XML の直列化: %w", err)
	}
	out := make([]byte, 0, len(bomUTF8)+len(xmlHeader)+len(body)+1)
	out = append(out, bomUTF8...)
	out = append(out, xmlHeader...)
	out = append(out, body...)
	out = append(out, '\n')
	return out, nil
}
