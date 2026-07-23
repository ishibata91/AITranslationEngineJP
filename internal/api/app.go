// Package api は Wails Bind の公開面を持つ。翻訳ジョブの起動、抽出子プロセスの起動、結果取得を frontend へ公開する。
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"aitranslationenginejp/internal/core/dictionary"
	"aitranslationenginejp/internal/core/prompt"
	"aitranslationenginejp/internal/core/runtimetag"
	"aitranslationenginejp/internal/engine"
	"aitranslationenginejp/internal/model"
	"aitranslationenginejp/internal/provider"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// progressEventName は本文翻訳の進捗を frontend へ push する runtime event 名。
const progressEventName = "translation:progress"

// defaultPageLimit は ListResultsPage の limit が未指定（0 以下）のときの既定ページ件数。
const defaultPageLimit = 50

// Store は api が結果一覧の keyset ページング、プロンプトテンプレートの CRUD、
// 指示文（directive）の編集と REC:FIELD 種別の参照に使う中心データアクセスの interface（使う分だけ宣言する）。
type Store interface {
	CountNarrations(ctx context.Context, plugin string) (int, error)
	CountLines(ctx context.Context, plugin string) (int, error)
	CountProperNouns(ctx context.Context, plugin string) (int, error)
	NarrationsAfter(ctx context.Context, plugin string, afterID int64, limit int) ([]model.Narration, error)
	LinesAfter(ctx context.Context, plugin string, afterID int64, limit int) ([]model.Line, error)
	ProperNounsAfter(ctx context.Context, plugin string, afterID int64, limit int) ([]model.ProperNoun, error)
	GetPromptTemplate(ctx context.Context) (model.PromptTemplate, error)
	SavePromptTemplate(ctx context.Context, t model.PromptTemplate) error
	LoadLineSpeakers(ctx context.Context, lineIDs []int64) (map[int64]model.LineSpeaker, error)
	ListDirectives(ctx context.Context) ([]model.Directive, error)
	SaveDirectiveInstruction(ctx context.Context, key, instruction string) error
	ListRecordTypeMaster(ctx context.Context) ([]model.RecordType, error)
	UpsertTargetPlugin(ctx context.Context, plugin, sourcePath string) error
	ListTargetPlugins(ctx context.Context) ([]model.TargetPlugin, error)
	DeleteTargetPlugin(ctx context.Context, plugin string) error
}

// ExtractorConfig は抽出子プロセス（C#）の起動設定。dotnet 起動に要するパスだけを持つ。
// 固有名派生で読む XML ディレクトリは抽出の入力ではないため、ここに含めず api.New が直接受け取る。
type ExtractorConfig struct {
	// DLLPath は事前 publish 済みの抽出子 DLL のパス（dev: tools/extractor/bin/publish/extractor.dll）。
	// 抽出のたびに MSBuild を評価する dotnet run を避け、ビルド済み DLL を dotnet で直接実行する。
	// ビルドは dev 起動 script（scripts/dev/run-wails.sh）が起動前に 1 度だけ行う。
	DLLPath   string
	SchemaDir string // db/migrations のパス
	DBPath    string // 中心 DB のパス。store と同じファイルを指す。
}

// Extractor は plugin を抽出し、原文を extracted_field（と話者 staging）へ書き込む境界。
// 本番は dotnet 子プロセス起動（DotnetExtractor）、テストは DB を直接 seed する fake を注入する。
// これは provider のような多態 port ではなく、子プロセス起動を呼び出し側から切り離してテスト可能にするための
// consumer 側 interface。api コンポーネント内に閉じ、.go-arch-lint.yml に新コンポーネントや新依存規則を足さない。
type Extractor interface {
	// Extract は pluginPath の plugin を抽出し、結果を中心 DB へ書き込む。書き込み先 DB は実装が保持する。
	Extract(ctx context.Context, pluginPath string) error
}

// DotnetExtractor は C# 抽出器（Mutagen）を publish 済み DLL の dotnet 直実行で起動する本番 Extractor。
// 起動に要するパス（DLL・schema・出力先 DB）は ExtractorConfig から受け取って保持する。
type DotnetExtractor struct {
	ext ExtractorConfig
}

// NewDotnetExtractor は本番の dotnet 抽出子を生成する。bootstrap が唯一の生成点。
func NewDotnetExtractor(ext ExtractorConfig) *DotnetExtractor {
	return &DotnetExtractor{ext: ext}
}

// Extract は publish 済み DLL を dotnet で直接実行し、pluginPath の plugin を抽出して extracted_field へ書き込む。
// DLL は起動前に build 済みである前提で、抽出のたびに MSBuild を通さない。DLL が無ければ無言でビルドせずエラーを返す。
func (d *DotnetExtractor) Extract(ctx context.Context, pluginPath string) error {
	// 抽出時にビルドを起こさない方針のため、DLL 不在はここで「未ビルド」として明示的に落とす。
	if _, err := os.Stat(d.ext.DLLPath); err != nil {
		return fmt.Errorf("抽出子が未ビルド（%s）。dev 起動 script が publish していない可能性がある: %w", d.ext.DLLPath, err)
	}
	dataFolder := filepath.Dir(pluginPath)
	pluginName := filepath.Base(pluginPath)
	args := buildExtractorArgs(d.ext.DLLPath, dataFolder, pluginName, d.ext.DBPath, d.ext.SchemaDir)
	// dotnet は固定コマンド、引数は内部生成のパスのみ。利用者が選んだ plugin を抽出するための意図的な子プロセス起動。
	out, err := exec.CommandContext(ctx, "dotnet", args...).CombinedOutput() //nolint:gosec // 固定コマンド dotnet・内部生成引数
	if err != nil {
		// 元実装と同じく dotnet の stdout/stderr をエラーへ含める。接頭は付けず、呼び出し側（RunExtractAndTranslate）が
		// "抽出に失敗" を付ける。最終メッセージは "抽出に失敗: <exec error>: <出力>" でリファクタ前と一致する。
		return fmt.Errorf("%w: %s", err, string(out))
	}
	return nil
}

// App は Wails へ Bind する公開面。
type App struct {
	store     Store
	engine    *engine.Engine
	batch     *engine.BatchRunner // xAI batch 翻訳の送信・反映。batch を使わない配線では nil。
	provider  provider.Translator
	extractor Extractor
	ctx       context.Context
}

// New は App を生成する。extractor は抽出子の注入点（本番は DotnetExtractor。抽出に要するパスは extractor が保持する）。
// batch は xAI batch 翻訳のオーケストレーション。batch を使わない配線（テスト用 harness など）では nil を渡してよい。
func New(store Store, eng *engine.Engine, batch *engine.BatchRunner, p provider.Translator, extractor Extractor) *App {
	if extractor == nil {
		// 抽出子は必須の注入物。nil interface（リテラル nil や未設定の interface 変数）を渡す配線ミスを起動時に弾く。
		// なお typed-nil（nil の concrete ポインタを interface に入れた値）は Go の制約で検出できないが、
		// composition root は concrete を直接 new して渡すため該当しない。
		panic("api.New: extractor は nil にできない")
	}
	return &App{store: store, engine: eng, batch: batch, provider: p, extractor: extractor}
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

// BatchPluginRequest は plugin 単位の batch 操作（状態確認・取り込み）の要求。
// Plugin は対象 plugin ファイル名（narration.plugin と同値・filepath.Base）。接続情報は都度受ける（永続化しない）。
type BatchPluginRequest struct {
	Plugin   string `json:"plugin"`
	Endpoint string `json:"endpoint"`
	APIKey   string `json:"apiKey"`
}

// BatchProgressView は xAI batch の進行状況（状態確認の応答）。Present=false は進行なし。
// Stage は "proper" / "body" / "done"。件数は現段 batch 由来。CanApply は取り込める完了段があるか。
type BatchProgressView struct {
	Present   bool   `json:"present"`
	Stage     string `json:"stage"`
	Total     int    `json:"total"`
	Pending   int    `json:"pending"`
	Succeeded int    `json:"succeeded"`
	Failed    int    `json:"failed"`
	CanApply  bool   `json:"canApply"`
}

// TermView は結果行の機械置換内訳 1 件（原語 → 確定訳語）。
// 結果取得時に各行の原文へ辞書を当て直して再構成する（保存はしない）。
type TermView struct {
	Source string `json:"source"`
	Dest   string `json:"dest"`
}

// PersonaView は結果行の口調メタデータ。口調が付いた台詞の判定結果と根拠。
// 結果行を展開したときに、判定結果（Cell・Trait）を強調し、根拠（段階・印・経路・性別）を小さく出す。
// 口調が付いた台詞だけ持つ（叙述文や口調なしの台詞は nil で省略）。
// 名指し話者は Cell・AttitudeBand・Marked を持つ。汎用台詞・PC 発話はそれらが空で、EmotionBand・Sex を持つ。
type PersonaView struct {
	Cell         string `json:"cell"`         // 基底口調セル名（判定結果）。汎用・PC は空
	Trait        string `json:"trait"`        // 基底口調の性質文（口調を普通の言葉で説明した一文）。汎用・PC は空
	AttitudeBand int    `json:"attitudeBand"` // 対人段階 0尊大/1中立/2丁寧。汎用・PC は持たない
	EmotionBand  int    `json:"emotionBand"`  // 感情段階 0抑制/1中/2激情
	Marked       int    `json:"marked"`       // 印（信頼度の目安）。汎用・PC は持たない
	DecisionPath string `json:"decisionPath"` // 本文/voice/保留/汎用/PC
	Sex          string `json:"sex"`          // 性別ラベル（女性/男性/空）。汎用・PC の一人称・語尾の根拠
}

// SpeakerView は結果行の話者識別と属性。誰の台詞かと、口調指示の根拠（性別・年齢・声型）を結果詳細で出す。
// 話者を解決できた台詞だけ持つ（叙述文や話者なしの台詞は nil で省略）。
type SpeakerView struct {
	EDID  string `json:"edid"`  // 話者の EditorID（例 AventusAretino）
	Sex   string `json:"sex"`   // 性別ラベル（女性/男性/空）
	Age   string `json:"age"`   // 年齢区分ラベル（老人/子供/成人/空）
	Voice string `json:"voice"` // 声型 EditorID（例 FemaleOldGrumpy。役割語でなく対人 prior の根拠）
}

// RecordTypeView は結果行の元レコード種別バッジ（箱 ・ REC:FIELD）。
// Box は概念モデルの箱（叙述文/固有名/定型句/台詞）、RecField は "REC:FIELD"（固有名は種別＝rec）。
type RecordTypeView struct {
	Box      string `json:"box"`
	RecField string `json:"recField"`
}

// ResultView は結果一覧の表示用 DTO。叙述文・定型句・台詞・固有名を共通の行として表す。
// Speaker・Directive・PersonaLabel・Persona は話者を解決できた台詞だけ持つ（叙述文や話者なしの台詞は省略）。
// Terms は本文で辞書から確定訳語へ置換した固有名の内訳。置換が無い行は空で省略する。
// Prompt は実際に翻訳 AI へ投げた完成プロンプト（base 指示＋directive＋機械置換済み原文）を取得時に再構成した全文。
// RecordType は元レコード種別バッジ。
type ResultView struct {
	EDID         string          `json:"edid"`
	Source       string          `json:"source"`
	Dest         string          `json:"dest"`
	StatusLabel  string          `json:"statusLabel"`
	RecordType   *RecordTypeView `json:"recordType,omitempty"`
	Speaker      *SpeakerView    `json:"speaker,omitempty"`
	Directive    string          `json:"directive,omitempty"`
	PersonaLabel string          `json:"personaLabel,omitempty"`
	Persona      *PersonaView    `json:"persona,omitempty"`
	Terms        []TermView      `json:"terms,omitempty"`
	Prompt       string          `json:"prompt,omitempty"`
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

// ProgressEvent は本文翻訳の進捗 payload。Stage は "extract"（翻訳前区間、不定）と "translate"（本文翻訳）。
// Step は extract 段内のサブ段（"extract" 台詞抽出 / "derive" 辞書準備 / "reference" 既存訳取込 / "ingest" 取込段）で、
// 翻訳前区間の無音を無くすために段ごとの見出しを frontend が出し分ける根拠にする。translate 段では空にする。
type ProgressEvent struct {
	Stage string `json:"stage"`
	Step  string `json:"step,omitempty"`
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
// plugin が空でなければその対象 plugin（plugin ファイル名）の結果だけに絞る。空なら全 plugin。
// cursor は ""（先頭）/ "n:<id>"（叙述文区間）/ "l:<id>"（台詞区間）。limit が 0 以下なら既定件数を使う。
func (a *App) ListResultsPage(plugin, cursor string, limit int) (ResultPage, error) {
	if limit <= 0 {
		limit = defaultPageLimit
	}
	return a.buildResultsPage(a.baseCtx(), plugin, cursor, limit)
}

// TargetPluginView は翻訳対象プラグイン一覧の 1 行の表示用 DTO。
// Plugin は plugin ファイル名（識別子）、SourcePath は選んだ plugin のフルパス（再実行に使う）、
// CreatedAt は初回登録時刻の文字列、Total はこの plugin が束ねる翻訳対象の総数、Translated は訳済み数。
type TargetPluginView struct {
	Plugin     string `json:"plugin"`
	SourcePath string `json:"sourcePath"`
	CreatedAt  string `json:"createdAt"`
	Total      int    `json:"total"`
	Translated int    `json:"translated"`
}

// ListTargetPlugins は翻訳した対象 plugin を新しい順で返す（一覧・進捗・削除・再実行の入口）。
func (a *App) ListTargetPlugins() ([]TargetPluginView, error) {
	rows, err := a.store.ListTargetPlugins(a.baseCtx())
	if err != nil {
		return nil, fmt.Errorf("翻訳対象プラグインの一覧取得: %w", err)
	}
	views := make([]TargetPluginView, len(rows))
	for i, r := range rows {
		views[i] = TargetPluginView{
			Plugin:     r.Plugin,
			SourcePath: r.SourcePath,
			CreatedAt:  r.CreatedAt,
			Total:      r.Total,
			Translated: r.Translated,
		}
	}
	return views, nil
}

// DeleteTargetPlugin は対象 plugin の翻訳成果を消す（対象スコープの行と連関を消し、共有資産は残す）。
// 消した後の再実行で、その plugin だけまっさら再翻訳される。
func (a *App) DeleteTargetPlugin(plugin string) error {
	if err := a.store.DeleteTargetPlugin(a.baseCtx(), plugin); err != nil {
		return fmt.Errorf("翻訳対象プラグインの削除: %w", err)
	}
	return nil
}

// PromptTemplateView はプロンプトテンプレート編集画面の表示用 DTO。
// BaseDirective は叙述文・台詞の両方に付く base 翻訳指示文、PersonaTemplate は話者のいる台詞に付く口調指示の雛形。
// GenericToneText・PcToneText・PcSex は話者なし台詞（汎用・PC）の口調設定（自由記述 2 つと PC 性別）。
type PromptTemplateView struct {
	BaseDirective   string `json:"baseDirective"`
	PersonaTemplate string `json:"personaTemplate"`
	GenericToneText string `json:"genericToneText"`
	PcToneText      string `json:"pcToneText"`
	PcSex           string `json:"pcSex"`
}

// GetPromptTemplate は編集画面の初期表示用に、現在保存されているプロンプトテンプレートと口調設定を返す。
func (a *App) GetPromptTemplate() (PromptTemplateView, error) {
	t, err := a.store.GetPromptTemplate(a.baseCtx())
	if err != nil {
		return PromptTemplateView{}, fmt.Errorf("プロンプトテンプレートの取得: %w", err)
	}
	return PromptTemplateView{
		BaseDirective:   t.BaseDirective,
		PersonaTemplate: t.PersonaTemplate,
		GenericToneText: t.GenericToneText,
		PcToneText:      t.PcToneText,
		PcSex:           t.PcSex,
	}, nil
}

// SavePromptTemplate は編集したプロンプトテンプレートと口調設定を保存する。
// 保存後の翻訳実行と、結果行の実プロンプト再構成は、この保存値を読んで反映する。
func (a *App) SavePromptTemplate(req PromptTemplateView) error {
	if err := a.store.SavePromptTemplate(a.baseCtx(), model.PromptTemplate{
		BaseDirective:   req.BaseDirective,
		PersonaTemplate: req.PersonaTemplate,
		GenericToneText: req.GenericToneText,
		PcToneText:      req.PcToneText,
		PcSex:           req.PcSex,
	}); err != nil {
		return fmt.Errorf("プロンプトテンプレートの保存: %w", err)
	}
	return nil
}

// TemplateVariableView は directive が実行時に埋める変数 1 件（プロンプトテンプレート画面のレコード別タブ用）。
type TemplateVariableView struct {
	Token       string `json:"token"`
	Description string `json:"description"`
}

// DirectiveView は指示文 1 件（key・指示文・変数宣言）。レコード別タブで指示文を編集する単位。
type DirectiveView struct {
	Key         string                 `json:"key"`
	Instruction string                 `json:"instruction"`
	Variables   []TemplateVariableView `json:"variables"`
}

// RecordAssignmentView は REC:FIELD → directive の割り当て 1 件（読み取り専用の対象一覧）。
type RecordAssignmentView struct {
	RecField    string `json:"recField"`
	LogicalName string `json:"logicalName"`
	Directive   string `json:"directive"`
}

// DirectiveEditingView はプロンプトテンプレート画面のレコード別タブの表示用 DTO。
// Directives は編集できる指示文一覧、Assignments は各 directive の対象 REC:FIELD（割り当ては固定で編集不可）。
type DirectiveEditingView struct {
	Directives  []DirectiveView        `json:"directives"`
	Assignments []RecordAssignmentView `json:"assignments"`
}

// GetDirectiveEditing は編集画面のレコード別タブ初期表示用に、指示文一覧と REC:FIELD 割り当てを返す。
func (a *App) GetDirectiveEditing() (DirectiveEditingView, error) {
	ctx := a.baseCtx()
	directives, err := a.store.ListDirectives(ctx)
	if err != nil {
		return DirectiveEditingView{}, fmt.Errorf("directive の取得: %w", err)
	}
	assignments, err := a.store.ListRecordTypeMaster(ctx)
	if err != nil {
		return DirectiveEditingView{}, fmt.Errorf("REC:FIELD 割り当ての取得: %w", err)
	}
	return DirectiveEditingView{
		Directives:  directiveViews(directives),
		Assignments: assignmentViews(assignments),
	}, nil
}

// SaveDirective は編集した指示文を保存する。割り当て（record_type_master）は固定のため変えず、指示文だけ書き換える。
func (a *App) SaveDirective(key, instruction string) error {
	if err := a.store.SaveDirectiveInstruction(a.baseCtx(), key, instruction); err != nil {
		return fmt.Errorf("指示文の保存: %w", err)
	}
	return nil
}

// directiveViews は directive 行を表示 DTO へ写し、変数宣言（JSON）を配列へ解く。
func directiveViews(rows []model.Directive) []DirectiveView {
	out := make([]DirectiveView, len(rows))
	for i, d := range rows {
		var vars []TemplateVariableView
		if strings.TrimSpace(d.Variables) != "" {
			// 変数宣言が壊れていても画面を止めない。解けなければ変数なし扱いにする。
			_ = json.Unmarshal([]byte(d.Variables), &vars)
		}
		out[i] = DirectiveView{Key: d.Key, Instruction: d.Instruction, Variables: vars}
	}
	return out
}

// assignmentViews は record_type_master 行を表示 DTO へ写す（REC:FIELD は "REC:FIELD" に組む）。
func assignmentViews(rows []model.RecordType) []RecordAssignmentView {
	out := make([]RecordAssignmentView, len(rows))
	for i, r := range rows {
		out[i] = RecordAssignmentView{
			RecField:    r.Rec + ":" + r.Field,
			LogicalName: r.LogicalName,
			Directive:   r.Directive,
		}
	}
	return out
}

// RunExtractAndTranslate は plugin を抽出し、未訳の叙述文と台詞を翻訳し、翻訳件数の要約を返す。
// 結果一覧は数万件になりうるためここでは返さず、frontend が ListResultsPage で取得する。
// 抽出中は extract 進捗を、本文翻訳中は translate 進捗を runtime event で push する。
func (a *App) RunExtractAndTranslate(req RunRequest) (RunResult, error) {
	ctx := a.baseCtx()

	if err := a.prepareForTranslation(ctx, req.PluginPath); err != nil {
		return RunResult{}, err
	}

	// 抽出した対象 plugin だけを翻訳する。plugin 列と同じ値（filepath.Base）で絞り、他 plugin の未訳を巻き込まない。
	count, runErr := a.engine.Run(ctx, provider.Connection{Endpoint: req.Endpoint, APIKey: req.APIKey}, req.Model, filepath.Base(req.PluginPath),
		func(done, total int) {
			a.emitProgress(ProgressEvent{Stage: "translate", Done: done, Total: total})
		})

	result := RunResult{TranslatedCount: count}
	if runErr != nil {
		return result, fmt.Errorf("翻訳に失敗: %w", runErr)
	}
	return result, nil
}

// prepareForTranslation は翻訳前区間（対象 plugin 登録 → 抽出 → 固有名派生 → 既存訳取込 → 取込段）を実行する。
// 同期翻訳（RunExtractAndTranslate）と batch 送信（SubmitBatchTranslation）が共有する。各段は冪等。
// 4 サブ段は無音で進むため段ごとに 1 回進捗を push し、frontend が見出しを出し分けられるようにする。
func (a *App) prepareForTranslation(ctx context.Context, pluginPath string) error {
	// 抽出の前に、翻訳した対象 plugin を登録する（plugin ファイル名で 1 対 1、フルパスを再実行用に保つ）。
	// 開始時に作るため、途中で失敗しても登録が残り、一覧からの削除でやり直せる。plugin 列（narration.plugin 等）と
	// 同じ値（filepath.Base）で束ねる。冪等（2 回目以降は source_path だけ更新）。
	if err := a.store.UpsertTargetPlugin(ctx, filepath.Base(pluginPath), pluginPath); err != nil {
		return fmt.Errorf("翻訳対象プラグインの登録に失敗: %w", err)
	}

	a.emitProgress(ProgressEvent{Stage: "extract", Step: "extract"})
	// 抽出子は注入点（本番は dotnet 子プロセス、テストは DB 直 seed の fake）。利用者が選んだ plugin を抽出する。
	if err := a.extractor.Extract(ctx, pluginPath); err != nil {
		return fmt.Errorf("抽出に失敗: %w", err)
	}

	// 抽出後、extracted_field の英日対から固有名の確定訳語（FULL 完全形）を書き、続けて人名の部分形
	//（名のみ・短名）を派生して master_term へ追記する。DB を作り直しても単独名（例 Aventus→アベンタス）が
	// 機械置換へ載るようにする。冪等。
	a.emitProgress(ProgressEvent{Stage: "extract", Step: "derive"})
	derivedCount, err := a.engine.DeriveMasterTerms(ctx)
	if err != nil {
		return fmt.Errorf("固有名の派生に失敗: %w", err)
	}

	// 同じ extracted_field の英日対から record 単位の既存訳（参照訳）を取り込む。翻訳で原文が完全一致する
	// 叙述文・台詞は AI を呼ばず既訳を流用する（known-issues 項目7）。冪等。
	a.emitProgress(ProgressEvent{Stage: "extract", Step: "reference"})
	referenceCount, err := a.engine.LoadReferenceTranslations(ctx)
	if err != nil {
		return fmt.Errorf("既存訳の取り込みに失敗: %w", err)
	}
	// 供給の集約件数を残す（英日対が無い時に確定訳語・参照訳が 0 のまま進む状態を切り分ける）。
	slog.InfoContext(ctx, "extracted_field の英日対から既存訳と確定訳語を組んだ",
		slog.String("event", "reference_supply_built"),
		slog.String("where", "api.prepareForTranslation"),
		slog.Int("masterTerms", derivedCount),
		slog.Int("references", referenceCount),
	)

	// 取込段: C# が素朴吸い出しした extracted_field を record_type_master で振り分け、
	// narration（叙述文・定型句）・proper_noun（固有名）・line（台詞）へ投入し、話者連関を解決する。
	// 本文・固有名フェーズより前に 1 度実行する。冪等。
	a.emitProgress(ProgressEvent{Stage: "extract", Step: "ingest"})
	if _, err := a.engine.Ingest(ctx); err != nil {
		return fmt.Errorf("取込段に失敗: %w", err)
	}
	return nil
}

// SubmitBatchTranslation は対象 plugin を抽出・取込した上で、未訳の固有名を xAI batch として送信する。
// 同期翻訳と同じ翻訳前区間を通した後、engine.Run の代わりに batch 送信を行う。結果は後日 RefreshBatchTranslations で反映する。
// 送信後は固有名 batch の反映待ちで、この呼び出しでは dest を確定しない。接続情報は都度受ける（永続化しない）。
func (a *App) SubmitBatchTranslation(req RunRequest) error {
	ctx := a.baseCtx()
	if a.batch == nil {
		return fmt.Errorf("batch 翻訳が未配線")
	}
	if err := a.prepareForTranslation(ctx, req.PluginPath); err != nil {
		return err
	}
	conn := provider.Connection{Endpoint: req.Endpoint, APIKey: req.APIKey}
	if err := a.batch.SubmitBatch(ctx, conn, req.Model, filepath.Base(req.PluginPath)); err != nil {
		return fmt.Errorf("batch 送信に失敗: %w", err)
	}
	return nil
}

// GetBatchProgress は対象 plugin の xAI batch 進行状況を副作用なしで返す（状態確認）。
// dest 更新・batch 送信・段更新を一切せず、現段 batch の状態を確認するだけ。進行が無ければ Present=false。
// 接続情報は都度受ける（永続化しない）。取り込みの前に押して進行段・件数・取り込み可否を得るために使う。
func (a *App) GetBatchProgress(req BatchPluginRequest) (BatchProgressView, error) {
	ctx := a.baseCtx()
	if a.batch == nil {
		return BatchProgressView{}, fmt.Errorf("batch 翻訳が未配線")
	}
	conn := provider.Connection{Endpoint: req.Endpoint, APIKey: req.APIKey}
	prog, ok, err := a.batch.ProgressStatus(ctx, conn, req.Plugin)
	if err != nil {
		return BatchProgressView{}, fmt.Errorf("batch 状態確認に失敗: %w", err)
	}
	if !ok {
		return BatchProgressView{Present: false}, nil
	}
	return BatchProgressView{
		Present:   true,
		Stage:     batchStageToken(prog.Stage),
		Total:     prog.Total,
		Pending:   prog.Pending,
		Succeeded: prog.Succeeded,
		Failed:    prog.Failed,
		CanApply:  prog.CanApply,
	}, nil
}

// RefreshBatchTranslations は対象 plugin の batch を取り込む（前進）。現段が完了なら結果を対象行へ書き戻し、
// 固有名段が完了すれば同じ呼び出しで本文 batch を送る。画面操作の時点だけ呼ぶ（常駐ポーリングはしない）。
// 接続情報は都度受ける（永続化しないため取り込みのたびに UI から渡す）。
func (a *App) RefreshBatchTranslations(req BatchPluginRequest) error {
	ctx := a.baseCtx()
	if a.batch == nil {
		return fmt.Errorf("batch 翻訳が未配線")
	}
	conn := provider.Connection{Endpoint: req.Endpoint, APIKey: req.APIKey}
	if err := a.batch.RefreshPlugin(ctx, conn, req.Plugin); err != nil {
		return fmt.Errorf("batch 取り込みに失敗: %w", err)
	}
	return nil
}

// batchStageToken は永続の進行段（proper_noun / body / done）を frontend の段トークン（proper / body / done）へ写す。
func batchStageToken(stage string) string {
	switch stage {
	case model.BatchStageProperNoun:
		return "proper"
	case model.BatchStageBody:
		return "body"
	default:
		return "done"
	}
}

// GetXAIModels は xAI の利用可能モデル一覧を返す（xAI 選択時のモデル選択用）。
// xAI の /v1/models は OpenAI 互換のため、同期 provider の ListModels を xAI endpoint へ向けて引く（batch 送信とは別経路）。
// Endpoint 未指定なら xAI の既定 base を使う。batch 非対応モデル（grok-4.5 系）は一覧から除く。
func (a *App) GetXAIModels(req ConnRequest) ([]string, error) {
	endpoint := req.Endpoint
	if strings.TrimSpace(endpoint) == "" {
		endpoint = provider.DefaultXAIEndpoint
	}
	models, err := a.provider.ListModels(a.baseCtx(), provider.Connection{Endpoint: endpoint, APIKey: req.APIKey})
	if err != nil {
		return nil, fmt.Errorf("xAI モデル一覧の取得: %w", err)
	}
	return filterBatchSupportedModels(models), nil
}

// filterBatchSupportedModels は batch 非対応モデルを一覧から除く。
// xAI の doc は grok-4.5 を batch 非対応と明記する。判明した非対応モデルだけを除く（他は素通し）。
func filterBatchSupportedModels(models []string) []string {
	out := make([]string, 0, len(models))
	for _, m := range models {
		if isBatchUnsupportedModel(m) {
			continue // batch 非対応（doc 明記）。
		}
		out = append(out, m)
	}
	return out
}

// isBatchUnsupportedModel は batch 非対応モデルかを判定する純粋関数。
// xAI のモデル名はドット表記（grok-4.5）とハイフン表記（grok-4-5）が混在しうるため、両表記を除外する。
// 大文字小文字は無視する（表記ゆれで除外漏れを起こさない）。
func isBatchUnsupportedModel(m string) bool {
	l := strings.ToLower(m)
	return strings.Contains(l, "grok-4.5") || strings.Contains(l, "grok-4-5")
}

// emitProgress は進捗を runtime event で frontend へ push する。runtime context が無い（テスト等）なら何もしない。
func (a *App) emitProgress(ev ProgressEvent) {
	if a.ctx == nil {
		return
	}
	wailsruntime.EventsEmit(a.ctx, progressEventName, ev)
}

// cursor の区間種別。叙述文（id 昇順）→ 台詞（id 昇順）→ 固有名（id 昇順）の連結列上の位置を表す。
const (
	cursorNarration = "n"
	cursorLine      = "l"
	cursorProper    = "p"
)

// parseCursor は cursor を区間種別と afterID へ分解する。""（先頭）や不正な cursor は叙述文先頭（id 0）とみなす。
func parseCursor(cursor string) (section string, afterID int64) {
	parts := strings.SplitN(cursor, ":", 2)
	if len(parts) == 2 && (parts[0] == cursorNarration || parts[0] == cursorLine || parts[0] == cursorProper) {
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

// buildResultsPage は cursor の指すページの叙述文・台詞・固有名を取り、台詞へ口調を一括で付けて ResultPage を返す。
// plugin が空でなければその対象 plugin の結果だけに絞る。
// 各行の原文へ機械置換辞書を当て直し、置換内訳（terms）と元レコード種別バッジ（recordType）を再構成して供給する。
func (a *App) buildResultsPage(ctx context.Context, plugin, cursor string, limit int) (ResultPage, error) {
	narrations, lines, propers, total, nextCursor, hasMore, err := a.pageRows(ctx, plugin, cursor, limit)
	if err != nil {
		return ResultPage{}, err
	}

	// base 指示と、(rec, field) ごとの directive 指示文をページ単位で 1 度だけ読み、実プロンプトを実行時と同じ式で再構成する。
	tmpl, err := a.store.GetPromptTemplate(ctx)
	if err != nil {
		return ResultPage{}, fmt.Errorf("プロンプトテンプレートの取得: %w", err)
	}
	directives, err := a.store.ListDirectives(ctx)
	if err != nil {
		return ResultPage{}, fmt.Errorf("directive の取得: %w", err)
	}
	masterRows, err := a.store.ListRecordTypeMaster(ctx)
	if err != nil {
		return ResultPage{}, fmt.Errorf("record_type_master の取得: %w", err)
	}
	instructionByKey, boxByRF, directiveByRF := resultLookups(directives, masterRows)

	// ページの台詞ぶんだけ口調を 1 度に一括生成する（台詞ごとの DB 問い合わせを避ける）。口調 directive の指示文を雛形にする。
	// 話者なし台詞（汎用・PC）は prompt_template 由来の自由記述口調へ感情と性別を重ねる。
	personas, err := a.engine.LinePersonas(ctx, lines, instructionByKey["口調"], engine.ToneDefaults{
		Generic: tmpl.GenericToneText,
		PC:      tmpl.PcToneText,
		PcSex:   tmpl.PcSex,
	})
	if err != nil {
		return ResultPage{}, fmt.Errorf("口調の一括生成: %w", err)
	}

	// 「誰の台詞か」と口調指示の根拠（性別・年齢・声型）を出すため、話者識別もページ単位で一括取得する。
	speakers, err := a.store.LoadLineSpeakers(ctx, lineIDs(lines))
	if err != nil {
		return ResultPage{}, fmt.Errorf("話者識別の一括取得: %w", err)
	}

	// 機械置換辞書をページ単位で 1 度だけ組み、各行の原文へ当て直して内訳を再構成する（master_term ∪ proper_noun）。
	dict, err := a.engine.LoadDictionary(ctx)
	if err != nil {
		return ResultPage{}, fmt.Errorf("機械置換辞書の構築: %w", err)
	}

	views := make([]ResultView, 0, len(narrations)+len(lines)+len(propers))
	for _, n := range narrations {
		view := narrationResultView(n)
		view.RecordType = recordTypeView(boxByRF, n.Rec, n.Field)
		// 送信経路と同じ順（タグ退避 → 辞書機械置換 → 生タグへ復元）で組み直し、置換内訳（terms）と実プロンプトを再構成する。
		masked, tags := runtimetag.Mask(n.Source)
		replaced, used := dict.Apply(masked)
		source := runtimetag.Restore(replaced, tags)
		view.Terms = termViews(used)
		// 叙述文・定型句は、その REC:FIELD の文体・定型句 directive を base へ合成した実プロンプトを再構成する。
		instruction := instructionByKey[directiveByRF[RecordKey{Rec: n.Rec, Field: n.Field}]]
		view.Prompt = prompt.RenderPrompt(prompt.ComposePrompt(tmpl.BaseDirective, instruction, source))
		views = append(views, view)
	}
	for _, l := range lines {
		p, hasPersona := personas[l.ID]
		masked, tags := runtimetag.Mask(l.Source)
		replaced, used := dict.Apply(masked)
		lineSource := runtimetag.Restore(replaced, tags)
		view := ResultView{
			EDID:         l.EDID,
			Source:       l.Source,
			Dest:         l.Dest,
			StatusLabel:  statusLabel(l.Status),
			RecordType:   recordTypeView(boxByRF, l.Rec, l.Field),
			Directive:    p.Directive,
			PersonaLabel: p.Label,
			Terms:        termViews(used),
			// 台詞は話者の口調指示を合成した実プロンプトを再構成する（口調指示の合成を目視で確かめる）。
			Prompt: prompt.RenderPrompt(prompt.ComposePrompt(tmpl.BaseDirective, p.Directive, lineSource)),
		}
		// 話者を解決できた台詞は、誰の台詞かと属性（性別・年齢・声型）を付ける。
		if sp, ok := speakers[l.ID]; ok {
			view.Speaker = &SpeakerView{
				EDID:  sp.EDID,
				Sex:   sexLabel(sp.Sex),
				Age:   ageLabel(sp.RaceEDID),
				Voice: sp.VoiceEDID,
			}
		}
		// 口調が付いた台詞は、展開時に出す口調メタ（判定結果と根拠）を付ける。
		// 汎用・PC は Cell・AttitudeBand・Marked が空で、感情段階と性別（条件由来 or 利用者選択）を出す。
		if hasPersona {
			view.Persona = &PersonaView{
				Cell:         p.Cell,
				Trait:        p.Trait,
				AttitudeBand: p.AttitudeBand,
				EmotionBand:  p.EmotionBand,
				Marked:       p.Marked,
				DecisionPath: p.DecisionPath,
				Sex:          sexLabel(p.Sex),
			}
		}
		views = append(views, view)
	}
	for _, pn := range propers {
		// 固有名は本文より先に確定した訳の単位。種別バッジは箱＝固有名、REC:FIELD は種別（category＝rec）。
		view := ResultView{
			EDID:        pn.Source,
			Source:      pn.Source,
			Dest:        pn.Dest,
			StatusLabel: statusLabel(pn.Status),
			RecordType:  &RecordTypeView{Box: "固有名", RecField: pn.Category},
		}
		views = append(views, view)
	}
	return ResultPage{Total: total, Results: views, NextCursor: nextCursor, HasMore: hasMore}, nil
}

// RecordKey は record_type_master の引きキー (rec, field)。結果行の種別バッジ・directive 再構成に使う。
type RecordKey struct {
	Rec   string
	Field string
}

// resultLookups は directive と record_type_master から結果再構成用の 3 つの引きを作る。
// instructionByKey は directive キー → 指示文、boxByRF は (rec, field) → box、directiveByRF は (rec, field) → directive キー。
func resultLookups(directives []model.Directive, masterRows []model.RecordType) (
	instructionByKey map[string]string, boxByRF map[RecordKey]string, directiveByRF map[RecordKey]string) {
	instructionByKey = make(map[string]string, len(directives))
	for _, d := range directives {
		instructionByKey[d.Key] = d.Instruction
	}
	boxByRF = make(map[RecordKey]string, len(masterRows))
	directiveByRF = make(map[RecordKey]string, len(masterRows))
	for _, r := range masterRows {
		key := RecordKey{Rec: r.Rec, Field: r.Field}
		boxByRF[key] = r.Box
		directiveByRF[key] = r.Directive
	}
	return instructionByKey, boxByRF, directiveByRF
}

// recordTypeView は (rec, field) から元レコード種別バッジを組む。master に無い (rec, field) は nil（バッジを出さない）。
func recordTypeView(boxByRF map[RecordKey]string, rec, field string) *RecordTypeView {
	box, ok := boxByRF[RecordKey{Rec: rec, Field: field}]
	if !ok {
		return nil
	}
	return &RecordTypeView{Box: box, RecField: rec + ":" + field}
}

// termViews は辞書の置換内訳（dictionary.Term）を結果行の表示 DTO へ写す。置換が無ければ nil を返す。
func termViews(used []dictionary.Term) []TermView {
	if len(used) == 0 {
		return nil
	}
	out := make([]TermView, len(used))
	for i, t := range used {
		out[i] = TermView{Source: t.Source, Dest: t.Dest}
	}
	return out
}

// pageRows は keyset cursor から当該ページの叙述文・台詞・固有名、総件数、次 cursor、続きの有無を決める。
// 区間は叙述文 → 台詞 → 固有名 の順に連結し、ある区間がページに満たなければ次区間の先頭から補充する。
// 口調・種別バッジは付けない（呼び出し側が付ける）。
func (a *App) pageRows(ctx context.Context, plugin, cursor string, limit int) (
	narrations []model.Narration, lines []model.Line, propers []model.ProperNoun,
	total int, nextCursor string, hasMore bool, err error) {
	total, err = a.countAll(ctx, plugin)
	if err != nil {
		return nil, nil, nil, 0, "", false, err
	}

	section, afterID := parseCursor(cursor)
	pb := &pageBuilder{remaining: limit, afterID: afterID, cur: section, plugin: plugin}
	// 叙述文 → 台詞 → 固有名 の順に区間を埋める。ある区間でページが確定（done）したらそこで打ち切る。
	for _, fill := range []sectionFiller{a.fillNarrations, a.fillLines, a.fillPropers} {
		var done bool
		done, err = fill(ctx, pb)
		if err != nil {
			return nil, nil, nil, 0, "", false, err
		}
		if done {
			break
		}
	}
	return pb.narrations, pb.lines, pb.propers, total, pb.nextCursor, pb.hasMore, nil
}

// countAll は叙述文・台詞・固有名の総件数を合算して返す。ページングの total に使う。plugin が空でなければその対象 plugin に絞る。
func (a *App) countAll(ctx context.Context, plugin string) (int, error) {
	nTotal, err := a.store.CountNarrations(ctx, plugin)
	if err != nil {
		return 0, fmt.Errorf("叙述文の件数: %w", err)
	}
	lTotal, err := a.store.CountLines(ctx, plugin)
	if err != nil {
		return 0, fmt.Errorf("台詞の件数: %w", err)
	}
	pTotal, err := a.store.CountProperNouns(ctx, plugin)
	if err != nil {
		return 0, fmt.Errorf("固有名の件数: %w", err)
	}
	return nTotal + lTotal + pTotal, nil
}

// pageBuilder は pageRows が区間をまたいで持ち回る蓄積と進行状態。各 fill メソッドが変異させる。
// remaining は残り取得枠、afterID は現区間の keyset カーソル位置、cur は現在の区間種別。
// nextCursor / hasMore はページ確定時に各 fill メソッドがセットする。
type pageBuilder struct {
	narrations []model.Narration
	lines      []model.Line
	propers    []model.ProperNoun
	remaining  int
	afterID    int64
	cur        string
	plugin     string // 空でなければ絞り込む対象 plugin（plugin ファイル名）。各 fill が store 問い合わせへ渡す。
	nextCursor string
	hasMore    bool
}

// sectionFiller は 1 区間を埋める関数の型。ページが確定したら done=true を返し、pageRows はそこで打ち切る。
type sectionFiller func(ctx context.Context, pb *pageBuilder) (done bool, err error)

// fillNarrations は叙述文区間を埋める。cursor がここから始まるときだけ取り、limit を超えたらこの区間で切る。
func (a *App) fillNarrations(ctx context.Context, pb *pageBuilder) (bool, error) {
	if pb.cur != cursorNarration {
		return false, nil
	}
	rows, err := a.store.NarrationsAfter(ctx, pb.plugin, pb.afterID, pb.remaining+1)
	if err != nil {
		return false, fmt.Errorf("叙述文ページの取得: %w", err)
	}
	if len(rows) > pb.remaining {
		pb.narrations = rows[:pb.remaining]
		pb.nextCursor = makeCursor(cursorNarration, pb.narrations[len(pb.narrations)-1].ID)
		pb.hasMore = true
		return true, nil
	}
	pb.narrations = rows
	pb.remaining -= len(rows)
	pb.cur, pb.afterID = cursorLine, 0
	return false, nil
}

// fillLines は台詞区間を埋める。叙述文の残りを台詞で補充する。remaining が 0 ちょうどなら、台詞があれば次ページへ送る。
func (a *App) fillLines(ctx context.Context, pb *pageBuilder) (bool, error) {
	if pb.cur != cursorLine {
		return false, nil
	}
	rows, err := a.store.LinesAfter(ctx, pb.plugin, pb.afterID, pb.remaining+1)
	if err != nil {
		return false, fmt.Errorf("台詞ページの取得: %w", err)
	}
	switch {
	case pb.remaining == 0:
		if len(rows) > 0 {
			pb.nextCursor = makeCursor(cursorLine, pb.afterID)
			pb.hasMore = true
			return true, nil
		}
	case len(rows) > pb.remaining:
		pb.lines = rows[:pb.remaining]
		pb.nextCursor = makeCursor(cursorLine, pb.lines[len(pb.lines)-1].ID)
		pb.hasMore = true
		return true, nil
	default:
		pb.lines = rows
		pb.remaining -= len(rows)
	}
	pb.cur, pb.afterID = cursorProper, 0
	return false, nil
}

// fillPropers は固有名区間（末尾）を埋める。台詞の残りを固有名で補充する。ここで尽きれば続きなし。
func (a *App) fillPropers(ctx context.Context, pb *pageBuilder) (bool, error) {
	if pb.cur != cursorProper {
		return false, nil
	}
	rows, err := a.store.ProperNounsAfter(ctx, pb.plugin, pb.afterID, pb.remaining+1)
	if err != nil {
		return false, fmt.Errorf("固有名ページの取得: %w", err)
	}
	if pb.remaining == 0 {
		if len(rows) > 0 {
			pb.nextCursor = makeCursor(cursorProper, pb.afterID)
			pb.hasMore = true
		}
		return true, nil
	}
	if len(rows) > pb.remaining {
		pb.propers = rows[:pb.remaining]
		pb.nextCursor = makeCursor(cursorProper, pb.propers[len(pb.propers)-1].ID)
		pb.hasMore = true
		return true, nil
	}
	pb.propers = rows
	return true, nil
}

// lineIDs は台詞行から id 列を取り出す（口調の一括生成の入力用）。
func lineIDs(lines []model.Line) []int64 {
	ids := make([]int64, len(lines))
	for i, l := range lines {
		ids[i] = l.ID
	}
	return ids
}

// buildExtractorArgs は publish 済み DLL を dotnet で直接実行する引数列を組む。
// 先頭は DLL パスで、dotnet run のような MSBuild 評価（run --project）を通さない。
func buildExtractorArgs(dllPath, dataFolder, plugin, dbPath, schemaDir string) []string {
	return []string{
		dllPath,
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

// sexLabel は性別コード（Female/Male）を表示ラベルへ写す。性別が取れない話者は空。
func sexLabel(sex string) string {
	switch sex {
	case "Female":
		return "女性"
	case "Male":
		return "男性"
	default:
		return ""
	}
}

// ageLabel は種族 EditorID から年齢区分の表示ラベルを引く。ElderRace=老人、*Child=子供、他=成人。種族不明は空。
func ageLabel(raceEDID string) string {
	switch {
	case raceEDID == "":
		return ""
	case strings.Contains(raceEDID, "Child"):
		return "子供"
	case strings.Contains(raceEDID, "Elder"):
		return "老人"
	default:
		return "成人"
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
