package service

import "testing"

// TestIsAllowedImportREC_only13Types は isAllowedImportREC が
// 13 種別だけ true を返し、DOOR:FULL / FLOR:FULL / FURN:FULL を含む
// 13 種別外では false を返すことを証明する。
func TestIsAllowedImportREC_only13Types(t *testing.T) {
	// isAllowedImportREC は recclassification.IsTermTarget の結果をそのまま返す責務境界を持つ。
	// DOOR:FULL / FLOR:FULL / FURN:FULL は旧仕様で許可されていたが、
	// master-dictionary-REQ-004 で対象外になった。

	cases := []struct {
		rec  string
		want bool
	}{
		// 13 種別 — すべて true
		{"BOOK:FULL", true},
		{"NPC_:FULL", true},
		{"NPC_:SHRT", true},
		{"ARMO:FULL", true},
		{"WEAP:FULL", true},
		{"LCTN:FULL", true},
		{"CELL:FULL", true},
		{"CONT:FULL", true},
		{"MISC:FULL", true},
		{"ALCH:FULL", true},
		{"RACE:FULL", true},
		{"INGR:FULL", true},
		{"SHOU:FULL", true},
		// 旧仕様で許可されていた 3 種別 — false に変わった
		{"DOOR:FULL", false},
		{"FLOR:FULL", false},
		{"FURN:FULL", false},
		// 形式違反
		{"", false},
		{"NPC_", false},
		{":FULL", false},
		{"npc_:full", false},
	}

	for _, c := range cases {
		t.Run(c.rec, func(t *testing.T) {
			got := isAllowedImportREC(c.rec)
			if got != c.want {
				t.Errorf("isAllowedImportREC(%q) = %v, want %v", c.rec, got, c.want)
			}
		})
	}
}

// TestCategoryFromREC_noDOORFLORFURNBranch は categoryFromREC が
// DOOR:FULL / FLOR:FULL / FURN:FULL に対して専用分岐を持たず、
// default の "その他" を返すことを証明する。
// これは master-dictionary-REQ-004 に基づき DOOR/FLOR/FURN 分岐が削除されたことの確認テスト。
func TestCategoryFromREC_noDOORFLORFURNBranch(t *testing.T) {
	// DOOR / FLOR / FURN が "その他" にフォールスルーすることで、
	// categoryFromREC から専用 case が消えていることを確認する。

	cases := []struct {
		rec  string
		want string
	}{
		{"DOOR:FULL", "その他"},
		{"FLOR:FULL", "その他"},
		{"FURN:FULL", "その他"},
	}

	for _, c := range cases {
		t.Run(c.rec, func(t *testing.T) {
			got := categoryFromREC(c.rec)
			if got != c.want {
				t.Errorf("categoryFromREC(%q) = %q, want %q (DOOR/FLOR/FURN 専用分岐が残っている可能性がある)", c.rec, got, c.want)
			}
		})
	}
}

// TestCategoryFromREC_13Types は categoryFromREC が
// 13 種別に対して仕様通りのカテゴリ文字列を返すことをテーブル駆動で証明する。
func TestCategoryFromREC_13Types(t *testing.T) {
	// master-dictionary-REQ-004 で定義された 13 種別のカテゴリ対応を検証する。

	cases := []struct {
		rec  string
		want string
	}{
		{"BOOK:FULL", "書籍"},
		{"NPC_:FULL", "NPC"},
		{"NPC_:SHRT", "NPC"},
		{"RACE:FULL", "NPC"},
		{"ARMO:FULL", "装備"},
		{"WEAP:FULL", "装備"},
		{"LCTN:FULL", "地名"},
		{"CELL:FULL", "地名"},
		{"CONT:FULL", "アイテム"},
		{"MISC:FULL", "アイテム"},
		{"INGR:FULL", "アイテム"},
		{"ALCH:FULL", "アイテム"},
		{"SHOU:FULL", "シャウト"},
	}

	for _, c := range cases {
		t.Run(c.rec, func(t *testing.T) {
			got := categoryFromREC(c.rec)
			if got != c.want {
				t.Errorf("categoryFromREC(%q) = %q, want %q", c.rec, got, c.want)
			}
		})
	}
}
