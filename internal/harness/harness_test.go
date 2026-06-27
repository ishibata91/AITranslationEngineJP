package harness

import (
	"flag"
	"os"
	"path/filepath"
	"testing"
)

// update は golden を再生成するフラグ。`go test -run TestSyntheticNonRegression -update` で更新する。
// 旧挙動へ凍結した比較基準を、意図したコード変更のときだけ更新する（不用意な上書きを避ける）。
var update = flag.Bool("update", false, "golden ファイルを再生成する")

// goldenPath は合成入力の比較基準。著作物を含まない合成出力のため commit して CI 恒久の回帰網にする。
const goldenPath = "testdata/synthetic.golden"

// 合成入力で RunExtractAndTranslate を端から端まで通し、送信プロンプト列と DB 最終状態が
// 凍結した golden と一致することを確かめる。リファクタ前後でこの一致が崩れなければ決定的振る舞いは非劣化。
func TestSyntheticNonRegression(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "synthetic.sqlite3")
	xmlDir := filepath.Join(dir, "xml")
	if err := os.MkdirAll(xmlDir, 0o750); err != nil {
		t.Fatalf("XML ディレクトリの作成: %v", err)
	}

	result, err := SyntheticRun(dbPath, xmlDir)
	if err != nil {
		t.Fatalf("合成 harness の実行: %v", err)
	}
	got := result.Serialize()

	if *update {
		// CI で -update を渡すと壊れた実装を比較基準として凍結し直してしまう。CI 環境では更新を禁じ、
		// golden の更新は手元での意図した操作に限る（CI は環境変数 CI を設定する慣行に従う）。
		if os.Getenv("CI") != "" {
			t.Fatal("CI 環境では -update を使えない（golden の誤上書き防止）。golden 更新は手元で行う")
		}
		if writeErr := os.WriteFile(goldenPath, []byte(got), 0o600); writeErr != nil {
			t.Fatalf("golden の書き出し: %v", writeErr)
		}
		t.Logf("golden を再生成した: %s", goldenPath)
		return
	}

	want, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("golden の読み込み: %v（初回は -update で生成する）", err)
	}
	if got != string(want) {
		t.Errorf("合成出力が golden と不一致。意図した変更なら -update で更新する。\n--- got ---\n%s", got)
	}
}

// 同一入力なら 2 回の実行で同じ出力になること（決定性そのものを確かめる）。
// 非劣化判定は出力の決定性が前提のため、fixture 由来の非決定（map 反復順の漏れ等）を検出する。
func TestSyntheticDeterministic(t *testing.T) {
	run := func() string {
		dir := t.TempDir()
		xmlDir := filepath.Join(dir, "xml")
		if err := os.MkdirAll(xmlDir, 0o750); err != nil {
			t.Fatalf("XML ディレクトリの作成: %v", err)
		}
		result, err := SyntheticRun(filepath.Join(dir, "synthetic.sqlite3"), xmlDir)
		if err != nil {
			t.Fatalf("合成 harness の実行: %v", err)
		}
		return result.Serialize()
	}
	if a, b := run(), run(); a != b {
		t.Errorf("同一入力で出力が揺れた（非決定）。\n--- 1 回目 ---\n%s\n--- 2 回目 ---\n%s", a, b)
	}
}
