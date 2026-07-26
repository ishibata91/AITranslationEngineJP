package harness

import (
	"path/filepath"
	"strings"
	"testing"

	"aitranslationenginejp/internal/core/dictionary"
	"aitranslationenginejp/internal/core/rolespeech"
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

// R-1-2: 既訳の収集は、対象 plugin が置かれた Data フォルダを走査対象にする。
// 対象 plugin もそのフォルダの中にあるため、対象 plugin 自身が日本語訳を持つ場合もその英日対が集まる。
// 収集は抽出より前に走る（横断辞書の派生が既訳を入力にするため、抽出前に既訳が揃っている必要がある）。
func TestReferenceCollectionScansPluginDataFolder(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "scan.sqlite3")
	f := SyntheticFixture()
	// 実運用と同じく Data フォルダ配下の plugin パスを渡す（plugin の束ねキーは filepath.Base）。
	pluginPath := filepath.Join("SyntheticData", f.PluginName)
	ex := &SeedExtractor{DBPath: dbPath, Fixture: f}

	roleSpeech, err := rolespeech.ParseRoleSpeech(strings.NewReader(syntheticRoleSpeech))
	if err != nil {
		t.Fatalf("合成役割語の構築: %v", err)
	}
	stop, err := dictionary.ParseStoplist(strings.NewReader(syntheticStopwords))
	if err != nil {
		t.Fatalf("合成 stoplist の構築: %v", err)
	}
	if _, err := Run(RunConfig{
		DBPath: dbPath, Extractor: ex, Lexicon: fakeLexicon{}, RoleSpeech: roleSpeech, Stoplist: stop,
		PluginPath: pluginPath, Model: "fake-model", Fixed: SyntheticAITranslations(),
	}); err != nil {
		t.Fatalf("合成 harness の実行: %v", err)
	}

	want := []string{"references:SyntheticData", "extract:" + pluginPath}
	if len(ex.Calls) != len(want) || ex.Calls[0] != want[0] || ex.Calls[1] != want[1] {
		t.Fatalf("抽出子の呼び出し = %v, want %v", ex.Calls, want)
	}
}
