package engine

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// 本ファイルは xTranslator 英日 XML から人名の部分形（名のみ・短名）を派生する I/O を持つ。
// 純粋ルール（termderive.go の DeriveTerms）へ XML 解析と用法集計を与える橋渡しで、
// 翻訳 Run の固有名派生（engine.DeriveMasterTerms）が使う。

// baseGamePrefixes は base ゲーム（Bethesda 公式本体・DLC）の XML ファイル名接頭。姓名分割（two）を
// base ゲーム限定にするための判定に使う。第三者 mod は Source/Dest 対応が信頼できないため two を派生しない。
var baseGamePrefixes = []string{"Skyrim", "Dawnguard", "Dragonborn", "HearthFires", "Update"}

// isBaseGame は XML ファイル名が base ゲーム由来かを返す。
func isBaseGame(filename string) bool {
	for _, p := range baseGamePrefixes {
		if strings.HasPrefix(filename, p) {
			return true
		}
	}
	return false
}

// DeriveTermsFromXMLDir はディレクトリ内の全 XML を解析し、用法分布を作り、純粋ルールで派生対を返す。
// baseSources は base 辞書の既存原語集合で、既出原語との衝突を避ける除外判定に使う。
func DeriveTermsFromXMLDir(xmlDir string, baseSources map[string]bool) ([]DerivedTerm, error) {
	entries, err := filepath.Glob(filepath.Join(xmlDir, "*.xml"))
	if err != nil {
		return nil, fmt.Errorf("XML の列挙: %w", err)
	}
	if len(entries) == 0 {
		return nil, fmt.Errorf("XML が無い: %s", xmlDir)
	}
	sort.Strings(entries)

	var fulls, shrts []NamePair
	var dialogues []string
	for _, path := range entries {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("%s の読み込み: %w", path, err)
		}
		f, s, d, err := parseTermXML(data, isBaseGame(filepath.Base(path)))
		if err != nil {
			return nil, fmt.Errorf("%s の解析: %w", path, err)
		}
		fulls = append(fulls, f...)
		shrts = append(shrts, s...)
		dialogues = append(dialogues, d...)
	}
	usage := BuildUsage(dialogues)
	return DeriveTerms(fulls, shrts, usage, baseSources, DefaultDeriveConfig()), nil
}

// xmlString は xTranslator XML の 1 レコード（<String>）。REC（種別）と原語・確定訳語を読む。
type xmlString struct {
	REC    string `xml:"REC"`
	Source string `xml:"Source"`
	Dest   string `xml:"Dest"`
}

// parseTermXML は 1 つの XML から NPC のフルネーム対・短名対と会話文の英語原文を取り出す。
// baseGame は供給元が base ゲームかを示し、フルネーム対へ印を付ける（two の base ゲーム限定判定用）。
// encoding/xml はエンティティ（&lt; など）を復号する。先頭 UTF-8 BOM は除いてから解析する。
func parseTermXML(data []byte, baseGame bool) (fulls, shrts []NamePair, dialogues []string, err error) { //nolint:gocognit // TODO(refactor): XML トークン走査の状態分岐（要素種別×収集先）。リファクタ本体で分割する。
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
				fulls = append(fulls, NamePair{Source: src, Dest: dst, BaseGame: baseGame})
			}
		case "NPC_:SHRT":
			if src != "" && dst != "" {
				shrts = append(shrts, NamePair{Source: src, Dest: dst, BaseGame: baseGame})
			}
		case "INFO:NAM1":
			if src != "" {
				dialogues = append(dialogues, src)
			}
		}
	}
	return fulls, shrts, dialogues, nil
}
