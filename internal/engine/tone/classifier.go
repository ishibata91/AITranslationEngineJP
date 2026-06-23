// Package tone は口調分類の不変ルールを閉じた純粋 IO クラスを提供する。
// 入力はキャッシュ済みの行特徴（line_analysis 由来の Features 列）と voice 気質ラベル、
// 出力は基底口調（対人段階・感情段階・セル名）と品質指標（印・決定経路）。
// DB 読み書き・prose の品詞解析・ハッシュ計算には依存しない（境界は persona-design.md）。
// 採点・印集計・帯分け・voice 気質 prior 融合の規則はすべて本パッケージへ閉じ、単体テストで固定する（R10）。
package tone

import "sort"

// Features は 1 台詞の機械検出特徴量。prose の品詞解析・NRC 辞書・正規表現で抽出した結果で、
// line_analysis に本文ハッシュでキャッシュされる。Classifier はこのキャッシュ済み特徴だけを
// 入力に取り、prose に依存せず純粋に分類する。
type Features struct {
	Sentences  int // 文数。0 のときは 1 文として扱う
	Polite     int // 丁寧定型の出現数
	Insult     int // 罵倒語の出現数
	Imperative int // 威圧の命令文か（0 か 1。誘導命令は line_analysis 段で除外済み）
	Exclaim    int // 感嘆符の数
	Elong      int // 引き伸ばし（同字 3 連以上）の数
	Emotion    int // 強感情語の数（差し替え可能な感情辞書で数える）
}

// 対人段階。基底口調の対人態度軸の値。
const (
	AttitudeArrogant = 0 // 尊大
	AttitudeNeutral  = 1 // 中立
	AttitudePolite   = 2 // 丁寧
)

// 感情段階。基底口調の感情表出軸の値。
const (
	EmotionSuppressed = 0 // 抑制
	EmotionMid        = 1 // 中
	EmotionIntense    = 2 // 激情
)

// 決定経路。対人段階をどの信号で決めたかを表す（persona-signal-map.md の縮退合成）。
const (
	PathDialogue = "本文"    // 印が閾値以上。本文 2 軸で対人を決めた
	PathVoice    = "voice" // 印不足。voice 気質 prior で対人を決めた
	PathHold     = "保留"    // 印不足かつ固有 voice で prior も無い。薄い本文値は当てにならず対人を中立へ寄せる
)

// Persona は Classifier の出力。基底口調（対人段階・感情段階・セル名）と、品質指標（印・決定経路・voice 気質名）。
type Persona struct {
	AttitudeBand  int     // 対人段階 0尊大/1中立/2丁寧
	EmotionBand   int     // 感情段階 0抑制/1中/2激情
	AttitudeScore float64 // 対人スコア（負＝尊大、正＝丁寧）
	EmotionScore  float64 // 感情スコア（高いほど激情）
	Marked        int     // 印。対人マーカーを含む台詞数で、信頼度の目安
	DecisionPath  string  // 本文/voice/保留
	VoiceLabel    string  // voice 気質の表示名（横柄・温厚老 等。固有 voice は「固有」、気質無しは「—」）
	Cell          string  // 基底口調セル名（感情段階×対人段階）
}

// Classifier は口調分類の不変ルールを閉じた純粋 IO クラス（R10）。
// しきい値と辞書を持ち、較正で差し替える。DB・prose・ハッシュには依存しない。
type Classifier struct {
	markedThreshold int          // 本文経路に要る印の下限。未満は voice prior へ畳む
	attitudePolite  float64      // 丁寧帯の下限スコア（これを超えたら丁寧）
	attitudeArrog   float64      // 尊大帯の上限スコア（これを下回ったら尊大）
	arousalIntense  float64      // 激情帯の下限スコア
	arousalMid      float64      // 中帯の下限スコア
	voiceTraits     []voiceTrait // voice 気質 prior の辞書
	cellNames       [3][3]string // 感情段階×対人段階 → 基底口調セル名
}

// defaultCellNames は 感情段階×対人段階 → 基底口調セル名。行が感情段階、列が対人段階。
var defaultCellNames = [3][3]string{
	{"冷然・見下し", "淡々・実務", "慇懃・端正"},  // 抑制 × 尊大/中立/丁寧
	{"ぞんざい", "平明", "物腰やわ"},        // 中 × 尊大/中立/丁寧
	{"居丈高・罵倒", "率直・興奮", "情に厚い懇願"}, // 激情 × 尊大/中立/丁寧
}

// CellName は基底口調の段階（対人段階・感情段階）からセル名を返す。段階が範囲外なら空文字。
// 生成済みの persona（段階だけを持つ）からセル名を引き直す用途に使う。
func CellName(attitudeBand, emotionBand int) string {
	if emotionBand < 0 || emotionBand >= len(defaultCellNames) {
		return ""
	}
	row := defaultCellNames[emotionBand]
	if attitudeBand < 0 || attitudeBand >= len(row) {
		return ""
	}
	return row[attitudeBand]
}

// NewClassifier は PoC（poc-tone-report.md）で較正したしきい値と気質辞書で Classifier を作る。
func NewClassifier() *Classifier {
	return &Classifier{
		markedThreshold: 10,
		attitudePolite:  0.15,
		attitudeArrog:   -0.15,
		arousalIntense:  1.0,
		arousalMid:      0.4,
		voiceTraits:     defaultVoiceTraits,
		cellNames:       defaultCellNames,
	}
}

// Classify は話者の行特徴群と voice EditorID から基底口調を出す。
// 対人軸は印が十分なら本文、足りなければ voice 気質 prior、固有 voice で prior も無ければ保留にする。
// 感情軸は印に依らず常に本文の感情スコアの中央値から決める。
func (c *Classifier) Classify(lines []Features, voice string) Persona {
	politeness := make([]float64, 0, len(lines))
	arousal := make([]float64, 0, len(lines))
	for _, f := range lines {
		p, a := scoreAxes(f)
		politeness = append(politeness, p)
		arousal = append(arousal, a)
	}
	dialogueAtt, marked := attitudeRatio(politeness)
	emo := median(arousal)
	att, attBand, path := c.fuseAttitude(dialogueAtt, marked, voice)
	emoBand := c.arousalBand(emo)
	return Persona{
		AttitudeBand:  attBand,
		EmotionBand:   emoBand,
		AttitudeScore: att,
		EmotionScore:  emo,
		Marked:        marked,
		DecisionPath:  path,
		VoiceLabel:    c.voiceLabel(voice),
		Cell:          c.cellNames[emoBand][attBand],
	}
}

// scoreAxes は 1 台詞の特徴量を 2 軸スコアへ写す。対人態度は命令文・罵倒で負、丁寧定型で正。
// 感情表出は感嘆符・引き伸ばし・強感情語の密度。文数 0 は 1 文として割る。
func scoreAxes(f Features) (politeness, arousal float64) {
	politeness = float64(f.Polite - f.Imperative - f.Insult)
	sents := f.Sentences
	if sents <= 0 {
		sents = 1
	}
	arousal = float64(f.Exclaim+f.Elong+f.Emotion) / float64(sents)
	return
}

// attitudeRatio は対人マーカーを含む台詞だけで、丁寧割合 − 粗野割合（-1..1）を返す。
// 中立台詞（命令も丁寧もない大多数）を分母から除き、台詞数で薄まる問題を避ける。
// marked はマーカーを含む台詞数で、信頼度の目安にする。
func attitudeRatio(politeness []float64) (ratio float64, marked int) {
	rude, polite := 0, 0
	for _, p := range politeness {
		switch {
		case p < 0:
			rude++
		case p > 0:
			polite++
		}
	}
	marked = rude + polite
	if marked == 0 {
		return 0, 0
	}
	return float64(polite-rude) / float64(marked), marked
}

// median は中央値を返す。空列は 0 を返す。元の列は変更しない。
func median(xs []float64) float64 {
	if len(xs) == 0 {
		return 0
	}
	s := append([]float64(nil), xs...)
	sort.Float64s(s)
	n := len(s)
	if n%2 == 1 {
		return s[n/2]
	}
	return (s[n/2-1] + s[n/2]) / 2
}

// attitudeBand は対人スコアを段階へ離散化する。
func (c *Classifier) attitudeBand(p float64) int {
	switch {
	case p > c.attitudePolite:
		return AttitudePolite
	case p < c.attitudeArrog:
		return AttitudeArrogant
	default:
		return AttitudeNeutral
	}
}

// arousalBand は感情スコアを段階へ離散化する。
func (c *Classifier) arousalBand(a float64) int {
	switch {
	case a > c.arousalIntense:
		return EmotionIntense
	case a > c.arousalMid:
		return EmotionMid
	default:
		return EmotionSuppressed
	}
}

// fuseAttitude は対人軸を融合する。印が十分なら本文、足りなければ voice 気質 prior、
// 固有 voice で prior も無ければ保留にする。保留では薄い本文値が実像とずれやすいため、
// 対人を中立（スコア 0）へ寄せる。感情軸は呼び出し側で常に本文から測るため影響しない。
func (c *Classifier) fuseAttitude(dialogueAtt float64, marked int, voice string) (att float64, band int, path string) {
	if marked >= c.markedThreshold {
		return dialogueAtt, c.attitudeBand(dialogueAtt), PathDialogue
	}
	if prior, ok := c.voicePrior(voice); ok {
		return prior, c.attitudeBand(prior), PathVoice
	}
	return 0, AttitudeNeutral, PathHold
}
