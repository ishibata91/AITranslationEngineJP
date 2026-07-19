package lexicon

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// strongThreshold は「感情表出の強さ」の二値ゲートの閾値。
// VADER の valence（快・不快の平均評定、概ね -4..+4）の絶対値がこの値以上の語を強感情語とする。
// 1.5 は angry(-2.3)・fear(-2.2)・terrible(-2.1) を含め、中立寄りの弱い語を落とす初期値。調整可能。
const strongThreshold = 1.5

// VADER は VADER lexicon（MIT ライセンス）由来の強感情語の集合。engine.EmotionLexicon を満たす。
// 感情のカテゴリ内訳は持たず、valence 絶対値の閾値で「強い感情表出の語か」の二値だけを判定する。
type VADER struct{ words map[string]struct{} }

// LoadVADER は VADER lexicon のファイル（タブ区切り: 語・valence 平均・標準偏差・生評定）を読み、
// valence 絶対値が strongThreshold 以上の語を集合にする。ファイルが読めなければエラーを返す。
// 照合を安定させるため語は小文字化して格納する（IsStrongEmotion は小文字化済みの語を前提にする）。
func LoadVADER(path string) (*VADER, error) {
	f, err := os.Open(path) //nolint:gosec // G304: path は感情辞書（VADER lexicon）の参照データのパス。差し替え可能な境界の意図的な読み込みのため限定許可する。
	if err != nil {
		return nil, fmt.Errorf("VADER 辞書を開けない (%s): %w", path, err)
	}
	defer func() { _ = f.Close() }()

	words := make(map[string]struct{})
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		p := strings.Split(sc.Text(), "\t")
		if len(p) < 2 {
			continue
		}
		mean, err := strconv.ParseFloat(p[1], 64)
		if err != nil {
			continue
		}
		if math.Abs(mean) >= strongThreshold {
			words[strings.ToLower(p[0])] = struct{}{}
		}
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("VADER 辞書の読み取り: %w", err)
	}
	return &VADER{words: words}, nil
}

// IsStrongEmotion は小文字化済みの語が強い感情語かを返す。
func (v *VADER) IsStrongEmotion(word string) bool {
	_, ok := v.words[word]
	return ok
}
