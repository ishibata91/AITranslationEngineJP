package api

import (
	"strings"
	"testing"

	"aitranslationenginejp/internal/model"
)

// 訳状態コードを画面表示ラベルへ写すこと（xTranslator の Status 値域に従う）。
func TestStatusLabel(t *testing.T) {
	cases := map[int]string{
		0: "未訳",
		1: "訳済",
		2: "部分",
		3: "仮訳",
		4: "承認",
		9: "未訳", // 未知は未訳扱い
	}
	for code, want := range cases {
		if got := statusLabel(code); got != want {
			t.Errorf("statusLabel(%d) = %q, want %q", code, got, want)
		}
	}
}

// narration を結果一覧の行へ写すこと。叙述文は口調指示を持たない。
func TestNarrationResultView(t *testing.T) {
	v := narrationResultView(model.Narration{
		EDID: "DLC1BookSerana", Source: "halls", Dest: "広間", Status: 3,
	})
	if v.EDID != "DLC1BookSerana" || v.Source != "halls" || v.Dest != "広間" || v.StatusLabel != "仮訳" {
		t.Errorf("narrationResultView = %+v", v)
	}
	if v.Directive != "" || v.PersonaLabel != "" {
		t.Errorf("叙述文に口調指示が付いた: %+v", v)
	}
}

// extractor 子プロセスの引数を組み立てること（dotnet run で extractor を起動）。
func TestBuildExtractorArgs(t *testing.T) {
	args := buildExtractorArgs("tools/extractor", "/sky/Data", "Dawnguard.esm", "db/c.sqlite3", "db/migrations")
	joined := strings.Join(args, " ")
	want := "run --project tools/extractor -- --data /sky/Data --plugin Dawnguard.esm --sqlite db/c.sqlite3 --schema db/migrations"
	if joined != want {
		t.Errorf("args = %q, want %q", joined, want)
	}
}
