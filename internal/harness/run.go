package harness

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"aitranslationenginejp/internal/api"
	"aitranslationenginejp/internal/core/linefeatures"
	"aitranslationenginejp/internal/core/rolespeech"
	"aitranslationenginejp/internal/engine"
	"aitranslationenginejp/internal/store"

	"github.com/jmoiron/sqlx"
)

// RunConfig は非劣化 harness の 1 回の実行に要る注入物。
// 抽出子・感情辞書・役割語・XML ディレクトリ・plugin パスを差し替え可能にし、
// 合成入力（SeedExtractor・最小辞書）と実 .esm 入力（dotnet 抽出子・実辞書）の両方を同じ Run で回す。
type RunConfig struct {
	DBPath      string        // 中心 DB（temp ファイル）のパス
	Extractor   api.Extractor // 抽出段の注入（合成は SeedExtractor、実データは api.DotnetExtractor）
	Lexicon     linefeatures.EmotionLexicon
	RoleSpeech  *rolespeech.Table
	TermsXMLDir string // 固有名派生が読む xTranslator XML ディレクトリ
	PluginPath  string // 抽出対象 plugin のパス
	Model       string // 送信モデル名（fake provider は記録のみ）
}

// Run は store・engine・provider・api を束ね、RunExtractAndTranslate を端から端まで通し、観測結果を捕獲する。
// provider は常に決定的 fake（RecordingProvider）で固定し、AI 出力の非決定性を排す。
func Run(cfg RunConfig) (Capture, error) {
	s, err := store.Open(cfg.DBPath)
	if err != nil {
		return Capture{}, fmt.Errorf("中心 DB を開けない: %w", err)
	}
	defer s.Close() //nolint:errcheck // 実行後の後始末。観測は別接続で行う。

	rec := &RecordingProvider{}
	eng := engine.New(s, rec, cfg.Lexicon, cfg.RoleSpeech)
	app := api.New(s, eng, rec, cfg.TermsXMLDir, cfg.Extractor)

	runResult, runErr := app.RunExtractAndTranslate(api.RunRequest{PluginPath: cfg.PluginPath, Model: cfg.Model})
	if runErr != nil {
		return Capture{}, fmt.Errorf("抽出＋翻訳の実行: %w", runErr)
	}

	db, err := sqlx.Connect("sqlite", cfg.DBPath)
	if err != nil {
		return Capture{}, fmt.Errorf("観測用 DB を開けない: %w", err)
	}
	defer db.Close() //nolint:errcheck // 観測後の後始末。
	dbState, err := captureDBState(db)
	if err != nil {
		return Capture{}, err
	}
	return Capture{Model: cfg.Model, TranslatedCount: runResult.TranslatedCount, Prompts: rec.Prompts(), DBState: dbState}, nil
}

// SyntheticRun は合成 fixture で非劣化 harness を回す。XML を xmlDir へ書き、SeedExtractor・最小辞書を組んで Run へ渡す。
// CI 恒久の回帰網（著作物を含まない自作入力）の実行点。
func SyntheticRun(dbPath, xmlDir string) (Capture, error) {
	f := SyntheticFixture()
	// ファイル名接頭を base ゲーム（Skyrim*）にして、固有名派生の姓名分割（two）経路まで通す。
	// 派生は base ゲーム由来の XML 限定のため、合成入力でもこの経路を golden で守るには接頭を合わせる。
	if err := os.WriteFile(filepath.Join(xmlDir, "Skyrim_Synthetic.xml"), []byte(f.TermsXML), 0o600); err != nil {
		return Capture{}, fmt.Errorf("合成 XML の書き出し: %w", err)
	}
	roleSpeech, err := rolespeech.ParseRoleSpeech(strings.NewReader(syntheticRoleSpeech))
	if err != nil {
		return Capture{}, fmt.Errorf("合成役割語の構築: %w", err)
	}
	return Run(RunConfig{
		DBPath:      dbPath,
		Extractor:   &SeedExtractor{DBPath: dbPath, Fixture: f},
		Lexicon:     fakeLexicon{},
		RoleSpeech:  roleSpeech,
		TermsXMLDir: xmlDir,
		PluginPath:  f.PluginName,
		Model:       "fake-model",
	})
}

// syntheticRoleSpeech は合成入力用の最小役割語表（タブ区切り 5 列、全ワイルドカード 1 行）。
// 実 assets/role-speech.tsv の内容変化から harness を切り離し、口調注入が決定的に通ることだけを保証する。
var syntheticRoleSpeech = strings.Join([]string{"*", "*", "*", "わたし", "落ち着いた言い回しにする。"}, "\t") + "\n"

// fakeLexicon は合成入力用の最小感情辞書。少数の強感情語だけを真とし、口調生成の感情経路を決定的に通す。
type fakeLexicon struct{}

// IsStrongEmotion は固定集合の語だけを強感情語とする。
func (fakeLexicon) IsStrongEmotion(word string) bool {
	switch word {
	case "hate", "fear", "kill", "rage", "afraid":
		return true
	}
	return false
}
