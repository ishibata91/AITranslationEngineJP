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

// Mask は 1 つの実行時タグを退避トークン ⟦0⟧ へ置き、原タグを退避列へ残す（機械置換からタグ内部を守る）。
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

// Mask→Restore の往復は、機械置換を通した後でも原文のタグを原形へ戻す（退避は一時措置、AI へは生タグを送る）。
func TestMaskRestoreRoundTrip(t *testing.T) {
	src := "Give <Global=Amount> to <Alias=Player> at <Alias=Location>."
	masked, tags := Mask(src)
	// masked は機械置換を通してもプレースホルダ ⟦i⟧ は当たらない（辞書語でない）。ここでは素通しを仮定。
	restored := Restore(masked, tags)
	if restored != src {
		t.Fatalf("往復で原文へ戻らない: got=%q want=%q", restored, src)
	}
}

// Restore は機械置換で周囲が変わっても、退避トークンだけを原タグへ戻す。
func TestRestoreAfterReplacementAroundToken(t *testing.T) {
	// "The Riften guard <Alias=NPC>" を Mask し、Riften→リフテンの置換が起きた後を模す。
	_, tags := Mask("The Riften guard <Alias=NPC>")
	restored := Restore("The リフテン guard ⟦0⟧", tags)
	if restored != "The リフテン guard <Alias=NPC>" {
		t.Fatalf("周囲置換後の復元が期待と違う: %q", restored)
	}
}

// Tags は本文中の生タグを出現順に返す（fake provider のタグ保持模倣・照合に使う）。
func TestTags(t *testing.T) {
	got := Tags("Give <Global=Amount> to <Alias=Player>")
	want := []string{"<Global=Amount>", "<Alias=Player>"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("生タグの抽出が期待と違う: %v", got)
	}
	if p := Tags("タグ無し"); len(p) != 0 {
		t.Fatalf("タグが無いのに返った: %v", p)
	}
}

// HasTag は生タグの有無を返す（プロンプト保護指示を付けるかの判定に使う）。
func TestHasTag(t *testing.T) {
	if !HasTag("届け先は <Alias=Player> だ") {
		t.Fatalf("生タグありを検出できない")
	}
	if HasTag("タグの無い本文") {
		t.Fatalf("生タグ無しを誤検出した")
	}
}

// CountMissing は AI 出力に元タグが原形で残っていない件数を返す（削除・改変の検出）。
func TestCountMissing(t *testing.T) {
	_, tags := Mask("<Global=Gold> gold for <Alias=Merchant>.")
	// モデルが 2 つ目のタグ <Alias=Merchant> を落とした出力。
	if got := CountMissing("<Global=Gold> の金貨を商人へ。", tags); got != 1 {
		t.Fatalf("欠落数が期待と違う: got=%d want=1", got)
	}
	// 全タグが原形で残る出力は欠落 0。
	if got := CountMissing("<Global=Gold> の金貨を <Alias=Merchant> へ。", tags); got != 0 {
		t.Fatalf("欠落が誤検出された: got=%d want=0", got)
	}
}

// CountMissing は同じタグが複数回出る本文で、出力側の回数不足を欠落として数える。
func TestCountMissingCountsOccurrences(t *testing.T) {
	_, tags := Mask("<Alias=NPC> と <Alias=NPC>")
	// 出力にタグが 1 回しかない（2 回必要）→ 欠落 1。
	if got := CountMissing("<Alias=NPC> と彼", tags); got != 1 {
		t.Fatalf("出現回数の欠落数が期待と違う: got=%d want=1", got)
	}
}

// GuardInstruction は生タグの保持を求める非空の指示文を返す。
func TestGuardInstruction(t *testing.T) {
	if GuardInstruction() == "" {
		t.Fatalf("保護指示文が空")
	}
}
