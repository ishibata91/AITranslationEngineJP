package termderive

import (
	"testing"
)

// derivedKey は派生結果を比較しやすい "source=>dest:kind" 文字列にする。
func derivedKey(t DerivedTerm) string {
	return t.Source + "=>" + t.Dest + ":" + t.Kind
}

func derivedKeys(ts []DerivedTerm) []string {
	out := make([]string, len(ts))
	for i, t := range ts {
		out[i] = derivedKey(t)
	}
	return out
}

func sameKeys(got []DerivedTerm, want []string) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if derivedKey(got[i]) != want[i] {
			return false
		}
	}
	return true
}

// 人名の部分形を 3 種（shrt / byname / two）で派生し、安全フィルタの各分岐で採否が変わること。
// 期待値は検証済みルール（_probe8/_probe9 と同じ構造判定）から読める結果にする。
func TestDeriveTerms(t *testing.T) {
	cfg := DefaultDeriveConfig()
	empty := Usage{LC: map[string]int{}, UC: map[string]int{}}

	cases := []struct {
		name        string
		fullPairs   []NamePair
		shrtPairs   []NamePair
		usage       Usage
		baseSources map[string]bool
		want        []string // "source=>dest:kind" の出現順
	}{
		{
			name:      "shrt は短縮別名をそのまま採る",
			shrtPairs: []NamePair{{Source: "Maro", Dest: "マロ"}},
			usage:     empty,
			want:      []string{"Maro=>マロ:shrt"},
		},
		{
			name:      "shrt の称号一致は捨てる",
			shrtPairs: []NamePair{{Source: "Master", Dest: "マスター"}},
			usage:     empty,
			want:      nil,
		},
		{
			name:      "shrt の一般語（用法比）は捨てる",
			shrtPairs: []NamePair{{Source: "Blood", Dest: "ブラッド"}},
			usage:     Usage{LC: map[string]int{"blood": 177}, UC: map[string]int{"blood": 38}},
			want:      nil,
		},
		{
			name:      "byname は二つ名の前部と Dest 末尾カタカナ連を採る",
			fullPairs: []NamePair{{Source: "Grelod the Kind", Dest: "親切者のグレロッド"}},
			usage:     empty,
			want:      []string{"Grelod=>グレロッド:byname"},
		},
		{
			name:      "byname は Dest 末尾にカタカナが無ければ捨てる",
			fullPairs: []NamePair{{Source: "Xeo the Brave", Dest: "勇敢な者"}},
			usage:     empty,
			want:      nil,
		},
		{
			name:      "byname は Dest 全体がカタカナ（部分形が取れない）なら捨てる",
			fullPairs: []NamePair{{Source: "Yuvo the Bold", Dest: "ボールド"}},
			usage:     empty,
			want:      nil,
		},
		{
			name:      "byname は前部が安全フィルタ不成立なら捨てる",
			fullPairs: []NamePair{{Source: "Are the Brave", Dest: "勇者アー"}},
			usage:     empty,
			want:      nil,
		},
		{
			name:      "two は空白2語姓名を中黒2語で整列して採る",
			fullPairs: []NamePair{{Source: "Mercer Frey", Dest: "メルセル・フレイ"}},
			usage:     empty,
			want:      []string{"Mercer=>メルセル:two", "Frey=>フレイ:two"},
		},
		{
			name:      "two は英語と中黒の語数が一致しなければ捨てる",
			fullPairs: []NamePair{{Source: "Mercer Frey Extra", Dest: "メルセル・フレイ"}},
			usage:     empty,
			want:      nil,
		},
		{
			name:      "two は Dest に漢字（役割語）を含めば捨てる",
			fullPairs: []NamePair{{Source: "Foobar Guard", Dest: "フーバー・衛兵"}},
			usage:     empty,
			want:      nil,
		},
		{
			name:      "two は語単位で安全フィルタを通った語だけ採る",
			fullPairs: []NamePair{{Source: "Ab Frey", Dest: "エ・フレイ"}},
			usage:     empty,
			want:      []string{"Frey=>フレイ:two"},
		},
		{
			name:        "base 辞書に既出の原語は派生しない",
			shrtPairs:   []NamePair{{Source: "Maro", Dest: "マロ"}},
			usage:       empty,
			baseSources: map[string]bool{"Maro": true},
			want:        nil,
		},
		{
			name:      "同一原語は最初に採れた由来を保つ",
			shrtPairs: []NamePair{{Source: "Grelod", Dest: "グレロド"}},
			fullPairs: []NamePair{{Source: "Grelod the Kind", Dest: "親切者のグレロッド"}},
			usage:     empty,
			want:      []string{"Grelod=>グレロド:shrt"},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			base := c.baseSources
			if base == nil {
				base = map[string]bool{}
			}
			got := DeriveTerms(c.fullPairs, c.shrtPairs, c.usage, base, cfg)
			if !sameKeys(got, c.want) {
				t.Errorf("DeriveTerms = %v, want %v", derivedKeys(got), c.want)
			}
		})
	}
}

// 安全フィルタ safePair の各棄却条件と採用条件を 1 つずつ突く。
func TestSafePair(t *testing.T) {
	cfg := DefaultDeriveConfig()
	empty := Usage{LC: map[string]int{}, UC: map[string]int{}}
	cases := []struct {
		name   string
		source string
		dest   string
		usage  Usage
		want   bool
	}{
		{"原語が空は不可", "", "グレロッド", empty, false},
		{"確定訳語が空は不可", "Grelod", "", empty, false},
		{"小文字始まりは不可", "grelod", "グレロッド", empty, false},
		{"最小長未満は不可", "Are", "アー", empty, false},
		{"英字以外を含む原語は不可", "Gr3lod", "グレロッド", empty, false},
		{"非カタカナの確定訳語は不可", "Grelod", "Grelod", empty, false},
		{"確定訳語が1文字は不可", "Grelod", "ア", empty, false},
		{"称号は不可（landmine）", "Master", "マスター", empty, false},
		{"通常の固有名は採用", "Grelod", "グレロッド", empty, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := safePair(c.source, c.dest, c.usage, cfg); got != c.want {
				t.Errorf("safePair(%q,%q) = %v, want %v", c.source, c.dest, got, c.want)
			}
		})
	}
}

// 一般語判定 landmine の各経路（縮約・称号・種族・用法比）を突く。
func TestLandmine(t *testing.T) {
	cfg := DefaultDeriveConfig()
	cases := []struct {
		name   string
		source string
		usage  Usage
		want   bool
	}{
		{"縮約素体は一般語", "Aren", Usage{LC: map[string]int{}, UC: map[string]int{}}, true},
		{"称号は一般語", "Lord", Usage{LC: map[string]int{}, UC: map[string]int{}}, true},
		{"種族・ハウス語は一般語", "Nord", Usage{LC: map[string]int{}, UC: map[string]int{}}, true},
		{"小文字用法が下限以上かつ大文字用法以上は一般語", "Blood", Usage{LC: map[string]int{"blood": 177}, UC: map[string]int{"blood": 38}}, true},
		{"小文字用法が下限未満は固有名扱い", "Snow", Usage{LC: map[string]int{"snow": 2}, UC: map[string]int{"snow": 1}}, false},
		{"大文字用法が小文字用法を上回れば固有名扱い", "Storm", Usage{LC: map[string]int{"storm": 5}, UC: map[string]int{"storm": 126}}, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := landmine(c.source, c.usage, cfg); got != c.want {
				t.Errorf("landmine(%q) = %v, want %v", c.source, got, c.want)
			}
		})
	}
}

// 日本語側のカタカナ語判定 isKanaToken の境界を突く。
func TestIsKanaToken(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"グレロッド", true},
		{"メルセル・フレイ", true},
		{"マーロ", true}, // 長音記号を含む
		{"", false},
		{"   ", false},
		{"Grelod", false},
		{"親切者", false},
		{"・", false}, // 中黒だけはカタカナを含まない
	}
	for _, c := range cases {
		if got := isKanaToken(c.in); got != c.want {
			t.Errorf("isKanaToken(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}

// 末尾カタカナ連 trailingKanaRun の取り出しを突く。
func TestTrailingKanaRun(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"親切者のグレロッド", "グレロッド"},
		{"勇者マーロ", "マーロ"},
		{"勇者", ""},
		{"メルセル・フレイ", "メルセル・フレイ"},
		{"", ""},
	}
	for _, c := range cases {
		if got := trailingKanaRun(c.in); got != c.want {
			t.Errorf("trailingKanaRun(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

// 漢字を含むかの判定 hasKanji を突く。
func TestHasKanji(t *testing.T) {
	if !hasKanji("衛兵") {
		t.Error("hasKanji(衛兵) = false, want true")
	}
	if hasKanji("メルセル") {
		t.Error("hasKanji(メルセル) = true, want false")
	}
}
