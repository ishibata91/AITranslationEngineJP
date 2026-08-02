// Package mention は本文中の既知固有名の言及検出を行う純粋ルール。
// 叙述文・台詞の原文へ語彙（master_term ∪ proper_noun の原語）を照合し、副作用なく一致した語彙を返す。
// 照合規則は機械置換（internal/core/dictionary）と同じ（貪欲最長一致・語境界・大小を区別しない・同一原語は先勝ち）に
// 揃える。言及レコードと注入語がずれると注入の事後検証が成り立たないため、規則の一致が不変条件。
package mention

import (
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"
)

// Term は言及検出の語彙 1 件。Kind と ID は呼び出し側の識別子（供給源の種別と行 id）で、
// 検出は解釈せず一致結果へそのまま返す。Source（原語）だけを照合に使う。
type Term struct {
	Kind   string // 語彙の供給源種別（呼び出し側が定める。例: master_term / proper_noun）
	ID     int64  // 供給源の行 id
	Source string // 原語（英語）。照合キー。
}

// Detector は言及の検出器。原文に対し語彙の原語を貪欲最長一致で照合し、一致した語彙を返す。
// 照合は最長一致優先・語境界を保ち、大文字小文字を区別しない。副作用を持たない。
type Detector struct {
	byLower map[string][]Term
	terms   []Term
	re      *regexp.Regexp
}

// NewDetector は語彙から検出器を作る。原語が空白だけの語彙は捨てる。
// 同一原語に複数の語彙が来た場合は最初の 1 つを採る（機械置換辞書の先勝ちと同じ規則。
// 呼び出し側が master_term → proper_noun の順で渡せば、注入と同じ側が言及の相手になる）。
func NewDetector(terms []Term) *Detector {
	byLower := make(map[string][]Term, len(terms))
	registered := make([]Term, 0, len(terms))
	sources := make([]string, 0, len(terms))
	for _, t := range terms {
		source := strings.TrimSpace(t.Source)
		if source == "" {
			continue
		}
		if _, ok := findEqualFoldTerm(registered, source); ok {
			continue // 既出の原語は最初の語彙を保つ。
		}
		t.Source = source
		lower := strings.ToLower(source)
		byLower[lower] = append(byLower[lower], t)
		registered = append(registered, t)
		sources = append(sources, source)
	}
	d := &Detector{byLower: byLower, terms: registered}
	if len(sources) == 0 {
		return d
	}
	// 最長一致を優先するため原語を長い順に並べ、alternation の先に長い語を置く。
	// Go の regexp は leftmost-first なので、同じ開始位置では先に並んだ（長い）語が当たる。
	sort.Slice(sources, func(i, j int) bool {
		leftLength := utf8.RuneCountInString(sources[i])
		rightLength := utf8.RuneCountInString(sources[j])
		if leftLength != rightLength {
			return leftLength > rightLength
		}
		return sources[i] < sources[j]
	})
	quoted := make([]string, len(sources))
	for i, source := range sources {
		quoted[i] = regexp.QuoteMeta(source)
	}
	// 語境界で囲み、原語の途中一致（Sword が Swordsman の内側に当たる等）を防ぐ。
	d.re = regexp.MustCompile(`(?i)\b(?:` + strings.Join(quoted, "|") + `)\b`)
	return d
}

// Detect は原文中に出現した語彙を出現順で返す。同じ原語の重複出現は 1 件に畳む。
// 語彙が空、または一致が無ければ nil を返す。
func (d *Detector) Detect(text string) []Term {
	if d.re == nil {
		return nil
	}
	var found []Term
	seen := make(map[string]bool)
	for _, match := range d.re.FindAllString(text, -1) {
		term, ok := findEqualFoldTerm(d.byLower[strings.ToLower(match)], match)
		if !ok {
			// regexp の大小無視と ToLower の同値関係が異なる場合だけ、全登録語から選ぶ。
			term, ok = findEqualFoldTerm(d.terms, match)
		}
		if !ok || seen[term.Source] {
			continue
		}
		seen[term.Source] = true
		found = append(found, term)
	}
	return found
}

func findEqualFoldTerm(terms []Term, source string) (Term, bool) {
	for _, term := range terms {
		if strings.EqualFold(term.Source, source) {
			return term, true
		}
	}
	return Term{}, false
}
