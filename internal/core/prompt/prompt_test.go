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

// FillVariables は directive の指示文中の変数トークンを vars の値へ差し込むこと。
// 口調 directive の {traits} に話者の性質を埋める経路と、変数なしの directive（固有名・定型句・文体）で
// 指示文をそのまま返す経路の両方を確かめる。
func TestFillVariables(t *testing.T) {
	cases := []struct {
		name        string
		instruction string
		vars        map[string]string
		want        string
	}{
		{
			name:        "口調 directive の {traits} を性質で埋める",
			instruction: "この台詞の話者の人物像:\n{traits}\nこの人物像に合う口調で訳すこと。",
			vars:        map[string]string{"{traits}": "- 口調: 柔らかく丁寧"},
			want:        "この台詞の話者の人物像:\n- 口調: 柔らかく丁寧\nこの人物像に合う口調で訳すこと。",
		},
		{
			name:        "変数なし（vars 空）は指示文をそのまま返す",
			instruction: "これは固有名詞です。簡潔に訳すこと。",
			vars:        nil,
			want:        "これは固有名詞です。簡潔に訳すこと。",
		},
		{
			name:        "宣言外のトークンはそのまま残す",
			instruction: "説明文を訳すこと。{unknown} は残す。",
			vars:        map[string]string{"{traits}": "x"},
			want:        "説明文を訳すこと。{unknown} は残す。",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := FillVariables(c.instruction, c.vars); got != c.want {
				t.Errorf("FillVariables = %q, want %q", got, c.want)
			}
		})
	}
}

// プロンプト合成（MECE）: Base 指示 ＋ 変数を埋めた directive 指示文 ＋ 機械置換済み原文。
// 口調 directive（変数あり）と固有名 directive（変数なし）の両経路で完成プロンプトを確かめる。
func TestComposeWithFilledDirective(t *testing.T) {
	const base = "あなたは翻訳者です。"

	// 口調 directive: {traits} を埋めてから base へ合成する。
	tone := FillVariables("人物像:\n{traits}", map[string]string{"{traits}": "- 丁寧"})
	got := ComposePrompt(base, tone, "Hello")
	if got.System != base+"\n\n人物像:\n- 丁寧" {
		t.Errorf("口調合成の System = %q", got.System)
	}

	// 固有名 directive: 変数なしで、base へ指示文をそのまま合成する。
	proper := FillVariables("これは固有名詞です。", nil)
	got = ComposePrompt(base, proper, "Dragonbane")
	if got.System != base+"\n\nこれは固有名詞です。" || got.User != "Dragonbane" {
		t.Errorf("固有名合成 = %+v", got)
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
