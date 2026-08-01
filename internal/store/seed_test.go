package store

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"aitranslationenginejp/db"

	"github.com/jmoiron/sqlx"
)

const acceptedToneDirective = `この台詞の話者の人物像:
{traits}
この人物像に合う口調と人称で訳すこと。
台詞は話し手と聞き手が決まっているので、主語を書かなくても誰の話かは伝わる。英語の I と you を日本語の主語へ置き換えず、原則として書かない。
台詞は常体で書く。ただし、丁寧な依頼や断りを命令へ変えない。「～してくれないか」「すまないが～」など、常体のまま原文の丁寧さを保つ。
原文の内容をすべて訳す。日本語で不要な主語と重複だけを省く。依頼、遠慮、推量、程度、話者の態度は短縮しない。
台詞は文末に句点を打たない。原文が疑問符・感嘆符で終わる時だけ ？ ！ を置く。
1 つの台詞に文が 2 つ以上ある時は、句点でつなげず全角空白で区切る。

英語の文法上の形ではなく、台詞全体が伝える内容と話者の意図を日本語で表す。

次の変更は、本文だけから意味と働きを確定できる場合に限って行う。

- 修辞疑問は、答えを求める質問として訳さない。原文が示す断定、皮肉、呆れ、非難を日本語で表す。例: Why am I not surprised? → やっぱり驚かない
- if you'll excuse me、I'm afraid、would you mind などの定型的な丁寧表現は、条件、恐怖、質問として直訳しない。退出の断り、残念の表明、依頼、拒否など、その箇所が果たす働きを日本語で表す。例: if you'll excuse me, I need to leave → すまないが、もう行かなければ。例: would you mind helping me? → 手伝ってくれないか
- 比喩や否定構文に対応する一般的な日本語表現が明確にある場合は、英語と同じ構造を残さず、その意味を直接表す。例: nothing to joke about → 軽く考えていい話ではない。例: the last word I received → 最後の便り
- 原文の複数の節が同じ内容を説明している場合は、日本語で同じ内容を重ねて言わない。ただし、条件、理由、程度、話者の態度は省かない。

次の内容は必ず保つ。

- 誰が行動するか。
- 誰に行動を求めているか。
- 質問、依頼、勧誘、拒否、警告、抗議、断言の区別。
- 丁寧さ、親しさ、敵意、皮肉の強さ。
- 肯定と否定、可能、義務、推量、条件、理由。
- 原文に明示された情報量。

本文だけでは意味を確定できない省略や指示対象は推測しない。対応する日本語表現に確信がない場合は、無理に意訳せず意味を保つ訳を選ぶ。

意訳する場合も、原文にない態度、経緯、指示対象を作らない。肯定と否定、頻度、程度、行為者と対象の関係を変えない。文をつなぎ直す場合も必要な述語と補語を残し、日本語として意味が完結する形にする。`

const migration17ToneDirective = `この台詞の話者の人物像:
{traits}
この人物像に合う口調と人称で訳すこと。
台詞は話し手と聞き手が決まっているので、主語を書かなくても誰の話かは伝わる。英語の I と you を日本語の主語へ置き換えず、原則として書かない。
台詞は常体で書く。話者が女性でも年配でも礼儀正しくても、です・ます体にしない。
原文の内容をすべて訳す。そのうえで英語の語順と品詞をなぞらず、説明を足さずに短く言い切る。
台詞は文末に句点を打たない。原文が疑問符・感嘆符で終わる時だけ ？ ！ を置く。
1 つの台詞に文が 2 つ以上ある時は、句点でつなげず全角空白で区切る。`

// migration の seed 整合を確かめる。directive は 9 指示文、record_type_master は全 REC:FIELD が
// directive を 1 つだけ持ち（排他）、参照先 directive が必ず存在し（網羅）、box と directive が対応すること。
func TestSeedRecordTypeMasterIntegrity(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "seed.sqlite3")
	s, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer func() { _ = s.Close() }()
	ctx := context.Background()

	directives, err := s.ListDirectives(ctx)
	if err != nil {
		t.Fatalf("ListDirectives: %v", err)
	}
	wantKeys := map[string]bool{
		"物品説明": true, "効果説明": true, "世界観断片": true, "書物体": true, "日記体": true,
		"固有名": true, "操作名": true, "語義": true, "口調": true,
	}
	if len(directives) != len(wantKeys) {
		t.Errorf("directive 数 = %d, want %d", len(directives), len(wantKeys))
	}
	keys := make(map[string]bool, len(directives))
	for _, d := range directives {
		keys[d.Key] = true
		if !wantKeys[d.Key] {
			t.Errorf("想定外の directive キー: %q", d.Key)
		}
	}
	for k := range wantKeys {
		if !keys[k] {
			t.Errorf("directive キー %q が seed に無い", k)
		}
	}

	rows, err := s.ListRecordTypeMaster(ctx)
	if err != nil {
		t.Fatalf("ListRecordTypeMaster: %v", err)
	}
	if len(rows) == 0 {
		t.Fatal("record_type_master の seed が空")
	}

	validBox := map[string]bool{"叙述文": true, "固有名": true, "定型句": true, "台詞": true}
	// box ごとに許される directive。叙述文は 5 文体のいずれか、定型句は短文 2 種のいずれか、
	// 固有名と台詞は単一 directive（箱と 1 対 1）。
	narrationDirectives := map[string]bool{
		"物品説明": true, "効果説明": true, "世界観断片": true, "書物体": true, "日記体": true,
	}
	setPhraseDirectives := map[string]bool{"操作名": true, "語義": true}

	seen := make(map[string]bool, len(rows))
	// directiveOf は REC:FIELD → 割り当てた指示文キー、usedDirectives は 1 件以上から参照された指示文キー。
	directiveOf := make(map[string]string, len(rows))
	usedDirectives := make(map[string]bool, len(rows))
	for _, r := range rows {
		rf := r.Rec + ":" + r.Field
		// 排他: 同じ (rec, field) が 2 度現れない（PK の担保を読み出し側でも確認）。
		if seen[rf] {
			t.Errorf("REC:FIELD %q が重複している（排他違反）", rf)
		}
		seen[rf] = true
		directiveOf[rf] = r.Directive
		usedDirectives[r.Directive] = true

		if !validBox[r.Box] {
			t.Errorf("%s: 未知の box %q", rf, r.Box)
		}
		// 網羅: 割り当て先 directive が必ず存在する（孤立参照なし）。
		if !keys[r.Directive] {
			t.Errorf("%s: 参照先 directive %q が存在しない", rf, r.Directive)
		}
		// box と directive の対応。
		switch r.Box {
		case "台詞":
			if r.Directive != "口調" {
				t.Errorf("%s: 台詞 box の directive = %q, want 口調", rf, r.Directive)
			}
		case "固有名":
			if r.Directive != "固有名" {
				t.Errorf("%s: 固有名 box の directive = %q, want 固有名", rf, r.Directive)
			}
		case "定型句":
			if !setPhraseDirectives[r.Directive] {
				t.Errorf("%s: 定型句 box の directive = %q, want 操作名・語義のいずれか", rf, r.Directive)
			}
		case "叙述文":
			if !narrationDirectives[r.Directive] {
				t.Errorf("%s: 叙述文 box の directive = %q, want 文体のいずれか", rf, r.Directive)
			}
		}
	}

	// 代表 REC:FIELD が網羅されていること（会話・本以外への拡張の最低保証）。
	// SPEL:DESC・WOOP:TNAM・RACE:DESC は指示文を割り直した対象で、分割後の割り当てを固定する。
	for _, rf := range []string{"BOOK:DESC", "INFO:NAM1", "WEAP:DESC", "WEAP:FULL", "ACTI:RNAM", "DIAL:FULL", "QUST:CNAM", "LSCR:DESC", "SPEL:DESC", "WOOP:TNAM", "RACE:DESC"} {
		if !seen[rf] {
			t.Errorf("代表 REC:FIELD %q が record_type_master に無い", rf)
		}
	}

	// 分割後の割り当て。旧「説明体」を物品説明と効果説明へ、旧「定型句」を操作名と語義へ割り、
	// RACE:DESC を世界観断片へ移した結果を固定する。
	for _, tc := range []struct{ recField, directive string }{
		{"WEAP:DESC", "物品説明"},
		{"SPEL:DESC", "効果説明"},
		{"MGEF:DNAM", "効果説明"},
		{"RACE:DESC", "世界観断片"},
		{"ACTI:RNAM", "操作名"},
		{"WOOP:TNAM", "語義"},
	} {
		if got := directiveOf[tc.recField]; got != tc.directive {
			t.Errorf("%s の directive = %q, want %q", tc.recField, got, tc.directive)
		}
	}

	// 9 指示文すべてが 1 件以上の REC:FIELD から参照されること（使われない指示文を残さない）。
	for k := range wantKeys {
		if !usedDirectives[k] {
			t.Errorf("指示文 %q を参照する REC:FIELD が無い", k)
		}
	}
}

// 新しい DB の口調 directive は採用した全文と一致し、人物像を埋める変数を 1 つだけ持つこと。
func TestFreshDatabaseUsesAcceptedToneDirective(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "seed2.sqlite3")
	s, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer func() { _ = s.Close() }()

	directives, err := s.ListDirectives(context.Background())
	if err != nil {
		t.Fatalf("ListDirectives: %v", err)
	}
	for _, d := range directives {
		if d.Key == "口調" {
			if d.Instruction != acceptedToneDirective {
				t.Errorf("口調 directive の instruction が採用した全文と一致しない:\ngot  %q\nwant %q", d.Instruction, acceptedToneDirective)
			}
			if count := strings.Count(d.Instruction, "{traits}"); count != 1 {
				t.Errorf("口調 directive の instruction にある {traits} = %d 個, want 1", count)
			}
			if !strings.Contains(d.Variables, "{traits}") {
				t.Errorf("口調 directive の variables に {traits} が無い: %q", d.Variables)
			}
			return
		}
	}
	t.Error("口調 directive が seed に無い")
}

// migration 18 の未編集の既定値は、migration 19 で採用した全文へ更新されること。
func TestMigration19UpdatesUneditedToneDirective(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "migration19-unedited.sqlite3")
	conn := openVersion18Database(t, dbPath)
	assertToneDirective(t, conn, migration17ToneDirective)
	closeDatabase(t, conn)

	s, err := Open(dbPath)
	if err != nil {
		t.Fatalf("migration 19 を適用する Open: %v", err)
	}
	defer func() { _ = s.Close() }()
	assertToneDirective(t, s.db, acceptedToneDirective)
}

// migration 18 までに利用者が編集した口調 directive は、空文字を含めて migration 19 が変更しないこと。
func TestMigration19PreservesEditedToneDirective(t *testing.T) {
	for _, tc := range []struct {
		name        string
		instruction string
	}{
		{name: "独自の文面", instruction: "利用者が編集した口調指示"},
		{name: "空文字", instruction: ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dbPath := filepath.Join(t.TempDir(), "migration19-edited.sqlite3")
			conn := openVersion18Database(t, dbPath)
			if _, err := conn.Exec(`UPDATE directive SET instruction = ? WHERE key = '口調'`, tc.instruction); err != nil {
				t.Fatalf("編集済み口調 directive の準備: %v", err)
			}
			closeDatabase(t, conn)

			s, err := Open(dbPath)
			if err != nil {
				t.Fatalf("migration 19 を適用する Open: %v", err)
			}
			defer func() { _ = s.Close() }()
			assertToneDirective(t, s.db, tc.instruction)
		})
	}
}

// migration 19 は口調 directive 以外の指示文と、全翻訳に使う base 指示を変更しないこと。
func TestMigration19LeavesOtherDirectivesAndBaseDirectiveUntouched(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "migration19-unrelated.sqlite3")
	conn := openVersion18Database(t, dbPath)

	otherDirectivesBefore := directiveInstructionsExceptTone(t, conn)
	var baseDirectiveBefore string
	if err := conn.Get(&baseDirectiveBefore, `SELECT base_directive FROM prompt_template WHERE id = 1`); err != nil {
		t.Fatalf("migration 19 適用前の base 指示の取得: %v", err)
	}
	closeDatabase(t, conn)

	s, err := Open(dbPath)
	if err != nil {
		t.Fatalf("migration 19 を適用する Open: %v", err)
	}
	defer func() { _ = s.Close() }()

	otherDirectivesAfter := directiveInstructionsExceptTone(t, s.db)
	if len(otherDirectivesAfter) != len(otherDirectivesBefore) {
		t.Errorf("口調以外の directive 数が変わった: got %d, want %d", len(otherDirectivesAfter), len(otherDirectivesBefore))
	}
	for key, before := range otherDirectivesBefore {
		if after, ok := otherDirectivesAfter[key]; !ok || after != before {
			t.Errorf("directive %q が変わった:\ngot  %q\nwant %q", key, after, before)
		}
	}
	var baseDirectiveAfter string
	if err := s.db.Get(&baseDirectiveAfter, `SELECT base_directive FROM prompt_template WHERE id = 1`); err != nil {
		t.Fatalf("migration 19 適用後の base 指示の取得: %v", err)
	}
	if baseDirectiveAfter != baseDirectiveBefore {
		t.Errorf("base 指示が変わった:\ngot  %q\nwant %q", baseDirectiveAfter, baseDirectiveBefore)
	}
}

// openVersion18Database は migration 19 の移行境界を独立して検証するため、18 本目までを明示的に適用する。
func openVersion18Database(t *testing.T, dbPath string) *sqlx.DB {
	t.Helper()
	conn, err := sqlx.Connect("sqlite", dbPath)
	if err != nil {
		t.Fatalf("version 18 DB を開く: %v", err)
	}
	migrations, err := db.Migrations()
	if err != nil {
		closeDatabase(t, conn)
		t.Fatalf("migration 一覧の取得: %v", err)
	}
	if len(migrations) < 19 {
		closeDatabase(t, conn)
		t.Fatalf("migration 数 = %d, want 19 以上", len(migrations))
	}
	for _, migration := range migrations[:18] {
		if _, err := conn.Exec(migration.SQL); err != nil {
			closeDatabase(t, conn)
			t.Fatalf("version 18 DB の準備で migration %s を適用: %v", migration.Name, err)
		}
	}
	if _, err := conn.Exec("PRAGMA user_version = 18"); err != nil {
		closeDatabase(t, conn)
		t.Fatalf("version 18 の記録: %v", err)
	}
	return conn
}

func assertToneDirective(t *testing.T, conn *sqlx.DB, want string) {
	t.Helper()
	var got string
	if err := conn.Get(&got, `SELECT instruction FROM directive WHERE key = '口調'`); err != nil {
		t.Fatalf("口調 directive の取得: %v", err)
	}
	if got != want {
		t.Errorf("口調 directive が想定外:\ngot  %q\nwant %q", got, want)
	}
}

func directiveInstructionsExceptTone(t *testing.T, conn *sqlx.DB) map[string]string {
	t.Helper()
	rows, err := conn.Queryx(`SELECT key, instruction FROM directive WHERE key <> '口調'`)
	if err != nil {
		t.Fatalf("口調以外の directive の取得: %v", err)
	}
	defer func() { _ = rows.Close() }()

	got := make(map[string]string)
	for rows.Next() {
		var key, instruction string
		if err := rows.Scan(&key, &instruction); err != nil {
			t.Fatalf("口調以外の directive の読み取り: %v", err)
		}
		got[key] = instruction
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("口調以外の directive の列挙: %v", err)
	}
	return got
}

func closeDatabase(t *testing.T, conn *sqlx.DB) {
	t.Helper()
	if err := conn.Close(); err != nil {
		t.Fatalf("DB のクローズ: %v", err)
	}
}

// GetPromptTemplate は prompt_template（0004）と tone_default（0007）を結合して読む。
// migration 0007 の口調設定 seed が入り、SavePromptTemplate で往復保存できることを確かめる。
func TestPromptTemplateToneDefaultsRoundTrip(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "seed3.sqlite3")
	s, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer func() { _ = s.Close() }()
	ctx := context.Background()

	got, err := s.GetPromptTemplate(ctx)
	if err != nil {
		t.Fatalf("GetPromptTemplate: %v", err)
	}
	if !strings.Contains(got.GenericToneText, "衛兵") {
		t.Errorf("汎用口調の seed が想定外: %q", got.GenericToneText)
	}
	if !strings.Contains(got.PcToneText, "プレイヤー") {
		t.Errorf("PC 口調の seed が想定外: %q", got.PcToneText)
	}
	if got.PcSex != "" {
		t.Errorf("PC 性別の既定 = %q, want 空", got.PcSex)
	}

	// 口調設定を書き換えて保存し、読み直して反映を確かめる（base 指示は据え置き）。
	got.GenericToneText = "威厳のある口調で訳す。"
	got.PcSex = "Male"
	if err = s.SavePromptTemplate(ctx, got); err != nil {
		t.Fatalf("SavePromptTemplate: %v", err)
	}
	after, err := s.GetPromptTemplate(ctx)
	if err != nil {
		t.Fatalf("GetPromptTemplate(再取得): %v", err)
	}
	if after.GenericToneText != "威厳のある口調で訳す。" || after.PcSex != "Male" {
		t.Errorf("保存が反映されない: generic=%q pcSex=%q", after.GenericToneText, after.PcSex)
	}
	if after.BaseDirective != got.BaseDirective {
		t.Errorf("base 指示が変わった: %q != %q", after.BaseDirective, got.BaseDirective)
	}
}
