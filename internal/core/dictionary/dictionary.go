// Package dictionary は固有名の機械置換（原語→確定訳語）を行う純粋ルール。
// 叙述文・台詞の原文へ辞書を貪欲最長一致で当て、副作用なく置換結果と内訳を返す。
package dictionary

import (
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"
)

// Term は固有名辞書の 1 件。原語（英語 FULL）を確定訳語（日本語の公式既訳）へ
// 機械置換するための対。
type Term struct {
	Source string // 原語（英語）
	Dest   string // 確定訳語（日本語）
}

// Dictionary は固有名の機械置換器。叙述文・台詞の原文に対し、辞書の原語を貪欲最長一致で照合し
// 確定訳語へ置換する。照合は最長一致優先・語境界を保ち、大文字小文字を区別しない。
// 副作用を持たない。
type Dictionary struct {
	byLower map[string][]Term
	terms   []Term
	re      *regexp.Regexp
}

// NewDictionary は辞書語の対から置換器を作る。原語または確定訳語が空の対は捨てる。
// 同一原語に複数の確定訳語が来た場合は最初の 1 つを採る（同綴り異義は概念モデル弱点1 として許容）。
func NewDictionary(terms []Term) *Dictionary {
	byLower := make(map[string][]Term, len(terms))
	registered := make([]Term, 0, len(terms))
	sources := make([]string, 0, len(terms))
	for _, t := range terms {
		source := strings.TrimSpace(t.Source)
		if source == "" || strings.TrimSpace(t.Dest) == "" {
			continue
		}
		if _, ok := findEqualFoldTerm(registered, source); ok {
			continue // 既出の原語は最初の訳を保つ。
		}
		term := Term{Source: source, Dest: t.Dest}
		lower := strings.ToLower(source)
		byLower[lower] = append(byLower[lower], term)
		registered = append(registered, term)
		sources = append(sources, source)
	}
	d := &Dictionary{byLower: byLower, terms: registered}
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

// Apply は原文中の固有名を確定訳語へ置換し、置換後の原文と置換した語の一覧を返す。
// 置換した語は出現順で、同じ語の重複は畳む（結果行への併記用）。辞書が空なら原文をそのまま返す。
func (d *Dictionary) Apply(source string) (string, []Term) {
	if d.re == nil {
		return source, nil
	}
	var used []Term
	seen := make(map[string]bool)
	replaced := d.re.ReplaceAllStringFunc(source, func(match string) string {
		term, ok := findEqualFoldTerm(d.byLower[strings.ToLower(match)], match)
		if !ok {
			// regexp の大小無視と ToLower の同値関係が異なる場合だけ、全登録語から選ぶ。
			term, _ = findEqualFoldTerm(d.terms, match)
		}
		if !seen[term.Source] {
			seen[term.Source] = true
			used = append(used, term)
		}
		return term.Dest
	})
	return replaced, used
}

func findEqualFoldTerm(terms []Term, source string) (Term, bool) {
	for _, term := range terms {
		if strings.EqualFold(term.Source, source) {
			return term, true
		}
	}
	return Term{}, false
}
