// Package lexicon は engine の EmotionLexicon 境界の concrete 実装を持つ。
// NRC EmoLex は研究用ライセンスのため、製品化時に MIT ライセンスの実装へ差し替える（差し替え可能な境界）。
package lexicon

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// strongEmotions は対象とする強い感情（怒り・恐怖・嫌悪・驚き・喜び）。感情表出軸の高さに効く。
var strongEmotions = map[string]bool{"anger": true, "fear": true, "disgust": true, "surprise": true, "joy": true}

// NRC は NRC EmoLex 由来の強感情語の集合。engine.EmotionLexicon を満たす。
type NRC struct{ words map[string]struct{} }

// LoadNRC は NRC EmoLex のファイル（タブ区切り: 語・感情・0|1）を読み、
// 強い感情のいずれかが 1 の語を集合にする。ファイルが読めなければエラーを返す。
func LoadNRC(path string) (*NRC, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("NRC 辞書を開けない (%s): %w", path, err)
	}
	defer f.Close()

	words := make(map[string]struct{})
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		p := strings.Split(sc.Text(), "\t")
		if len(p) != 3 || p[2] != "1" {
			continue
		}
		if strongEmotions[p[1]] {
			words[p[0]] = struct{}{}
		}
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("NRC 辞書の読み取り: %w", err)
	}
	return &NRC{words: words}, nil
}

// IsStrongEmotion は小文字化済みの語が強い感情語かを返す。
func (n *NRC) IsStrongEmotion(word string) bool {
	_, ok := n.words[word]
	return ok
}
