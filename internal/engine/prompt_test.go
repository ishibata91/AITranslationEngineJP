package engine

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
