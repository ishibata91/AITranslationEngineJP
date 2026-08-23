// Package japanesetext は翻訳対象の原文に含まれる日本語文字を判定する。
package japanesetext

import "unicode"

// Contains はひらがな、カタカナ、漢字のいずれかを含む原文を判定する。
func Contains(source string) bool {
	for _, r := range source {
		if unicode.Is(unicode.Hiragana, r) || unicode.Is(unicode.Katakana, r) || unicode.Is(unicode.Han, r) {
			return true
		}
	}
	return false
}
