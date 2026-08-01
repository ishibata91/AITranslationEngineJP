package harness

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

// 共有オラクル（test-oracle/specs.json）の integration 段（Go ツールの継ぎ目）を、入口→出口の結合テストで照合する。
// 書き方はフォルダの CLAUDE.md に従う: 入口 SyntheticRun を 1 回通し、read-only の出口を spec ごとに照合する。
// 守るのは単体で守れない継ぎ目だけ（件数保存・未知除外・固有名一貫・感情結線・話者結線・stoplist 一貫）。
// 単段で純粋に閉じるルール（口調閾値・派生・stoplist 判定・役割語引き等）は core package の単体テストが守る。

// translationTargetProperNoun は proper_noun のうち翻訳対象（抽出由来）に絞る条件。
// 判別の正本は model.OriginDerived だが、harness は architecture 上 model へ依存できないため
// 本 package ではこの 1 定数へ写し、各 SQL には同じ条件文を書かない。
const translationTargetProperNoun = "origin = ''"

// probe は 1 回の入口→出口の出口。Capture（プロンプト・件数）と最終 DB を持つ。read-only。
type probe struct {
	cap Capture
	db  *sqlx.DB
}

// 1 spec = 1 関数。左が id、右が出口への期待値（Assert）。given は fixture 側にある。
var goOracles = map[string]func(t *testing.T, p probe){
	// 件数保存: 翻訳対象に数えた件数と、出力した件数が一致する（取込→翻訳で落ちない）。
	// 人名の部分形（origin が derive の行）は取込段を通らない機械派生で、翻訳対象に数えない。
	// 数えると派生の増減が件数保存の崩れに見えるため、取り込んだ固有名だけを対象にする。
	"count-parity": func(t *testing.T, p probe) {
		out := countRows(t, p.db, `SELECT count(*) FROM proper_noun WHERE `+translationTargetProperNoun) +
			countRows(t, p.db, `SELECT count(*) FROM narration`) +
			countRows(t, p.db, `SELECT count(*) FROM line`)
		if p.cap.TranslatedCount != out {
			t.Fatalf("件数保存が崩れた: 翻訳件数=%d, 出力行数=%d", p.cap.TranslatedCount, out)
		}
	},

	// 未知除外: 分類表に無い REC:FIELD は、どの箱へも出ず、件数にも含めない。
	"box-routing-unknown-skipped": func(t *testing.T, p probe) {
		for _, tbl := range []string{"narration", "line", "proper_noun"} {
			assertCount(t, p.db, 0, `SELECT count(*) FROM `+tbl+` WHERE source LIKE '%should be skipped%'`)
		}
	},

	// 固有名一貫: 同じ固有名（Aventus Aretino）が叙述文と台詞の両方に出て、同一訳で本文へ入る。
	"proper-noun-consistent": func(t *testing.T, p probe) {
		var dest string
		if err := p.db.Get(&dest, `SELECT dest FROM proper_noun WHERE source='Aventus Aretino'`); err != nil {
			t.Fatalf("固有名 Aventus Aretino の訳が無い: %v", err)
		}
		// 叙述文（WEAP:DESC）と台詞（0x520）の両プロンプトで、原文が同一の確定訳へ置換されている。
		narr := promptContainingUser(t, p, "once held by")
		line := promptContainingUser(t, p, "Have you seen")
		if !strings.Contains(narr.User, dest) || !strings.Contains(line.User, dest) {
			t.Fatalf("固有名の訳が叙述文と台詞で一致しない: dest=%q\n  叙述文=%q\n  台詞=%q", dest, narr.User, line.User)
		}
	},

	// 感情結線: 抽出した台詞感情（Fear）が、その台詞のプロンプトへ乗る（感情 staging→プロンプト）。
	"line-emotion-injected": func(t *testing.T, p probe) {
		pr := promptContainingUser(t, p, "come back to town")
		if !strings.Contains(pr.System, "恐れ") {
			t.Fatalf("台詞感情がプロンプトへ乗っていない（Fear=恐れ を期待）:\n%s", pr.System)
		}
	},

	// 話者結線: 話者が結ばれた台詞のプロンプトへ、その話者の人物像・口調が乗る（話者→ペルソナ→プロンプト）。
	// 人物像には役割語（一人称）と口調の例文も含む。役割語は3キー、例文は4キーで別々に引き、
	// 1つの口調指示へ束ねる。ここは合成表で配線だけを守る。
	"speaker-tone-injected": func(t *testing.T, p probe) {
		pr := promptContainingUser(t, p, "trouble in town")
		if !strings.Contains(pr.System, "人物像") {
			t.Fatalf("話者の人物像がプロンプトへ乗っていない:\n%s", pr.System)
		}
		if !strings.Contains(pr.System, "一人称は「わたし」") {
			t.Fatalf("役割語の一人称がプロンプトへ乗っていない:\n%s", pr.System)
		}
		if !strings.Contains(pr.System, "- 例: I will go. → わたしが行きます。") {
			t.Fatalf("口調の例文がプロンプトへ乗っていない:\n%s", pr.System)
		}
	},

	// 指示文結線: REC:FIELD へ割り当てた指示文が、その本文のプロンプトへ乗る（割り当て→取込→プロンプト）。
	// 割り当ての正本は migration の record_type_master seed で、harness は実 seed を通す。
	// 叙述文（WEAP:DESC＝物品説明）と定型句（ACTI:RNAM＝操作名）で、別の箱が別の指示文を受けることを見る。
	"narration-directive-injected": func(t *testing.T, p probe) {
		narr := promptContainingUser(t, p, "once held by")
		if !strings.Contains(narr.System, "品物の説明文") {
			t.Fatalf("WEAP:DESC へ物品説明の指示文が乗っていない:\n%s", narr.System)
		}
		phrase := promptContainingUser(t, p, "Open")
		if !strings.Contains(phrase.System, "動詞の終止形") {
			t.Fatalf("ACTI:RNAM へ操作名の指示文が乗っていない:\n%s", phrase.System)
		}
		// 箱が違えば違う指示文を受けること（1 指示文へ畳まれていない）。
		if narr.System == phrase.System {
			t.Fatalf("叙述文と定型句が同じ指示文を受けている:\n%s", narr.System)
		}
	},

	// stoplist 一貫: 供給から外した一般語（Yes）は本文で機械置換されず、本文の一般語出現を壊さない。
	"stoplist-preserved-in-body": func(t *testing.T, p probe) {
		pr := promptContainingUser(t, p, "the road is clear")
		if !strings.Contains(pr.User, "Yes, the road is clear") {
			t.Fatalf("stoplist 語 Yes が本文で置換された（原文保持を期待）:\n%s", pr.User)
		}
	},

	// 実行時タグ保護: <Alias=...> を含む叙述文は、AI へ生タグを見せ（user に原形）、system にタグ保護指示が乗り、出力にタグが残る。
	"runtime-tag-preserved": func(t *testing.T, p probe) {
		// 送信プロンプトの user は生タグ <Alias=Player> を原形で持つこと（AI に意味の分かるタグを見せる）。
		pr := promptContainingUser(t, p, "Deliver this letter to")
		if !strings.Contains(pr.User, "<Alias=Player>") {
			t.Fatalf("生タグが原形で AI へ渡っていない:\n%s", pr.User)
		}
		// タグを持つ本文には system へタグ保護指示が乗ること。
		if !strings.Contains(pr.System, "一字一句変えずに残すこと") {
			t.Fatalf("タグ保護指示が system に乗っていない:\n%s", pr.System)
		}
		// 出力にタグが原形で残る行は、そのまま書き戻される（欠落なし → 仮訳として保存）。
		var dest string
		if err := p.db.Get(&dest, `SELECT dest FROM narration WHERE source LIKE '%Deliver this letter to%'`); err != nil {
			t.Fatalf("タグ入り叙述文の訳が無い: %v", err)
		}
		if !strings.Contains(dest, "<Alias=Player>") {
			t.Fatalf("実行時タグが最終出力に残っていない: dest=%q", dest)
		}
	},

	// 既存訳の横断供給: 翻訳対象が日本語 Strings を持たない plugin でも、同じ Data フォルダの別 plugin の
	// 英日対から既存訳と固有名の訳が組まれ、対象 plugin の本文へその訳が当たる。
	// 合成 fixture の Synthetic.esm は日本語を持たず、既訳はすべて別 plugin 由来（References）で供給する。
	"reference-supply-cross-plugin": func(t *testing.T, p probe) {
		// 別 plugin の英日対から横断辞書へ固有名の完全形が入る。
		var dest string
		if err := p.db.Get(&dest,
			`SELECT dest FROM master_term WHERE source='Aventus Aretino'`); err != nil {
			t.Fatalf("別 plugin 由来の固有名が横断辞書に無い: %v", err)
		}
		if dest != "アベンタス・アレティノ" {
			t.Fatalf("別 plugin 由来の訳語 = %q, want アベンタス・アレティノ", dest)
		}
		// 対象 plugin の本文へ、その訳が機械置換で当たる。
		narr := promptContainingUser(t, p, "once held by")
		if !strings.Contains(narr.User, dest) {
			t.Fatalf("別 plugin 由来の訳が本文へ当たっていない:\n%s", narr.User)
		}
	},

	// 実行内で確定した氏名からの部分形（姓名分割）: 既存訳を持たない mod NPC の氏名は AI 訳で確定する。
	// その確定訳から名だけ・苗字だけの訳を作り、本文の該当語へ当てる。
	"run-name-part-derived": func(t *testing.T, p probe) {
		var full string
		if err := p.db.Get(&full, `SELECT dest FROM proper_noun WHERE source='Sorine Trueblade'`); err != nil {
			t.Fatalf("mod NPC の氏名が確定していない: %v", err)
		}
		given, family, ok := strings.Cut(full, "・")
		if !ok {
			t.Fatalf("氏名の確定訳が中黒区切りでない: %q", full)
		}
		pr := promptContainingUser(t, p, "guards the gate")
		if strings.Contains(pr.User, "Sorine") || !strings.Contains(pr.User, given) {
			t.Fatalf("名だけの人名が訳語へ置き換わっていない（want %q）:\n%s", given, pr.User)
		}
		if strings.Contains(pr.User, "Trueblade") || !strings.Contains(pr.User, family) {
			t.Fatalf("苗字だけの人名が訳語へ置き換わっていない（want %q）:\n%s", family, pr.User)
		}
	},

	// 実行内で確定した氏名からの部分形（二つ名）: 二つ名を除いた名だけが本文へ出た時も訳語へ置き換わる。
	"run-byname-part-derived": func(t *testing.T, p probe) {
		var full string
		if err := p.db.Get(&full, `SELECT dest FROM proper_noun WHERE source='Ulfarr the Grim'`); err != nil {
			t.Fatalf("二つ名を持つ mod NPC の氏名が確定していない: %v", err)
		}
		pr := promptContainingUser(t, p, "old shrine")
		if strings.Contains(pr.User, "Ulfarr") {
			t.Fatalf("二つ名を除いた名が原文のまま残った（確定訳 %q）:\n%s", full, pr.User)
		}
	},

	// 一般語との衝突: mod NPC の名が本文の一般語と同綴りなら、部分形として辞書へ載せない。
	// 載せると一般語の出現（文頭の大文字を含む）が訳語へ置き換わり、原文が壊れる。
	"run-common-word-part-skipped": func(t *testing.T, p probe) {
		assertCount(t, p.db, 0, `SELECT count(*) FROM proper_noun WHERE source='Hunter' AND NOT `+translationTargetProperNoun)
		// 文頭で大文字になる一般語の出現が原文のまま残る（実データで観測した破壊置換の形）。
		head := promptContainingUser(t, p, "tracks are fresh")
		if !strings.Contains(head.User, "Hunter tracks are fresh") {
			t.Fatalf("文頭の一般語 Hunter が置換された:\n%s", head.User)
		}
		// 小文字始まりの出現も原文のまま残る。
		low := promptContainingUser(t, p, "near the ridge")
		if !strings.Contains(low.User, "a hunter near the ridge") {
			t.Fatalf("小文字始まりの一般語 hunter が置換された:\n%s", low.User)
		}
	},

	// 既存訳の流用: 原文が既訳と完全一致する台詞は、AI を呼ばず既訳を確定訳（status=1）で書き戻す。
	"existing-translation-reused": func(t *testing.T, p probe) {
		var dest string
		var status int
		if err := p.db.QueryRow(`SELECT dest, status FROM line WHERE source='Well met, traveler.'`).Scan(&dest, &status); err != nil {
			t.Fatalf("既訳一致の台詞が無い: %v", err)
		}
		if dest != "ようこそ、旅の方。" || status != 1 {
			t.Fatalf("既訳が確定訳で流用されていない: dest=%q status=%d（既訳・status=1 を期待）", dest, status)
		}
		// この原文は AI へ渡っていないこと（provider を呼ばずに流用した）。
		for _, pr := range p.cap.Prompts {
			if strings.Contains(pr.User, "Well met, traveler.") {
				t.Fatalf("既訳一致の台詞が AI へ渡った（流用されていない）:\n%s", pr.User)
			}
		}
	},
}

// TestGoOracles は入口→出口を 1 回通し、read-only の出口を各 spec 関数へ渡す（パラメタライズド）。
func TestGoOracles(t *testing.T) {
	p := runSynthetic(t)
	for id, fn := range goOracles {
		t.Run(id, func(t *testing.T) { fn(t, p) })
	}
}

// 網羅番人: 登録した id と specs.json の integration 段 id が完全一致すること。
// spec を足して関数を書き忘れたら（または余計な id を登録したら）ここで落ちる。
func TestGoOracleCoverage(t *testing.T) {
	want := integrationSpecIDs(t)
	got := make([]string, 0, len(goOracles))
	for id := range goOracles {
		got = append(got, id)
	}
	sort.Strings(got)
	if strings.Join(want, ",") != strings.Join(got, ",") {
		t.Fatalf("登録オラクルと specs.json が不一致:\n  specs=%v\n  登録 =%v", want, got)
	}
}

// runSynthetic は合成入力で入口→出口を 1 回通し、出口（Capture と最終 DB）を返す。
func runSynthetic(t *testing.T) probe {
	t.Helper()
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "synthetic.sqlite3")
	captured, err := SyntheticRun(dbPath)
	if err != nil {
		t.Fatalf("合成 harness の実行: %v", err)
	}
	db, err := sqlx.Connect("sqlite", dbPath)
	if err != nil {
		t.Fatalf("最終 DB を開けない: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return probe{cap: captured, db: db}
}

// promptContainingUser は user 本文に substr を含む送信プロンプトを 1 件返す。無ければ落とす。
func promptContainingUser(t *testing.T, p probe, substr string) RecordedPrompt {
	t.Helper()
	for _, pr := range p.cap.Prompts {
		if strings.Contains(pr.User, substr) {
			return pr
		}
	}
	t.Fatalf("user 本文に %q を含むプロンプトが無い", substr)
	return RecordedPrompt{}
}

// countRows は count クエリの結果を返す。
func countRows(t *testing.T, db *sqlx.DB, query string) int {
	t.Helper()
	var got int
	if err := db.Get(&got, query); err != nil {
		t.Fatalf("件数取得（%s）: %v", query, err)
	}
	return got
}

// assertCount は count クエリの結果が want と一致することを確かめる。
func assertCount(t *testing.T, db *sqlx.DB, want int, query string) {
	t.Helper()
	if got := countRows(t, db, query); got != want {
		t.Fatalf("件数が期待と違う: got=%d want=%d\n  query=%s", got, want, query)
	}
}

// integrationSpecIDs は specs.json の integration 段 id を返す（Go ツールの継ぎ目）。
func integrationSpecIDs(t *testing.T) []string {
	t.Helper()
	var doc struct {
		Specs []struct {
			ID    string `json:"id"`
			Stage string `json:"stage"`
		} `json:"specs"`
	}
	raw, err := os.ReadFile(filepath.Join(repoRoot(t), "test-oracle", "specs.json"))
	if err != nil {
		t.Fatalf("specs.json の読み込み: %v", err)
	}
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatalf("specs.json の解析: %v", err)
	}
	ids := make([]string, 0)
	for _, s := range doc.Specs {
		if s.Stage == "integration" {
			ids = append(ids, s.ID)
		}
	}
	sort.Strings(ids)
	return ids
}

// repoRoot は go.mod を含む repo root を、テスト実行 dir から遡って探す。
func repoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("作業 dir の取得: %v", err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("go.mod を含む repo root が見つからない")
		}
		dir = parent
	}
}
