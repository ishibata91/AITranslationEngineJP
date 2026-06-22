// Command poc-tone は口調抽出 PoC（v2）。db/poc-skyrim.sqlite3 の実台詞から、
// 基底口調（対人態度×感情表出）と分散、語彙マーカーを機械抽出して表示する。
// v2 の変更: 対人態度を命令文（prose の品詞解析）と丁寧定型で測り、集計を割合へ変更。
// 感情を NRC EmoLex（外部辞書）で拡充。計画は poc-tone-extraction.md。
//
// 実行モード:
//
//	go run ./cmd/poc-tone          手で選んだ 13 NPC（枠・破れ疑い）を出す。
//	go run ./cmd/poc-tone verify   台詞 100 以上の全話者へかける無作為検証。選択バイアスを外す。
package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"math"
	"os"
	"regexp"
	"sort"
	"strings"
	"unicode"

	"github.com/jdkato/prose/v2"
	_ "modernc.org/sqlite"
)

const nrcPath = "dictionaries/nrc-emolex.txt"

// dbPathFromEnv は中心 DB のパスを返す。POC_DB で別プラグインの抽出 DB（inigo 等）へ向けられる。
func dbPathFromEnv() string {
	if p := os.Getenv("POC_DB"); p != "" {
		return p
	}
	return "db/poc-skyrim.sqlite3"
}

type target struct{ edid, group string }

// fixedTargets は v2 レポートで使った手選びの 13 NPC。再現用に残す。
var fixedTargets = []target{
	{"Maven", "枠"},
	{"Ulfric", "枠"},
	{"GeneralTullius", "枠"},
	{"Astrid", "枠"},
	{"FarengarSecretFire", "枠"},
	{"AventusAretino", "枠"},
	{"Nazeem", "枠少"},
	{"Grelod", "枠少"},
	{"ConstanceMichel", "枠少"},
	{"Paarthurnax", "破れ疑い"},
	{"Cicero", "破れ疑い"},
	{"DA15Sheogorath", "破れ疑い"},
	{"MaiqTheLiar", "破れ疑い"},
}

// sparseTargets は voice 気質 prior の経路を実証する検証集合。
// 台詞が少なく印が 10 未満で、汎用気質 voice を持つ固有キャラ。先頭に対照の多弁キャラを置く。
var sparseTargets = []target{
	{"Tolfdir", "多弁対照"},
	{"PhinisGestor", "多弁対照"},
	{"Nazeem", "台詞少"},
	{"Grelod", "台詞少"},
	{"ConstanceMichel", "台詞少"},
	{"Edorfin", "1台詞"},
	{"Curwe", "1台詞"},
	{"Vilod", "1台詞"},
	{"Edith", "1台詞"},
	{"Tsavani", "1台詞"},
}

// 丁寧定型（politeness 研究の緩衝・間接表現を参考にした小辞書）。
var politeWords = []string{"please", "sorry", "i'm afraid", "thank", "would you", "could you", "perhaps", "my dear", "must ask", "forgive", "pardon", "if you would", "i think", "i believe"}

// Skyrim 固有の罵倒（NRC で拾えない世界固有語）。
var insultWords = []string{"riff-raff", "milk-drinker", "fool", "wretch", "defiler", "debaser", "n'wah", "s'wit", "over your head"}

var (
	sentSplit = regexp.MustCompile(`[.!?]+`)
	wordRe    = regexp.MustCompile(`[a-z']+`)
)

// NRC EmoLex（外部辞書）。strong は怒り・恐怖・嫌悪・驚き・喜びの強い感情語。
var nrcEmotion = map[string]bool{}

func loadNRC(path string) {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	strong := map[string]bool{"anger": true, "fear": true, "disgust": true, "surprise": true, "joy": true}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		p := strings.Split(sc.Text(), "\t")
		if len(p) != 3 || p[2] != "1" {
			continue
		}
		if strong[p[1]] {
			nrcEmotion[p[0]] = true
		}
	}
	if err := sc.Err(); err != nil {
		panic(err)
	}
}

func countOccur(text string, words []string) int {
	n := 0
	for _, w := range words {
		n += strings.Count(text, w)
	}
	return n
}

func wordTokens(text string) []string {
	return wordRe.FindAllString(strings.ToLower(text), -1)
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

// guidanceLeads は誘導・協調の命令で始まる定型。護衛 NPC の道案内を威圧から外すために使う。
var guidanceLeads = []string{
	"come on", "come, ", "this way", "follow me", "follow ", "stay close",
	"stay with me", "stay behind", "stay back", "keep up", "keep close",
	"keep moving", "keep an eye", "keep your", "keep quiet", "look around",
	"look out", "wait here", "wait for", "get down", "get back", "get inside",
	"hurry", "over here", "after me", "right behind", "watch out",
	"take cover", "take it slow",
}

// isGuidanceImperative は誘導・協調の命令かを判定する。包括命令（let's / let us）と
// 道案内の定型で始まる文を、威圧の命令から外す。Hadvar・Ralof の道案内を尊大と数える誤りへの対応。
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

// Features は 1 台詞の機械検出特徴量。
type Features struct {
	Sentences, Polite, Insult, Imperative, Question, Exclaim, Elong, NRCEmo int
}

func extractFeatures(text string) Features {
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
	q := 0
	if strings.Contains(text, "?") {
		q = 1
	}
	emo := 0
	for _, w := range wordTokens(text) {
		if nrcEmotion[w] {
			emo++
		}
	}
	return Features{
		Sentences:  sents,
		Polite:     countOccur(tl, politeWords),
		Insult:     countOccur(tl, insultWords),
		Imperative: imp,
		Question:   q,
		Exclaim:    strings.Count(text, "!"),
		Elong:      countElong(tl),
		NRCEmo:     emo,
	}
}

// scoreAxes は特徴量を 2 軸スコアへ写す。対人態度は命令文・罵倒で負、丁寧定型で正。
// 疑問文は丁寧とは限らない（詰問・皮肉もある）ため加点しない。
func scoreAxes(f Features) (politeness, arousal float64) {
	politeness = float64(f.Polite) - float64(f.Imperative) - float64(f.Insult)
	arousal = float64(f.Exclaim+f.Elong+f.NRCEmo) / float64(f.Sentences)
	return
}

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

func stddev(xs []float64) float64 {
	if len(xs) < 2 {
		return 0
	}
	var m float64
	for _, x := range xs {
		m += x
	}
	m /= float64(len(xs))
	var v float64
	for _, x := range xs {
		v += (x - m) * (x - m)
	}
	return math.Sqrt(v / float64(len(xs)))
}

// attitudeRatio は対人マーカーを含む台詞だけで、丁寧割合 − 粗野割合（-1..1）を返す。
// 中立台詞（命令も丁寧もない大多数）を分母から除き、台詞数で薄まる問題を避ける。
// marked はマーカーを含む台詞数で、信頼度の目安にする。
func attitudeRatio(ps []float64) (ratio float64, marked int) {
	rude, polite := 0, 0
	for _, p := range ps {
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

func attitudeBand(p float64) int { // 0 尊大 / 1 中立 / 2 丁寧
	switch {
	case p > 0.15:
		return 2
	case p < -0.15:
		return 0
	default:
		return 1
	}
}

func arousalBand(a float64) int { // 0 抑制 / 1 中 / 2 激情
	switch {
	case a > 1.0:
		return 2
	case a > 0.4:
		return 1
	default:
		return 0
	}
}

var cellNames = [3][3]string{
	{"冷然・見下し", "淡々・実務", "慇懃・端正"},
	{"ぞんざい", "平明", "物腰やわ"},
	{"居丈高・罵倒", "率直・興奮", "情に厚い懇願"},
}

func markersFromMeta(voice, race string) []string {
	var m []string
	switch {
	case strings.Contains(voice, "Old"):
		m = append(m, "老人(声型)")
	case strings.Contains(voice, "Child"):
		m = append(m, "子供(声型)")
	case strings.Contains(voice, "Young"):
		m = append(m, "若年(声型)")
	}
	if strings.Contains(voice, "Condescending") {
		m = append(m, "横柄(声型)")
	}
	switch {
	case strings.Contains(race, "Khajiit"):
		m = append(m, "Khajiit訛り")
	case strings.Contains(race, "Argonian"):
		m = append(m, "Argonian訛り")
	case strings.Contains(race, "Dragon"):
		m = append(m, "龍(種族)")
	case strings.Contains(race, "Elder"):
		m = append(m, "老人(種族)")
	case strings.Contains(race, "Child"):
		m = append(m, "子供(種族)")
	}
	return m
}

// lexicalFromText は本文から語彙マーカーを検出する。古英語は単語境界一致で誤爆を防ぐ。
// Khajiit 三人称は種族レコードが Khajiit のときだけ出す。Breton 等の `this one`（指示語）誤爆を防ぐ。
func lexicalFromText(texts []string, race string) []string {
	words := map[string]bool{}
	var sb strings.Builder
	for _, t := range texts {
		sb.WriteByte(' ')
		sb.WriteString(strings.ToLower(t))
		for _, w := range wordTokens(t) {
			words[w] = true
		}
	}
	joined := sb.String()
	var m []string
	for _, a := range []string{"thee", "thou", "thy", "thine", "hast", "doth", "tis"} {
		if words[a] {
			m = append(m, "古英語(本文)")
			break
		}
	}
	if strings.Contains(race, "Khajiit") && strings.Contains(joined, "this one") {
		m = append(m, "Khajiit三人称(本文)")
	}
	return m
}

// voiceTraits は汎用 voice_type の気質辞書。token は voice EditorID の気質部分（Male/Female を除く）。
// prior は対人軸の事前値（負＝尊大、正＝丁寧）。label は表示用の短い気質名。
// 出典は voice_type の EditorID 命名（制作側が付けた気質ラベル）。本文と独立な源。
var voiceTraits = []struct {
	token string
	prior float64
	label string
}{
	{"Condescending", -0.6, "横柄"},
	{"ElfHaughty", -0.5, "高慢"},
	{"Brute", -0.5, "粗暴"},
	{"Bandit", -0.5, "粗暴"},
	{"SlyCynical", -0.4, "皮肉"},
	{"OldGrumpy", -0.3, "気難老"},
	{"Shrill", -0.3, "刺々"},
	{"Commander", -0.1, "指揮"},
	{"Guard", -0.1, "衛兵"},
	{"Coward", 0.4, "臆病"},
	{"OldKindly", 0.4, "温厚老"},
	{"EvenToned", 0.0, "平静"},
	{"Sultry", 0.0, "色気"},
	{"YoungEager", 0.0, "前のめり"},
	{"Drunk", 0.0, "酔"},
	{"Commoner", 0.0, "庶民"},
	{"Soldier", 0.0, "兵"},
	{"Warlock", 0.0, "術士"},
}

func isUniqueVoice(voice string) bool {
	return voice == "" || strings.Contains(voice, "Unique") ||
		strings.HasPrefix(voice, "Cr") || strings.HasPrefix(voice, "SPECIAL")
}

// voicePrior は voice 気質から対人軸の事前値を返す。固有・特殊 voice は prior 無し（ok=false）。
// 汎用だが気質中立の voice（Nord 等の種族系）は中立 0 を返す。
func voicePrior(voice string) (prior float64, ok bool) {
	if isUniqueVoice(voice) {
		return 0, false
	}
	for _, t := range voiceTraits {
		if strings.Contains(voice, t.token) {
			return t.prior, true
		}
	}
	return 0, true
}

func voiceLabel(voice string) string {
	if isUniqueVoice(voice) {
		return "固有"
	}
	for _, t := range voiceTraits {
		if strings.Contains(voice, t.token) {
			return t.label
		}
	}
	return "—"
}

// fuseAttitude は対人軸を融合する。印が十分なら本文、足りなければ voice 気質 prior、
// 固有 voice で prior も無ければ保留（本文の薄い値をそのまま使うが信頼低い）。
func fuseAttitude(dialogueAtt float64, marked int, voice string) (att float64, band int, path string) {
	if marked >= 10 {
		return dialogueAtt, attitudeBand(dialogueAtt), "本文"
	}
	if vp, ok := voicePrior(voice); ok {
		return vp, attitudeBand(vp), "voice"
	}
	return dialogueAtt, attitudeBand(dialogueAtt), "保留"
}

// result は 1 話者の抽出結果。verify モードでソート・集計に使う。
type result struct {
	edid, group        string
	lines, marked      int
	att, emoMed, emoSd float64
	attBand            int
	cell               string
	path, vlabel       string
	markers            []string
}

func processSpeaker(db *sql.DB, tg target) (result, bool) {
	rows, err := db.Query(`
SELECT l.source, COALESCE(v.edid,''), COALESCE(r.edid,'')
FROM line l
JOIN line_speaker ls ON ls.line_id=l.id
JOIN speaker s ON s.id=ls.speaker_id
LEFT JOIN voice_type v ON s.voice_type_id=v.id
LEFT JOIN race r ON s.race_id=r.id
WHERE s.edid=? AND TRIM(l.source)<>'' AND LENGTH(l.source)>15`, tg.edid)
	if err != nil {
		panic(err)
	}
	var voice, race string
	var texts []string
	var ps, as []float64
	for rows.Next() {
		var src, v, r string
		if err := rows.Scan(&src, &v, &r); err != nil {
			panic(err)
		}
		voice, race = v, r
		texts = append(texts, src)
		p, a := scoreAxes(extractFeatures(src))
		ps = append(ps, p)
		as = append(as, a)
	}
	rows.Close()
	if len(ps) == 0 {
		return result{}, false
	}
	attVal, marked := attitudeRatio(ps)
	amed := median(as)
	att, ab, path := fuseAttitude(attVal, marked, voice)
	cell := cellNames[arousalBand(amed)][ab]
	markers := append(markersFromMeta(voice, race), lexicalFromText(texts, race)...)
	return result{
		edid: tg.edid, group: tg.group, lines: len(ps), marked: marked,
		att: att, emoMed: amed, emoSd: stddev(as), attBand: ab,
		cell: cell, path: path, vlabel: voiceLabel(voice), markers: markers,
	}, true
}

// speakerTargets は台詞 minLines 以上の全話者を返す。手選びを排し選択バイアスを外す母集団。
func speakerTargets(db *sql.DB, minLines int) []target {
	rows, err := db.Query(`
SELECT s.edid, count(*) AS c
FROM speaker s
JOIN line_speaker ls ON ls.speaker_id=s.id
JOIN line l ON l.id=ls.line_id
WHERE TRIM(s.edid)<>'' AND TRIM(l.source)<>'' AND LENGTH(l.source)>15
GROUP BY s.id
HAVING c>=?
ORDER BY c DESC`, minLines)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var ts []target
	for rows.Next() {
		var edid string
		var c int
		if err := rows.Scan(&edid, &c); err != nil {
			panic(err)
		}
		ts = append(ts, target{edid, fmt.Sprintf("%d台詞", c)})
	}
	return ts
}

func printTable(results []result) {
	fmt.Printf("%-20s %4s %4s %-6s %-13s %6s %6s  %-8s %s\n", "話者", "台詞", "印", "経路", "基底口調", "対人", "感情", "voice気質", "語彙マーカー")
	fmt.Println(strings.Repeat("-", 104))
	for _, r := range results {
		fmt.Printf("%-20s %4d %4d %-6s %-13s %6.2f %6.2f  %-8s %s\n",
			r.edid, r.lines, r.marked, r.path, r.cell, r.att, r.emoMed, r.vlabel, strings.Join(r.markers, ","))
	}
}

// summarize は無作為検証の所見を出す。対人軸が中立へ潰れず判別したかを確かめる。
func summarize(results []result) {
	reliable := 0
	band := [3]int{} // 0 尊大 / 1 中立 / 2 丁寧
	cellHist := map[string]int{}
	for _, r := range results {
		if r.marked < 10 {
			continue
		}
		reliable++
		band[r.attBand]++
		cellHist[r.cell]++
	}
	fmt.Printf("\n## 無作為検証の所見（台詞100以上の全%d話者）\n", len(results))
	fmt.Printf("印10以上で評価可能: %d 話者\n", reliable)
	fmt.Printf("対人帯の分布(印10以上): 尊大 %d / 中立 %d / 丁寧 %d\n", band[0], band[1], band[2])
	if reliable > 0 {
		fmt.Printf("中立に潰れた割合: %.0f%%（v1 は全員潰れた。低いほど判別している）\n",
			100*float64(band[1])/float64(reliable))
	}
	type kv struct {
		k string
		v int
	}
	var cells []kv
	for k, v := range cellHist {
		cells = append(cells, kv{k, v})
	}
	sort.Slice(cells, func(i, j int) bool { return cells[i].v > cells[j].v })
	fmt.Print("基底口調セルの分布(印10以上): ")
	var parts []string
	for _, c := range cells {
		parts = append(parts, fmt.Sprintf("%s %d", c.k, c.v))
	}
	fmt.Println(strings.Join(parts, " / "))
}

func main() {
	loadNRC(nrcPath)

	db, err := sql.Open("sqlite", dbPathFromEnv())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	mode := ""
	if len(os.Args) > 1 {
		mode = os.Args[1]
	}

	var tgs []target
	sortByAtt := true
	switch mode {
	case "verify":
		tgs = speakerTargets(db, 100)
		fmt.Printf("NRC 強感情語 %d 件を読み込み\n", len(nrcEmotion))
		fmt.Printf("## 無作為検証（台詞100以上の全話者、融合・対人軸でソート）\n\n")
	case "all":
		tgs = speakerTargets(db, 1)
		fmt.Printf("NRC 強感情語 %d 件を読み込み\n", len(nrcEmotion))
		fmt.Printf("## 全話者（台詞1以上、融合・対人軸でソート）\n\n")
	case "sparse":
		tgs = sparseTargets
		sortByAtt = false // 経路の違いを見せるため入力順を保つ
		fmt.Printf("NRC 強感情語 %d 件を読み込み\n", len(nrcEmotion))
		fmt.Printf("## 融合の検証（voice 気質 prior の経路、入力順）\n\n")
	default:
		tgs = fixedTargets
		fmt.Printf("NRC 強感情語 %d 件を読み込み\n\n", len(nrcEmotion))
		fmt.Println("## 基底口調と分散（v2: 命令文＋NRC＋割合集計）")
	}

	var results []result
	for _, tg := range tgs {
		r, ok := processSpeaker(db, tg)
		if ok {
			results = append(results, r)
		}
	}
	// 対人スコア昇順（尊大→丁寧）に並べ、判別の広がりを見やすくする。
	if sortByAtt {
		sort.Slice(results, func(i, j int) bool { return results[i].att < results[j].att })
	}
	printTable(results)
	if mode == "verify" {
		summarize(results)
	}
}
