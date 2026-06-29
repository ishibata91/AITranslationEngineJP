// Package termusage は会話文から各英単語の用法分布（一般語 LC・固有名 UC）を集計する純粋ルール。
// termderive の安全フィルタが、一般語（Master・Blood など）の派生を捨てる根拠に使う。
package termusage

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"aitranslationenginejp/internal/core/termderive"
)

// 本ファイルは派生規則の安全フィルタが使う用法集計を作る純粋関数を持つ。会話文（INFO:NAM1 の英語原文）から、
// 各英単語が一般語として小文字で出た回数（LC）と、固有名として文頭以外で大文字始まりに出た回数（UC）を数える。
// 用法分布で「一般語か固有名か」を構造判定し、Master や Blood のような一般語の派生を捨てる根拠にする。

var (
	// sentSplitRe は文末記号（. ! ?）と続く空白で文を分ける。文頭語の大文字始まりを固有名と誤認しないため。
	sentSplitRe = regexp.MustCompile(`[.!?]\s+`)
	// wordRe は英字・アポストロフィからなる語を取り出す。
	wordRe = regexp.MustCompile(`[A-Za-z'][A-Za-z']*`)
	// lowerRunRe は語内の連続する小文字英字を取り出す（縮約 aren't を aren・t に割って数える）。
	lowerRunRe = regexp.MustCompile(`[a-z]+`)
)

// BuildUsage は会話文の英語原文から各英単語の用法分布（LC / UC）を作る。
// 文ごとに語を取り、小文字始まりの語は一般語用法として小文字連を LC に積み、
// 文頭以外で大文字始まりの語は固有名用法として小文字化した語を UC に積む。文頭語（i==0）の大文字は数えない。
func BuildUsage(dialogues []string) termderive.Usage {
	u := termderive.Usage{LC: map[string]int{}, UC: map[string]int{}}
	for _, line := range dialogues {
		for _, sent := range sentSplitRe.Split(line, -1) {
			accumulateSentence(sent, u)
		}
	}
	return u
}

// accumulateSentence は 1 文のトークンを走査し、各トークンの用法を分布 u へ積む。
// u は map を持つため値渡しでも積み足しは呼び出し元へ反映される。
func accumulateSentence(sent string, u termderive.Usage) {
	for i, tok := range wordRe.FindAllString(sent, -1) {
		accumulateToken(i, tok, u)
	}
}

// accumulateToken は 1 トークンを大小・文頭で振り分けて用法分布 u へ積む。
// 小文字始まりは一般語用法として小文字連を LC に積み、文頭以外（i>0）の大文字始まりは固有名用法として UC に積む。
// 文頭語（i==0）の大文字は数えない。
func accumulateToken(i int, tok string, u termderive.Usage) {
	first, _ := utf8.DecodeRuneInString(tok)
	if unicode.IsLower(first) {
		for _, w := range lowerRunRe.FindAllString(tok, -1) {
			u.LC[w]++
		}
		return
	}
	if i > 0 {
		u.UC[strings.Trim(strings.ToLower(tok), "'")]++
	}
}
