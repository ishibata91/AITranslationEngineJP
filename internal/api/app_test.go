package api

import (
	"context"
	"strconv"
	"strings"
	"testing"

	"aitranslationenginejp/internal/model"
)

// 訳状態コードを画面表示ラベルへ写すこと（xTranslator の Status 値域に従う）。
func TestStatusLabel(t *testing.T) {
	cases := map[int]string{
		0: "未訳",
		1: "訳済",
		2: "部分",
		3: "仮訳",
		4: "承認",
		9: "未訳", // 未知は未訳扱い
	}
	for code, want := range cases {
		if got := statusLabel(code); got != want {
			t.Errorf("statusLabel(%d) = %q, want %q", code, got, want)
		}
	}
}

// batch 非対応モデル（grok-4.5）を、ドット表記・ハイフン表記・大文字混在のいずれでも除外し、
// 対応モデルは素通しすること。xAI のモデル名表記ゆれで除外漏れを起こさないため。
func TestFilterBatchSupportedModels(t *testing.T) {
	got := filterBatchSupportedModels([]string{
		"grok-4.5", "grok-4-5", "GROK-4.5-mini", "grok-4", "grok-3",
	})
	want := []string{"grok-4", "grok-3"}
	if len(got) != len(want) {
		t.Fatalf("filterBatchSupportedModels = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("filterBatchSupportedModels[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

// narration を結果一覧の行へ写すこと。叙述文は口調指示を持たない。
func TestNarrationResultView(t *testing.T) {
	v := narrationResultView(model.Narration{
		EDID: "DLC1BookSerana", Source: "halls", Dest: "広間", Status: 3,
	})
	if v.EDID != "DLC1BookSerana" || v.Source != "halls" || v.Dest != "広間" || v.StatusLabel != "仮訳" {
		t.Errorf("narrationResultView = %+v", v)
	}
	if v.Directive != "" || v.PersonaLabel != "" {
		t.Errorf("叙述文に口調指示が付いた: %+v", v)
	}
}

// extractor 子プロセスの引数を組み立てること（publish 済み DLL を dotnet で直接実行する）。
// dotnet run（run --project）で毎回 MSBuild を評価するのを避けるため、先頭は DLL パスで
// run・--project を含まないことを固定する。
func TestBuildExtractorArgs(t *testing.T) {
	dll := "tools/extractor/bin/publish/extractor.dll"
	args := buildExtractorArgs(dll, "/sky/Data", "Dawnguard.esm", "db/c.sqlite3", "db/migrations")
	joined := strings.Join(args, " ")
	want := dll + " --data /sky/Data --plugin Dawnguard.esm --sqlite db/c.sqlite3 --schema db/migrations"
	if joined != want {
		t.Errorf("args = %q, want %q", joined, want)
	}
	if args[0] != dll {
		t.Errorf("先頭は DLL パスであるべき: args[0] = %q", args[0])
	}
	for _, a := range args {
		if a == "run" || a == "--project" {
			t.Errorf("dotnet run の引数（run / --project）が残っている: %q", joined)
		}
	}
}

// fakePageStore は api.Store を満たす最小の fake。id 昇順の叙述文・台詞・固有名を保持し、keyset 範囲取得を模す。
// DB を使わず、pageRows の cursor 境界ロジック（連結列の順次送り・区間跨ぎ・末尾）だけを確かめる。
type fakePageStore struct {
	narrations []model.Narration
	lines      []model.Line
	propers    []model.ProperNoun
}

// plugin フィルタは pageRows の cursor 境界ロジック確認には使わないため無視する（呼び出しは "" を渡す）。
func (f *fakePageStore) CountNarrations(_ context.Context, _ string) (int, error) {
	return len(f.narrations), nil
}
func (f *fakePageStore) CountLines(_ context.Context, _ string) (int, error) {
	return len(f.lines), nil
}
func (f *fakePageStore) CountProperNouns(_ context.Context, _ string) (int, error) {
	return len(f.propers), nil
}

func (f *fakePageStore) NarrationsAfter(_ context.Context, _ string, afterID int64, limit int) ([]model.Narration, error) {
	var out []model.Narration
	for _, n := range f.narrations {
		if n.ID > afterID {
			out = append(out, n)
			if len(out) == limit {
				break
			}
		}
	}
	return out, nil
}

func (f *fakePageStore) LinesAfter(_ context.Context, _ string, afterID int64, limit int) ([]model.Line, error) {
	var out []model.Line
	for _, l := range f.lines {
		if l.ID > afterID {
			out = append(out, l)
			if len(out) == limit {
				break
			}
		}
	}
	return out, nil
}

func (f *fakePageStore) ProperNounsAfter(_ context.Context, _ string, afterID int64, limit int) ([]model.ProperNoun, error) {
	var out []model.ProperNoun
	for _, p := range f.propers {
		if p.ID > afterID {
			out = append(out, p)
			if len(out) == limit {
				break
			}
		}
	}
	return out, nil
}

// 翻訳対象プラグインの CRUD は pageRows テストで未使用のスタブ（interface を満たすためだけ）。
func (f *fakePageStore) UpsertTargetPlugin(_ context.Context, _, _ string) error { return nil }
func (f *fakePageStore) ListTargetPlugins(_ context.Context) ([]model.TargetPlugin, error) {
	return nil, nil
}
func (f *fakePageStore) DeleteTargetPlugin(_ context.Context, _ string) error { return nil }

// pageRows の cursor 境界ロジックだけを確かめる fake のため、テンプレート・directive CRUD は未使用のスタブにする。
func (f *fakePageStore) GetPromptTemplate(_ context.Context) (model.PromptTemplate, error) {
	return model.PromptTemplate{}, nil
}
func (f *fakePageStore) SavePromptTemplate(_ context.Context, _ model.PromptTemplate) error {
	return nil
}
func (f *fakePageStore) LoadLineSpeakers(_ context.Context, _ []int64) (map[int64]model.LineSpeaker, error) {
	return map[int64]model.LineSpeaker{}, nil
}
func (f *fakePageStore) ListDirectives(_ context.Context) ([]model.Directive, error) {
	return nil, nil
}
func (f *fakePageStore) SaveDirectiveInstruction(_ context.Context, _, _ string) error {
	return nil
}
func (f *fakePageStore) ListRecordTypeMaster(_ context.Context) ([]model.RecordType, error) {
	return nil, nil
}

func narrationsWithIDs(ids ...int64) []model.Narration {
	rows := make([]model.Narration, len(ids))
	for i, id := range ids {
		rows[i] = model.Narration{ID: id}
	}
	return rows
}

func linesWithIDs(ids ...int64) []model.Line {
	rows := make([]model.Line, len(ids))
	for i, id := range ids {
		rows[i] = model.Line{ID: id}
	}
	return rows
}

// rowKeys は叙述文・台詞を "n:<id>" / "l:<id>" で表し、ページ走査の順序確認に使う。
func rowKeys(narrations []model.Narration, lines []model.Line) []string {
	keys := make([]string, 0, len(narrations)+len(lines))
	for _, n := range narrations {
		keys = append(keys, "n:"+strconv.FormatInt(n.ID, 10))
	}
	for _, l := range lines {
		keys = append(keys, "l:"+strconv.FormatInt(l.ID, 10))
	}
	return keys
}

// cursor "" から hasMore=false まで pageRows を辿り、連結列（叙述文→台詞、各 id 昇順）を
// ページサイズ 2 で順に取り切ること。総件数が一定で、末尾でちょうど止まること。
func TestPageRowsWalk(t *testing.T) {
	app := &App{store: &fakePageStore{
		narrations: narrationsWithIDs(1, 2, 3),
		lines:      linesWithIDs(1, 2, 3, 4, 5),
	}}

	var got []string
	cursor := ""
	pages := 0
	for {
		narrations, lines, _, total, next, hasMore, err := app.pageRows(context.Background(), "", cursor, 2)
		if err != nil {
			t.Fatalf("pageRows: %v", err)
		}
		if total != 8 {
			t.Errorf("total = %d, want 8", total)
		}
		got = append(got, rowKeys(narrations, lines)...)
		pages++
		if pages > 10 {
			t.Fatalf("ページが終わらない（無限ループ）")
		}
		if !hasMore {
			break
		}
		cursor = next
	}

	want := []string{"n:1", "n:2", "n:3", "l:1", "l:2", "l:3", "l:4", "l:5"}
	if len(got) != len(want) {
		t.Fatalf("走査結果 = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("走査[%d] = %s, want %s", i, got[i], want[i])
		}
	}
	if pages != 4 {
		t.Errorf("ページ数 = %d, want 4（8 件 ÷ 2）", pages)
	}
}

// 叙述文区間から台詞区間へ跨ぐページで、叙述文の残り＋台詞の先頭を返し、次 cursor が台詞区間（l:<id>）になること。
func TestPageRowsCrossBoundary(t *testing.T) {
	app := &App{store: &fakePageStore{
		narrations: narrationsWithIDs(1, 2, 3),
		lines:      linesWithIDs(1, 2, 3, 4, 5),
	}}

	narrations, lines, _, _, next, hasMore, err := app.pageRows(context.Background(), "", "n:2", 2)
	if err != nil {
		t.Fatalf("pageRows: %v", err)
	}
	if len(narrations) != 1 || narrations[0].ID != 3 {
		t.Errorf("叙述文 = %v, want [3]", narrations)
	}
	if len(lines) != 1 || lines[0].ID != 1 {
		t.Errorf("台詞 = %v, want [1]", lines)
	}
	if next != "l:1" {
		t.Errorf("next cursor = %q, want l:1", next)
	}
	if !hasMore {
		t.Errorf("hasMore = false, want true")
	}
}

// 叙述文がちょうどページを埋め、台詞が残る場合、次 cursor が台詞先頭（l:0）になり hasMore=true であること。
func TestPageRowsNarrationsExactlyFill(t *testing.T) {
	app := &App{store: &fakePageStore{
		narrations: narrationsWithIDs(1, 2),
		lines:      linesWithIDs(1, 2, 3),
	}}

	narrations, lines, _, total, next, hasMore, err := app.pageRows(context.Background(), "", "", 2)
	if err != nil {
		t.Fatalf("pageRows: %v", err)
	}
	if len(narrations) != 2 || len(lines) != 0 {
		t.Errorf("叙述文=%v 台詞=%v, want 叙述文 2・台詞 0", narrations, lines)
	}
	if next != "l:0" {
		t.Errorf("next cursor = %q, want l:0", next)
	}
	if !hasMore {
		t.Errorf("hasMore = false, want true")
	}
	if total != 5 {
		t.Errorf("total = %d, want 5", total)
	}

	// 続く l:0 から台詞先頭が取れること。
	_, lines2, _, _, _, _, err := app.pageRows(context.Background(), "", next, 2)
	if err != nil {
		t.Fatalf("pageRows(l:0): %v", err)
	}
	if len(lines2) != 2 || lines2[0].ID != 1 {
		t.Errorf("台詞先頭 = %v, want [1,2]", lines2)
	}
}

// 結果が無いとき、総件数 0・hasMore=false・next 空であること。
func TestPageRowsEmpty(t *testing.T) {
	app := &App{store: &fakePageStore{}}

	narrations, lines, propers, total, next, hasMore, err := app.pageRows(context.Background(), "", "", 2)
	if err != nil {
		t.Fatalf("pageRows: %v", err)
	}
	if len(narrations) != 0 || len(lines) != 0 || len(propers) != 0 || total != 0 || hasMore || next != "" {
		t.Errorf("空のはず: 叙述文=%v 台詞=%v 固有名=%v total=%d hasMore=%v next=%q",
			narrations, lines, propers, total, hasMore, next)
	}
}

// 台詞のあと固有名区間へ跨ぐページで、台詞の残り＋固有名の先頭を返し、次 cursor が固有名区間（p:<id>）になること。
// 連結列は叙述文 → 台詞 → 固有名。固有名が末尾区間で、ここで尽きれば hasMore=false になることも確かめる。
func TestPageRowsProperNounSection(t *testing.T) {
	app := &App{store: &fakePageStore{
		narrations: narrationsWithIDs(1),
		lines:      linesWithIDs(1, 2),
		propers:    propersWithIDs(1, 2),
	}}

	// "" から走査して、連結列 n:1 → l:1,l:2 → p:1,p:2 を順に取り切ること。
	var got []string
	cursor := ""
	for pages := 0; ; pages++ {
		narrations, lines, propers, total, next, hasMore, err := app.pageRows(context.Background(), "", cursor, 2)
		if err != nil {
			t.Fatalf("pageRows: %v", err)
		}
		if total != 5 {
			t.Errorf("total = %d, want 5", total)
		}
		for _, n := range narrations {
			got = append(got, "n:"+strconv.FormatInt(n.ID, 10))
		}
		for _, l := range lines {
			got = append(got, "l:"+strconv.FormatInt(l.ID, 10))
		}
		for _, p := range propers {
			got = append(got, "p:"+strconv.FormatInt(p.ID, 10))
		}
		if !hasMore {
			break
		}
		cursor = next
		if pages > 10 {
			t.Fatalf("ページが終わらない（無限ループ）")
		}
	}
	want := []string{"n:1", "l:1", "l:2", "p:1", "p:2"}
	if len(got) != len(want) {
		t.Fatalf("走査結果 = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("走査[%d] = %s, want %s", i, got[i], want[i])
		}
	}
}

func propersWithIDs(ids ...int64) []model.ProperNoun {
	rows := make([]model.ProperNoun, len(ids))
	for i, id := range ids {
		rows[i] = model.ProperNoun{ID: id}
	}
	return rows
}

// R-1-2: 翻訳対象に選んだ plugin 自身が日本語訳を持つ時も、その plugin の英日対が既訳として集まること。
// 既訳の収集は対象 plugin が置かれた Data フォルダ全体を走査対象にするため、対象 plugin も走査に入る。
// 走査を 1 本へ絞る --plugin を渡さないこと（渡すと対象以外の plugin が供給から落ちる）を併せて固定する。
func TestBuildReferenceScanArgs(t *testing.T) {
	dll := "tools/extractor/bin/publish/extractor.dll"
	args := buildReferenceScanArgs(dll, "/sky/Data", "db/c.sqlite3", "db/migrations")
	joined := strings.Join(args, " ")
	want := dll + " --data /sky/Data --sqlite db/c.sqlite3 --schema db/migrations --references"
	if joined != want {
		t.Errorf("args = %q, want %q", joined, want)
	}
	for _, a := range args {
		if a == "--plugin" {
			t.Errorf("走査を 1 本へ絞る --plugin が渡っている: %q", joined)
		}
	}
}
