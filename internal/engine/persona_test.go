package engine

import (
	"strings"
	"testing"

	"aitranslationenginejp/internal/model"
)

// 口調 traits が揃うと、テンプレートの {traits} へ声質・種族・所属を箇条書きで差し込んだ口調指示文を組むこと。
func TestBuildPersonaDirective(t *testing.T) {
	got := buildPersonaDirective(testPersonaTemplate, model.SpeakerPersona{
		VoiceNature:    "幼い少年の声",
		RaceNature:     "ノルド（実直で粘り強い北方の気質）",
		FactionNatures: []string{"闇の一党の気風（冷徹さ）"},
	})
	for _, want := range []string{"声質: 幼い少年の声", "種族の気質: ノルド", "所属の気風: 闇の一党", "口調と人称で訳すこと"} {
		if !strings.Contains(got, want) {
			t.Errorf("directive に %q が無い:\n%s", want, got)
		}
	}
	// {traits} の差し込み口は置換され、テンプレートにそのまま残らないこと。
	if strings.Contains(got, traitsToken) {
		t.Errorf("差し込み口 %q が未置換で残った:\n%s", traitsToken, got)
	}
}

// traits が 1 つも無ければ空文字を返し、ペルソナ指示を注入しないこと。
func TestBuildPersonaDirectiveEmpty(t *testing.T) {
	if got := buildPersonaDirective(testPersonaTemplate, model.SpeakerPersona{}); got != "" {
		t.Errorf("空の persona で directive=%q、空を期待", got)
	}
}

// 口調チップの短い要約は、声質を最優先し、無ければ種族、所属の順で採ること。
func TestBuildPersonaLabel(t *testing.T) {
	cases := []struct {
		name string
		in   model.SpeakerPersona
		want string
	}{
		{"声質優先", model.SpeakerPersona{VoiceNature: "幼い少年の声", RaceNature: "ノルド"}, "声質: 幼い少年の声"},
		{"一般的な声は種族へ", model.SpeakerPersona{VoiceNature: "男性の声", RaceNature: "ノルドの子供（北方の気質）"}, "種族: ノルドの子供（北方の気質）"},
		{"種族へ後退", model.SpeakerPersona{RaceNature: "ノルド"}, "種族: ノルド"},
		{"所属へ後退", model.SpeakerPersona{FactionNatures: []string{"戦士団の気風"}}, "所属: 戦士団の気風"},
		{"無し", model.SpeakerPersona{}, ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := buildPersonaLabel(c.in); got != c.want {
				t.Errorf("label = %q, want %q", got, c.want)
			}
		})
	}
}
