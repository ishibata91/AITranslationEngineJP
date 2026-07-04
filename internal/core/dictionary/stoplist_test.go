package dictionary

import (
	"errors"
	"strings"
	"testing"
	"testing/iotest"
)

// 1 行 1 語のリストを読み、原語全体の小文字一語一致だけを選別すること。
// 大小の正規化・空白の除去・複数語の不適用・空入力の扱いを確かめる。
func TestStoplistBlocks(t *testing.T) {
	list, err := ParseStoplist(strings.NewReader("yes\n No \n\nOPEN\n"))
	if err != nil {
		t.Fatalf("ParseStoplist: %v", err)
	}
	cases := []struct {
		name   string
		source string
		want   bool
	}{
		{"一語一致は大小を問わず除く", "Yes", true},
		{"リスト側の空白・大文字は正規化して除く", "no", true},
		{"原語側の前後空白は落として判定する", " Open ", true},
		{"リストに無い語は除かない", "Riften", false},
		{"複数語の固有名には当てない", "Yes Man", false},
		{"空の原語は除かない", "  ", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := list.Blocks(c.source); got != c.want {
				t.Errorf("Blocks(%q) = %v, want %v", c.source, got, c.want)
			}
		})
	}
}

// nil の Stoplist は選別なし（常に false）で、辞書・言及の供給をそのまま通すこと。
func TestStoplistNilBlocksNothing(t *testing.T) {
	var list *Stoplist
	if list.Blocks("Yes") {
		t.Error("nil Stoplist が語を除いた（選別なしで通すべき）")
	}
}

// リストの読み取りが失敗したら error を返すこと。
func TestParseStoplistReadError(t *testing.T) {
	readErr := errors.New("read broken")
	if _, err := ParseStoplist(iotest.ErrReader(readErr)); !errors.Is(err, readErr) {
		t.Errorf("ParseStoplist error = %v, want wrapped %v", err, readErr)
	}
}
