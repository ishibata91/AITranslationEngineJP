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

// Store は api が結果一覧の keyset ページングとプロンプトテンプレートの CRUD に使う
// 中心データアクセスの interface（使う分だけ宣言する）。
type Store interface {
	CountNarrations(ctx context.Context) (int, error)
	CountLines(ctx context.Context) (int, error)
	NarrationsAfter(ctx context.Context, afterID int64, limit int) ([]model.Narration, error)
	LinesAfter(ctx context.Context, afterID int64, limit int) ([]model.Line, error)
	GetPromptTemplate(ctx context.Context) (model.PromptTemplate, error)
	SavePromptTemplate(ctx context.Context, t model.PromptTemplate) error
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

// TermView は結果行の機械置換内訳 1 件（原語 → 確定訳語）。
// 結果取得時に各行の原文へ辞書を当て直して再構成する（保存はしない）。
type TermView struct {
	Source string `json:"source"`
	Dest   string `json:"dest"`
}

// ResultView は結果一覧の表示用 DTO。叙述文と台詞を共通の行として表す。
// Directive と PersonaLabel は話者を解決できた台詞だけ持つ（叙述文や話者なしの台詞は空で省略）。
// Terms は本文で辞書から確定訳語へ置換した固有名の内訳。置換が無い行は空で省略する。
// Prompt は実際に翻訳 AI へ投げた完成プロンプト（base 指示＋口調指示＋機械置換済み原文）を取得時に再構成した全文。
type ResultView struct {
	EDID         string     `json:"edid"`
	Source       string     `json:"source"`
	Dest         string     `json:"dest"`
	StatusLabel  string     `json:"statusLabel"`
	Directive    string     `json:"directive,omitempty"`
	PersonaLabel string     `json:"personaLabel,omitempty"`
	Terms        []TermView `json:"terms,omitempty"`
	Prompt       string     `json:"prompt,omitempty"`
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

// PromptTemplateView はプロンプトテンプレート編集画面の表示用 DTO。
// BaseDirective は叙述文・台詞の両方に付く base 翻訳指示文、PersonaTemplate は話者のいる台詞に付く口調指示の雛形。
type PromptTemplateView struct {
	BaseDirective   string `json:"baseDirective"`
	PersonaTemplate string `json:"personaTemplate"`
}

// GetPromptTemplate は編集画面の初期表示用に、現在保存されているプロンプトテンプレートを返す。
func (a *App) GetPromptTemplate() (PromptTemplateView, error) {
	t, err := a.store.GetPromptTemplate(a.baseCtx())
	if err != nil {
		return PromptTemplateView{}, fmt.Errorf("プロンプトテンプレートの取得: %w", err)
	}
	return PromptTemplateView{BaseDirective: t.BaseDirective, PersonaTemplate: t.PersonaTemplate}, nil
}

// SavePromptTemplate は編集したプロンプトテンプレートを保存する。
// 保存後の翻訳実行と、結果行の実プロンプト再構成は、この保存値を読んで反映する。
func (a *App) SavePromptTemplate(req PromptTemplateView) error {
	if err := a.store.SavePromptTemplate(a.baseCtx(), model.PromptTemplate{
		BaseDirective:   req.BaseDirective,
		PersonaTemplate: req.PersonaTemplate,
	}); err != nil {
		return fmt.Errorf("プロンプトテンプレートの保存: %w", err)
	}
	return nil
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
// 各行の原文へ機械置換辞書を当て直し、置換内訳（terms）を再構成して結果行へ供給する。
func (a *App) buildResultsPage(ctx context.Context, cursor string, limit int) (ResultPage, error) {
	narrations, lines, total, nextCursor, hasMore, err := a.pageRows(ctx, cursor, limit)
	if err != nil {
		return ResultPage{}, err
	}

	// プロンプトテンプレート（base 指示・口調指示の雛形）をページ単位で 1 度だけ読む。
	// 実行時と同じ雛形で実プロンプトを再構成し、保存済みテンプレートを参照へ反映する。
	tmpl, err := a.store.GetPromptTemplate(ctx)
	if err != nil {
		return ResultPage{}, fmt.Errorf("プロンプトテンプレートの取得: %w", err)
	}

	// ページの台詞ぶんだけ口調を 1 度に一括生成する（台詞ごとの DB 問い合わせを避ける）。
	personas, err := a.engine.LinePersonas(ctx, lineIDs(lines), tmpl.PersonaTemplate)
	if err != nil {
		return ResultPage{}, fmt.Errorf("口調の一括生成: %w", err)
	}

	// 機械置換辞書をページ単位で 1 度だけ組み、各行の原文へ当て直して内訳を再構成する。
	dict, err := a.engine.LoadDictionary(ctx)
	if err != nil {
		return ResultPage{}, fmt.Errorf("機械置換辞書の構築: %w", err)
	}

	views := make([]ResultView, 0, len(narrations)+len(lines))
	for _, n := range narrations {
		view := narrationResultView(n)
		// 原文へ辞書を当て直し、置換内訳（terms）と、置換済み原文から組んだ実プロンプトを再構成する。
		replaced, used := dict.Apply(n.Source)
		view.Terms = termViews(used)
		// 叙述文は口調指示なし。実行時と同じ構築関数で完成プロンプトを組み、表示用全文へ描く。
		view.Prompt = engine.RenderPrompt(engine.ComposePrompt(tmpl.BaseDirective, "", replaced))
		views = append(views, view)
	}
	for _, l := range lines {
		p := personas[l.ID]
		replaced, used := dict.Apply(l.Source)
		views = append(views, ResultView{
			EDID:         l.EDID,
			Source:       l.Source,
			Dest:         l.Dest,
			StatusLabel:  statusLabel(l.Status),
			Directive:    p.Directive,
			PersonaLabel: p.Label,
			Terms:        termViews(used),
			// 台詞は話者の口調指示を合成した実プロンプトを再構成する（口調指示の合成を目視で確かめる）。
			Prompt: engine.RenderPrompt(engine.ComposePrompt(tmpl.BaseDirective, p.Directive, replaced)),
		})
	}
	return ResultPage{Total: total, Results: views, NextCursor: nextCursor, HasMore: hasMore}, nil
}

// termViews は辞書の置換内訳（engine.DictionaryTerm）を結果行の表示 DTO へ写す。置換が無ければ nil を返す。
func termViews(used []engine.DictionaryTerm) []TermView {
	if len(used) == 0 {
		return nil
	}
	out := make([]TermView, len(used))
	for i, t := range used {
		out[i] = TermView{Source: t.Source, Dest: t.Dest}
	}
	return out
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
