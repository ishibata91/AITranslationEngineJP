package engine

import (
	"crypto/sha256"
	"encoding/hex"
	"regexp"
	"strings"
	"unicode"

	"aitranslationenginejp/internal/engine/tone"

	"github.com/jdkato/prose/v2"
)

// EmotionLexicon は英語の強感情語を判定する辞書の境界。行特徴抽出の感情語カウントに使う。
// dev は NRC 実装（研究用ライセンス）、製品化時に MIT ライセンスの実装へ差し替える差し替え可能な境界。
type EmotionLexicon interface {
	// IsStrongEmotion は小文字化済みの語が強い感情語かを返す。
	IsStrongEmotion(word string) bool
}

// 丁寧定型（politeness 研究の緩衝・間接表現を参考にした小辞書）。
var politeWords = []string{"please", "sorry", "i'm afraid", "thank", "would you", "could you", "perhaps", "my dear", "must ask", "forgive", "pardon", "if you would", "i think", "i believe"}

// Skyrim 固有の罵倒（感情辞書で拾えない世界固有語）。
var insultWords = []string{"riff-raff", "milk-drinker", "fool", "wretch", "defiler", "debaser", "n'wah", "s'wit", "over your head"}

// guidanceLeads は誘導・協調の命令で始まる定型。護衛 NPC の道案内を威圧から外すために使う。
var guidanceLeads = []string{
	"come on", "come, ", "this way", "follow me", "follow ", "stay close",
	"stay with me", "stay behind", "stay back", "keep up", "keep close",
	"keep moving", "keep an eye", "keep your", "keep quiet", "look around",
	"look out", "wait here", "wait for", "get down", "get back", "get inside",
	"hurry", "over here", "after me", "right behind", "watch out",
	"take cover", "take it slow",
}

var (
	sentSplit  = regexp.MustCompile(`[.!?]+`)
	lineWordRe = regexp.MustCompile(`[a-z']+`)
)

// SourceHash は本文テキストの内容アドレス（line_analysis のキー）。同一英文は同一ハッシュに集約する。
func SourceHash(text string) string {
	sum := sha256.Sum256([]byte(text))
	return hex.EncodeToString(sum[:])
}

// ExtractFeatures は 1 本文から機械検出特徴量（tone.Features）を抽出する。最も重い処理で、
// line_analysis に本文ハッシュでキャッシュする（純粋 IO の Classifier の外）。prose の品詞解析と
// 差し替え可能な感情辞書（lex）に依存する。文数 0 は 1 文として扱う側（Classifier）に委ねる。
func ExtractFeatures(text string, lex EmotionLexicon) tone.Features {
	tl := strings.ToLower(text)
	sents := 0
	for _, s := range sentSplit.Split(tl, -1) {
		if strings.TrimSpace(s) != "" {
			sents++
		}
	}
	if sents == 0 {
		sents = 1
	}
	imp := 0
	if isImperative(text) && !isGuidanceImperative(text) {
		imp = 1
	}
	emo := 0
	for _, w := range wordTokens(text) {
		if lex.IsStrongEmotion(w) {
			emo++
		}
	}
	return tone.Features{
		Sentences:  sents,
		Polite:     countOccur(tl, politeWords),
		Insult:     countOccur(tl, insultWords),
		Imperative: imp,
		Exclaim:    strings.Count(text, "!"),
		Elong:      countElong(tl),
		Emotion:    emo,
	}
}

// countOccur は語リストの各語の出現数の合計を返す。
func countOccur(text string, words []string) int {
	n := 0
	for _, w := range words {
		n += strings.Count(text, w)
	}
	return n
}

// wordTokens は小文字化した英単語トークン列を返す。
func wordTokens(text string) []string {
	return lineWordRe.FindAllString(strings.ToLower(text), -1)
}

// countElong は同字 3 連続以上（引き伸ばし）の数を数える。
func countElong(text string) int {
	n := 0
	runes := []rune(strings.ToLower(text))
	for i := 0; i < len(runes); {
		j := i
		for j+1 < len(runes) && runes[j+1] == runes[i] {
			j++
		}
		if j-i+1 >= 3 && unicode.IsLetter(runes[i]) {
			n++
		}
		i = j + 1
	}
	return n
}

// isGuidanceImperative は誘導・協調の命令かを判定する。包括命令（let's / let us）と
// 道案内の定型で始まる文を威圧の命令から外す（護衛 NPC の道案内を尊大と数える誤りへの対応）。
func isGuidanceImperative(text string) bool {
	t := strings.TrimLeft(strings.ToLower(strings.TrimSpace(text)), `"'-– `)
	if strings.HasPrefix(t, "let's ") || strings.HasPrefix(t, "let us ") {
		return true
	}
	for _, g := range guidanceLeads {
		if strings.HasPrefix(t, g) {
			return true
		}
	}
	return false
}

// isImperative は文頭の意味トークンが動詞原形(VB)なら命令文とみなす。
// prose の品詞解析（Penn Treebank の外部学習モデル）を参照する文構造判定。
func isImperative(text string) bool {
	doc, err := prose.NewDocument(text,
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
