// Package termderive は固有名（人名）の部分形（名のみ・短名）の訳を派生する純粋ルール。
// 用法分布（一般語か固有名か）で安全フィルタを掛け、確定訳語の対だけを返す。
package termderive

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

// 本ファイルは人名の部分形を辞書へ派生する純粋ルールを持つ。副作用を持たず、入力（NPC 対・用法集計・
// base 原語集合・設定）から出力（安全な派生対）を決める。XML 解析や master_term 書き込みは本ルールの外側
// （ビルド時コマンド）で行い、本ルールは判定だけを担う。置換器 dictionary.go と同じ engine 配下に置く。
//
// 問題: master 辞書は固有名（NPC_:FULL のフルネーム）を持つが、台詞・叙述文では人名が「名のみ」や「短名」で
// 出るため辞書とかみ合わず、名のみで出た人名が AI 任せになって訳が揺れる（例: 辞書に Grelod the Kind は
// あるが、台詞の Grelod 単独は当たらない）。本ルールは部分形の訳を機械派生し、安全フィルタで一般語衝突を
// 防いで辞書へ足せる対だけを返す。

// 派生の由来種別。master_term・proper_noun の category 印（"derive:"+Kind）にも使う。
const (
	KindShrt   = "shrt"   // NPC_:SHRT（作者記述の短縮別名）の通過。
	KindByname = "byname" // " the " を含むフルネームの前部（二つ名の前の名）。
	KindTwo    = "two"    // 空白 2 語・中黒 2 語で割れる姓名分割。
)

// NamePair は NPC レコードの「原語（英語）→ 訳語（日本語）」の 1 対。派生規則の入力。
// 訳の作者（人間の既訳か、実行内の AI 訳か）で扱いを変えない。採否は語形の防御（safePair・landmine と
// 各派生の構造判定）だけで決める。
type NamePair struct {
	Source string
	Dest   string
	Index  int // 入力列内の位置。派生元のメタデータを呼び出し側で復元する。
}

// Usage は本文中の各英単語の用法分布。語は小文字化して数える。
// LC は文中で小文字始まりに出た回数（一般語用法）、UC は文頭以外で大文字始まりに出た回数（固有名用法）。
type Usage struct {
	LC map[string]int
	UC map[string]int
}

// DerivedTerm は派生で得た「原語（部分形）→ 確定訳語」と由来種別。
type DerivedTerm struct {
	Source      string
	Dest        string
	Kind        string
	SourceIndex int // 派生元 NamePair の Index。
}

// DeriveConfig は派生規則のしきい値と除外集合。実データ検証で決めた既定は DefaultDeriveConfig が返す。
type DeriveConfig struct {
	Contractions        map[string]bool // 縮約素体（aren など）。「't」分割で誤当たりする語を捨てる。
	Titles              map[string]bool // 称号（master, lord など）。役割語を名と誤認しない。
	Factions            map[string]bool // 種族・ハウス語（dunmer, nord など）。
	MinSourceLen        int             // 原語の最小長（大文字始まりの短い一般語を弾く）。
	LandmineLCMin       int             // 一般語判定の小文字出現の下限（lc>=本値 かつ lc>=uc で一般語）。
	HyphenatedNameParts bool            // ハイフンを含む姓を中黒で区切る部分と対応付けるか。
}

// bynameSepRe は二つ名の区切り（前後に空白を伴う the）。前部の名を取り出す分割に使う。
var bynameSepRe = regexp.MustCompile(`(?i)\s+the\s+`)

// DefaultDeriveConfig は実データ（base + DLC）で検証した既定設定を返す。
// 集合と数値は _probe8/_probe9 の検証で過剰置換 0（破壊型）・被覆 99.9%・held-out 汎用性を確かめた値。
func DefaultDeriveConfig() DeriveConfig {
	return DeriveConfig{
		Contractions:  wordSet("aren isn wasn weren hasn haven hadn doesn didn don won couldn wouldn shouldn mustn needn daren ain shan mightn oughtn"),
		Titles:        wordSet("master lord lady captain king queen sir madam jarl thane general commander chief elder matron scourge"),
		Factions:      wordSet("dunmer dwemer altmer bosmer breton nord orc khajiit argonian redguard imperial falmer forsworn redoran hlaalu telvanni indoril stormcloak"),
		MinSourceLen:  4,
		LandmineLCMin: 3,
	}
}

// wordSet は空白区切りの語を集合に変換する（設定の可読な記述用）。
func wordSet(s string) map[string]bool {
	out := map[string]bool{}
	for _, w := range strings.Fields(s) {
		out[w] = true
	}
	return out
}

// DeriveTerms は NPC 対から人名の部分形の訳を派生し、安全フィルタを通った対だけを返す。
// 同一原語に複数の由来が当たる場合は最初に採れた由来を保つ（shrt → byname → two の順）。
// base 辞書に既にある原語は base を優先して派生に含めない。
func DeriveTerms(fullPairs, shrtPairs []NamePair, usage Usage, baseSources map[string]bool, cfg DeriveConfig) []DerivedTerm {
	out := []DerivedTerm{}
	seen := map[string]bool{}
	add := func(source, dest, kind string, sourceIndex int) {
		if baseSources[source] || seen[source] {
			return
		}
		seen[source] = true
		out = append(out, DerivedTerm{Source: source, Dest: dest, Kind: kind, SourceIndex: sourceIndex})
	}

	// 派生規則ごとに helper を呼ぶ。同一原語は最初に採れた由来を保つため shrt → byname → two の順で適用する。
	deriveShrt(shrtPairs, usage, cfg, add)
	deriveByname(fullPairs, usage, cfg, add)
	deriveTwo(fullPairs, usage, cfg, add)
	return out
}

// deriveShrt は shrt 由来（作者記述の短縮別名）の安全な対を add へ渡す。
func deriveShrt(shrtPairs []NamePair, usage Usage, cfg DeriveConfig, add func(source, dest, kind string, sourceIndex int)) {
	for _, p := range shrtPairs {
		if safePair(p.Source, p.Dest, usage, cfg) {
			add(p.Source, p.Dest, KindShrt, p.Index)
		}
	}
}

// deriveByname は byname 由来の安全な対を add へ渡す。
// " the " を含む名の前部（英）と、Dest 末尾のカタカナ連（日）を対にする。
func deriveByname(fullPairs []NamePair, usage Usage, cfg DeriveConfig, add func(source, dest, kind string, sourceIndex int)) {
	for _, p := range fullPairs {
		if !hasByname(p.Source) {
			continue
		}
		en := strings.TrimSpace(bynameSepRe.Split(p.Source, 2)[0])
		ja := trailingKanaRun(p.Dest)
		if ja == "" || ja == p.Dest {
			continue
		}
		if safePair(en, ja, usage, cfg) {
			add(en, ja, KindByname, p.Index)
		}
	}
}

// deriveTwo は two 由来（姓名分割）の安全な対を add へ渡す。
// 空白 2 語の姓名を、漢字を含まない中黒区切りの Dest と同数で整列する。語数一致・中黒区切り・漢字なしが
// 対応ずれの防御になるため、供給元（base ゲームか mod か、人間の既訳か AI 訳か）では絞らない。
func deriveTwo(fullPairs []NamePair, usage Usage, cfg DeriveConfig, add func(source, dest, kind string, sourceIndex int)) {
	for _, p := range fullPairs {
		if hasByname(p.Source) {
			continue
		}
		toks := strings.Fields(p.Source)
		parts := strings.Split(p.Dest, "・")
		partCounts, ok := namePartCounts(toks, cfg.HyphenatedNameParts)
		if !ok || len(toks) < 2 || sum(partCounts) != len(parts) || hasKanji(p.Dest) {
			continue
		}
		partOffset := 0
		for i := range toks {
			en := strings.TrimSpace(toks[i])
			partEnd := partOffset + partCounts[i]
			ja := strings.TrimSpace(strings.Join(parts[partOffset:partEnd], "・"))
			partOffset = partEnd
			if safePair(en, ja, usage, cfg) {
				add(en, ja, KindTwo, p.Index)
			}
		}
	}
}

func namePartCounts(tokens []string, allowHyphenated bool) ([]int, bool) {
	counts := make([]int, len(tokens))
	for i, token := range tokens {
		counts[i] = 1
		if !allowHyphenated || !strings.Contains(token, "-") {
			continue
		}
		segments := strings.Split(token, "-")
		for _, segment := range segments {
			if strings.TrimSpace(segment) == "" {
				return nil, false
			}
		}
		counts[i] = len(segments)
	}
	return counts, true
}

func sum(values []int) int {
	total := 0
	for _, value := range values {
		total += value
	}
	return total
}

// safePair は派生候補（英原語 source、日訳 dest）を採るかを判定する。観測した失敗の手書き除外でなく、
// 語形と用法分布による構造判定にする。1 つでも条件を満たさなければ捨てる。
func safePair(source, dest string, usage Usage, cfg DeriveConfig) bool {
	if source == "" || dest == "" {
		return false
	}
	first, _ := utf8.DecodeRuneInString(source)
	if !unicode.IsUpper(first) { // 大文字始まりの固有名だけを採る。
		return false
	}
	if utf8.RuneCountInString(source) < cfg.MinSourceLen { // 短い一般語（Are 等）を弾く。
		return false
	}
	for _, r := range source { // ASCII の英字・アポストロフィ・ハイフンだけを許す。
		allowed := (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || r == '\'' || r == '-'
		if !allowed {
			return false
		}
	}
	if !isKanaToken(dest) { // 日本語側は純カタカナ（＋中黒）に限る。漢字の役割語を名と誤認しない。
		return false
	}
	if utf8.RuneCountInString(dest) < 2 {
		return false
	}
	if landmine(source, usage, cfg) { // 用法分布・称号・種族・縮約で一般語を弾く。
		return false
	}
	return true
}

// landmine は英原語を一般語として捨てるかを判定する。称号・種族・縮約素体に一致するか、
// 本文で小文字（一般語用法）が下限以上かつ大文字始まり（固有名用法）以上のとき一般語とみなす。
func landmine(source string, usage Usage, cfg DeriveConfig) bool {
	st := strings.ToLower(strings.Trim(source, "'"))
	if cfg.Contractions[st] || cfg.Titles[st] || cfg.Factions[st] {
		return true
	}
	lc, uc := usage.LC[st], usage.UC[st]
	return lc >= cfg.LandmineLCMin && lc >= uc
}

// hasByname は名が二つ名（前後に空白を伴う the）を含むかを返す。byname 派生の対象判定に使う。
func hasByname(source string) bool {
	return strings.Contains(strings.ToLower(source), " the ")
}

// isKanaToken は文字列がカタカナ（と区切りの中黒）だけで構成され、カタカナを 1 文字以上含むかを返す。
func isKanaToken(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}
	hasKana := false
	for _, r := range s {
		switch {
		case r == '・': // 中黒は姓名区切りとして許すが、カタカナとは数えない。
		case isKatakana(r):
			hasKana = true
		default:
			return false
		}
	}
	return hasKana
}

// isKatakana はカタカナ（U+30A0..U+30FF）または長音記号かを返す。
func isKatakana(r rune) bool {
	return (r >= '゠' && r <= 'ヿ') || r == 'ー'
}

// hasKanji は文字列が漢字（CJK 統合漢字）を含むかを返す。姓名分割で役割語入りの Dest を弾くために使う。
func hasKanji(s string) bool {
	for _, r := range s {
		if r >= '一' && r <= '鿿' {
			return true
		}
	}
	return false
}

// trailingKanaRun は文字列末尾から続くカタカナ連（中黒を含む）を元の順で返す。
// 二つ名訳（例 親切者のグレロッド）から名のカタカナ部（グレロッド）を取り出すために使う。
func trailingKanaRun(d string) string {
	runes := []rune(d)
	i := len(runes)
	for i > 0 && (isKatakana(runes[i-1]) || runes[i-1] == '・') {
		i--
	}
	return string(runes[i:])
}
