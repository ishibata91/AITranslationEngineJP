package prompt

import (
	"strings"
	"testing"
)

// ComposePrompt は base 指示・口調指示・原文から完成プロンプトを組むこと。
// 口調指示があれば base 指示の後ろへ 1 行空けて合成し、無ければ base 指示だけにする。原文は user へそのまま入れる。
func TestComposePrompt(t *testing.T) {
	const base = "base 指示"
	cases := []struct {
		name       string
		persona    string
		source     string
		wantSystem string
		wantUser   string
	}{
		{
			name:       "口調指示ありは base の後ろへ合成する",
			persona:    "この台詞の話者の人物像:\n- 声質: 幼い少年の声",
			source:     "The リフテン guard.",
			wantSystem: "base 指示\n\nこの台詞の話者の人物像:\n- 声質: 幼い少年の声",
			wantUser:   "The リフテン guard.",
		},
		{
			name:       "口調指示なしは base 指示だけにする",
			persona:    "",
			source:     "halls",
			wantSystem: "base 指示",
			wantUser:   "halls",
		},
		{
			name:       "空白だけの口調指示は無し扱いにする",
			persona:    "   \n  ",
			source:     "cairn",
			wantSystem: "base 指示",
			wantUser:   "cairn",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := ComposePrompt(base, c.persona, c.source)
			if got.System != c.wantSystem {
				t.Errorf("System = %q, want %q", got.System, c.wantSystem)
			}
			if got.User != c.wantUser {
				t.Errorf("User = %q, want %q", got.User, c.wantUser)
			}
		})
	}
}

// ComposePrompt は原文に実行時タグ（<...>）があれば system へタグ保護指示を足し、
// 無ければ足さない（タグを持つ本文にだけモデルへタグの原形保持を求める）。
func TestComposePromptAddsTagGuardWhenTagPresent(t *testing.T) {
	const base = "base 指示"
	// 生タグありは保護指示が system へ足される。
	withTag := ComposePrompt(base, "", "届け先は <Alias=Player> だ")
	if !strings.Contains(withTag.System, "残すこと") {
		t.Fatalf("生タグありでタグ保護指示が付いていない: %q", withTag.System)
	}
	// user はタグを退避せず生のまま入れる。
	if !strings.Contains(withTag.User, "<Alias=Player>") {
		t.Fatalf("user に生タグが入っていない: %q", withTag.User)
	}
	// タグなしは base のまま（保護指示を足さない）。
	noTag := ComposePrompt(base, "", "普通の本文")
	if noTag.System != base {
		t.Fatalf("タグなしで system が変わった: %q", noTag.System)
	}
}

// プロンプト合成（MECE）: Base 指示 ＋ 差し込み済みの directive 指示文 ＋ 機械置換済み原文。
// 口調 directive（差し込み済み）と固有名 directive（差し込みなし）の両経路で完成プロンプトを確かめる。
// {traits} の差し込みは personatone.BuildToneDirective が担うため、ここは差し込み後の文字列を入力にする。
func TestComposeWithFilledDirective(t *testing.T) {
	const base = "あなたは翻訳者です。"

	// 口調 directive: 話者の性質を差し込んだ状態で base へ合成する。
	got := ComposePrompt(base, "人物像:\n- 丁寧", "Hello")
	if got.System != base+"\n\n人物像:\n- 丁寧" {
		t.Errorf("口調合成の System = %q", got.System)
	}

	// 固有名 directive: 差し込む変数を持たず、base へ指示文をそのまま合成する。
	got = ComposePrompt(base, "これは固有名詞です。", "Dragonbane")
	if got.System != base+"\n\nこれは固有名詞です。" || got.User != "Dragonbane" {
		t.Errorf("固有名合成 = %+v", got)
	}
}

// ComposeBodyPrompt は本文を英語のまま user へ置き、辞書候補の表示に meaning を混ぜない。
func TestComposeBodyPromptKeepsSourceAndOmitsMeaning(t *testing.T) {
	p := ComposeBodyPrompt("base", "directive", "The Riften guard waited.", []BodyReference{{
		Source: "Riften", Dest: "リフテン", PartOfSpeech: "noun", SkyrimCategory: "city", Origin: "事前作成済み翻訳辞書",
	}})
	if p.User != "The Riften guard waited." {
		t.Fatalf("本文が変わった: %q", p.User)
	}
	for _, want := range []string{"Riften -> リフテン", "品詞: noun", "Skyrimカテゴリ: city", "出どころ: 事前作成済み翻訳辞書"} {
		if !strings.Contains(p.System, want) {
			t.Errorf("system に %q が無い: %q", want, p.System)
		}
	}
	if strings.Contains(p.System, "城塞を守る都市") {
		t.Errorf("meaning が system に出た: %q", p.System)
	}
}

// RenderPrompt は完成プロンプトを役割見出し付きの 1 文字列へ描き、system と user の両方を全文含めること。
func TestRenderPrompt(t *testing.T) {
	got := RenderPrompt(ComposePrompt("base 指示", "口調指示", "原文"))
	want := "system:\nbase 指示\n\n口調指示\n\nuser:\n原文"
	if got != want {
		t.Errorf("RenderPrompt = %q, want %q", got, want)
	}
	// system と user の両方の中身が表示文字列に残ること（実プロンプト目視確認の前提）。
	if !strings.Contains(got, "口調指示") || !strings.Contains(got, "原文") {
		t.Errorf("RenderPrompt に system/user の中身が欠けた: %q", got)
	}
}
