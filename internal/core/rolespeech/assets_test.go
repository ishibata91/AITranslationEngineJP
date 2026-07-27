package rolespeech

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// 本 file だけは repo の実 asset（assets/role-speech.tsv・assets/role-speech-examples.tsv）を読む。
// 照合規則そのものは rolespeech_test.go が合成入力で守り、ここは実表に行の欠落が無いことだけを守る。
// harness は合成表を使うため実表の欠落を検出しない（internal/harness/run.go）。その穴をここで埋める。

// 実 asset を読み、役割語と例文を束ねた Table を返す。
func realTable(t *testing.T) *Table {
	t.Helper()
	roleFile, err := os.Open(filepath.Join("..", "..", "..", "assets", "role-speech.tsv"))
	if err != nil {
		t.Fatalf("役割語テンプレートを開けない: %v", err)
	}
	defer roleFile.Close() //nolint:errcheck // 読み取り後の後始末。
	tbl, err := ParseRoleSpeech(roleFile)
	if err != nil {
		t.Fatalf("役割語テンプレートの読み込み: %v", err)
	}
	exampleFile, err := os.Open(filepath.Join("..", "..", "..", "assets", "role-speech-examples.tsv"))
	if err != nil {
		t.Fatalf("口調例文テンプレートを開けない: %v", err)
	}
	defer exampleFile.Close() //nolint:errcheck // 読み取り後の後始末。
	tbl, err = ParseRoleSpeechExamples(tbl, exampleFile)
	if err != nil {
		t.Fatalf("口調例文テンプレートの読み込み: %v", err)
	}
	return tbl
}

// allCells は基底口調の 9 セル名。正本は internal/core/tone/classifier.go:72 の defaultCellNames で、
// 並びも同じ（感情段階 抑制→中→激情 の各行に、対人段階 尊大→中立→丁寧 を並べる）。
// rolespeech は tone に依存しない純粋ルールのため、テストでも import せず literal で持つ。
// 正本のセル名を変えたら、この配列と assets/role-speech-examples.tsv の cell 列も揃える。
var allCells = []string{
	"冷然・見下し", "淡々・実務", "慇懃・端正", // 抑制 × 尊大/中立/丁寧
	"ぞんざい", "平明", "物腰やわ", // 中 × 尊大/中立/丁寧
	"居丈高・罵倒", "率直・興奮", "情に厚い懇願", // 激情 × 尊大/中立/丁寧
}

// 名指し話者の全区分（役割区分 3 × 性別 2 × セル 9 の 54 通り）が、役割語と例文の両方で一致を返すこと。
// 名指し話者の性別は NPC の Female flag から決まり空にならないため、この 54 通りで全て覆う。
func TestRealTableCoversAllNamedSpeakerKeys(t *testing.T) {
	tbl := realTable(t)
	for _, race := range []string{"child", "adult", "elder"} {
		for _, sex := range []string{"male", "female"} {
			for _, cell := range allCells {
				tmpl, ok := tbl.Lookup(race, sex, cell)
				if !ok {
					t.Errorf("%s/%s/%s: 一致する行が無い", race, sex, cell)
					continue
				}
				if tmpl.FirstPerson == "" && tmpl.Register == "" {
					t.Errorf("%s/%s/%s: 一人称も言い回しも空（役割語が付かない）", race, sex, cell)
				}
				if tmpl.Example.IsZero() {
					t.Errorf("%s/%s/%s: 例文が無い", race, sex, cell)
				}
			}
		}
	}
}

// 性別を取れない話者（汎用台詞、PC 性別が未設定の PC 発話）の経路。
// personatone.freeRoleSpeechLines が Lookup("adult", sex, "") を呼ぶため、セル列のワイルドカード行で受ける。
func TestRealTableCoversFreeLinePath(t *testing.T) {
	tbl := realTable(t)
	for _, sex := range []string{"male", "female", ""} {
		tmpl, ok := tbl.Lookup("adult", sex, "")
		if !ok {
			t.Errorf("汎用経路 adult/%q/（セルなし）: 一致する行が無い", sex)
			continue
		}
		if tmpl.FirstPerson == "" {
			t.Errorf("汎用経路 adult/%q: 一人称が空", sex)
		}
		if tmpl.Example.IsZero() {
			t.Errorf("汎用経路 adult/%q: 例文が無い", sex)
		}
	}
}

// 成人男性の一人称は、対人段階が尊大の 3 セルと、感情段階が激情で対人段階が中立の 1 セルだけ「俺」。
// 残る 5 セルは「私」。
//
// 例文訳は一人称を書かないことがある。日本語は場面から分かる主語を落とすので、例文の側も
// 落とした形にしてある（2026-07-28）。よって「一人称を含むこと」は要求せず、含む場合に
// 役割語と食い違わないことだけを見る。「私」のセルへ「俺」が現れる取り違えはここで落ちる。
func TestRealTableAdultMaleFirstPerson(t *testing.T) {
	tbl := realTable(t)
	// 対人段階が尊大の 3 セル（抑制・中・激情）と、感情段階が激情で対人段階が中立の 1 セル。
	assertive := map[string]bool{
		"冷然・見下し": true,
		"ぞんざい":   true,
		"居丈高・罵倒": true,
		"率直・興奮":  true,
	}
	for _, cell := range allCells {
		want := "私"
		if assertive[cell] {
			want = "俺"
		}
		tmpl, ok := tbl.Lookup("adult", "male", cell)
		if !ok {
			t.Fatalf("adult/male/%s: 一致する行が無い", cell)
		}
		if tmpl.FirstPerson != want {
			t.Errorf("adult/male/%s の一人称 = %q, want %q", cell, tmpl.FirstPerson, want)
		}
		other := "俺"
		if want == "俺" {
			other = "私"
		}
		if strings.Contains(tmpl.Example.Dest, other) {
			t.Errorf("adult/male/%s の例文訳が役割語と違う一人称 %q を含む: %q", cell, other, tmpl.Example.Dest)
		}
	}
}
