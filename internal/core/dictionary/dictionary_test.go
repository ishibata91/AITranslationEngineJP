package dictionary

import (
	"testing"
)

// 貪欲最長一致で固有名を確定訳語へ置換すること。最長一致・語境界・大小を区別しない照合・重複の畳みを確かめる。
func TestDictionaryApply(t *testing.T) {
	cases := []struct {
		name     string
		terms    []Term
		source   string
		want     string
		wantUsed []string // 置換した原語の一覧（出現順、重複なし）
	}{
		{
			name:     "最長一致を優先する",
			terms:    []Term{{"Iron", "鉄"}, {"Iron Sword", "鉄の剣"}},
			source:   "An Iron Sword and Iron ore.",
			want:     "An 鉄の剣 and 鉄 ore.",
			wantUsed: []string{"Iron Sword", "Iron"},
		},
		{
			name:     "語境界の内側には当てない",
			terms:    []Term{{"Sword", "剣"}},
			source:   "The Swordsman drew a Sword.",
			want:     "The Swordsman drew a 剣.",
			wantUsed: []string{"Sword"},
		},
		{
			name:     "大小を区別せず同じ固有名を置換する",
			terms:    []Term{{"Storm", "嵐"}},
			source:   "Storm hit during the storm.",
			want:     "嵐 hit during the 嵐.",
			wantUsed: []string{"Storm"},
		},
		{
			name:     "同じ固有名の複数出現は同一訳へ置換し used は畳む",
			terms:    []Term{{"Riften", "リフテン"}},
			source:   "The Riften guard left Riften.",
			want:     "The リフテン guard left リフテン.",
			wantUsed: []string{"Riften"},
		},
		{
			name:     "辞書に無い語は素通しする",
			terms:    []Term{{"Riften", "リフテン"}},
			source:   "The guard left town.",
			want:     "The guard left town.",
			wantUsed: nil,
		},
		{
			name:     "空辞書は原文をそのまま返す",
			terms:    nil,
			source:   "Iron Sword.",
			want:     "Iron Sword.",
			wantUsed: nil,
		},
		{
			// 原語が空・訳語が空の対は NewDictionary が捨てる。残る有効語だけが置換される。
			name: "原語空・訳語空の対は捨てる",
			terms: []Term{
				{"", "捨てる訳"},       // 原語が空 → 捨てる
				{"Whiterun", ""},   // 訳語が空 → 捨てる
				{"Riften", "リフテン"}, // 有効
			},
			source:   "From Whiterun to Riften.",
			want:     "From Whiterun to リフテン.",
			wantUsed: []string{"Riften"},
		},
		{
			// 同じ長さの別の原語が 2 つあると、並べ替えの同長分岐（辞書順）を通る。
			// どちらも語境界を保って独立に置換されることを確かめる。
			name:     "同長の別原語は辞書順で並べどちらも置換する",
			terms:    []Term{{"Storm", "嵐"}, {"Frost", "霜"}},
			source:   "Frost follows Storm.",
			want:     "霜 follows 嵐.",
			wantUsed: []string{"Frost", "Storm"},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d := NewDictionary(c.terms)
			gotText, gotUsed := d.Apply(c.source)
			if gotText != c.want {
				t.Errorf("Apply text = %q, want %q", gotText, c.want)
			}
			if !sameSources(gotUsed, c.wantUsed) {
				t.Errorf("Apply used = %v, want %v", sources(gotUsed), c.wantUsed)
			}
		})
	}
}

// 同綴り異義（同一原語に複数の訳）は最初の訳を保つこと。
func TestDictionaryFirstWinsOnDuplicateSource(t *testing.T) {
	d := NewDictionary([]Term{
		{"Chest", "宝箱"},
		{"Chest", "胸"},
	})
	got, _ := d.Apply("A Chest here.")
	if got != "A 宝箱 here." {
		t.Errorf("Apply = %q, want %q", got, "A 宝箱 here.")
	}
}

// R-1 の大小無視、語境界、置換内訳の仕様を公開 entrypoint から検証する。
func TestDictionaryApplyCaseInsensitiveSpecifications(t *testing.T) {
	tests := []struct {
		name     string
		terms    []Term
		source   string
		want     string
		wantUsed []string
	}{
		{
			name:     "本文から固有名を見つける時は大文字小文字を区別せず機械置換辞書の同じ訳語へ置き換えること",
			terms:    []Term{{Source: "Inigo", Dest: "イニーゴ"}},
			source:   "I asked iNiGo to wait.",
			want:     "I asked イニーゴ to wait.",
			wantUsed: []string{"Inigo"},
		},
		{
			name:     "本文の固有名がすべて小文字の場合とすべて大文字の場合も機械置換辞書の同じ訳語へ置き換えること",
			terms:    []Term{{Source: "Inigo", Dest: "イニーゴ"}},
			source:   "inigo met INIGO.",
			want:     "イニーゴ met イニーゴ.",
			wantUsed: []string{"Inigo"},
		},
		{
			name:     "機械置換辞書にない固有名の内側で一致した部分の先頭または末尾に語境界がない場合は置き換えないこと",
			terms:    []Term{{Source: "Inigo", Dest: "イニーゴ"}},
			source:   "MiniGopher stayed.",
			want:     "MiniGopher stayed.",
			wantUsed: nil,
		},
		{
			name:     "本文に大文字小文字が異なる同じ固有名が複数ある場合置換内訳には機械置換辞書の固有名と訳語を1件表示すること",
			terms:    []Term{{Source: "Inigo", Dest: "イニーゴ"}},
			source:   "Inigo met inigo and INIGO.",
			want:     "イニーゴ met イニーゴ and イニーゴ.",
			wantUsed: []string{"Inigo"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDictionary(tt.terms)
			gotText, gotUsed := d.Apply(tt.source)
			if gotText != tt.want {
				t.Errorf("Apply text = %q, want %q", gotText, tt.want)
			}
			if !sameSources(gotUsed, tt.wantUsed) {
				t.Errorf("Apply used = %v, want %v", sources(gotUsed), tt.wantUsed)
			}
		})
	}
}

// ToLower の key が同じでも EqualFold で異なる登録語は候補から失わない。
func TestDictionaryRetainsTermsWithSameLowerKeyButDifferentEqualFold(t *testing.T) {
	// İ と i は ToLower の key が同じでも EqualFold では異なるため、別の固有名として保持する。
	d := NewDictionary([]Term{{Source: "İ", Dest: "点付き"}, {Source: "i", Dest: "ASCII"}})

	got, used := d.Apply("i")

	if got != "ASCII" || !sameSources(used, []string{"i"}) {
		t.Errorf("Apply = %q, used = %v, want ASCII / [i]", got, sources(used))
	}
}

// ToLower の key が異なる SimpleFold 同値語も全登録語から取得する。
func TestDictionaryFallsBackWhenEqualFoldTermsHaveDifferentLowerKeys(t *testing.T) {
	// Σ と ς は EqualFold で同じだが ToLower の key が異なるため、候補 bucket に無い時は全登録語から選ぶ。
	d := NewDictionary([]Term{{Source: "AΣA", Dest: "シグマ"}})

	got, used := d.Apply("aςa")

	if got != "シグマ" || !sameSources(used, []string{"AΣA"}) {
		t.Errorf("Apply = %q, used = %v, want シグマ / [AΣA]", got, sources(used))
	}
}

// ToLower の key が異なる EqualFold 同値語にも辞書データの先勝ちを適用する。
func TestDictionaryFirstWinsForEqualFoldTermsWithDifferentLowerKeys(t *testing.T) {
	d := NewDictionary([]Term{{Source: "AΣA", Dest: "先の訳"}, {Source: "AςA", Dest: "後の訳"}})

	got, used := d.Apply("aςa")

	if got != "先の訳" || !sameSources(used, []string{"AΣA"}) {
		t.Errorf("Apply = %q, used = %v, want 先の訳 / [AΣA]", got, sources(used))
	}
}

// 大小無視で短い語にも一致する場合は UTF-8 byte 数ではなく rune 数による最長一致を優先する。
func TestDictionaryPrefersLongestRuneCountWithCaseInsensitiveRegexp(t *testing.T) {
	// K は K と EqualFold で同じだが UTF-8 byte 数が多い。rune 数が長い Kx を先に照合する。
	d := NewDictionary([]Term{{Source: "Kx", Dest: "長い語"}, {Source: "K", Dest: "短い語"}})

	got, used := d.Apply("kx")

	if got != "長い語" || !sameSources(used, []string{"Kx"}) {
		t.Errorf("Apply = %q, used = %v, want 長い語 / [Kx]", got, sources(used))
	}
}

func sources(terms []Term) []string {
	out := make([]string, len(terms))
	for i, t := range terms {
		out[i] = t.Source
	}
	return out
}

func sameSources(got []Term, want []string) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i].Source != want[i] {
			return false
		}
	}
	return true
}
