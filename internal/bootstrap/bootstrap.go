// Package bootstrap は composition root。store・provider を生成し、engine・api へ注入する唯一の場所。
package bootstrap

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"aitranslationenginejp/internal/api"
	"aitranslationenginejp/internal/core/dictionary"
	"aitranslationenginejp/internal/core/rolespeech"
	"aitranslationenginejp/internal/engine"
	"aitranslationenginejp/internal/lexicon"
	"aitranslationenginejp/internal/provider"
	"aitranslationenginejp/internal/store"
)

// dev 既定のパス。wails dev は repo root を作業ディレクトリにするため相対パスで足りる。
const (
	devDBPath = "db/aitranslation.dev.sqlite3"
	// extractorDLL は起動前に publish 済みの抽出子 DLL。dev 起動 script（scripts/dev/run-wails.sh）が
	// wails dev 起動の前に dotnet publish で生成する。抽出のたびに dotnet run で MSBuild を評価しない。
	extractorDLL  = "tools/extractor/bin/publish/extractor.dll"
	migrationsDir = "db/migrations"
	// vaderDictPath は口調生成の感情辞書（VADER lexicon、MIT ライセンス）。git 同梱可のため assets に置く。
	// 出典・checksum は併置の assets/vader_lexicon.LICENSE に記録する。
	vaderDictPath = "assets/vader_lexicon.txt"
	// roleSpeechPath は注入時に引く一人称・語尾テンプレート。中身は実画面確認で見直すため外部ファイルに置く。
	roleSpeechPath = "assets/role-speech.tsv"
	// roleSpeechExamplePath は口調の例文（英語原文と日本語訳文の 1 対）。役割語と同じキーで引く。
	// 説明文だけでは同じ話者の台詞ごとに一人称・語尾が揺れるため、訳例で固定する。
	roleSpeechExamplePath = "assets/role-speech-examples.tsv"
	// stopwordsPath は機械置換辞書・言及語彙の供給から除く一般語リスト（stopwords-iso 配布、MIT）。
	// 出典・checksum は併置の assets/stopwords-en.LICENSE に記録する。
	stopwordsPath = "assets/stopwords-en.txt"
)

// NewApp は中心 DB を開き、全層を配線して api.App を返す。Close 用に store も返す。
func NewApp() (*api.App, *store.Store, error) {
	s, err := store.Open(devDBPath)
	if err != nil {
		return nil, nil, fmt.Errorf("中心 DB を開けない: %w", err)
	}

	// 翻訳は本文が長く時間がかかるため、HTTP の per-request timeout を長めに取る。
	client := &http.Client{Timeout: 10 * time.Minute}
	p := provider.NewOpenAICompatible(client)

	// 口調生成の感情辞書を読む。差し替え可能な境界（engine.EmotionLexicon）の concrete 実装。
	lex, err := lexicon.LoadVADER(vaderDictPath)
	if err != nil {
		_ = s.Close()
		return nil, nil, fmt.Errorf("感情辞書の読み込み: %w", err)
	}

	// 注入時に引く一人称・語尾テンプレートを読む。性別・年齢・基底口調で役割語を引く参照データ。
	// asset のファイル open は composition root の責務。純粋な解析は core/rolespeech が行う。
	rolesFile, err := os.Open(roleSpeechPath)
	if err != nil {
		_ = s.Close()
		return nil, nil, fmt.Errorf("役割語テンプレートを開けない (%s): %w", roleSpeechPath, err)
	}
	defer rolesFile.Close() //nolint:errcheck // 読み取り後の後始末。
	roles, err := rolespeech.ParseRoleSpeech(rolesFile)
	if err != nil {
		_ = s.Close()
		return nil, nil, fmt.Errorf("役割語テンプレートの読み込み: %w", err)
	}

	// 6列の口調例文を読み、役割語とは独立した4キーの行集合として同じ Table へ束ねる。
	exampleFile, err := os.Open(roleSpeechExamplePath)
	if err != nil {
		_ = s.Close()
		return nil, nil, fmt.Errorf("口調例文テンプレートを開けない (%s): %w", roleSpeechExamplePath, err)
	}
	defer exampleFile.Close() //nolint:errcheck // 読み取り後の後始末。
	roles, err = rolespeech.ParseRoleSpeechExamples(roles, exampleFile)
	if err != nil {
		_ = s.Close()
		return nil, nil, fmt.Errorf("口調例文テンプレートの読み込み: %w", err)
	}

	// 一般語 stoplist を読む。固有名の供給が本文の通常語（Yes・No 等）へ誤って当たるのを供給側で止める。
	stopFile, err := os.Open(stopwordsPath)
	if err != nil {
		_ = s.Close()
		return nil, nil, fmt.Errorf("一般語 stoplist を開けない (%s): %w", stopwordsPath, err)
	}
	defer stopFile.Close() //nolint:errcheck // 読み取り後の後始末。
	stop, err := dictionary.ParseStoplist(stopFile)
	if err != nil {
		_ = s.Close()
		return nil, nil, fmt.Errorf("一般語 stoplist の読み込み: %w", err)
	}

	eng := engine.New(s, p, lex, roles, stop)
	// OpenAI と xAI の batch クライアントを提供元に対応付け、送信・反映のオーケストレーションへ配線する。
	// batch は同じ長め timeout の HTTP client を使う（file upload・結果取得も時間がかかりうる）。
	openAIBatch := provider.NewOpenAIBatch(client)
	xaiBatch := provider.NewXAIBatch(client)
	batchRunner := engine.NewBatchRunner(eng, map[string]provider.BatchTranslator{
		provider.BatchProviderOpenAI: openAIBatch,
		provider.BatchProviderXAI:    xaiBatch,
	}, s)
	ext := api.ExtractorConfig{
		DLLPath:   extractorDLL,
		SchemaDir: migrationsDir,
		DBPath:    devDBPath,
	}
	// 抽出子は本番の dotnet 子プロセス起動。composition root が concrete を生成して注入する唯一の場所。
	// 抽出に要するパスは DotnetExtractor が保持する。
	app := api.New(s, eng, batchRunner, p, api.NewDotnetExtractor(ext))
	return app, s, nil
}
