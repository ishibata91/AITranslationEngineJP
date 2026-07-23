package harness

import (
	"path/filepath"
	"testing"
)

// 合成入力の非劣化は、golden 文字列比較でなくオラクル照合（oracle_test.go）が担う。
// 本 file は決定性そのもの（同一入力で出力が揺れない）だけを守る。オラクルの前提となる決定的振る舞いを検出する。

// 同一入力なら 2 回の実行で同じ出力になること（決定性そのものを確かめる）。
// 非劣化判定は出力の決定性が前提のため、fixture 由来の非決定（map 反復順の漏れ等）を検出する。
func TestSyntheticDeterministic(t *testing.T) {
	run := func() string {
		dir := t.TempDir()
		result, err := SyntheticRun(filepath.Join(dir, "synthetic.sqlite3"))
		if err != nil {
			t.Fatalf("合成 harness の実行: %v", err)
		}
		return result.Serialize()
	}
	if a, b := run(), run(); a != b {
		t.Errorf("同一入力で出力が揺れた（非決定）。\n--- 1 回目 ---\n%s\n--- 2 回目 ---\n%s", a, b)
	}
}
