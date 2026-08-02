package mention

import (
	"reflect"
	"testing"
)

// 本文中の語彙が出現順で検出され、Kind と ID が語彙のまま返ること（基本経路）。
func TestDetectReturnsTermsInOrderOfAppearance(t *testing.T) {
	d := NewDetector([]Term{
		{Kind: "master_term", ID: 1, Source: "Whiterun"},
		{Kind: "proper_noun", ID: 2, Source: "Dragonbane"},
	})

	got := d.Detect("Take Dragonbane to Whiterun.")

	want := []Term{
		{Kind: "proper_noun", ID: 2, Source: "Dragonbane"},
		{Kind: "master_term", ID: 1, Source: "Whiterun"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Detect = %+v, want %+v", got, want)
	}
}

// 同じ原語が本文に複数回出ても 1 件に畳むこと（言及は原語単位の関連で、出現回数は持たない）。
func TestDetectCollapsesRepeatedOccurrences(t *testing.T) {
	d := NewDetector([]Term{{Kind: "master_term", ID: 1, Source: "Whiterun"}})

	got := d.Detect("Whiterun guards protect Whiterun.")

	if len(got) != 1 || got[0].ID != 1 {
		t.Errorf("Detect = %+v, want Whiterun 1 件", got)
	}
}

// 重なる語彙は最長一致を優先し、内側の短い語を別の言及として重複検出しないこと
// （機械置換 dictionary の貪欲最長一致と同じ規則）。
func TestDetectPrefersLongestMatch(t *testing.T) {
	d := NewDetector([]Term{
		{Kind: "master_term", ID: 1, Source: "Sword"},
		{Kind: "master_term", ID: 2, Source: "Iron Sword"},
	})

	got := d.Detect("An Iron Sword shines.")

	if len(got) != 1 || got[0].ID != 2 {
		t.Errorf("Detect = %+v, want Iron Sword（最長一致）だけ", got)
	}
}

// 同じ長さの語彙が並んでも語順（アルファベット順）で決定的に組まれ、両方検出できること（並びの決定性）。
func TestDetectSameLengthSourcesAreDeterministic(t *testing.T) {
	d := NewDetector([]Term{
		{Kind: "master_term", ID: 1, Source: "Riften"},
		{Kind: "master_term", ID: 2, Source: "Aventu"},
	})

	got := d.Detect("Aventu left Riften.")

	if len(got) != 2 || got[0].ID != 2 || got[1].ID != 1 {
		t.Errorf("Detect = %+v, want Aventu, Riften の順", got)
	}
}

// 語境界を跨ぐ途中一致（Sword が Swordsman の内側）には当たらないこと。
func TestDetectRespectsWordBoundary(t *testing.T) {
	d := NewDetector([]Term{{Kind: "master_term", ID: 1, Source: "Sword"}})

	if got := d.Detect("The Swordsman walks."); got != nil {
		t.Errorf("Detect = %+v, want nil（途中一致は言及にしない）", got)
	}
}

// 大小だけが異なる表記も登録済みの固有名（Sword）の言及として返すこと。
func TestDetectIsCaseInsensitive(t *testing.T) {
	d := NewDetector([]Term{{Kind: "master_term", ID: 1, Source: "Sword"}})

	got := d.Detect("a sword on the table")
	if len(got) != 1 || got[0].ID != 1 || got[0].Source != "Sword" {
		t.Errorf("Detect = %+v, want registered Sword", got)
	}
}

// 同一原語が複数の語彙（master_term と proper_noun）から来た場合は最初の 1 つが勝つこと
// （機械置換辞書の先勝ちと同じ規則。master_term を先に渡せば注入と同じ側が言及の相手になる）。
func TestNewDetectorFirstTermWinsOnDuplicateSource(t *testing.T) {
	d := NewDetector([]Term{
		{Kind: "master_term", ID: 1, Source: "Dragonbane"},
		{Kind: "proper_noun", ID: 9, Source: "Dragonbane"},
	})

	got := d.Detect("Dragonbane hums.")

	if len(got) != 1 || got[0].Kind != "master_term" || got[0].ID != 1 {
		t.Errorf("Detect = %+v, want master_term ID=1（先勝ち）", got)
	}
}

// 原語が空白だけの語彙は捨てること（照合キーにならない入力の除外）。
func TestNewDetectorDropsBlankSources(t *testing.T) {
	d := NewDetector([]Term{
		{Kind: "master_term", ID: 1, Source: "   "},
		{Kind: "master_term", ID: 2, Source: "Riften"},
	})

	got := d.Detect("Riften stands.")

	if len(got) != 1 || got[0].ID != 2 {
		t.Errorf("Detect = %+v, want Riften だけ", got)
	}
}

// 語彙が空（全部空白で捨てた場合を含む）なら、どの本文にも一致しないこと（境界）。
func TestDetectWithEmptyVocabulary(t *testing.T) {
	for name, d := range map[string]*Detector{
		"語彙なし":    NewDetector(nil),
		"空白だけの語彙": NewDetector([]Term{{Kind: "master_term", ID: 1, Source: " "}}),
	} {
		if got := d.Detect("Whiterun"); got != nil {
			t.Errorf("%s: Detect = %+v, want nil", name, got)
		}
	}
}

// 原語に正規表現のメタ文字（.）が含まれても、任意一文字ではなく文字どおりに照合すること（QuoteMeta の経路）。
func TestDetectQuotesRegexMetaCharacters(t *testing.T) {
	d := NewDetector([]Term{{Kind: "master_term", ID: 1, Source: "St. Sword"}})

	// '.' が任意一文字として当たると StX Sword に誤一致する。文字どおりの照合なら一致しない。
	if got := d.Detect("The StX Sword shines."); got != nil {
		t.Errorf("Detect = %+v, want nil（メタ文字の誤一致）", got)
	}
	got := d.Detect("The St. Sword shines.")
	if len(got) != 1 || got[0].ID != 1 {
		t.Errorf("Detect = %+v, want St. Sword", got)
	}
}

// 中置の記号（アポストロフィ）を含む固有名（M'aiq）が語境界内で一致すること（実データ代表の名前形）。
func TestDetectMatchesNameWithApostrophe(t *testing.T) {
	d := NewDetector([]Term{{Kind: "master_term", ID: 1, Source: "M'aiq"}})

	got := d.Detect("M'aiq knows much.")

	if len(got) != 1 || got[0].ID != 1 {
		t.Errorf("Detect = %+v, want M'aiq", got)
	}
}

// 言及検出も機械置換と同じ大小無視の規則で登録済みの固有名を返す。
func TestDetectorCaseInsensitiveSpecifications(t *testing.T) {
	tests := []struct {
		name string
		text string
	}{
		{
			name: "本文から固有名を見つける時は大文字小文字を区別しないこと",
			text: "I asked iNiGo to wait.",
		},
		{
			name: "本文の固有名がすべて小文字の場合も見つけること",
			text: "I asked inigo to wait.",
		},
		{
			name: "本文の固有名がすべて大文字の場合も見つけること",
			text: "I asked INIGO to wait.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDetector([]Term{{Kind: "proper_noun", ID: 1, Source: "Inigo"}})

			got := d.Detect(tt.text)

			if len(got) != 1 || got[0].ID != 1 || got[0].Source != "Inigo" {
				t.Errorf("Detect = %+v, want registered Inigo", got)
			}
		})
	}
}

// 大小だけが異なる複数の出現を登録済みの同じ固有名1件へまとめる。
func TestDetectCollapsesCaseVariants(t *testing.T) {
	d := NewDetector([]Term{{Kind: "proper_noun", ID: 1, Source: "Inigo"}})

	got := d.Detect("Inigo met inigo and INIGO.")

	if len(got) != 1 || got[0].ID != 1 || got[0].Source != "Inigo" {
		t.Errorf("Detect = %+v, want registered Inigo 1 件", got)
	}
}

// ToLower の key が同じでも EqualFold で異なる登録語は候補から失わない。
func TestDetectorRetainsTermsWithSameLowerKeyButDifferentEqualFold(t *testing.T) {
	d := NewDetector([]Term{{Kind: "proper_noun", ID: 1, Source: "İ"}, {Kind: "proper_noun", ID: 2, Source: "i"}})

	got := d.Detect("i")

	if len(got) != 1 || got[0].ID != 2 || got[0].Source != "i" {
		t.Errorf("Detect = %+v, want ASCII i", got)
	}
}

// ToLower の key が異なる SimpleFold 同値語も全登録語から取得する。
func TestDetectorFallsBackWhenEqualFoldTermsHaveDifferentLowerKeys(t *testing.T) {
	d := NewDetector([]Term{{Kind: "proper_noun", ID: 1, Source: "AΣA"}})

	got := d.Detect("aςa")

	if len(got) != 1 || got[0].ID != 1 || got[0].Source != "AΣA" {
		t.Errorf("Detect = %+v, want registered AΣA", got)
	}
}

// ToLower の key が異なる EqualFold 同値語にも語彙の先勝ちを適用する。
func TestNewDetectorFirstTermWinsForEqualFoldTermsWithDifferentLowerKeys(t *testing.T) {
	d := NewDetector([]Term{{Kind: "master_term", ID: 1, Source: "AΣA"}, {Kind: "proper_noun", ID: 9, Source: "AςA"}})

	got := d.Detect("aςa")

	if len(got) != 1 || got[0].Kind != "master_term" || got[0].ID != 1 || got[0].Source != "AΣA" {
		t.Errorf("Detect = %+v, want first registered AΣA", got)
	}
}

// 大小無視で短い語にも一致する場合は UTF-8 byte 数ではなく rune 数による最長一致を優先する。
func TestDetectorPrefersLongestRuneCountWithCaseInsensitiveRegexp(t *testing.T) {
	d := NewDetector([]Term{{Kind: "proper_noun", ID: 1, Source: "Kx"}, {Kind: "proper_noun", ID: 2, Source: "K"}})

	got := d.Detect("kx")

	if len(got) != 1 || got[0].ID != 1 || got[0].Source != "Kx" {
		t.Errorf("Detect = %+v, want Kx（rune 数の最長一致）", got)
	}
}
