// Package bootstrap は composition root。store・provider を生成し、engine・api へ注入する唯一の場所。
package bootstrap

import (
	"fmt"
	"net/http"
	"time"

	"aitranslationenginejp/internal/api"
	"aitranslationenginejp/internal/engine"
	"aitranslationenginejp/internal/lexicon"
	"aitranslationenginejp/internal/provider"
	"aitranslationenginejp/internal/store"
)

// dev 既定のパス。wails dev は repo root を作業ディレクトリにするため相対パスで足りる。
const (
	devDBPath        = "db/aitranslation.dev.sqlite3"
	extractorProject = "tools/extractor"
	migrationsDir    = "db/migrations"
	// nrcDictPath は口調生成の感情辞書（研究用ライセンス）。製品化時に MIT 実装へ差し替える。
	nrcDictPath = "dictionaries/nrc-emolex.txt"
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
	lex, err := lexicon.LoadNRC(nrcDictPath)
	if err != nil {
		_ = s.Close()
		return nil, nil, fmt.Errorf("感情辞書の読み込み: %w", err)
	}

	eng := engine.New(s, p, lex)
	app := api.New(s, eng, p, api.ExtractorConfig{
		ProjectPath: extractorProject,
		SchemaDir:   migrationsDir,
		DBPath:      devDBPath,
	})
	return app, s, nil
}
