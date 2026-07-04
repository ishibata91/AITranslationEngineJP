// candidate.go は本文コーパスから辞書に無い固有名候補を検出する純粋ルール（言及検出の未知語側）。
// 既知語の言及検出（mention.go）が辞書に載った語を本文から拾うのに対し、こちらは辞書
// （master_term ∪ proper_noun の原語）に載っていない固有名を拾う。known-issues 1番
// （名前付きレコードに出ず本文・会話文にだけ現れる語）の候補検出を担う。
//
// 方式は決定的なヒューリスティックの組み合わせで、LLM を使わない（達成基準の決定性のため）。
//  1. 文頭以外で大文字始まりの語を固有名用法とみなし、隣接する大文字語・接続語（of / the / &）・
//     小文字接頭の姓（gro-Nolob 等）を 1 句に結合する。
//  2. 文頭など曖昧位置の句は、先頭語を含む解釈と外した解釈の両方を候補にする。閉クラスの
//     機能語（A・An・And 等）で始まる句と、品詞解析で動詞先頭と判定した 1 行文（命令形の
//     クエスト目標行）は、先頭語を外した解釈だけを候補にする。
//  3. 句の全体と、接続語で分けた成分、既知語を彫り出した残り、所有格の前後、小文字用法の
//     無い内側の語を候補に出す。辞書は最長一致を好み全体と部分の両方を語彙に持つため。
//  4. 一般語 stoplist（機械置換辞書の供給選別と同じ規則）に載る 1 語は捨てる。曖昧位置だけに
//     出た 1 語は、小文字用法が無い、または固有表現認識が entity と判定した場合だけ通す。
//
// 出力は正規化後表記で一意（重複なし）・辞書順で、同一入力に同一出力を返す。副作用を持たない。
// 評価結果と既知の限界は docs/exec-plans の dictionary-missing-term-detection を参照。

package mention

import (
	"regexp"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

// CandidateTerm は未知固有名候補 1 件。Term は正規化済み表記で、検出結果内で一意。
// 付随する回数は呼び出し側の診断（候補の確からしさの目視確認）用の根拠。
type CandidateTerm struct {
	Term        string // 正規化済み表記（NormalizeCandidate 適用後）
	Occurrences int    // コーパス内の出現回数
	MidSentence int    // 曖昧でない位置（文頭・引用符直後以外）での出現回数
	NERHits     int    // 固有表現認識が entity と判定した出現回数
}

// Stoplist は一般語 stoplist の境界（使う側で定義する小さい interface）。
// 実装は機械置換辞書の供給選別と同じ規則（internal/core/dictionary の Stoplist）を想定する。
type Stoplist interface {
	// Blocks は原語を候補から除くべき一般語かを返す。
	Blocks(source string) bool
}

// TextAnalyzer は外部学習モデルによる本文解析の境界。固有表現の表記列（Entities）と、
// 文頭の語が動詞かの判定（LeadingVerb。命令形のクエスト目標行の識別に使う）を提供する。
// prose 実装（ProseAnalyzer）とテスト用の偽実装を差し替えられるようにする。
type TextAnalyzer interface {
	Entities(text string) []string
	LeadingVerb(sentence string) bool
}

// CandidateDetector は未知固有名候補の検出器。既知語（辞書の原語）・一般語 stoplist・
// 本文解析（固有表現・品詞）を判定材料に持つ。副作用を持たない。
type CandidateDetector struct {
	known map[string]bool
	stop  Stoplist
	ner   TextAnalyzer
}

// NewCandidateDetector は検出器を作る。knownSources は辞書の原語（master_term ∪ proper_noun、
// stoplist 選別後）で、正規化して既知語集合に積む。stop は nil なら一般語の選別なし、
// ner は nil なら本文解析の補強を使わない。
func NewCandidateDetector(knownSources []string, stop Stoplist, ner TextAnalyzer) *CandidateDetector {
	known := make(map[string]bool, len(knownSources))
	for _, s := range knownSources {
		if n := NormalizeCandidate(s); n != "" {
			known[n] = true
		}
	}
	return &CandidateDetector{known: known, stop: stop, ner: ner}
}

// NormalizeCandidate は候補表記の正規化。前後空白の除去・語間空白の 1 個への畳み込み・
// 末尾の所有格（'s と複数形所有格の '）の除去を行い、大文字小文字は区別したまま保つ
// （達成基準の正規化基準）。語中のアポストロフィ（M'aiq、Mehrunes' Razor の内側）は保つ。
func NormalizeCandidate(s string) string {
	s = strings.Join(strings.Fields(s), " ")
	s = strings.TrimSuffix(s, "'s")
	s = strings.TrimSuffix(s, "’s")
	return strings.TrimRight(s, "'’-")
}

var (
	// candMarkupRe は動的タグ（<Alias=...> 等）の除去。タグ内の語（Alias 等）を固有名と誤認しないため、
	// トークン化の前に空白へ置き換える。タグの翻訳時保護は本ルールの外（known-issues 2番）。
	candMarkupRe = regexp.MustCompile(`<[^>]*>`)
	// candSentenceRe は文境界。文末記号（連続可）と閉じ引用符・閉じ括弧の後の空白、または改行で文を分ける。
	// 次の文の先頭語は曖昧位置（大文字が構文由来かもしれない）として扱う。
	candSentenceRe = regexp.MustCompile(`[.!?]+["')\]]*\s+|\r?\n+`)
	// candTokenRe は英字の語（内部に単独のアポストロフィ・ハイフンを許す。M'aiq、Riff-raff。
	// 連続するハイフン -- はダッシュなので語に含めない）と、複数語固有名の内側に現れる
	// 記号接続語 &（Kolb & the Dragon）。
	candTokenRe = regexp.MustCompile(`[A-Za-z]+(?:['’-][A-Za-z]+)*|&`)
	// candLowerNameRe は小文字接頭の姓（gro-Nolob、gra-Mutha 等のオーク姓）。小文字始まりでも
	// 固有名の一部として句へ取り込む。
	candLowerNameRe = regexp.MustCompile(`^[a-z'’]+[-'’][A-Z]`)
)

// candConnectors は複数語固有名の内側に現れる接続語（College of Winterhold、
// Potema the Wolf Queen、Kolb & the Dragon）。大文字語に挟まれた時だけ句へ取り込む。
var candConnectors = map[string]bool{"of": true, "the": true, "&": true}

// candPronouns は一人称代名詞とその縮約。英語の表記規則で常に大文字始まりのため、
// 大文字が固有名の根拠にならない。句の部材にしない。
var candPronouns = map[string]bool{"I": true, "I'm": true, "I've": true, "I'll": true, "I'd": true}

// candTitleFillers は title-case の見出し行（書名・章見出し）に混ざる小文字の機能語・記号。
var candTitleFillers = map[string]bool{
	"of": true, "the": true, "a": true, "an": true, "and": true,
	"by": true, "for": true, "in": true, "to": true, "on": true, "&": true,
}

// candLeadSuppress は曖昧位置の句頭で固有名の一部と考えにくい英語の閉クラスの機能語
// （冠詞・限定詞・接続詞・前置詞・疑問詞・談話的な間投詞）。文法上の閉クラスであり
// 語彙データセットではないため、コード内の固定集合で持つ（candConnectors・candPronouns と同じ扱い）。
// stoplist（stopwords-iso）は fire・clear のような固有名の先頭になり得る内容語も含むため、
// この判定には使わない。The は The Reach のように固有名の先頭が普通にあるため含めない。
var candLeadSuppress = map[string]bool{
	// 冠詞・限定詞
	"a": true, "an": true, "all": true, "any": true, "another": true, "some": true,
	"each": true, "every": true, "no": true, "this": true, "that": true, "these": true,
	"those": true, "such": true, "both": true, "either": true, "neither": true,
	"my": true, "your": true, "his": true, "her": true, "its": true, "our": true, "their": true,
	// 接続詞・従属詞
	"and": true, "but": true, "or": true, "nor": true, "so": true, "if": true,
	"when": true, "while": true, "although": true, "though": true, "because": true,
	"since": true, "as": true, "unless": true, "until": true, "once": true, "than": true,
	// 前置詞
	"of": true, "by": true, "for": true, "in": true, "to": true, "on": true, "at": true,
	"with": true, "from": true, "into": true, "onto": true, "upon": true, "over": true,
	"under": true, "before": true, "after": true, "during": true, "between": true,
	"among": true, "through": true, "across": true, "against": true, "about": true,
	"without": true, "within": true, "behind": true, "beyond": true, "near": true,
	// 疑問詞・指示副詞
	"what": true, "who": true, "whom": true, "whose": true, "which": true, "where": true,
	"why": true, "how": true, "there": true, "here": true, "then": true,
	// 談話の間投詞・応答詞
	"oh": true, "ah": true, "aye": true, "yes": true, "well": true, "hey": true,
	"now": true, "not": true,
}

// candToken はトークン化した語 1 件。位置の性質（曖昧位置か、直前と連結可能か）と、
// 直前の区切り文字列（引用符で括られた名前の判定用）を持つ。
type candToken struct {
	text      string
	sepBefore string // 直前トークン（または文頭）からこの語までの区切り文字列
	ambiguous bool   // 文頭、または引用符・コロン等の直後（大文字が構文由来かもしれない位置）
	joined    bool   // 直前トークンと 1 個の空白だけで隣接している（句の結合を許す）
}

// candResetsStart は語の直前の区切り文字列が「文頭に相当する曖昧位置」を作るかを返す。
// 引用符の開始・コロン・丸括弧・ダッシュ・省略記号の直後は、引用文や項目の頭で大文字になり得る。
func candResetsStart(sep string) bool {
	if strings.Contains(sep, "--") {
		return true
	}
	return strings.ContainsAny(sep, `"“”:(—…`)
}

// candTokens は 1 文をトークン列と末尾の残り文字列（最終トークンの後の区切り）に分ける。
// 連結（joined）は 1 個の空白だけの隣接に限る。2 個以上の空白は整形された一覧・見出しの
// 項目境界（Fire Salts␣␣Ruined Attempts）である見込みが高いため、句を切る。
func candTokens(sent string) ([]candToken, string) {
	locs := candTokenRe.FindAllStringIndex(sent, -1)
	toks := make([]candToken, 0, len(locs))
	prevEnd := -1
	for _, loc := range locs {
		// 先頭トークンの区切りは文頭からの前置文字列（開き引用符の検出用）。
		sep := sent[:loc[0]]
		if prevEnd >= 0 {
			sep = sent[prevEnd:loc[0]]
		}
		reset := candResetsStart(sep)
		toks = append(toks, candToken{
			text:      sent[loc[0]:loc[1]],
			sepBefore: sep,
			ambiguous: len(toks) == 0 || reset,
			joined:    len(toks) > 0 && sep == " " && !reset,
		})
		prevEnd = loc[1]
	}
	tail := ""
	if prevEnd >= 0 {
		tail = sent[prevEnd:]
	}
	return toks, tail
}

// candNameToken は語が固有名の部材になり得るか（大文字始まり、または小文字接頭の姓）を返す。
// 1 文字の大文字も部材にする（A Minor Maze の A、Chapter X の X）。
func candNameToken(text string) bool {
	first, _ := utf8.DecodeRuneInString(text)
	return unicode.IsUpper(first) || candLowerNameRe.MatchString(text)
}

// candUsage はコーパス全体の語の用法分布。lc は小文字用法（一般語として出た回数）。
// キーは小文字化した語。曖昧位置だけに出た 1 語の最終判定（acceptCandidate）と、
// 複数語成分の内側の 1 語候補化（componentSpans）が「一般語として使われない語」の根拠に使う。
type candUsage struct {
	lc map[string]int
}

// candCollectUsage はコーパス全体を 1 巡し、動的タグを除いた本文と語の用法分布を作る。
// 用法分布を句の組み立てより先に確定させ、本文の並び順に依存しない判定にする。
func candCollectUsage(texts []string) (candUsage, []string) {
	usage := candUsage{lc: make(map[string]int)}
	cleaned := make([]string, len(texts))
	for i, text := range texts {
		cleaned[i] = candMarkupRe.ReplaceAllString(text, " ")
		for _, sent := range candSentenceRe.Split(cleaned[i], -1) {
			toks, _ := candTokens(sent)
			for _, t := range toks {
				first, _ := utf8.DecodeRuneInString(t.text)
				if unicode.IsLower(first) && !candLowerNameRe.MatchString(t.text) {
					usage.lc[strings.ToLower(t.text)]++
				}
			}
		}
	}
	return usage, cleaned
}

// candPhrase は大文字語の連なり 1 句。tokens は大文字語と、間に挟まった接続語（小文字のまま）。
type candPhrase struct {
	tokens    []string
	ambiguous bool // 句の先頭トークンが曖昧位置にある
	weak      bool // title-case の見出し行の中にあり、文中大文字が固有名の根拠にならない
	quoted    bool // 句全体が引用符で括られている（"Clear Skies" のような名前確定の表記）
	soleLine  bool // 文の全トークンがこの句だけ（Kill Vittoria Vici のような 1 行文）
	verbLead  bool // soleLine かつ先頭語が動詞（命令形のクエスト目標行。全体を候補にしない）
}

// candTitleCaseSentence は文全体が title-case の見出し（書名・章見出し・整形ヘッダ）かを返す。
// 大文字始まりの語が 3 つ以上あり、残りが機能語・1 文字（間隔空け表記の断片）だけの文は
// 見出しとみなす。見出しの中の大文字は表記規則由来で、固有名の根拠（文中大文字用法）にしない。
func candTitleCaseSentence(toks []candToken) bool {
	caps := 0
	for _, t := range toks {
		first, _ := utf8.DecodeRuneInString(t.text)
		switch {
		case unicode.IsUpper(first):
			if utf8.RuneCountInString(t.text) >= 2 {
				caps++
			}
		case !candTitleFillers[t.text] && !candLowerNameRe.MatchString(t.text):
			return false
		}
	}
	return caps >= 3
}

// candQuote は区切り文字列が引用符を含むかを返す（引用符で括られた名前の判定用）。
func candQuote(sep string) bool {
	return strings.ContainsAny(sep, `"“”`)
}

// candPhrases はトークン列から固有名の句を組む。固有名の部材（candNameToken。ただし常に
// 大文字始まりの一人称代名詞は除く）を句に積み、空白だけで隣接した部材は同じ句へ結合する。
// 小文字の接続語（最大 2 連続。of the）は、直後に部材が続く時だけ句へ取り込む。
// カンマ等の区切りや小文字の通常語で句を閉じる。曖昧位置（文頭等）の先頭語の妥当性は
// ここでは判定せず、候補化（phraseSpans の先頭抜き変種）と最終判定（acceptCandidate の
// 用法分布・固有表現認識）に委ねる。weak は文全体が見出し（candTitleCaseSentence）の印。
// tail は文末の残り文字列で、文末で閉じた句の引用符括り（"Clear Skies"）の判定に使う。
func candPhrases(toks []candToken, tail string, weak bool) []candPhrase {
	var phrases []candPhrase
	var cur []string
	var pending []string
	curAmbiguous := false
	curFirstSep := ""
	flush := func(sepAfter string) {
		if len(cur) > 0 {
			quoted := candQuote(curFirstSep) && candQuote(sepAfter)
			phrases = append(phrases, candPhrase{
				tokens: cur, ambiguous: curAmbiguous && !quoted, weak: weak, quoted: quoted,
			})
		}
		cur = nil
		pending = nil
	}
	for _, t := range toks {
		switch {
		case candNameToken(t.text) && !candPronouns[t.text]:
			if len(cur) > 0 && t.joined {
				cur = append(cur, pending...)
				cur = append(cur, t.text)
				pending = nil
				continue
			}
			flush(t.sepBefore)
			cur = []string{t.text}
			curAmbiguous = t.ambiguous
			curFirstSep = t.sepBefore
		case len(cur) > 0 && t.joined && len(pending) < 2 && candConnectors[t.text]:
			pending = append(pending, t.text)
		default:
			flush(t.sepBefore)
		}
	}
	flush(tail)
	// 文の全トークンが 1 句に収まった場合（クエスト目標のような 1 行文）に印を付ける。
	if len(phrases) == 1 && len(phrases[0].tokens) == len(toks) {
		phrases[0].soleLine = true
	}
	return phrases
}

// carveRange は句のトークン範囲 [i, j) から既知語の最長部分列を彫り出し、残った未知の範囲を
// emit へ渡す。範囲は成分（接続語を含まない連なり）の内側なので、接続語の点検は要らない。
// 全体が既知なら何も残らない。最長・左端優先の走査で決定的に動く。
func (d *CandidateDetector) carveRange(tokens []string, i, j int, emit func(int, int)) {
	if i >= j {
		return
	}
	for l := j - i; l >= 1; l-- {
		for s := i; s+l <= j; s++ {
			if d.known[NormalizeCandidate(strings.Join(tokens[s:s+l], " "))] {
				d.carveRange(tokens, i, s, emit)
				d.carveRange(tokens, s+l, j, emit)
				return
			}
		}
	}
	emit(i, j)
}

// phraseSpans は 1 句から候補にするトークン範囲を列挙する。辞書は最長一致を好み、全体と部分の
// 両方を別語彙として持つ（College of Winterhold と Winterhold）ため、複数の解釈を重複なく出す。
func (d *CandidateDetector) phraseSpans(p candPhrase, usage candUsage, nerSet map[string]bool) [][2]int {
	var spans [][2]int
	seen := make(map[[2]int]bool)
	add := func(i, j int) {
		sp := [2]int{i, j}
		if i < j && !seen[sp] {
			seen[sp] = true
			spans = append(spans, sp)
		}
	}
	switch {
	// 引用符で括られた句（"Clear Skies"）は名前の境界が確定しているため、全体だけを出す。
	case p.quoted:
		d.rangeSpans(p.tokens, 0, usage, nerSet, add)
	// 命令形のクエスト目標行（Kill Vittoria Vici のような動詞先頭の 1 行文）は、動詞が
	// 固有名の一部でないため先頭抜きの解釈だけを出す。動詞判定は scanSentence が
	// 品詞解析（TextAnalyzer.LeadingVerb）で行い verbLead に印を付ける。
	case p.verbLead:
		d.rangeSpans(p.tokens, 1, usage, nerSet, add)
	// 閉クラスの機能語（A・An・And・Of 等。The は例外）で始まる曖昧句は、先頭語が固有名の
	// 一部である見込みが薄いため、先頭抜きの解釈だけを出す（A Bosmer → Bosmer）。
	// 見出し行（weak）は例外。書名は機能語で始まる形（A Minor Maze・Of Crossed Daggers）が
	// 普通にあるため、全体も出す。
	case p.ambiguous && !p.weak && len(p.tokens) >= 2 && candLeadSuppress[strings.ToLower(p.tokens[0])]:
		d.rangeSpans(p.tokens, 1, usage, nerSet, add)
	default:
		d.rangeSpans(p.tokens, 0, usage, nerSet, add)
		// 曖昧位置で始まる複数語の句は、先頭語が文の大文字（構文由来）かもしれないため、
		// 先頭語を外した解釈も出す（Reduces Health → Reduces Health と Health）。
		if p.ambiguous && len(p.tokens) >= 2 {
			d.rangeSpans(p.tokens, 1, usage, nerSet, add)
		}
	}
	return spans
}

// rangeSpans はトークン列の [lo, len) 範囲へ句の候補化規則（全体・成分ごとの各解釈）を適用する。
func (d *CandidateDetector) rangeSpans(tokens []string, lo int, usage candUsage, nerSet map[string]bool, add func(int, int)) {
	// 先頭の接続語を落とす（句は接続語で終わらないため、末尾の点検は要らない）。
	sub := tokens[lo:]
	for len(sub) > 0 && candConnectors[sub[0]] {
		lo++
		sub = tokens[lo:]
	}
	// 句全体は未知なら常に候補にする。成分（Pillar・Sacrifice）が個々に既知でも、全体
	// （Pillar of Sacrifice）が独立した語彙として辞書に載る形が普通にあるため、抑止しない。
	if !d.known[NormalizeCandidate(strings.Join(sub, " "))] {
		add(lo, len(tokens))
	}
	for _, c := range candComponents(tokens, lo) {
		if d.known[NormalizeCandidate(strings.Join(tokens[c[0]:c[1]], " "))] {
			continue
		}
		d.componentSpans(tokens, c, usage, nerSet, add)
	}
}

// componentSpans は接続語を含まない成分 1 つから、彫り出し残り・所有格の前後・内側の 1 語を出す。
func (d *CandidateDetector) componentSpans(tokens []string, c [2]int, usage candUsage, nerSet map[string]bool, add func(int, int)) {
	d.carveRange(tokens, c[0], c[1], add)
	// 成分の内側の所有格は名前の境界（Wylandriah's Ingot = 所有者 + 対象物）でもあるため、
	// 所有格の前後で割った解釈も出す（Azura's Star のような所有格込みの語彙は成分全体が担う）。
	for k := c[0]; k < c[1]-1; k++ {
		if candPossessive(tokens[k]) {
			d.carveRange(tokens, c[0], k+1, add)
			d.carveRange(tokens, k+1, c[1], add)
		}
	}
	if c[1]-c[0] < 2 {
		return
	}
	// 複数語成分の内側の語も、コーパス内に小文字用法が無い語（Poetic Edda の Edda、
	// Thane Bryling の Bryling のような名前確実の語）または固有表現認識が entity と
	// 判定した語に限り 1 語で候補にする。辞書の語彙は複合の一部（Winking Skeever の
	// Skeever）を単独項目でも持つため。小文字用法のある語（Guild Treasury Door の
	// Door 等）まで無条件に割ると一般語が混入するので出さない。
	for k := c[0]; k < c[1]; k++ {
		n := NormalizeCandidate(tokens[k])
		if (usage.lc[strings.ToLower(n)] == 0 || nerSet[n]) && !d.known[n] {
			add(k, k+1)
		}
	}
}

// candPossessive は語が所有格（'s・複数形の s'）で終わるかを返す。
func candPossessive(tok string) bool {
	for _, suf := range []string{"'s", "’s", "s'", "s’"} {
		if strings.HasSuffix(tok, suf) {
			return true
		}
	}
	return false
}

// candComponents は句のトークン列の [lo, len) を接続語で分けた成分（接続語を含まない連なり）の範囲に割る。
func candComponents(tokens []string, lo int) [][2]int {
	var comps [][2]int
	start := -1
	for i := lo; i < len(tokens); i++ {
		if candConnectors[tokens[i]] {
			if start >= 0 {
				comps = append(comps, [2]int{start, i})
				start = -1
			}
			continue
		}
		if start < 0 {
			start = i
		}
	}
	if start >= 0 {
		comps = append(comps, [2]int{start, len(tokens)})
	}
	return comps
}

// candEntitySet は固有表現の表記列を正規化済みの照合集合へ畳む。entity 全体と、その中の
// 大文字語トークン単体の両方を積む（句の彫り出し後の 1 語にも当たるようにする）。
func candEntitySet(entities []string) map[string]bool {
	set := make(map[string]bool)
	for _, e := range entities {
		if n := NormalizeCandidate(e); n != "" {
			set[n] = true
		}
		for _, tok := range candTokenRe.FindAllString(e, -1) {
			first, _ := utf8.DecodeRuneInString(tok)
			if unicode.IsUpper(first) {
				set[NormalizeCandidate(tok)] = true
			}
		}
	}
	return set
}

// candStats は候補 1 件のコーパス内の根拠集計。
type candStats struct {
	occurrences int
	midSentence int
	nerHits     int
	alone       int // 1 語だけの句として単独で出た回数（称号・種別語の分離判定に使う）
}

// DetectCandidates はコーパス（本文群）から未知固有名候補を検出する。出力は正規化後表記で
// 一意・辞書順。同一入力に同一出力を返し、副作用を持たない。
// 1 巡目で語の用法分布（candUsage）を集め、2 巡目で句を組んで候補を集計する。
func (d *CandidateDetector) DetectCandidates(texts []string) []CandidateTerm {
	usage, cleaned := candCollectUsage(texts)
	stats := make(map[string]*candStats)
	for _, clean := range cleaned {
		d.scanText(clean, usage, stats)
	}
	d.deriveSplits(stats)
	return d.finishCandidates(stats, usage)
}

// scanText は動的タグ除去済みの本文 1 件を走査し、候補の根拠を stats へ積む。
func (d *CandidateDetector) scanText(clean string, usage candUsage, stats map[string]*candStats) {
	var nerSet map[string]bool
	if d.ner != nil {
		nerSet = candEntitySet(d.ner.Entities(clean))
	}
	for _, sent := range candSentenceRe.Split(clean, -1) {
		d.scanSentence(sent, usage, nerSet, stats)
	}
}

// scanSentence は 1 文から句を組み、各句の候補範囲を stats へ積む。
func (d *CandidateDetector) scanSentence(sent string, usage candUsage, nerSet map[string]bool, stats map[string]*candStats) {
	toks, tail := candTokens(sent)
	phrases := candPhrases(toks, tail, candTitleCaseSentence(toks))
	d.markVerbLead(phrases, sent, usage)
	for _, p := range phrases {
		for _, sp := range d.phraseSpans(p, usage, nerSet) {
			recordSpan(stats, p, sp, nerSet)
		}
	}
}

// recordSpan は候補範囲 1 件の根拠を stats へ積む。範囲は必ず英字の語を含む
// （接続語だけの範囲は生成されない）ため、正規化後も空にならない。
func recordSpan(stats map[string]*candStats, p candPhrase, sp [2]int, nerSet map[string]bool) {
	n := NormalizeCandidate(strings.Join(p.tokens[sp[0]:sp[1]], " "))
	s := stats[n]
	if s == nil {
		s = &candStats{}
		stats[n] = s
	}
	s.occurrences++
	// 曖昧位置の句頭と見出し行の出現は、文中大文字（固有名用法）の根拠に数えない。
	if !p.weak && (!p.ambiguous || sp[0] != 0) {
		s.midSentence++
	}
	if nerSet[n] {
		s.nerHits++
	}
	if len(p.tokens) == 1 {
		s.alone++
	}
}

// markVerbLead は 1 行文の句の先頭語に小文字用法があるとき（Kill … と Mace Etiquette の両方が
// あり得る）、品詞解析で動詞と判定できた行だけを命令形のクエスト目標行として扱う。
// 見出し判定（weak）と重なっても命令形が勝つ（Kill Vittoria Vici は 3 語とも大文字始まりで
// 見出しに見えるが、書名でなく目標行）。
func (d *CandidateDetector) markVerbLead(phrases []candPhrase, sent string, usage candUsage) {
	if d.ner == nil || len(phrases) != 1 {
		return
	}
	p := &phrases[0]
	if p.soleLine && (p.ambiguous || p.weak) && !p.quoted && len(p.tokens) >= 2 &&
		usage.lc[strings.ToLower(NormalizeCandidate(p.tokens[0]))] > 0 &&
		d.ner.LeadingVerb(sent) {
		p.verbLead = true
	}
}

// deriveSplits は複合候補から称号・種別語を外した派生候補を stats へ足す。
// 先頭または末尾の語が 1 語だけの句としても頻出する（単独で立つ）語は、称号（Jarl Hrongar の
// Jarl）や種別語（Dragonrend Shout の Shout）の可能性が高く、辞書の語彙はそれを含まない形で
// 載ることが多い。外した残りが未知なら派生候補として足す（複合の解釈は残したまま増やす）。
// 派生は決定的（辞書順の走査・固定閾値）で、新しく足した候補からも再帰的に派生する。
func (d *CandidateDetector) deriveSplits(stats map[string]*candStats) {
	const aloneMin = 2 // 単独で立つとみなす単独句出現の下限
	queue := make([]string, 0, len(stats))
	for term := range stats {
		queue = append(queue, term)
	}
	sort.Strings(queue)
	for len(queue) > 0 {
		term := queue[0]
		queue = queue[1:]
		toks := strings.Fields(term)
		if len(toks) < 2 {
			continue
		}
		s := stats[term]
		if st := stats[toks[0]]; st != nil && st.alone >= aloneMin {
			// 先頭語を外した残りは句の内側なので、曖昧位置でない出現として数える。
			queue = d.addDerived(stats, queue, toks[1:], s, s.occurrences)
		}
		if st := stats[toks[len(toks)-1]]; st != nil && st.alone >= aloneMin {
			queue = d.addDerived(stats, queue, toks[:len(toks)-1], s, s.midSentence)
		}
	}
}

// addDerived は派生候補 1 件を stats へ足す。残りの端に出た接続語は固有名の一部でないため
// 落とす。派生先が既知語、または既に候補にある（直接の根拠を持つ）場合は足さない。
// 派生先は元の語より必ず短くなる（語を 1 つ外す）ため、元の語との同一は起きない。
func (d *CandidateDetector) addDerived(stats map[string]*candStats, queue []string, sub []string, src *candStats, mid int) []string {
	for len(sub) > 0 && candConnectors[sub[0]] {
		sub = sub[1:]
	}
	for len(sub) > 0 && candConnectors[sub[len(sub)-1]] {
		sub = sub[:len(sub)-1]
	}
	n := NormalizeCandidate(strings.Join(sub, " "))
	if d.known[n] || stats[n] != nil {
		return queue
	}
	stats[n] = &candStats{occurrences: src.occurrences, midSentence: mid, nerHits: src.nerHits}
	return append(queue, n)
}

// finishCandidates は集計済みの根拠へ最終判定（acceptCandidate）を掛け、辞書順の一意な
// 候補列に整える。
func (d *CandidateDetector) finishCandidates(stats map[string]*candStats, usage candUsage) []CandidateTerm {
	out := make([]CandidateTerm, 0, len(stats))
	for term, s := range stats {
		if !d.acceptCandidate(term, s, usage) {
			continue
		}
		out = append(out, CandidateTerm{
			Term: term, Occurrences: s.occurrences, MidSentence: s.midSentence, NERHits: s.nerHits,
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Term < out[j].Term })
	return out
}

// acceptCandidate は集計済みの候補 1 件を最終判定する。1 語候補は 2 文字以上かつ stoplist の
// 一般語でないこと。複数語の句は表記そのものが固有名の根拠なので通す（一般語の複合
// （Fire Salts・Giant's Toe）が辞書の item 名として普通にあるため、語の一般性では落とさない）。
// 曖昧位置（文頭・引用符直後）だけに出た 1 語は、コーパス内に小文字用法が無い、または
// 固有表現認識が entity と判定した場合だけ通す（文頭の一般語を落とす）。
func (d *CandidateDetector) acceptCandidate(term string, s *candStats, usage candUsage) bool {
	multi := strings.Contains(term, " ")
	if !multi {
		if utf8.RuneCountInString(term) < 2 {
			return false
		}
		if d.stop != nil && d.stop.Blocks(term) {
			return false
		}
	}
	if s.midSentence > 0 || multi {
		return true
	}
	return usage.lc[strings.ToLower(term)] == 0 || s.nerHits > 0
}
