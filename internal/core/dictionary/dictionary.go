// Package dictionary は固有名の機械置換（原語→確定訳語）を行う純粋ルール。
// 叙述文・台詞の原文へ辞書を貪欲最長一致で当て、副作用なく置換結果と内訳を返す。
package dictionary

import (
	"regexp"
	"sort"
	"strings"
)

// Term は固有名辞書の 1 件。原語（英語 FULL）を確定訳語（日本語の公式既訳）へ
// 機械置換するための対。
type Term struct {
	Source string // 原語（英語）
	Dest   string // 確定訳語（日本語）
}

// Dictionary は固有名の機械置換器。叙述文・台詞の原文に対し、辞書の原語を貪欲最長一致で照合し
// 確定訳語へ置換する。照合は最長一致優先・語境界・大小区別（原語は大文字始まりの固有名なので、
// 大小を区別すると一般語の小文字出現に当たりにくい）で行う。副作用を持たない。
type Dictionary struct {
	bySource map[string]string
	re       *regexp.Regexp
}

// NewDictionary は辞書語の対から置換器を作る。原語または確定訳語が空の対は捨てる。
// 同一原語に複数の確定訳語が来た場合は最初の 1 つを採る（同綴り異義は概念モデル弱点1 として許容）。
func NewDictionary(terms []Term) *Dictionary {
	bySource := make(map[string]string, len(terms))
	sources := make([]string, 0, len(terms))
	for _, t := range terms {
		source := strings.TrimSpace(t.Source)
		if source == "" || strings.TrimSpace(t.Dest) == "" {
			continue
		}
		if _, ok := bySource[source]; ok {
			continue // 既出の原語は最初の訳を保つ。
		}
		bySource[source] = t.Dest
		sources = append(sources, source)
	}
	d := &Dictionary{bySource: bySource}
	if len(sources) == 0 {
		return d
	}
	// 最長一致を優先するため原語を長い順に並べ、alternation の先に長い語を置く。
	// Go の regexp は leftmost-first なので、同じ開始位置では先に並んだ（長い）語が当たる。
	sort.Slice(sources, func(i, j int) bool {
		if len(sources[i]) != len(sources[j]) {
			return len(sources[i]) > len(sources[j])
		}
		return sources[i] < sources[j]
	})
	quoted := make([]string, len(sources))
	for i, source := range sources {
		quoted[i] = regexp.QuoteMeta(source)
	}
	// 語境界で囲み、原語の途中一致（Sword が Swordsman の内側に当たる等）を防ぐ。
	d.re = regexp.MustCompile(`\b(?:` + strings.Join(quoted, "|") + `)\b`)
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
		// d.re は bySource のキー（QuoteMeta 済み原語）だけを alternation に並べて組むため、
		// 一致した match は必ず bySource のキーになる。ここでの取得は欠落しない。
		dest := d.bySource[match]
		if !seen[match] {
			seen[match] = true
			used = append(used, Term{Source: match, Dest: dest})
		}
		return dest
	})
	return replaced, used
}

// Extract はApplyと同じ照合規則で原文に含まれる原語を抽出する。
// 本文を置換せず、出現順で同じ原語を一意に返す。
func (d *Dictionary) Extract(source string) []Term {
	if d.re == nil {
		return nil
	}
	seen := make(map[string]bool)
	used := make([]Term, 0)
	for _, match := range d.re.FindAllString(source, -1) {
		if seen[match] {
			continue
		}
		seen[match] = true
		used = append(used, Term{Source: match, Dest: d.bySource[match]})
	}
	return used
}
