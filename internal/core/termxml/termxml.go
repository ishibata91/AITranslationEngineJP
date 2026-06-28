// Package termxml は xTranslator 英日 XML を解析し、人名の部分形（名のみ・短名）を派生する純粋ルール。
// ファイル I/O は持たない。呼び出し側が読み込んだ XML の中身（[]byte）を XMLFile で渡す
// （ディレクトリ走査・os.ReadFile は engine 側の DeriveMasterTerms が行う）。
package termxml

import (
	"bytes"
	"encoding/xml"
	"io"
	"strings"

	"aitranslationenginejp/internal/core/termderive"
	"aitranslationenginejp/internal/core/termusage"
)

// baseGamePrefixes は base ゲーム（Bethesda 公式本体・DLC）の XML ファイル名接頭。姓名分割（two）を
// base ゲーム限定にするための判定に使う。第三者 mod は Source/Dest 対応が信頼できないため two を派生しない。
var baseGamePrefixes = []string{"Skyrim", "Dawnguard", "Dragonborn", "HearthFires", "Update"}

// IsBaseGame は XML ファイル名が base ゲーム由来かを返す。
func IsBaseGame(filename string) bool {
	for _, p := range baseGamePrefixes {
		if strings.HasPrefix(filename, p) {
			return true
		}
	}
	return false
}

// XMLFile は読み込み済みの 1 つの xTranslator XML。Name はファイル名（base ゲーム判定に使う）、
// Data は中身（os.ReadFile の結果）。ファイル open は呼び出し側が行い、純粋部はこの object だけ受け取る。
type XMLFile struct {
	Name string
	Data []byte
}

// DeriveTermsFromFiles は読み込み済みの XML 群を解析し、用法分布を作り、純粋ルールで派生対を返す。
// baseSources は base 辞書の既存原語集合で、既出原語との衝突を避ける除外判定に使う。
// 入力順は呼び出し側で決める（決定性のため engine はファイル名昇順で渡す）。
func DeriveTermsFromFiles(files []XMLFile, baseSources map[string]bool) ([]termderive.DerivedTerm, error) {
	var fulls, shrts []termderive.NamePair
	var dialogues []string
	for _, file := range files {
		f, s, d, err := ParseTermXML(file.Data, IsBaseGame(file.Name))
		if err != nil {
			return nil, err
		}
		fulls = append(fulls, f...)
		shrts = append(shrts, s...)
		dialogues = append(dialogues, d...)
	}
	usage := termusage.BuildUsage(dialogues)
	return termderive.DeriveTerms(fulls, shrts, usage, baseSources, termderive.DefaultDeriveConfig()), nil
}

// xmlString は xTranslator XML の 1 レコード（<String>）。REC（種別）と原語・確定訳語を読む。
type xmlString struct {
	REC    string `xml:"REC"`
	Source string `xml:"Source"`
	Dest   string `xml:"Dest"`
}

// ParseTermXML は 1 つの XML から NPC のフルネーム対・短名対と会話文の英語原文を取り出す。
// baseGame は供給元が base ゲームかを示し、フルネーム対へ印を付ける（two の base ゲーム限定判定用）。
// encoding/xml はエンティティ（&lt; など）を復号する。先頭 UTF-8 BOM は除いてから解析する。
func ParseTermXML(data []byte, baseGame bool) (fulls, shrts []termderive.NamePair, dialogues []string, err error) { //nolint:gocognit // TODO(refactor): XML トークン走査の状態分岐（要素種別×収集先）。リファクタ本体で分割する。
	data = bytes.TrimPrefix(data, []byte("\uFEFF"))
	dec := xml.NewDecoder(bytes.NewReader(data))
	for {
		tok, tokErr := dec.Token()
		if tokErr == io.EOF {
			break
		}
		if tokErr != nil {
			return nil, nil, nil, tokErr
		}
		start, ok := tok.(xml.StartElement)
		if !ok || start.Name.Local != "String" {
			continue
		}
		var s xmlString
		if decErr := dec.DecodeElement(&s, &start); decErr != nil {
			return nil, nil, nil, decErr
		}
		rec := strings.ToUpper(strings.TrimSpace(s.REC))
		src := strings.TrimSpace(s.Source)
		dst := strings.TrimSpace(s.Dest)
		switch rec {
		case "NPC_:FULL":
			if src != "" && dst != "" {
				fulls = append(fulls, termderive.NamePair{Source: src, Dest: dst, BaseGame: baseGame})
			}
		case "NPC_:SHRT":
			if src != "" && dst != "" {
				shrts = append(shrts, termderive.NamePair{Source: src, Dest: dst, BaseGame: baseGame})
			}
		case "INFO:NAM1":
			if src != "" {
				dialogues = append(dialogues, src)
			}
		}
	}
	return fulls, shrts, dialogues, nil
}
