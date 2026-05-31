package recclassification

import (
	"slices"
	"sort"
	"testing"
)

// TestIsTermTarget_includes13Types は IsTermTarget が仕様で定めた 13 種別すべてに
// true を返すことを証明する（term-translation-phase-REQ-008、master-dictionary-REQ-004）。
func TestIsTermTarget_includes13Types(t *testing.T) {
	t.Parallel()

	// 仕様で定められた 13 種別（detail-spec-diff.md term-translation-phase-REQ-008）
	cases := []struct {
		rec  string
		want bool
		desc string
	}{
		{"BOOK:FULL", true, "BOOK:FULL は翻訳対象 13 種別のひとつ"},
		{"NPC_:FULL", true, "NPC_:FULL は翻訳対象 13 種別のひとつ"},
		{"NPC_:SHRT", true, "NPC_:SHRT は翻訳対象 13 種別のひとつ"},
		{"ARMO:FULL", true, "ARMO:FULL は翻訳対象 13 種別のひとつ"},
		{"WEAP:FULL", true, "WEAP:FULL は翻訳対象 13 種別のひとつ"},
		{"LCTN:FULL", true, "LCTN:FULL は翻訳対象 13 種別のひとつ"},
		{"CELL:FULL", true, "CELL:FULL は翻訳対象 13 種別のひとつ"},
		{"CONT:FULL", true, "CONT:FULL は翻訳対象 13 種別のひとつ"},
		{"MISC:FULL", true, "MISC:FULL は翻訳対象 13 種別のひとつ"},
		{"ALCH:FULL", true, "ALCH:FULL は翻訳対象 13 種別のひとつ"},
		{"RACE:FULL", true, "RACE:FULL は翻訳対象 13 種別のひとつ"},
		{"INGR:FULL", true, "INGR:FULL は翻訳対象 13 種別のひとつ"},
		{"SHOU:FULL", true, "SHOU:FULL は翻訳対象 13 種別のひとつ"},
	}

	for _, tc := range cases {
		t.Run(tc.rec, func(t *testing.T) {
			t.Parallel()
			// Act
			got := IsTermTarget(tc.rec)
			// Assert: 13 種別の各 REC は true を返す
			if got != tc.want {
				t.Errorf("%s: IsTermTarget(%q) = %v, want %v", tc.desc, tc.rec, got, tc.want)
			}
		})
	}
}

// TestIsTermTarget_excludes13TypesOuter は IsTermTarget が仕様変更で除外された 3 種別および
// 形式違反入力に false を返すことを証明する（term-translation-phase-REQ-008）。
func TestIsTermTarget_excludes13TypesOuter(t *testing.T) {
	t.Parallel()

	cases := []struct {
		rec  string
		desc string
	}{
		// 仕様変更で除外された 3 種別
		{"DOOR:FULL", "DOOR:FULL は仕様変更で 13 種別から除外された REC"},
		{"FLOR:FULL", "FLOR:FULL は仕様変更で 13 種別から除外された REC"},
		{"FURN:FULL", "FURN:FULL は仕様変更で 13 種別から除外された REC"},
		// 形式違反入力
		{"", "空文字は REC 形式に違反する"},
		{"NPC_", "フィールド部が欠けた入力は REC 形式に違反する"},
		{":FULL", "レコード種別部が欠けた入力は REC 形式に違反する"},
		{"NPC_:", "フィールド値が空の入力は REC 形式に違反する"},
		{"npc_:full", "小文字の入力は REC 大文字表記規則に違反する"},
		{"UNKNOWN:FULL", "未知の REC は 13 種別に該当しない"},
	}

	for _, tc := range cases {
		t.Run(tc.rec+"_false", func(t *testing.T) {
			t.Parallel()
			// Act
			got := IsTermTarget(tc.rec)
			// Assert: 13 種別外入力は例外を送出せず false を返す
			if got != false {
				t.Errorf("%s: IsTermTarget(%q) = %v, want false", tc.desc, tc.rec, got)
			}
		})
	}
}

// TestIsTermTarget_NPC_FULL_and_NPC_SHRT_areIndependent は NPC_:FULL と NPC_:SHRT が
// それぞれ独立して true を返し、片方への判定が他方に影響しないことを証明する
// （term-translation-phase-REQ-008 の「別 REC として扱う」仕様）。
func TestIsTermTarget_NPC_FULL_and_NPC_SHRT_areIndependent(t *testing.T) {
	t.Parallel()

	// Arrange: NPC_ プレフィクスで 2 種の REC を用意する
	npcFull := "NPC_:FULL"
	npcShrt := "NPC_:SHRT"

	// Act
	gotFull := IsTermTarget(npcFull)
	gotShrt := IsTermTarget(npcShrt)

	// Assert: NPC_:FULL と NPC_:SHRT はいずれも独立して true になる
	if !gotFull {
		t.Errorf("NPC_:FULL は 13 種別に含まれるが IsTermTarget(%q) = false", npcFull)
	}
	if !gotShrt {
		t.Errorf("NPC_:SHRT は 13 種別に含まれるが IsTermTarget(%q) = false", npcShrt)
	}
}

// TestTermTargetRECList_returns13UniqueItems は TermTargetRECList が 13 種別を
// 重複なく返すことを証明する（SQL IN 句の単一情報源保証）。
func TestTermTargetRECList_returns13UniqueItems(t *testing.T) {
	t.Parallel()

	// Act
	list := TermTargetRECList()

	// Assert: 件数が 13 であること
	const wantCount = 13
	if len(list) != wantCount {
		t.Errorf("TermTargetRECList は %d 件を返す必要があるが %d 件だった", wantCount, len(list))
	}

	// Assert: 重複がないこと
	seen := make(map[string]struct{}, len(list))
	for _, rec := range list {
		if _, dup := seen[rec]; dup {
			t.Errorf("TermTargetRECList に重複した REC が含まれる: %q", rec)
		}
		seen[rec] = struct{}{}
	}
}

// TestTermTargetRECList_containsAll13Specs は TermTargetRECList が仕様で定めた
// 13 種別をすべて含むことを証明する（term-translation-phase-REQ-008）。
func TestTermTargetRECList_containsAll13Specs(t *testing.T) {
	t.Parallel()

	// Arrange: 仕様で定められた 13 種別（detail-spec-diff.md より）
	specRECs := []string{
		"BOOK:FULL",
		"NPC_:FULL",
		"NPC_:SHRT",
		"ARMO:FULL",
		"WEAP:FULL",
		"LCTN:FULL",
		"CELL:FULL",
		"CONT:FULL",
		"MISC:FULL",
		"ALCH:FULL",
		"RACE:FULL",
		"INGR:FULL",
		"SHOU:FULL",
	}

	// Act
	list := TermTargetRECList()

	// Assert: 仕様の各 REC がリストに含まれること
	for _, spec := range specRECs {
		if !slices.Contains(list, spec) {
			t.Errorf("TermTargetRECList に仕様 REC %q が含まれていない", spec)
		}
	}
}

// TestTermTargetRECList_excludesDoorFlorFurn は TermTargetRECList が
// 仕様変更で除外された DOOR:FULL、FLOR:FULL、FURN:FULL を含まないことを証明する。
func TestTermTargetRECList_excludesDoorFlorFurn(t *testing.T) {
	t.Parallel()

	// Arrange: 仕様変更で除外された 3 種別
	excluded := []string{"DOOR:FULL", "FLOR:FULL", "FURN:FULL"}

	// Act
	list := TermTargetRECList()

	// Assert: 除外対象が含まれないこと
	for _, rec := range excluded {
		if slices.Contains(list, rec) {
			t.Errorf("TermTargetRECList に除外 REC %q が含まれている（仕様変更で除外済み）", rec)
		}
	}
}

// TestTermTargetRECList_returnsDefensiveCopy は TermTargetRECList が呼び出しごとに
// 独立したスライスを返し、変更が内部集合に影響しないことを証明する。
func TestTermTargetRECList_returnsDefensiveCopy(t *testing.T) {
	t.Parallel()

	// Arrange: 1 回目の呼び出しで取得したリストを変更する
	list1 := TermTargetRECList()
	original := make([]string, len(list1))
	copy(original, list1)

	// Act: リストを破壊的にソートして内部集合が影響を受けるか確認する
	sort.Strings(list1)
	list2 := TermTargetRECList()

	// Assert: 2 回目の呼び出し結果が元の定義順と一致すること
	for i, rec := range original {
		if list2[i] != rec {
			t.Errorf(
				"TermTargetRECList は防御的コピーを返す必要があるが、内部集合が変更された: index %d: got %q, want %q",
				i, list2[i], rec,
			)
		}
	}
}
