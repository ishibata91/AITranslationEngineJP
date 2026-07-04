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

// 大文字小文字を区別し、一般語の小文字出現（sword）を固有名（Sword）の言及にしないこと。
func TestDetectIsCaseSensitive(t *testing.T) {
	d := NewDetector([]Term{{Kind: "master_term", ID: 1, Source: "Sword"}})

	if got := d.Detect("a sword on the table"); got != nil {
		t.Errorf("Detect = %+v, want nil（大小区別）", got)
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
