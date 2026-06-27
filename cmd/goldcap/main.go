// Command goldcap は実 .esm 由来の baseline golden を捕獲・比較する local 限定ツール。
// 実データ harness（local 限定・再利用可能）の実行点で、内部は internal/harness の Run を本番 Extractor と実辞書で回す。
// provider は harness の決定的 fake のまま（AI 出力の非決定を排す）。golden と入力は著作物を含むため gitignore へ置く。
//
// 使い方:
//
//	capture: go run ./cmd/goldcap -mode capture -plugin <実.esm> -golden <出力先(gitignore)>
//	compare: go run ./cmd/goldcap -mode compare -plugin <実.esm> -golden <比較基準(gitignore)>
//
// golden は指定 git ref を check out した worktree で capture すると、比較基準を旧挙動へ凍結できる（scripts/golden/capture.sh）。
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"aitranslationenginejp/internal/api"
	"aitranslationenginejp/internal/engine"
	"aitranslationenginejp/internal/harness"
	"aitranslationenginejp/internal/lexicon"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "goldcap:", err)
		os.Exit(1)
	}
}

// options は goldcap の実行設定。既定は dev のパス配置（bootstrap と同じ）に合わせる。
type options struct {
	mode      string
	plugin    string
	golden    string
	xmlDir    string
	nrcPath   string
	rolePath  string
	project   string
	schemaDir string
	model     string
}

func parseFlags() options {
	var o options
	flag.StringVar(&o.mode, "mode", "capture", "capture（golden 生成）または compare（golden 比較）")
	flag.StringVar(&o.plugin, "plugin", "", "抽出対象の実 plugin パス（.esm/.esp）。必須")
	flag.StringVar(&o.golden, "golden", "", "golden ファイルのパス（gitignore 配下を指定）。必須")
	flag.StringVar(&o.xmlDir, "xml", "dictionaries/xTranslatorXMLs", "固有名派生が読む xTranslator XML ディレクトリ")
	flag.StringVar(&o.nrcPath, "nrc", "dictionaries/nrc-emolex.txt", "感情辞書（NRC）のパス")
	flag.StringVar(&o.rolePath, "roles", "assets/role-speech.tsv", "役割語テンプレートのパス")
	flag.StringVar(&o.project, "extractor-project", "tools/extractor", "C# 抽出器の dotnet project パス")
	flag.StringVar(&o.schemaDir, "schema", "db/migrations", "schema migration のディレクトリ")
	flag.StringVar(&o.model, "model", "goldcap-fake", "送信モデル名（fake provider は記録のみ）")
	flag.Parse()
	return o
}

func run() error {
	o := parseFlags()
	if o.plugin == "" || o.golden == "" {
		flag.Usage()
		return fmt.Errorf("-plugin と -golden は必須")
	}

	// 実辞書・実役割語を読む（合成 harness と違い、実データの口調生成・感情経路を本番素材で通す）。
	lex, err := lexicon.LoadNRC(o.nrcPath)
	if err != nil {
		return fmt.Errorf("感情辞書の読み込み: %w", err)
	}
	roles, err := engine.LoadRoleSpeech(o.rolePath)
	if err != nil {
		return fmt.Errorf("役割語テンプレートの読み込み: %w", err)
	}

	// 実行ごとに使い捨ての temp DB を開く。抽出子は同じ temp DB へ書き込む（DBPath を揃える）。
	tmpDir, err := os.MkdirTemp("", "goldcap-")
	if err != nil {
		return fmt.Errorf("temp ディレクトリの作成: %w", err)
	}
	defer os.RemoveAll(tmpDir) //nolint:errcheck // 使い捨ての後始末。
	dbPath := filepath.Join(tmpDir, "capture.sqlite3")

	extractor := api.NewDotnetExtractor(api.ExtractorConfig{
		ProjectPath: o.project,
		SchemaDir:   o.schemaDir,
		DBPath:      dbPath,
	})

	captured, err := harness.Run(harness.RunConfig{
		DBPath:      dbPath,
		Extractor:   extractor,
		Lexicon:     lex,
		RoleSpeech:  roles,
		TermsXMLDir: o.xmlDir,
		PluginPath:  o.plugin,
		Model:       o.model,
	})
	if err != nil {
		return fmt.Errorf("実データ harness の実行: %w", err)
	}
	got := captured.Serialize()

	switch o.mode {
	case "capture":
		if err := os.WriteFile(o.golden, []byte(got), 0o600); err != nil {
			return fmt.Errorf("golden の書き出し: %w", err)
		}
		fmt.Printf("golden を捕獲した: %s（プロンプト %d 件）\n", o.golden, len(captured.Prompts))
		return nil
	case "compare":
		want, err := os.ReadFile(o.golden)
		if err != nil {
			return fmt.Errorf("golden の読み込み: %w", err)
		}
		if got != string(want) {
			fmt.Fprintln(os.Stderr, firstDiff(string(want), got))
			return fmt.Errorf("実データ出力が golden と不一致（リファクタで決定的振る舞いが変わった可能性）: %s", o.golden)
		}
		fmt.Printf("一致: 実データ出力は golden と非劣化（プロンプト %d 件）\n", len(captured.Prompts))
		return nil
	default:
		return fmt.Errorf("未知の mode: %q（capture か compare）", o.mode)
	}
}

// firstDiff は want と got の最初に食い違う行を、前後数行の文脈つきで返す（compare 不一致時の原因追跡用）。
func firstDiff(want, got string) string {
	wl := strings.Split(want, "\n")
	gl := strings.Split(got, "\n")
	n := min(len(wl), len(gl))
	for i := range n {
		if wl[i] != gl[i] {
			start := max(i-3, 0)
			var b strings.Builder
			fmt.Fprintf(&b, "--- 最初の相違: %d 行目 ---\n", i+1)
			for j := start; j < i; j++ {
				fmt.Fprintf(&b, "  %s\n", wl[j])
			}
			fmt.Fprintf(&b, "- golden: %s\n", wl[i])
			fmt.Fprintf(&b, "+ 実行値: %s\n", gl[i])
			return b.String()
		}
	}
	return fmt.Sprintf("--- 行数の相違: golden=%d 行 実行値=%d 行 ---", len(wl), len(gl))
}
