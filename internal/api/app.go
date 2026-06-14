// Package api は Wails Bind の公開面を持つ。翻訳ジョブの起動、抽出子プロセスの起動、結果取得を frontend へ公開する。
package api

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"aitranslationenginejp/internal/engine"
	"aitranslationenginejp/internal/model"
	"aitranslationenginejp/internal/provider"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// progressEventName は本文翻訳の進捗を frontend へ push する runtime event 名。
const progressEventName = "translation:progress"

// defaultPageLimit は ListResultsPage の limit が未指定（0 以下）のときの既定ページ件数。
const defaultPageLimit = 50

// Store は api が結果一覧の keyset ページングに使う中心データアクセスの interface（使う分だけ宣言する）。
type Store interface {
	CountNarrations(ctx context.Context) (int, error)
	CountLines(ctx context.Context) (int, error)
	NarrationsAfter(ctx context.Context, afterID int64, limit int) ([]model.Narration, error)
	LinesAfter(ctx context.Context, afterID int64, limit int) ([]model.Line, error)
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

// RunResult は実行結果の要約。結果一覧は数万件になりうるためここでは返さず、
// frontend が ListResultsPage で先頭ページから取得する。
type RunResult struct {
	TranslatedCount int `json:"translatedCount"`
}

// ResultPage は結果一覧の keyset cursor ページ。Total は叙述文＋台詞の総件数。
// NextCursor は次ページ取得用の cursor、HasMore は次ページの有無。
type ResultPage struct {
	Total      int          `json:"total"`
	Results    []ResultView `json:"results"`
	NextCursor string       `json:"nextCursor"`
	HasMore    bool         `json:"hasMore"`
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

// ListResultsPage は中心 DB の叙述文と台詞を keyset cursor ページで返す（起動時・ページ送り・実行後の取得を統一）。
// cursor は ""（先頭）/ "n:<id>"（叙述文区間）/ "l:<id>"（台詞区間）。limit が 0 以下なら既定件数を使う。
func (a *App) ListResultsPage(cursor string, limit int) (ResultPage, error) {
	if limit <= 0 {
		limit = defaultPageLimit
	}
	return a.buildResultsPage(a.baseCtx(), cursor, limit)
}

// RunExtractAndTranslate は plugin を抽出し、未訳の叙述文と台詞を翻訳し、翻訳件数の要約を返す。
// 結果一覧は数万件になりうるためここでは返さず、frontend が ListResultsPage で取得する。
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

	result := RunResult{TranslatedCount: count}
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

// cursor の区間種別。叙述文（id 昇順）を先頭に、台詞（id 昇順）を続けた連結列上の位置を表す。
const (
	cursorNarration = "n"
	cursorLine      = "l"
)

// parseCursor は cursor を区間種別と afterID へ分解する。""（先頭）や不正な cursor は叙述文先頭（id 0）とみなす。
func parseCursor(cursor string) (section string, afterID int64) {
	parts := strings.SplitN(cursor, ":", 2)
	if len(parts) == 2 && (parts[0] == cursorNarration || parts[0] == cursorLine) {
		if id, err := strconv.ParseInt(parts[1], 10, 64); err == nil {
			return parts[0], id
		}
	}
	return cursorNarration, 0
}

// makeCursor は区間種別と id を cursor 文字列へ組む。
func makeCursor(section string, id int64) string {
	return section + ":" + strconv.FormatInt(id, 10)
}

// buildResultsPage は cursor の指すページの叙述文・台詞を取り、台詞へ口調を一括で付けて ResultPage を返す。
func (a *App) buildResultsPage(ctx context.Context, cursor string, limit int) (ResultPage, error) {
	narrations, lines, total, nextCursor, hasMore, err := a.pageRows(ctx, cursor, limit)
	if err != nil {
		return ResultPage{}, err
	}

	// ページの台詞ぶんだけ口調を 1 度に一括生成する（台詞ごとの DB 問い合わせを避ける）。
	personas, err := a.engine.LinePersonas(ctx, lineIDs(lines))
	if err != nil {
		return ResultPage{}, fmt.Errorf("口調の一括生成: %w", err)
	}

	views := make([]ResultView, 0, len(narrations)+len(lines))
	for _, n := range narrations {
		views = append(views, narrationResultView(n))
	}
	for _, l := range lines {
		p := personas[l.ID]
		views = append(views, ResultView{
			EDID:         l.EDID,
			Source:       l.Source,
			Dest:         l.Dest,
			StatusLabel:  statusLabel(l.Status),
			Directive:    p.Directive,
			PersonaLabel: p.Label,
		})
	}
	return ResultPage{Total: total, Results: views, NextCursor: nextCursor, HasMore: hasMore}, nil
}

// pageRows は keyset cursor から当該ページの叙述文・台詞、総件数、次 cursor、続きの有無を決める。
// 叙述文区間でページに満たなければ台詞を先頭から補充して台詞区間へ移る。口調は付けない（呼び出し側が付ける）。
func (a *App) pageRows(ctx context.Context, cursor string, limit int) (narrations []model.Narration, lines []model.Line, total int, nextCursor string, hasMore bool, err error) {
	nTotal, err := a.store.CountNarrations(ctx)
	if err != nil {
		return nil, nil, 0, "", false, fmt.Errorf("叙述文の件数: %w", err)
	}
	lTotal, err := a.store.CountLines(ctx)
	if err != nil {
		return nil, nil, 0, "", false, fmt.Errorf("台詞の件数: %w", err)
	}
	total = nTotal + lTotal

	section, afterID := parseCursor(cursor)

	// 台詞区間: 台詞だけを afterID から取る。
	if section == cursorLine {
		lines, err = a.store.LinesAfter(ctx, afterID, limit+1)
		if err != nil {
			return nil, nil, 0, "", false, fmt.Errorf("台詞ページの取得: %w", err)
		}
		hasMore = len(lines) > limit
		if hasMore {
			lines = lines[:limit]
		}
		if len(lines) > 0 {
			nextCursor = makeCursor(cursorLine, lines[len(lines)-1].ID)
		}
		return nil, lines, total, nextCursor, hasMore, nil
	}

	// 叙述文区間: 叙述文で埋め、足りなければ台詞先頭から補充する。
	narrations, err = a.store.NarrationsAfter(ctx, afterID, limit+1)
	if err != nil {
		return nil, nil, 0, "", false, fmt.Errorf("叙述文ページの取得: %w", err)
	}
	if len(narrations) > limit {
		narrations = narrations[:limit]
		return narrations, nil, total, makeCursor(cursorNarration, narrations[len(narrations)-1].ID), true, nil
	}

	remaining := limit - len(narrations)
	lines, err = a.store.LinesAfter(ctx, 0, remaining+1)
	if err != nil {
		return nil, nil, 0, "", false, fmt.Errorf("台詞ページの取得: %w", err)
	}
	lineHasMore := len(lines) > remaining
	if lineHasMore {
		lines = lines[:remaining]
	}
	switch {
	case len(lines) > 0:
		return narrations, lines, total, makeCursor(cursorLine, lines[len(lines)-1].ID), lineHasMore, nil
	case lineHasMore:
		// 叙述文がちょうどページを埋め、台詞が残る。次ページは台詞先頭（afterID 0）から。
		return narrations, nil, total, makeCursor(cursorLine, 0), true, nil
	default:
		return narrations, nil, total, "", false, nil
	}
}

// lineIDs は台詞行から id 列を取り出す（口調の一括生成の入力用）。
func lineIDs(lines []model.Line) []int64 {
	ids := make([]int64, len(lines))
	for i, l := range lines {
		ids[i] = l.ID
	}
	return ids
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
