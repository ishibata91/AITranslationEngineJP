package engine

import "testing"

// SelectSupply は固有名の既訳の有無から供給源を選ぶこと。
// 既訳あり → 権威訳（Dest を埋め、書き込み宛先なし）、既訳なし → AI 訳（宛先は proper_noun）。
func TestSelectSupply(t *testing.T) {
	authoritative := map[string]string{
		"Riften":     "リフテン",
		"Dragonbane": "竜の災い",
	}

	// 既訳あり: 権威訳を返し、AI 訳しない（WriteTarget は空＝書き戻し不要）。
	hit := SelectSupply("Riften", authoritative)
	if hit.Kind != SupplyAuthoritative {
		t.Errorf("既訳ありの Kind = %d, want SupplyAuthoritative", hit.Kind)
	}
	if hit.Dest != "リフテン" {
		t.Errorf("既訳ありの Dest = %q, want リフテン", hit.Dest)
	}
	if hit.WriteTarget != "" {
		t.Errorf("既訳ありの WriteTarget = %q, want 空", hit.WriteTarget)
	}

	// 既訳なし: AI 訳対象で、宛先は proper_noun（master_term ではない）。
	miss := SelectSupply("Whiterun", authoritative)
	if miss.Kind != SupplyAITranslate {
		t.Errorf("既訳なしの Kind = %d, want SupplyAITranslate", miss.Kind)
	}
	if miss.Dest != "" {
		t.Errorf("既訳なしの Dest = %q, want 空", miss.Dest)
	}
	if miss.WriteTarget != properNounTable {
		t.Errorf("既訳なしの WriteTarget = %q, want %q", miss.WriteTarget, properNounTable)
	}
}

// 不変ルール: どの入力でも、AI 訳の書き込み宛先は proper_noun で、master_term を返さない（方針A）。
// 既訳あり・なしの両経路を網羅し、宛先が master_term になり得ないことを確かめる。
func TestSelectSupplyNeverWritesMasterTerm(t *testing.T) {
	authoritative := map[string]string{"Known": "既訳"}
	for _, source := range []string{"Known", "Unknown", "", "Another"} {
		sup := SelectSupply(source, authoritative)
		if sup.WriteTarget == masterTermTable {
			t.Errorf("source=%q で WriteTarget が master_term になった（不変ルール違反）", source)
		}
		if sup.Kind == SupplyAITranslate && sup.WriteTarget != properNounTable {
			t.Errorf("source=%q の AI 訳宛先 = %q, want %q", source, sup.WriteTarget, properNounTable)
		}
	}
}

// 空の既訳辞書では、すべて AI 訳（既訳ヒットなし）になること（境界）。
func TestSelectSupplyEmptyAuthoritative(t *testing.T) {
	sup := SelectSupply("Anything", map[string]string{})
	if sup.Kind != SupplyAITranslate || sup.WriteTarget != properNounTable {
		t.Errorf("空辞書での判定 = %+v, want AI 訳→proper_noun", sup)
	}
}
