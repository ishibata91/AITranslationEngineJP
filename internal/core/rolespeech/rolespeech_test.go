package rolespeech

import (
	"errors"
	"strings"
	"testing"
	"testing/iotest"
)

// RoleClassOfRace は race EditorID を child / elder / adult へ畳む。
func TestRoleClassOfRace(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"NordRaceChild", "child"},
		{"ImperialRaceChild", "child"},
		{"ElderRace", "elder"},
		{"ImperialRace", "adult"},
		{"", "adult"},
	}
	for _, tc := range cases {
		if got := RoleClassOfRace(tc.in); got != tc.want {
			t.Errorf("RoleClassOfRace(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

// ParseRoleSpeech はコメント・空行を読み飛ばし、5 列の行をテンプレートに読む。
func TestParseRoleSpeech(t *testing.T) {
	const data = "# コメント行\n\n" +
		"child\tmale\t*\tぼく\t平易に。\n" +
		"elder\tfemale\t*\tわたし\t落ち着いて。\n"
	tbl, err := ParseRoleSpeech(strings.NewReader(data))
	if err != nil {
		t.Fatalf("ParseRoleSpeech: %v", err)
	}
	if len(tbl.rows) != 2 {
		t.Fatalf("行数 = %d, want 2（コメントと空行を除く）", len(tbl.rows))
	}
	if tbl.rows[0] != (roleSpeechRow{race: "child", sex: "male", cell: "*", tmpl: Template{FirstPerson: "ぼく", Register: "平易に。"}}) {
		t.Errorf("先頭行が想定外: %+v", tbl.rows[0])
	}
}

// ParseRoleSpeech は列が 5 未満の行をエラーにする。
func TestParseRoleSpeechTooFewColumns(t *testing.T) {
	_, err := ParseRoleSpeech(strings.NewReader("child\tmale\t*\tぼく\n"))
	if err == nil {
		t.Fatal("列不足でエラーを期待")
	}
}

// ParseRoleSpeech は読み取りエラー（scanner エラー）をそのまま返す。
func TestParseRoleSpeechReadError(t *testing.T) {
	_, err := ParseRoleSpeech(iotest.ErrReader(errors.New("読み取り失敗")))
	if err == nil {
		t.Fatal("読み取りエラーを期待")
	}
}

// Lookup は完全一致・ワイルドカード・具体度優先・不一致・nil を分ける。
func TestRoleSpeechLookup(t *testing.T) {
	const data = "child\tmale\t*\tぼく\t幼く。\n" +
		"adult\tfemale\t*\tわたし\tやわらかく。\n" +
		"adult\tfemale\t冷然・見下し\tわたし\t冷たく。\n" +
		"*\t*\t*\t\t標準。\n"
	tbl, err := ParseRoleSpeech(strings.NewReader(data))
	if err != nil {
		t.Fatalf("ParseRoleSpeech: %v", err)
	}

	// 完全一致（性別・年齢でテンプレートを引く）。
	if got, ok := tbl.Lookup("child", "male", "平明"); !ok || got.FirstPerson != "ぼく" {
		t.Errorf("child/male = (%+v, %v), want ぼく", got, ok)
	}
	// セル指定の行が、セル無指定（*）の行より具体的で優先される。
	if got, ok := tbl.Lookup("adult", "female", "冷然・見下し"); !ok || got.Register != "冷たく。" {
		t.Errorf("adult/female/冷然 = (%+v, %v), want 冷たく。（具体度優先）", got, ok)
	}
	// セルが一致しない場合は cell=* の行へ落ちる。
	if got, ok := tbl.Lookup("adult", "female", "物腰やわ"); !ok || got.Register != "やわらかく。" {
		t.Errorf("adult/female/物腰やわ = (%+v, %v), want やわらかく。", got, ok)
	}
	// どの具体行にも当たらない組は、全ワイルドカードの既定行へ落ちる。
	if got, ok := tbl.Lookup("adult", "male", "平明"); !ok || got.Register != "標準。" {
		t.Errorf("adult/male = (%+v, %v), want 標準。（* * * へ）", got, ok)
	}

	// 既定行が無ければ不一致は ok=false。
	noDefault, err := ParseRoleSpeech(strings.NewReader("elder\tfemale\t*\tわたし\t落ち着いて。\n"))
	if err != nil {
		t.Fatalf("ParseRoleSpeech: %v", err)
	}
	if _, ok := noDefault.Lookup("adult", "male", "平明"); ok {
		t.Error("一致行が無いのに ok=true")
	}

	// nil レシーバは安全に ok=false。
	var nilTable *Table
	if _, ok := nilTable.Lookup("child", "male", "平明"); ok {
		t.Error("nil テーブルで ok=true")
	}
}
