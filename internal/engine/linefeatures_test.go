package engine

import "testing"

// fakeLexicon は語の集合だけを持つ感情辞書。差し替え可能な境界（EmotionLexicon）の代替。
type fakeLexicon struct{ words map[string]bool }

func (f fakeLexicon) IsStrongEmotion(w string) bool { return f.words[w] }

// SourceHash は同一本文で同一ハッシュ、別本文で別ハッシュを返す（line_analysis のキー）。
func TestSourceHash(t *testing.T) {
	a := SourceHash("Get out of my way.")
	b := SourceHash("Get out of my way.")
	c := SourceHash("Please come in.")
	if a == "" {
		t.Fatal("SourceHash が空")
	}
	if a != b {
		t.Fatalf("同一本文で別ハッシュ: %q vs %q", a, b)
	}
	if a == c {
		t.Fatal("別本文で同一ハッシュ")
	}
}

// ExtractFeatures は丁寧定型・罵倒・感嘆・引き伸ばし・感情語を数える（命令文に依らない経路）。
// 文頭が代名詞のため命令文判定は false で固定でき、prose の曖昧さを避ける。
func TestExtractFeaturesCounts(t *testing.T) {
	lex := fakeLexicon{words: map[string]bool{"afraid": true}}
	// 「my dear」「i'm afraid」が丁寧定型 2、「fool」が罵倒 1、感嘆 1、「noooo」が引き伸ばし 1、「afraid」が感情語 1。
	got := ExtractFeatures("My dear, I'm afraid. Noooo! You fool.", lex)
	if got.Sentences != 3 {
		t.Errorf("Sentences = %d, want 3", got.Sentences)
	}
	if got.Polite != 2 {
		t.Errorf("Polite = %d, want 2", got.Polite)
	}
	if got.Insult != 1 {
		t.Errorf("Insult = %d, want 1", got.Insult)
	}
	if got.Exclaim != 1 {
		t.Errorf("Exclaim = %d, want 1", got.Exclaim)
	}
	if got.Elong != 1 {
		t.Errorf("Elong = %d, want 1", got.Elong)
	}
	if got.Emotion != 1 {
		t.Errorf("Emotion = %d, want 1", got.Emotion)
	}
	if got.Imperative != 0 {
		t.Errorf("Imperative = %d, want 0（文頭が代名詞）", got.Imperative)
	}
}

// ExtractFeatures は威圧の命令を 1、誘導の命令（道案内）を 0 に分ける（修正1の保証）。
func TestExtractFeaturesImperativeVsGuidance(t *testing.T) {
	lex := fakeLexicon{}
	if got := ExtractFeatures("Get out of my way.", lex); got.Imperative != 1 {
		t.Errorf("威圧の命令 Imperative = %d, want 1", got.Imperative)
	}
	if got := ExtractFeatures("Come on, this way.", lex); got.Imperative != 0 {
		t.Errorf("誘導の命令 Imperative = %d, want 0（道案内は威圧から除外）", got.Imperative)
	}
	if got := ExtractFeatures("I am right here.", lex); got.Imperative != 0 {
		t.Errorf("平叙文 Imperative = %d, want 0", got.Imperative)
	}
}
