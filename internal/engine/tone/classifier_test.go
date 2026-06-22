package tone

import "testing"

// scoreAxes は特徴量を 2 軸スコアへ写す純粋関数。対人は丁寧−命令−罵倒、感情は密度。
func TestScoreAxes(t *testing.T) {
	cases := []struct {
		name           string
		in             Features
		wantPoliteness float64
		wantArousal    float64
	}{
		{
			// 丁寧定型 2・命令 0・罵倒 0 で対人は正、感嘆 1 を 2 文で割り感情は 0.5。
			name:           "丁寧定型は対人を正へ",
			in:             Features{Sentences: 2, Polite: 2, Exclaim: 1},
			wantPoliteness: 2,
			wantArousal:    0.5,
		},
		{
			// 威圧の命令 1・罵倒 1 で対人は負。
			name:           "命令と罵倒は対人を負へ",
			in:             Features{Sentences: 1, Imperative: 1, Insult: 1},
			wantPoliteness: -2,
			wantArousal:    0,
		},
		{
			// 文数 0 は 0 除算を避けて 1 文として割る（境界）。
			name:           "文数0は1文として割る",
			in:             Features{Sentences: 0, Exclaim: 2, Elong: 1, Emotion: 1},
			wantPoliteness: 0,
			wantArousal:    4,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p, a := scoreAxes(tc.in)
			if p != tc.wantPoliteness || a != tc.wantArousal {
				t.Fatalf("scoreAxes(%+v) = (%v, %v), want (%v, %v)", tc.in, p, a, tc.wantPoliteness, tc.wantArousal)
			}
		})
	}
}

// attitudeRatio は対人マーカーを含む台詞だけで丁寧割合−粗野割合を返し、印を数える。
func TestAttitudeRatio(t *testing.T) {
	cases := []struct {
		name       string
		in         []float64
		wantRatio  float64
		wantMarked int
	}{
		// 全て粗野（負）→ -1、印は件数。
		{"全て粗野は-1", []float64{-1, -2, -1}, -1, 3},
		// 全て丁寧（正）→ +1。
		{"全て丁寧は+1", []float64{1, 2}, 1, 2},
		// 丁寧3・粗野1 で (3-1)/4。中立(0)は印に数えない。
		{"中立は分母から除く", []float64{1, 1, 1, -1, 0, 0}, 0.5, 4},
		// マーカー無し（全て中立）→ 0、印 0。
		{"マーカー無しは印0", []float64{0, 0}, 0, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r, m := attitudeRatio(tc.in)
			if r != tc.wantRatio || m != tc.wantMarked {
				t.Fatalf("attitudeRatio(%v) = (%v, %d), want (%v, %d)", tc.in, r, m, tc.wantRatio, tc.wantMarked)
			}
		})
	}
}

// median は中央値を返す。空・奇数・偶数の各経路を固定する。
func TestMedian(t *testing.T) {
	cases := []struct {
		name string
		in   []float64
		want float64
	}{
		{"空は0", nil, 0},
		{"奇数は中央", []float64{3, 1, 2}, 2},
		{"偶数は中央2件の平均", []float64{4, 1, 3, 2}, 2.5},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := median(tc.in); got != tc.want {
				t.Fatalf("median(%v) = %v, want %v", tc.in, got, tc.want)
			}
		})
	}
}

// attitudeBand は対人スコアを段階へ離散化する。閾値の内外差を固定する。
func TestAttitudeBand(t *testing.T) {
	c := NewClassifier()
	cases := []struct {
		name string
		in   float64
		want int
	}{
		{"閾値超で丁寧", 0.16, AttitudePolite},
		{"丁寧閾値ちょうどは中立", 0.15, AttitudeNeutral},
		{"閾値未満で尊大", -0.16, AttitudeArrogant},
		{"尊大閾値ちょうどは中立", -0.15, AttitudeNeutral},
		{"中央は中立", 0, AttitudeNeutral},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := c.attitudeBand(tc.in); got != tc.want {
				t.Fatalf("attitudeBand(%v) = %d, want %d", tc.in, got, tc.want)
			}
		})
	}
}

// arousalBand は感情スコアを段階へ離散化する。閾値の内外差を固定する。
func TestArousalBand(t *testing.T) {
	c := NewClassifier()
	cases := []struct {
		name string
		in   float64
		want int
	}{
		{"激情閾値超で激情", 1.01, EmotionIntense},
		{"激情閾値ちょうどは中", 1.0, EmotionMid},
		{"中閾値超で中", 0.41, EmotionMid},
		{"中閾値ちょうどは抑制", 0.4, EmotionSuppressed},
		{"低密度は抑制", 0, EmotionSuppressed},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := c.arousalBand(tc.in); got != tc.want {
				t.Fatalf("arousalBand(%v) = %d, want %d", tc.in, got, tc.want)
			}
		})
	}
}

// isUniqueVoice は固有・特殊 voice を判定する。気質 prior を持たない voice を分ける。
func TestIsUniqueVoice(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want bool
	}{
		{"空は固有扱い", "", true},
		{"Unique を含むは固有", "MaleUniqueInigo", true},
		{"Cr 接頭は creature", "CrHorse", true},
		{"SPECIAL 接頭は特殊", "SPECIALMournfulNord", true},
		{"汎用気質 voice は固有でない", "MaleNord", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := isUniqueVoice(tc.in); got != tc.want {
				t.Fatalf("isUniqueVoice(%q) = %v, want %v", tc.in, got, tc.want)
			}
		})
	}
}

// voicePrior は voice 気質から対人 prior を返す。固有は無し、汎用気質は値、汎用中立は 0。
func TestVoicePrior(t *testing.T) {
	c := NewClassifier()
	cases := []struct {
		name      string
		in        string
		wantPrior float64
		wantOK    bool
	}{
		{"固有 voice は prior 無し", "MaleUniqueUlfric", 0, false},
		{"横柄 voice は負の prior", "MaleCondescending", -0.6, true},
		{"温厚老 voice は正の prior", "MaleOldKindly", 0.4, true},
		{"気質中立の汎用 voice は0かつ有り", "FemaleNord", 0, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p, ok := c.voicePrior(tc.in)
			if p != tc.wantPrior || ok != tc.wantOK {
				t.Fatalf("voicePrior(%q) = (%v, %v), want (%v, %v)", tc.in, p, ok, tc.wantPrior, tc.wantOK)
			}
		})
	}
}

// voiceLabel は voice 気質の表示名を返す。固有・汎用気質・汎用中立を分ける。
func TestVoiceLabel(t *testing.T) {
	c := NewClassifier()
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"固有 voice は固有", "FemaleUniqueAstrid", "固有"},
		{"横柄 voice は横柄", "MaleCondescending", "横柄"},
		{"気質中立の汎用 voice は—", "MaleNord", "—"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := c.voiceLabel(tc.in); got != tc.want {
				t.Fatalf("voiceLabel(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

// CellName は段階からセル名を引く。範囲内と範囲外を固定する。
func TestCellName(t *testing.T) {
	cases := []struct {
		name     string
		att, emo int
		want     string
	}{
		{"尊大×抑制", AttitudeArrogant, EmotionSuppressed, "冷然・見下し"},
		{"丁寧×中", AttitudePolite, EmotionMid, "物腰やわ"},
		{"中立×激情", AttitudeNeutral, EmotionIntense, "率直・興奮"},
		{"感情段階が範囲外は空", AttitudeNeutral, 3, ""},
		{"対人段階が範囲外は空", 3, EmotionMid, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := CellName(tc.att, tc.emo); got != tc.want {
				t.Fatalf("CellName(%d,%d) = %q, want %q", tc.att, tc.emo, got, tc.want)
			}
		})
	}
}

// rudeLine は威圧の命令 1 を持つ粗野な台詞特徴。
func rudeLine() Features { return Features{Sentences: 1, Imperative: 1} }

// politeLine は丁寧定型 1 を持つ丁寧な台詞特徴。
func politeLine() Features { return Features{Sentences: 1, Polite: 1} }

// repeat は同じ特徴を n 件並べた話者入力を作る。
func repeat(f Features, n int) []Features {
	out := make([]Features, n)
	for i := range out {
		out[i] = f
	}
	return out
}

// Classify は本文経路で、印が十分な尊大話者を尊大・本文経路へ分類する。
func TestClassifyDialoguePathArrogant(t *testing.T) {
	c := NewClassifier()
	// 粗野な命令を 12 件（印 12 ≥ 10）持つ汎用 voice 話者。本文 2 軸で対人を決める。
	got := c.Classify(repeat(rudeLine(), 12), "MaleNord")
	if got.DecisionPath != PathDialogue {
		t.Fatalf("DecisionPath = %q, want %q", got.DecisionPath, PathDialogue)
	}
	if got.AttitudeBand != AttitudeArrogant {
		t.Fatalf("AttitudeBand = %d, want %d(尊大)", got.AttitudeBand, AttitudeArrogant)
	}
	if got.Marked != 12 {
		t.Fatalf("Marked = %d, want 12", got.Marked)
	}
	if got.AttitudeScore != -1 {
		t.Fatalf("AttitudeScore = %v, want -1", got.AttitudeScore)
	}
	// 感情 0（抑制）× 尊大 のセル。感情入力が無いため中帯の「ぞんざい」ではない。
	if got.Cell != "冷然・見下し" {
		t.Fatalf("Cell = %q, want %q", got.Cell, "冷然・見下し")
	}
}

// Classify は本文経路で、印が十分な丁寧話者を丁寧・本文経路へ分類する。
func TestClassifyDialoguePathPolite(t *testing.T) {
	c := NewClassifier()
	got := c.Classify(repeat(politeLine(), 11), "FemaleNord")
	if got.DecisionPath != PathDialogue {
		t.Fatalf("DecisionPath = %q, want %q", got.DecisionPath, PathDialogue)
	}
	if got.AttitudeBand != AttitudePolite {
		t.Fatalf("AttitudeBand = %d, want %d(丁寧)", got.AttitudeBand, AttitudePolite)
	}
}

// Classify は印不足の話者を、voice 気質 prior で対人を決める（Nazeem 型・横柄）。
func TestClassifyVoicePath(t *testing.T) {
	c := NewClassifier()
	// 印不足（マーカー 2 < 10）で横柄 voice を持つ固有キャラ相当。voice 経路で尊大。
	got := c.Classify(repeat(rudeLine(), 2), "MaleCondescending")
	if got.DecisionPath != PathVoice {
		t.Fatalf("DecisionPath = %q, want %q", got.DecisionPath, PathVoice)
	}
	if got.AttitudeBand != AttitudeArrogant {
		t.Fatalf("AttitudeBand = %d, want %d(尊大)", got.AttitudeBand, AttitudeArrogant)
	}
	if got.AttitudeScore != -0.6 {
		t.Fatalf("AttitudeScore = %v, want -0.6", got.AttitudeScore)
	}
	if got.VoiceLabel != "横柄" {
		t.Fatalf("VoiceLabel = %q, want %q", got.VoiceLabel, "横柄")
	}
}

// Classify は印不足かつ固有 voice で prior も無い話者を保留にし、薄い本文値を保持する。
func TestClassifyHoldPath(t *testing.T) {
	c := NewClassifier()
	// 1 台詞・粗野（印 1）で固有 voice（prior 無し）。保留経路で本文の薄い値を保持する。
	got := c.Classify([]Features{rudeLine()}, "MaleUniqueEdorfin")
	if got.DecisionPath != PathHold {
		t.Fatalf("DecisionPath = %q, want %q", got.DecisionPath, PathHold)
	}
	if got.AttitudeScore != -1 {
		t.Fatalf("AttitudeScore = %v, want -1（本文の薄い値を保持）", got.AttitudeScore)
	}
	if got.VoiceLabel != "固有" {
		t.Fatalf("VoiceLabel = %q, want %q", got.VoiceLabel, "固有")
	}
}

// Classify は感情軸を、印に依らず本文の感情スコア中央値から決める（激情）。
func TestClassifyEmotionFromDialogue(t *testing.T) {
	c := NewClassifier()
	// 感嘆 2 を 1 文で割り感情 2.0（>1.0）。印不足でも感情は本文から測れる。
	loud := Features{Sentences: 1, Exclaim: 2}
	got := c.Classify(repeat(loud, 3), "FemaleNord")
	if got.EmotionBand != EmotionIntense {
		t.Fatalf("EmotionBand = %d, want %d(激情)", got.EmotionBand, EmotionIntense)
	}
	if got.EmotionScore != 2 {
		t.Fatalf("EmotionScore = %v, want 2", got.EmotionScore)
	}
}

// Classify は台詞が無い話者でも、voice 中立で破綻せず中立・voice 経路を返す。
func TestClassifyEmptyLines(t *testing.T) {
	c := NewClassifier()
	got := c.Classify(nil, "FemaleNord")
	if got.Marked != 0 {
		t.Fatalf("Marked = %d, want 0", got.Marked)
	}
	if got.DecisionPath != PathVoice {
		t.Fatalf("DecisionPath = %q, want %q", got.DecisionPath, PathVoice)
	}
	if got.AttitudeBand != AttitudeNeutral || got.EmotionBand != EmotionSuppressed {
		t.Fatalf("bands = (%d,%d), want (中立,抑制)", got.AttitudeBand, got.EmotionBand)
	}
	if got.Cell != "淡々・実務" {
		t.Fatalf("Cell = %q, want %q", got.Cell, "淡々・実務")
	}
}
