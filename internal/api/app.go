// Package api は Wails Bind の公開面を持つ。翻訳ジョブの起動、抽出子プロセスの起動、結果取得を frontend へ公開する。
package api

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"

	"aitranslationenginejp/internal/engine"
	"aitranslationenginejp/internal/model"
	"aitranslationenginejp/internal/provider"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// progressEventName は本文翻訳の進捗を frontend へ push する runtime event 名。
const progressEventName = "translation:progress"

// Store は api が結果一覧に使う中心データアクセスの interface（使う分だけ宣言する）。
type Store interface {
	ListNarrations(ctx context.Context) ([]model.Narration, error)
	ListLines(ctx context.Context) ([]model.Line, error)
}

// ExtractorConfig は抽出子プロセス（C#）の起動設定。
type ExtractorConfig struct {
	ProjectPath string // dotnet run --project のパス（dev: tools/extractor）
	SchemaDir   string // db/migrations のパス
	DBPath      string // 中心 DB のパス。store と同じファイルを指す。
}

// App は Wails へ Bind する公開面。
type App struct {
	store    Store
	engine   *engine.Engine
	provider provider.Translator
	ext      ExtractorConfig
	ctx      context.Context
}

// New は App を生成する。
func New(store Store, eng *engine.Engine, p provider.Translator, ext ExtractorConfig) *App {
	return &App{store: store, engine: eng, provider: p, ext: ext}
}

// Startup は Wails 起動時に runtime context を受け取る。ファイルダイアログと進捗 event の push に使う。
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) baseCtx() context.Context {
	if a.ctx != nil {
		return a.ctx
	}
	return context.Background()
}

// ConnRequest は AI サービスへの接続情報（画面から都度渡す）。
type ConnRequest struct {
	Endpoint string `json:"endpoint"`
	APIKey   string `json:"apiKey"`
}

// RunRequest は抽出＋翻訳の実行要求。
type RunRequest struct {
	PluginPath string `json:"pluginPath"`
	Endpoint   string `json:"endpoint"`
	APIKey     string `json:"apiKey"`
	Model      string `json:"model"`
}

// ResultView は結果一覧の表示用 DTO。叙述文と台詞を共通の行として表す。
// Directive と PersonaLabel は話者を解決できた台詞だけ持つ（叙述文や話者なしの台詞は空で省略）。
type ResultView struct {
	EDID         string `json:"edid"`
	Source       string `json:"source"`
	Dest         string `json:"dest"`
	StatusLabel  string `json:"statusLabel"`
	Directive    string `json:"directive,omitempty"`
	PersonaLabel string `json:"personaLabel,omitempty"`
}

// RunResult は実行結果。Results は叙述文に続けて台詞を並べる。
type RunResult struct {
	TranslatedCount int          `json:"translatedCount"`
	Results         []ResultView `json:"results"`
}

// ProgressEvent は本文翻訳の進捗 payload。Stage は "extract"（台詞抽出、不定）と "translate"（本文翻訳）。
type ProgressEvent struct {
	Stage string `json:"stage"`
	Done  int    `json:"done"`
	Total int    `json:"total"`
}

// GetModels は接続先の利用可能モデル一覧を返す（画面のモデル選択用）。
func (a *App) GetModels(req ConnRequest) ([]string, error) {
	models, err := a.provider.ListModels(a.baseCtx(), provider.Connection{Endpoint: req.Endpoint, APIKey: req.APIKey})
	if err != nil {
		return nil, fmt.Errorf("モデル一覧の取得: %w", err)
	}
	return models, nil
}

// SelectPluginFile は plugin ファイル選択ダイアログを開き、選んだフルパスを返す。
func (a *App) SelectPluginFile() (string, error) {
	path, err := wailsruntime.OpenFileDialog(a.baseCtx(), wailsruntime.OpenDialogOptions{
		Title: "plugin を選択",
		Filters: []wailsruntime.FileFilter{
			{DisplayName: "Skyrim plugin (*.esp;*.esm;*.esl)", Pattern: "*.esp;*.esm;*.esl"},
		},
	})
	if err != nil {
		return "", fmt.Errorf("ファイル選択: %w", err)
	}
	return path, nil
}

// ListResults は中心 DB の叙述文と台詞を結果一覧として返す（起動時の前回結果表示用）。
func (a *App) ListResults() ([]ResultView, error) {
	return a.buildResults(a.baseCtx())
}

// RunExtractAndTranslate は plugin を抽出し、未訳の叙述文と台詞を翻訳し、結果一覧を返す。
// 抽出中は extract 進捗を、本文翻訳中は translate 進捗を runtime event で push する。
func (a *App) RunExtractAndTranslate(req RunRequest) (RunResult, error) {
	ctx := a.baseCtx()
	dataFolder := filepath.Dir(req.PluginPath)
	pluginName := filepath.Base(req.PluginPath)

	a.emitProgress(ProgressEvent{Stage: "extract"})
	args := buildExtractorArgs(a.ext.ProjectPath, dataFolder, pluginName, a.ext.DBPath, a.ext.SchemaDir)
	// dotnet は固定コマンド、引数は内部生成のパスのみ。利用者が選んだ plugin を抽出するための意図的な子プロセス起動。
	out, err := exec.CommandContext(ctx, "dotnet", args...).CombinedOutput() //nolint:gosec // 固定コマンド dotnet・内部生成引数
	if err != nil {
		return RunResult{}, fmt.Errorf("抽出に失敗: %w: %s", err, string(out))
	}

	count, runErr := a.engine.Run(ctx, provider.Connection{Endpoint: req.Endpoint, APIKey: req.APIKey}, req.Model,
		func(done, total int) {
			a.emitProgress(ProgressEvent{Stage: "translate", Done: done, Total: total})
		})

	results, err := a.buildResults(ctx)
	if err != nil {
		return RunResult{}, err
	}
	result := RunResult{TranslatedCount: count, Results: results}
	if runErr != nil {
		return result, fmt.Errorf("翻訳に失敗: %w", runErr)
	}
	return result, nil
}

// emitProgress は進捗を runtime event で frontend へ push する。runtime context が無い（テスト等）なら何もしない。
func (a *App) emitProgress(ev ProgressEvent) {
	if a.ctx == nil {
		return
	}
	wailsruntime.EventsEmit(a.ctx, progressEventName, ev)
}

// buildResults は叙述文と台詞を結果一覧の行へ変換する。台詞は話者属性から口調指示文と短い要約を付ける。
func (a *App) buildResults(ctx context.Context) ([]ResultView, error) {
	narrations, err := a.store.ListNarrations(ctx)
	if err != nil {
		return nil, fmt.Errorf("叙述文の取得: %w", err)
	}
	lines, err := a.store.ListLines(ctx)
	if err != nil {
		return nil, fmt.Errorf("台詞の取得: %w", err)
	}

	views := make([]ResultView, 0, len(narrations)+len(lines))
	for _, n := range narrations {
		views = append(views, narrationResultView(n))
	}
	for _, l := range lines {
		directive, label, err := a.engine.LineDirective(ctx, l.ID)
		if err != nil {
			return nil, fmt.Errorf("台詞の口調指示の組み立て: %w", err)
		}
		views = append(views, ResultView{
			EDID:         l.EDID,
			Source:       l.Source,
			Dest:         l.Dest,
			StatusLabel:  statusLabel(l.Status),
			Directive:    directive,
			PersonaLabel: label,
		})
	}
	return views, nil
}

// buildExtractorArgs は dotnet run で extractor を起動する引数列を組む。
func buildExtractorArgs(projectPath, dataFolder, plugin, dbPath, schemaDir string) []string {
	return []string{
		"run", "--project", projectPath, "--",
		"--data", dataFolder,
		"--plugin", plugin,
		"--sqlite", dbPath,
		"--schema", schemaDir,
	}
}

// narrationResultView は叙述文を結果一覧の行へ写す。叙述文は話者を持たないため口調指示は付けない。
func narrationResultView(n model.Narration) ResultView {
	return ResultView{
		EDID:        n.EDID,
		Source:      n.Source,
		Dest:        n.Dest,
		StatusLabel: statusLabel(n.Status),
	}
}

// statusLabel は xTranslator の訳状態コードを表示ラベルへ写す。
func statusLabel(status int) string {
	switch status {
	case 1:
		return "訳済"
	case 2:
		return "部分"
	case 3:
		return "仮訳"
	case 4:
		return "承認"
	default:
		return "未訳"
	}
}
