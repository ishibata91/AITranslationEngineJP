package runtimetag

import (
	"reflect"
	"testing"
)

// Mask はタグの無い本文を素通しし、退避タグ列を空にする（保護対象が無ければ本文を変えない）。
func TestMaskNoTags(t *testing.T) {
	masked, tags := Mask("Hello world, no tags here.")
	if masked != "Hello world, no tags here." {
		t.Fatalf("タグ無し本文が書き換わった: %q", masked)
	}
	if len(tags) != 0 {
		t.Fatalf("退避タグが空でない: %v", tags)
	}
}

// Mask は 1 つの実行時タグを退避トークン ⟦0⟧ へ置き、原タグを退避列へ残す（AI へタグ内部を渡さない）。
func TestMaskSingleTag(t *testing.T) {
	masked, tags := Mask("Bring it to <Alias=Player>.")
	if masked != "Bring it to ⟦0⟧." {
		t.Fatalf("退避結果が期待と違う: %q", masked)
	}
	if !reflect.DeepEqual(tags, []string{"<Alias=Player>"}) {
		t.Fatalf("退避タグ列が期待と違う: %v", tags)
	}
}

// Mask は複数タグを出現順に ⟦0⟧⟦1⟧… へ退避する（順序を退避列の添字と一致させる）。
func TestMaskMultipleTagsInOrder(t *testing.T) {
	masked, tags := Mask("<Global=GoldAmount> gold for <Alias=Merchant>.")
	if masked != "⟦0⟧ gold for ⟦1⟧." {
		t.Fatalf("複数タグの退避結果が期待と違う: %q", masked)
	}
	want := []string{"<Global=GoldAmount>", "<Alias=Merchant>"}
	if !reflect.DeepEqual(tags, want) {
		t.Fatalf("退避タグ列が期待と違う: %v", tags)
	}
}

// Unmask はモデルがプレースホルダを保持した出力を原タグへ復元し、欠落 0 を返す（正常な往復）。
func TestUnmaskRestoresAndReportsNoLoss(t *testing.T) {
	_, tags := Mask("Bring it to <Alias=Player>.")
	// モデルが訳語の中でプレースホルダを保持した想定の出力。
	restored, lost := Unmask("⟦0⟧ に届けよ。", tags)
	if restored != "<Alias=Player> に届けよ。" {
		t.Fatalf("復元結果が期待と違う: %q", restored)
	}
	if lost != 0 {
		t.Fatalf("欠落が誤検出された: lost=%d", lost)
	}
}

// Unmask はモデルがプレースホルダを落とした出力で、そのタグを欠落として数える（改変・削除の検出）。
func TestUnmaskDetectsLostTag(t *testing.T) {
	_, tags := Mask("<Global=Gold> gold for <Alias=Merchant>.")
	// モデルが 2 つ目のプレースホルダ ⟦1⟧ を落とした（削除した）出力。
	restored, lost := Unmask("⟦0⟧ の金貨を商人へ。", tags)
	if lost != 1 {
		t.Fatalf("欠落数が期待と違う: got=%d want=1", lost)
	}
	// 残ったプレースホルダは復元し、落ちた分は復元されない（原文の一部が失われたことが結果に残る）。
	if restored != "<Global=Gold> の金貨を商人へ。" {
		t.Fatalf("復元結果が期待と違う: %q", restored)
	}
}

// Unmask は並び替えられたプレースホルダを、位置でなく連番で照合して正しく復元する。
func TestUnmaskHandlesReorderedPlaceholders(t *testing.T) {
	_, tags := Mask("<Global=Gold> for <Alias=Merchant>")
	// モデルが語順を変え、⟦1⟧ を先に出した出力。
	restored, lost := Unmask("⟦1⟧ へ ⟦0⟧", tags)
	if lost != 0 {
		t.Fatalf("欠落が誤検出された: lost=%d", lost)
	}
	if restored != "<Alias=Merchant> へ <Global=Gold>" {
		t.Fatalf("並び替え時の復元が期待と違う: %q", restored)
	}
}

// Unmask はモデルが同じプレースホルダを重複させた出力でも、全て原タグへ復元する。
func TestUnmaskHandlesDuplicatedPlaceholder(t *testing.T) {
	_, tags := Mask("Speak to <Alias=NPC>.")
	restored, lost := Unmask("⟦0⟧ と ⟦0⟧ に話せ。", tags)
	if lost != 0 {
		t.Fatalf("欠落が誤検出された: lost=%d", lost)
	}
	if restored != "<Alias=NPC> と <Alias=NPC> に話せ。" {
		t.Fatalf("重複時の復元が期待と違う: %q", restored)
	}
}

// Placeholders は本文中の退避トークンを出現順に返す（fake provider の保持模倣に使う）。
func TestPlaceholdersReturnsTokensInOrder(t *testing.T) {
	got := Placeholders("訳文 ⟦1⟧ と ⟦0⟧ を含む")
	want := []string{"⟦1⟧", "⟦0⟧"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("退避トークンの抽出が期待と違う: %v", got)
	}
	if p := Placeholders("トークン無し"); len(p) != 0 {
		t.Fatalf("退避トークンが無いのに返った: %v", p)
	}
}

// HasPlaceholder は退避トークンの有無を返す（プロンプト保護指示を付けるかの判定に使う）。
func TestHasPlaceholder(t *testing.T) {
	if !HasPlaceholder("届け先は ⟦0⟧ だ") {
		t.Fatalf("退避トークンありを検出できない")
	}
	if HasPlaceholder("タグの無い本文") {
		t.Fatalf("退避トークン無しを誤検出した")
	}
}

// GuardInstruction は退避トークンの保持を求める非空の指示文を返す。
func TestGuardInstruction(t *testing.T) {
	if GuardInstruction() == "" {
		t.Fatalf("保護指示文が空")
	}
}

// Mask→Unmask の往復は、全プレースホルダを保持した出力なら原文のタグを完全に復元する（不変条件）。
func TestMaskUnmaskRoundTrip(t *testing.T) {
	src := "Give <Global=Amount> to <Alias=Player> at <Alias=Location>."
	masked, tags := Mask(src)
	restored, lost := Unmask(masked, tags) // 退避テキストをそのまま戻せば原文に一致する。
	if lost != 0 {
		t.Fatalf("往復で欠落が出た: lost=%d", lost)
	}
	if restored != src {
		t.Fatalf("往復で原文へ戻らない: got=%q want=%q", restored, src)
	}
}
