package engine

import (
	"testing"
)

// 貪欲最長一致で固有名を確定訳語へ置換すること。最長一致・語境界・大小区別・重複の畳みを確かめる。
func TestDictionaryApply(t *testing.T) {
	cases := []struct {
		name     string
		terms    []DictionaryTerm
		source   string
		want     string
		wantUsed []string // 置換した原語の一覧（出現順、重複なし）
	}{
		{
			name:     "最長一致を優先する",
			terms:    []DictionaryTerm{{"Iron", "鉄"}, {"Iron Sword", "鉄の剣"}},
			source:   "An Iron Sword and Iron ore.",
			want:     "An 鉄の剣 and 鉄 ore.",
			wantUsed: []string{"Iron Sword", "Iron"},
		},
		{
			name:     "語境界の内側には当てない",
			terms:    []DictionaryTerm{{"Sword", "剣"}},
			source:   "The Swordsman drew a Sword.",
			want:     "The Swordsman drew a 剣.",
			wantUsed: []string{"Sword"},
		},
		{
			name:     "大小を区別し小文字の一般語は置換しない",
			terms:    []DictionaryTerm{{"Storm", "嵐"}},
			source:   "Storm hit during the storm.",
			want:     "嵐 hit during the storm.",
			wantUsed: []string{"Storm"},
		},
		{
			name:     "同じ固有名の複数出現は同一訳へ置換し used は畳む",
			terms:    []DictionaryTerm{{"Riften", "リフテン"}},
			source:   "The Riften guard left Riften.",
			want:     "The リフテン guard left リフテン.",
			wantUsed: []string{"Riften"},
		},
		{
			name:     "辞書に無い語は素通しする",
			terms:    []DictionaryTerm{{"Riften", "リフテン"}},
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
			terms: []DictionaryTerm{
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
			// どちらも語境界・大小区別で独立に置換されることを確かめる。
			name:     "同長の別原語は辞書順で並べどちらも置換する",
			terms:    []DictionaryTerm{{"Storm", "嵐"}, {"Frost", "霜"}},
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
	d := NewDictionary([]DictionaryTerm{
		{"Chest", "宝箱"},
		{"Chest", "胸"},
	})
	got, _ := d.Apply("A Chest here.")
	if got != "A 宝箱 here." {
		t.Errorf("Apply = %q, want %q", got, "A 宝箱 here.")
	}
}

func sources(terms []DictionaryTerm) []string {
	out := make([]string, len(terms))
	for i, t := range terms {
		out[i] = t.Source
	}
	return out
}

func sameSources(got []DictionaryTerm, want []string) bool {
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
