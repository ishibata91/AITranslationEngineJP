package engine

import "testing"

// xTranslator XML から NPC のフルネーム対・短名対と会話文を取り出すこと。
// 先頭 BOM の除去と、エンティティ（&lt;）の復号、対象 REC 以外の素通しを確かめる。
func TestParseTermXML(t *testing.T) {
	const fixture = "\uFEFF" + `<?xml version="1.0" encoding="UTF-8"?>
<SSTXMLRessources>
  <Content>
    <String><REC>NPC_:FULL</REC><Source>Grelod the Kind</Source><Dest>親切者のグレロッド</Dest></String>
    <String><REC>NPC_:SHRT</REC><Source>Maro</Source><Dest>マロ</Dest></String>
    <String><REC>INFO:NAM1</REC><Source>Bring it to &lt;Alias&gt;.</Source><Dest>x</Dest></String>
    <String><REC>WEAP:FULL</REC><Source>Iron Sword</Source><Dest>鉄の剣</Dest></String>
  </Content>
</SSTXMLRessources>`

	fulls, shrts, dialogues, err := parseTermXML([]byte(fixture), true)
	if err != nil {
		t.Fatalf("parseTermXML error: %v", err)
	}
	// NPC_:FULL だけを fulls に取り、base ゲーム印を付ける。
	if len(fulls) != 1 || fulls[0].Source != "Grelod the Kind" || !fulls[0].BaseGame {
		t.Errorf("fulls = %+v, want 1 件の Grelod the Kind（base ゲーム）", fulls)
	}
	// NPC_:SHRT だけを shrts に取る。
	if len(shrts) != 1 || shrts[0].Source != "Maro" || shrts[0].Dest != "マロ" {
		t.Errorf("shrts = %+v, want 1 件の Maro=>マロ", shrts)
	}
	// INFO:NAM1 の英語原文を会話文に取り、エンティティを復号する。
	if len(dialogues) != 1 || dialogues[0] != "Bring it to <Alias>." {
		t.Errorf("dialogues = %v, want 1 件の復号済み原文", dialogues)
	}
}
