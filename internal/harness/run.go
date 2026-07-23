package harness

import (
	"fmt"
	"strings"

	"aitranslationenginejp/internal/api"
	"aitranslationenginejp/internal/core/dictionary"
	"aitranslationenginejp/internal/core/linefeatures"
	"aitranslationenginejp/internal/core/rolespeech"
	"aitranslationenginejp/internal/engine"
	"aitranslationenginejp/internal/store"

	"github.com/jmoiron/sqlx"
)

// RunConfig は非劣化 harness の 1 回の実行に要る注入物。
// 抽出子・感情辞書・役割語・plugin パスを差し替え可能にし、
// 合成入力（SeedExtractor・最小辞書）と実 .esm 入力（dotnet 抽出子・実辞書）の両方を同じ Run で回す。
type RunConfig struct {
	DBPath     string        // 中心 DB（temp ファイル）のパス
	Extractor  api.Extractor // 抽出段の注入（合成は SeedExtractor、実データは api.DotnetExtractor）
	Lexicon    linefeatures.EmotionLexicon
	RoleSpeech *rolespeech.Table
	Stoplist   *dictionary.Stoplist // 機械置換辞書・言及語彙の供給から一般語を除く選別（nil なら選別なし）
	PluginPath string               // 抽出対象 plugin のパス
	Model      string               // 送信モデル名（fake provider は記録のみ）
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
	eng := engine.New(s, rec, cfg.Lexicon, cfg.RoleSpeech, cfg.Stoplist)
	// harness は同期経路（RunExtractAndTranslate）だけを通す結合テスト基盤のため、batch runner は nil。
	app := api.New(s, eng, nil, rec, cfg.Extractor)

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

// SyntheticRun は合成 fixture で非劣化 harness を回す。SeedExtractor・最小辞書を組んで Run へ渡す。
// 固有名の確定訳語・参照訳の供給は fixture の英日対（SeedField.Dest）が担う（XML 直読みは廃止）。
// CI 恒久の回帰網（著作物を含まない自作入力）の実行点。
func SyntheticRun(dbPath string) (Capture, error) {
	f := SyntheticFixture()
	roleSpeech, err := rolespeech.ParseRoleSpeech(strings.NewReader(syntheticRoleSpeech))
	if err != nil {
		return Capture{}, fmt.Errorf("合成役割語の構築: %w", err)
	}
	stop, err := dictionary.ParseStoplist(strings.NewReader(syntheticStopwords))
	if err != nil {
		return Capture{}, fmt.Errorf("合成 stoplist の構築: %w", err)
	}
	return Run(RunConfig{
		DBPath:     dbPath,
		Extractor:  &SeedExtractor{DBPath: dbPath, Fixture: f},
		Lexicon:    fakeLexicon{},
		RoleSpeech: roleSpeech,
		Stoplist:   stop,
		PluginPath: f.PluginName,
		Model:      "fake-model",
	})
}

// syntheticRoleSpeech は合成入力用の最小役割語表（タブ区切り 5 列、全ワイルドカード 1 行）。
// 実 assets/role-speech.tsv の内容変化から harness を切り離し、口調注入が決定的に通ることだけを保証する。
var syntheticRoleSpeech = strings.Join([]string{"*", "*", "*", "わたし", "落ち着いた言い回しにする。"}, "\t") + "\n"

// syntheticStopwords は合成入力用の最小一般語リスト。実 assets/stopwords-en.txt の内容変化から
// harness を切り離し、fixture の管理用文字列（Yes・No）が置換も言及もされないことだけを golden で凍結する。
const syntheticStopwords = "yes\nno\n"

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
