package termusage

import (
	"testing"
)

// 会話文から用法分布を作ること。文分割・文頭語の不算入・大文字始まり（固有名用法）の UC 算入・
// 小文字始まり（一般語用法）の LC 算入・縮約の小文字連分割・アポストロフィ始まり語の扱いを突く。
func TestBuildUsage(t *testing.T) {
	u := BuildUsage([]string{
		"I saw Grelod. The master left.", // 文分割と文頭語の不算入、固有名 Grelod の UC
		"You aren't Mercer here.",        // 縮約 aren't の小文字連分割、固有名 Mercer の UC
		"Well 'tis Brand.",               // アポストロフィ始まり語と固有名 Brand の UC
	})

	// 固有名用法（文頭以外の大文字始まり）は UC に積む。
	wantUC := map[string]int{"grelod": 1, "mercer": 1, "tis": 1, "brand": 1}
	for w, n := range wantUC {
		if u.UC[w] != n {
			t.Errorf("UC[%q] = %d, want %d", w, u.UC[w], n)
		}
	}
	// 一般語用法（小文字始まり）は LC に積む。縮約 aren't は aren と t に割れる。
	wantLC := map[string]int{"saw": 1, "master": 1, "left": 1, "aren": 1, "t": 1, "here": 1}
	for w, n := range wantLC {
		if u.LC[w] != n {
			t.Errorf("LC[%q] = %d, want %d", w, u.LC[w], n)
		}
	}
	// 文頭語の大文字始まり（I / The / You / Well）は固有名用法に数えない。
	for _, w := range []string{"i", "the", "you", "well"} {
		if u.UC[w] != 0 {
			t.Errorf("UC[%q] = %d, want 0（文頭語は数えない）", w, u.UC[w])
		}
	}
}
