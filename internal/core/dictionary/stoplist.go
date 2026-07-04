package dictionary

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Stoplist は機械置換辞書と言及語彙の供給源から一般語を除く選別規則。
// 本文の通常語と同綴りの固有名（FULL 名や階級称号が "Yes"・"No" 等の管理用文字列）が
// 辞書へ載ると、文頭の大文字出現に誤って当たり本文を壊すため、供給の段階で除く。
// 語集合は外部配布の stopword リスト（stopwords-iso、assets/stopwords-en.txt）を使い、独自拡張はしない。
type Stoplist struct {
	words map[string]bool
}

// ParseStoplist は 1 行 1 語の stopword リスト（stopwords-iso 配布形式）を読んで選別規則を作る。
// 語は小文字へ正規化し、前後の空白を落とし、空行は読み飛ばす。
func ParseStoplist(r io.Reader) (*Stoplist, error) {
	words := make(map[string]bool)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		word := strings.ToLower(strings.TrimSpace(scanner.Text()))
		if word == "" {
			continue
		}
		words[word] = true
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("stopword リストの読み取り: %w", err)
	}
	return &Stoplist{words: words}, nil
}

// Blocks は原語を機械置換辞書・言及語彙の供給から除くべきかを返す。
// 原語全体を小文字化した一語一致で判定する。複数語の固有名（空白区切りの名前）は
// 語の一部が一般語でも全体は固有の並びなので当てない。nil の Stoplist は選別なし（常に false）。
func (s *Stoplist) Blocks(source string) bool {
	if s == nil {
		return false
	}
	trimmed := strings.TrimSpace(source)
	if trimmed == "" || len(strings.Fields(trimmed)) != 1 {
		return false
	}
	return s.words[strings.ToLower(trimmed)]
}
