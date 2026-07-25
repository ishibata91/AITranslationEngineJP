package store

import (
	"context"
	"path/filepath"
	"strings"
	"testing"
)

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

// 口調 directive は変数 {traits} を宣言していること（本文フェーズの台詞プロンプトで話者の性質を埋めるため）。
func TestSeedToneDirectiveHasTraitsVariable(t *testing.T) {
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
			if !strings.Contains(d.Instruction, "{traits}") {
				t.Errorf("口調 directive の instruction に {traits} が無い: %q", d.Instruction)
			}
			if !strings.Contains(d.Variables, "{traits}") {
				t.Errorf("口調 directive の variables に {traits} が無い: %q", d.Variables)
			}
			return
		}
	}
	t.Error("口調 directive が seed に無い")
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
