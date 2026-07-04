// prose.go は TextAnalyzer の prose 実装。外部学習モデルによる固有表現抽出（NER）と
// 品詞解析で、候補検出（candidate.go）の判定を補強する。学習モデルの推論は同一入力に
// 同一出力を返す（決定的）。学習データ（ニュース英語）由来のため fantasy 語への判定品質は
// 保証せず、単独の根拠でなく補強信号としてだけ使う。

package mention

import (
	"unicode"

	"github.com/jdkato/prose/v2"
)

// ProseAnalyzer は prose の固有表現抽出と品詞解析で本文を解析する。
type ProseAnalyzer struct{}

// Entities は本文中の固有表現の表記列を返す。解析に失敗した本文は補強なし（nil）として扱う。
func (ProseAnalyzer) Entities(text string) []string {
	doc, err := prose.NewDocument(text)
	if err != nil {
		return nil
	}
	ents := doc.Entities()
	out := make([]string, 0, len(ents))
	for _, e := range ents {
		out = append(out, e.Text)
	}
	return out
}

// LeadingVerb は文頭の意味トークンが動詞原形（VB）かを返す。命令形のクエスト目標行
// （Kill Vittoria Vici）の識別に使う。判定方法は口調分類の命令文判定
// （internal/core/linefeatures の isImperative）と同じ。
func (ProseAnalyzer) LeadingVerb(sentence string) bool {
	doc, err := prose.NewDocument(sentence,
		prose.WithExtraction(false), prose.WithSegmentation(false))
	if err != nil {
		return false
	}
	for _, t := range doc.Tokens() {
		if len(t.Tag) == 0 || !unicode.IsLetter(rune(t.Tag[0])) {
			continue // 記号トークンを飛ばす
		}
		return t.Tag == "VB"
	}
	return false
}
