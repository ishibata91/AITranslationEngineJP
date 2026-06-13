// Package bootstrap は composition root。store・provider を生成し、engine・api へ注入する唯一の場所。
package bootstrap

import (
	"fmt"
	"net/http"
	"time"

	"aitranslationenginejp/internal/api"
	"aitranslationenginejp/internal/engine"
	"aitranslationenginejp/internal/provider"
	"aitranslationenginejp/internal/store"
)

// dev 既定のパス。wails dev は repo root を作業ディレクトリにするため相対パスで足りる。
const (
	devDBPath        = "db/aitranslation.dev.sqlite3"
	extractorProject = "tools/extractor"
	migrationsDir    = "db/migrations"
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

	eng := engine.New(s, p)
	app := api.New(s, eng, p, api.ExtractorConfig{
		ProjectPath: extractorProject,
		SchemaDir:   migrationsDir,
		DBPath:      devDBPath,
	})
	return app, s, nil
}
