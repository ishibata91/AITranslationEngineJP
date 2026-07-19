package termxml

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"
)

// ReferenceEntry は 1 レコード分の既存訳（参照訳）。xTranslator 英日 XML の <String> の REC を
// rec / field に分け、原文（Source）と既訳（Dest）を持つ。項目7 の完全一致置換の照合単位。
// FormID は xTranslator XML が持たないため、照合は (rec, field, source) で行う。
type ReferenceEntry struct {
	Rec    string // レコード種別（例 WEAP）。REC の ":" 前。
	Field  string // フィールド（例 DESC）。REC の ":" 後。
	Source string // 原文（英語）。照合キー。XML の復号結果をそのまま持つ（前後空白も保つ）。
	Dest   string // 既訳（日本語）。置換で書き戻す確定訳。
}

// ReferencesFromFiles は読み込み済みの XML 群から、record 単位の既存訳を全て集める。
// 入力順は呼び出し側で決める（決定性のため engine はファイル名昇順で渡す）。
func ReferencesFromFiles(files []XMLFile) ([]ReferenceEntry, error) {
	var out []ReferenceEntry
	for _, file := range files {
		entries, err := ParseReferences(file.Data)
		if err != nil {
			return nil, err
		}
		out = append(out, entries...)
	}
	return out, nil
}

// ParseReferences は 1 つの XML から record 単位の既存訳を取り出す。
// REC を ":" で rec / field に分け、rec・field・原文・既訳がすべて非空の <String> だけを採る
// （rec / field は前後空白を除き大文字化して、抽出レコードの rec / field 表記へ合わせる）。
// encoding/xml はエンティティ（&lt; など）を復号する。先頭 UTF-8 BOM は除いてから解析する。
func ParseReferences(data []byte) ([]ReferenceEntry, error) {
	data = bytes.TrimPrefix(data, []byte("\uFEFF"))
	dec := xml.NewDecoder(bytes.NewReader(data))
	var out []ReferenceEntry
	for {
		tok, tokErr := dec.Token()
		if errors.Is(tokErr, io.EOF) {
			break
		}
		if tokErr != nil {
			return nil, fmt.Errorf("XML トークンの取得に失敗した: %w", tokErr)
		}
		start, ok := tok.(xml.StartElement)
		if !ok || start.Name.Local != "String" {
			continue
		}
		var s xmlString
		if decErr := dec.DecodeElement(&s, &start); decErr != nil {
			return nil, fmt.Errorf("XML の String 要素の復号に失敗した: %w", decErr)
		}
		if e, ok := referenceOf(s); ok {
			out = append(out, e)
		}
	}
	return out, nil
}

// referenceOf は 1 つの <String> を参照訳 1 件へ変換する。
// rec / field / 原文 / 既訳のいずれかが空なら採らない（照合に使えないため）。
// 原文・既訳は復号したまま（前後空白を保つ）。抽出レコードの source と厳密一致で照合するため。
func referenceOf(s xmlString) (ReferenceEntry, bool) {
	rec, field, found := strings.Cut(strings.ToUpper(strings.TrimSpace(s.REC)), ":")
	if !found || rec == "" || field == "" {
		return ReferenceEntry{}, false
	}
	if strings.TrimSpace(s.Source) == "" || strings.TrimSpace(s.Dest) == "" {
		return ReferenceEntry{}, false
	}
	return ReferenceEntry{Rec: rec, Field: field, Source: s.Source, Dest: s.Dest}, true
}
