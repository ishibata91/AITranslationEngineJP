package personatone

import (
	"strings"
	"testing"

	"aitranslationenginejp/internal/core/rolespeech"
	"aitranslationenginejp/internal/core/tone"
	"aitranslationenginejp/internal/model"
)

// toneTraitOf は 9 セル全てで性質文を持ち、範囲外は空を返す。
func TestToneTraitOf(t *testing.T) {
	for emo := range 3 {
		for att := range 3 {
			if toneTraitOf(att, emo) == "" {
				t.Errorf("基底口調セル (対人%d, 感情%d) の性質文が空", att, emo)
			}
		}
	}
	if toneTraitOf(tone.AttitudeNeutral, 3) != "" {
		t.Error("感情段階が範囲外で性質文が空でない")
	}
	if toneTraitOf(3, tone.EmotionMid) != "" {
		t.Error("対人段階が範囲外で性質文が空でない")
	}
}

// raceMarkerTrait は種族訛りを持つ種族にだけ注記を返す。
func TestRaceMarkerTrait(t *testing.T) {
	if m := raceMarkerTrait("KhajiitRace"); !strings.Contains(m, "三人称") {
		t.Errorf("Khajiit の種族訛り = %q, 三人称の注記を期待", m)
	}
	if raceMarkerTrait("ArgonianRace") == "" {
		t.Error("Argonian の種族訛りが空")
	}
	if raceMarkerTrait("NordRace") != "" {
		t.Error("種族訛りを持たない Nord で注記が空でない")
	}
}

// buildToneTraits は性質文に種族訛りを重ね、範囲外の段階では空を返す（役割語テンプレート無し）。
func TestBuildToneTraits(t *testing.T) {
	// 丁寧×中（物腰やわ）＋ Khajiit → 口調行 ＋ 種族訛り行。役割語テンプレート nil。
	traits := BuildToneTraits(model.LinePersonaInput{
		AttitudeBand: tone.AttitudePolite, EmotionBand: tone.EmotionMid, RaceEDID: "KhajiitRace",
	}, nil)
	if len(traits) != 2 {
		t.Fatalf("Khajiit の口調指示行 = %d, want 2", len(traits))
	}
	if !strings.HasPrefix(traits[0], "- 口調: ") || !strings.HasPrefix(traits[1], "- 種族訛り: ") {
		t.Errorf("行の組み立てが想定外: %v", traits)
	}

	// 種族訛りを持たない種族は口調行のみ。
	plain := BuildToneTraits(model.LinePersonaInput{
		AttitudeBand: tone.AttitudeArrogant, EmotionBand: tone.EmotionSuppressed, RaceEDID: "NordRace",
	}, nil)
	if len(plain) != 1 {
		t.Fatalf("Nord の口調指示行 = %d, want 1", len(plain))
	}

	// 段階が範囲外なら口調指示なし（空）。
	if got := BuildToneTraits(model.LinePersonaInput{AttitudeBand: 9, EmotionBand: 9}, nil); got != nil {
		t.Errorf("範囲外段階で traits = %v, want nil", got)
	}
}

// buildToneTraits は役割語テンプレートがあれば、性質文と種族訛りの間へ一人称・言い回しの行を挟む。
func TestBuildToneTraitsWithRoleSpeech(t *testing.T) {
	roles, err := rolespeech.ParseRoleSpeech(strings.NewReader(
		"elder\tfemale\t*\tわたし\t年配の女性らしく落ち着いて。\n"))
	if err != nil {
		t.Fatalf("ParseRoleSpeech: %v", err)
	}
	// 老女（ElderRace + Female）・物腰やわ（丁寧×中）。口調 ＋ 性別 ＋ 人称と言い回し の 3 行（種族訛り無し）。
	traits := BuildToneTraits(model.LinePersonaInput{
		AttitudeBand: tone.AttitudePolite, EmotionBand: tone.EmotionMid,
		Sex: "Female", RaceEDID: "ElderRace",
	}, roles)
	if len(traits) != 3 {
		t.Fatalf("老女の口調指示行 = %d, want 3: %v", len(traits), traits)
	}
	if !strings.HasPrefix(traits[2], "- 人称と言い回し: ") || !strings.Contains(traits[2], "わたし") {
		t.Errorf("役割語の行が想定外: %q", traits[2])
	}

	// テンプレートに当たらない成人男（adult/male、行なし）は役割語を付けず、口調 ＋ 性別 の 2 行。
	none := BuildToneTraits(model.LinePersonaInput{
		AttitudeBand: tone.AttitudeNeutral, EmotionBand: tone.EmotionMid,
		Sex: "Male", RaceEDID: "NordRace",
	}, roles)
	if len(none) != 2 {
		t.Fatalf("成人男の口調指示行 = %d, want 2: %v", len(none), none)
	}
}

// buildFreeToneTraits は自由記述の口調へ、感情段階の助言と性別の一人称・語尾を重ねる（汎用・PC 用）。
func TestBuildGenericToneTraits(t *testing.T) {
	roles, err := rolespeech.ParseRoleSpeech(strings.NewReader(
		"adult\tfemale\t*\tわたし\t女性らしいやわらかな言い回し。\n"))
	if err != nil {
		t.Fatalf("ParseRoleSpeech: %v", err)
	}

	// 自由記述が空なら口調指示なし（空）。
	if got := BuildGenericToneTraits("", tone.EmotionMid, "Female", roles, ""); got != nil {
		t.Errorf("自由記述が空で traits = %v, want nil", got)
	}

	// 感情が中（1）は助言を出さない。女性は性別行と一人称・語尾が付く。口調 ＋ 性別 ＋ 人称 の 3 行。
	mid := BuildGenericToneTraits("汎用台詞。", tone.EmotionMid, "Female", roles, "")
	if len(mid) != 3 {
		t.Fatalf("中・女性の口調指示行 = %d, want 3: %v", len(mid), mid)
	}
	if !strings.HasPrefix(mid[0], "- 口調: ") || !strings.Contains(mid[2], "わたし") {
		t.Errorf("行の組み立てが想定外: %v", mid)
	}

	// 感情が激情（2）は助言を出す。性別なしは一人称・語尾を付けない。口調 ＋ 感情 の 2 行。
	intense := BuildGenericToneTraits("汎用台詞。", tone.EmotionIntense, "", roles, "")
	if len(intense) != 2 {
		t.Fatalf("激情・性別なしの口調指示行 = %d, want 2: %v", len(intense), intense)
	}
	if !strings.HasPrefix(intense[1], "- 感情: ") {
		t.Errorf("感情助言の行が想定外: %v", intense)
	}

	// 抑制（0）＋ 男性 ＋ roles nil は口調 ＋ 性別 ＋ 感情 の 3 行（一人称・語尾なし）。
	suppressed := BuildPCToneTraits("PC の選択肢。", tone.EmotionSuppressed, "Male", nil, "")
	if len(suppressed) != 3 {
		t.Fatalf("抑制の口調指示行 = %d, want 3: %v", len(suppressed), suppressed)
	}

	// テンプレートに当たらない成人男（adult/male、行なし）は一人称・語尾を付けず、中なら口調 ＋ 性別 の 2 行。
	maleMid := BuildGenericToneTraits("汎用台詞。", tone.EmotionMid, "Male", roles, "")
	if len(maleMid) != 2 {
		t.Fatalf("中・成人男の口調指示行 = %d, want 2: %v", len(maleMid), maleMid)
	}
}

// lineEmotionLine は TRDT 感情型を助言行へ写す。Neutral・空・未知型は空（加算なし）。
func TestLineEmotionLine(t *testing.T) {
	cases := map[string]string{
		tone.LineEmotionAnger:    "怒り",
		tone.LineEmotionDisgust:  "嫌悪",
		tone.LineEmotionFear:     "恐れ",
		tone.LineEmotionSad:      "悲しみ",
		tone.LineEmotionHappy:    "喜び",
		tone.LineEmotionSurprise: "驚き",
		tone.LineEmotionPuzzled:  "戸惑い",
	}
	for typ, word := range cases {
		got := lineEmotionLine(typ)
		if !strings.HasPrefix(got, "- 感情: ") || !strings.Contains(got, word) {
			t.Errorf("感情型 %q の行 = %q, 感情語 %q を期待", typ, got, word)
		}
	}
	for _, none := range []string{tone.LineEmotionNeutral, "", "Unknown"} {
		if got := lineEmotionLine(none); got != "" {
			t.Errorf("加算しない型 %q で行 = %q, want 空", none, got)
		}
	}
}

// BuildToneTraits は非 Neutral の台詞感情を口調指示の末尾へ 1 行足し、Neutral では足さない（ペルソナ基底は不変）。
func TestBuildToneTraitsLineEmotion(t *testing.T) {
	base := model.LinePersonaInput{AttitudeBand: tone.AttitudeNeutral, EmotionBand: tone.EmotionMid}
	plain := BuildToneTraits(base, nil)

	angry := base
	angry.EmotionType = tone.LineEmotionAnger
	withEmo := BuildToneTraits(angry, nil)
	if len(withEmo) != len(plain)+1 {
		t.Fatalf("感情行の追加が想定外: plain=%v withEmo=%v", plain, withEmo)
	}
	last := withEmo[len(withEmo)-1]
	if !strings.HasPrefix(last, "- 感情: ") || !strings.Contains(last, "怒り") {
		t.Errorf("感情行が想定外: %q", last)
	}

	neutral := base
	neutral.EmotionType = tone.LineEmotionNeutral
	if got := BuildToneTraits(neutral, nil); len(got) != len(plain) {
		t.Errorf("Neutral で感情行が増えた: %v", got)
	}
}

// BuildGenericToneTraits は TRDT 種別があれば台詞感情行を出し、本文推定の強度助言を出さない（二重回避）。
// TRDT が無ければ従来の強度助言へ落ちる。
func TestBuildGenericToneTraitsLineEmotion(t *testing.T) {
	// 本文は激情だが TRDT=悲しみ → TRDT 種別を出し、強度助言（「高ぶった」）は出さない。
	got := BuildGenericToneTraits("汎用台詞。", tone.EmotionIntense, "", nil, tone.LineEmotionSad)
	if len(got) != 2 {
		t.Fatalf("口調 ＋ 感情の 2 行を期待: %v", got)
	}
	if !strings.Contains(got[1], "悲しみ") || strings.Contains(got[1], "高ぶった") {
		t.Errorf("TRDT 種別でなく強度助言が出た: %q", got[1])
	}

	// TRDT 無し（空）＋ 激情 → 従来の強度助言へ落ちる。
	fallback := BuildGenericToneTraits("汎用台詞。", tone.EmotionIntense, "", nil, "")
	if len(fallback) != 2 || !strings.Contains(fallback[1], "高ぶった") {
		t.Errorf("TRDT 無しで強度助言に落ちない: %v", fallback)
	}
}

// freeRoleSpeechLines（formatRoleSpeech）は一人称のみ・言い回しのみ・両方空を出し分ける。
func TestFreeRoleSpeechFormatting(t *testing.T) {
	firstOnly, _ := rolespeech.ParseRoleSpeech(strings.NewReader("adult\tmale\t*\t俺\t\n"))
	if got := freeRoleSpeechLines("Male", firstOnly, false); len(got) != 1 || got[0] != "- 人称と言い回し: 一人称は「俺」。" {
		t.Errorf("一人称のみ = %v", got)
	}
	registerOnly, _ := rolespeech.ParseRoleSpeech(strings.NewReader("adult\tmale\t*\t\t男性らしい言い回し。\n"))
	if got := freeRoleSpeechLines("Male", registerOnly, false); len(got) != 1 || got[0] != "- 言い回し: 男性らしい言い回し。" {
		t.Errorf("言い回しのみ = %v", got)
	}
	bothEmpty, _ := rolespeech.ParseRoleSpeech(strings.NewReader("adult\tmale\t*\t\t\n"))
	if got := freeRoleSpeechLines("Male", bothEmpty, false); len(got) != 0 {
		t.Errorf("一人称も言い回しも空で = %v, want 行なし", got)
	}
}

// 性別を取れない話者（汎用台詞、PC 性別が未設定の PC 発話）でも役割語を引き、
// 性別列のワイルドカード行へ落ちること。早期 return で打ち切らない。
func TestFreeRoleSpeechUnknownSexFallsBackToWildcard(t *testing.T) {
	roles, err := rolespeech.ParseRoleSpeech(strings.NewReader(
		"adult\tmale\t*\t俺\t\nadult\t*\t*\t私\t\n"))
	if err != nil {
		t.Fatalf("役割語表の解析: %v", err)
	}
	got := freeRoleSpeechLines("", roles, true)
	if len(got) != 1 || !strings.Contains(got[0], "「私」") {
		t.Fatalf("性別不明でワイルドカード行へ落ちない: %v", got)
	}
	// 口調指示の組み立てまで通ること（PC 発話・汎用台詞の経路）。
	traits := BuildGenericToneTraits("汎用台詞。", tone.EmotionMid, "", roles, "")
	if len(traits) != 2 || !strings.Contains(traits[1], "「私」") {
		t.Fatalf("性別不明の口調指示に一人称が乗らない: %v", traits)
	}
}

// 名指し話者の口調指示にも、人称の行に続けて例文の行が乗ること（セル別に引く経路）。
func TestBuildToneTraitsIncludesExample(t *testing.T) {
	roles, err := rolespeech.ParseRoleSpeech(strings.NewReader("adult\tmale\t*\t俺\t\n"))
	if err != nil {
		t.Fatalf("役割語表の解析: %v", err)
	}
	cell := tone.CellName(tone.AttitudeArrogant, tone.EmotionMid)
	roles, err = rolespeech.ParseRoleSpeechExamples(roles,
		strings.NewReader("adult\tdefault\tmale\t"+cell+"\tMove.\tどけ。俺の邪魔だ。\n"))
	if err != nil {
		t.Fatalf("例文表の解析: %v", err)
	}
	got := BuildToneTraits(model.LinePersonaInput{
		AttitudeBand: tone.AttitudeArrogant,
		EmotionBand:  tone.EmotionMid,
		Sex:          "Male",
		RaceEDID:     "NordRace",
	}, roles)
	if len(got) != 4 {
		t.Fatalf("口調・性別・人称・例文の 4 行を期待: %v", got)
	}
	if got[3] != "- 例: Move. → どけ。俺の邪魔だ。" {
		t.Errorf("例文行 = %q", got[3])
	}
}

// 例文を持つテンプレートは、人称の行に続けて「- 例:」の行を口調指示へ足す。
func TestRoleSpeechExampleLine(t *testing.T) {
	roles, err := rolespeech.ParseRoleSpeech(strings.NewReader("adult\tmale\t*\t俺\t\n"))
	if err != nil {
		t.Fatalf("役割語表の解析: %v", err)
	}
	roles, err = rolespeech.ParseRoleSpeechExamples(roles,
		strings.NewReader("adult\tdefault\tmale\t*\tI got it!\tやったぞ、俺がやった。\n"))
	if err != nil {
		t.Fatalf("例文表の解析: %v", err)
	}
	got := freeRoleSpeechLines("Male", roles, true)
	if len(got) != 2 {
		t.Fatalf("人称と例文の 2 行を期待: %v", got)
	}
	if got[1] != "- 例: I got it! → やったぞ、俺がやった。" {
		t.Errorf("例文行 = %q", got[1])
	}
}

// R-1-1・R-1-2・R-1-4: 名指しKhajiitには組み合わせ別3例と使い方の指示が入力順で入る。
func TestBuildToneTraitsIncludesKhajiitExamplesAndUsageInstruction(t *testing.T) {
	roles, err := rolespeech.ParseRoleSpeech(strings.NewReader("adult\tfemale\t*\t\t落ち着いて。\n"))
	if err != nil {
		t.Fatalf("役割語表の解析: %v", err)
	}
	cell := tone.CellName(tone.AttitudeNeutral, tone.EmotionMid)
	data := strings.Join([]string{
		"adult\tkhajiit\tfemale\t" + cell + "\tF1\tこの者はF1",
		"adult\tkhajiit\tfemale\t" + cell + "\tF2\tF2には種族表現なし",
		"adult\tkhajiit\tfemale\t" + cell + "\tF3\tこの者にF3",
	}, "\n") + "\n"
	roles, err = rolespeech.ParseRoleSpeechExamples(roles, strings.NewReader(data))
	if err != nil {
		t.Fatalf("例文表の解析: %v", err)
	}
	got := BuildToneTraits(model.LinePersonaInput{
		AttitudeBand: tone.AttitudeNeutral, EmotionBand: tone.EmotionMid,
		Sex: "Female", RaceEDID: "KhajiitRace",
	}, roles)
	joined := strings.Join(got, "\n")
	if strings.Index(joined, "F1 →") >= strings.Index(joined, "F2 →") || strings.Index(joined, "F2 →") >= strings.Index(joined, "F3 →") {
		t.Errorf("3例が入力順でない: %v", got)
	}
	exampleText := strings.Join(got[3:6], "\n")
	if strings.Count(exampleText, "この者") != 2 || strings.Contains(got[4], "この者") {
		t.Errorf("Khajiit例の種族表現が想定外: %v", got)
	}
	if got[6] != exampleUsageInstruction {
		t.Errorf("3例直後の使い方指示 = %q", got[6])
	}
}

// R-2-1〜R-2-3: 汎用台詞だけが性別別3例を持ち、PC発話と性別不明の汎用台詞は例を持たない。
func TestGenericAndPCToneTraitsSeparateExamples(t *testing.T) {
	roles, err := rolespeech.ParseRoleSpeech(strings.NewReader("adult\t*\t*\t\t自然に。\n"))
	if err != nil {
		t.Fatalf("役割語表の解析: %v", err)
	}
	var rows []string
	for _, sex := range []string{"male", "female"} {
		for i := 1; i <= 3; i++ {
			rows = append(rows, "adult\tdefault\t"+sex+"\t*\tF"+string(rune('0'+i))+"\t"+sex+"-F"+string(rune('0'+i)))
		}
	}
	roles, err = rolespeech.ParseRoleSpeechExamples(roles, strings.NewReader(strings.Join(rows, "\n")+"\n"))
	if err != nil {
		t.Fatalf("例文表の解析: %v", err)
	}
	male := strings.Join(BuildGenericToneTraits("汎用台詞。", tone.EmotionMid, "Male", roles, ""), "\n")
	female := strings.Join(BuildGenericToneTraits("汎用台詞。", tone.EmotionMid, "Female", roles, ""), "\n")
	if strings.Count(male, "- 例:") != 3 || !strings.Contains(male, "male-F1") || strings.Contains(male, "female-F1") {
		t.Errorf("男性の汎用例が想定外: %q", male)
	}
	if strings.Count(female, "- 例:") != 3 || !strings.Contains(female, "→ female-F1") || strings.Contains(female, "→ male-F1") {
		t.Errorf("女性の汎用例が想定外: %q", female)
	}
	unknown := strings.Join(BuildGenericToneTraits("汎用台詞。", tone.EmotionMid, "", roles, ""), "\n")
	pc := strings.Join(BuildPCToneTraits("PC発話。", tone.EmotionMid, "Female", roles, tone.LineEmotionHappy), "\n")
	if strings.Contains(unknown, "- 例:") {
		t.Errorf("性別不明の汎用台詞に例がある: %q", unknown)
	}
	if strings.Contains(pc, "- 例:") || !strings.Contains(pc, "- 性別: 女性") || !strings.Contains(pc, "喜び") || !strings.Contains(pc, "自然に") {
		t.Errorf("PC発話の属性維持または例除外が想定外: %q", pc)
	}
}

// buildToneDirective は traits があればテンプレートの {traits} を置換し、無ければ空を返す。
func TestBuildToneDirective(t *testing.T) {
	tmpl := "口調指示:\n{traits}"
	got := BuildToneDirective(tmpl, []string{"- 口調: 柔らかく丁寧"})
	if !strings.Contains(got, "- 口調: 柔らかく丁寧") || strings.Contains(got, "{traits}") {
		t.Errorf("{traits} の置換が想定外: %q", got)
	}
	if BuildToneDirective(tmpl, nil) != "" {
		t.Error("traits が空で空文字を返さない")
	}
}

// buildToneLabel は基底口調セル名のチップ要約を返し、範囲外は空。
func TestBuildToneLabel(t *testing.T) {
	got := BuildToneLabel(model.LinePersonaInput{AttitudeBand: tone.AttitudePolite, EmotionBand: tone.EmotionMid})
	if got != "口調: 物腰やわ" {
		t.Errorf("buildToneLabel = %q, want %q", got, "口調: 物腰やわ")
	}
	if BuildToneLabel(model.LinePersonaInput{AttitudeBand: 9, EmotionBand: 9}) != "" {
		t.Error("範囲外段階で空文字を返さない")
	}
}

// R-4-1: 性別が取れる名指し話者の台詞について、口調指示に性別を示す行が出ること。
// 行は一人称・言い回しの行とは別に立てる。役割語は成人男性と性別不明が同じ出力になるため、
// 引いた結果だけでは男女の区別がプロンプトに現れない。
//
// 性別の行は 2026-07-28 に「- 性別: 男性」だけへ短縮した。変える前は「男性の話者として訳す。」を
// 続けていたが、実験 task dialogue-tone-naturalness が指示文の長さと破損行の件数の関係を測り、
// 口調指示が長いほど JSON が壊れる行が増えることを確かめた。事実だけを置く形へ縮めた。
func TestBuildToneTraitsHasSexLine(t *testing.T) {
	roles, err := rolespeech.ParseRoleSpeech(strings.NewReader("adult\t*\t*\t私\t\n"))
	if err != nil {
		t.Fatalf("役割語表の解析: %v", err)
	}
	for _, tc := range []struct{ sex, want string }{
		{"Male", "- 性別: 男性"},
		{"Female", "- 性別: 女性"},
	} {
		traits := BuildToneTraits(model.LinePersonaInput{
			AttitudeBand: tone.AttitudeNeutral, EmotionBand: tone.EmotionMid,
			Sex: tc.sex, RaceEDID: "NordRace",
		}, roles)
		if !hasLine(traits, tc.want) {
			t.Errorf("性別 %q の口調指示に性別の行が無い: %v", tc.sex, traits)
		}
		// 性別の行と一人称・言い回しの行が、それぞれ独立した行として並ぶこと（1 行へ畳まない）。
		if !hasPrefix(traits, "- 人称と言い回し: ") {
			t.Errorf("性別 %q の口調指示に一人称の行が無い: %v", tc.sex, traits)
		}
	}
}

// R-4-2: 話者を解決できない汎用台詞とプレイヤーの選択肢についても、性別が取れる時は
// 名指し話者と同じ形の行が出ること。
func TestBuildGenericToneTraitsHasSameSexLine(t *testing.T) {
	roles, err := rolespeech.ParseRoleSpeech(strings.NewReader("adult\t*\t*\t私\t\n"))
	if err != nil {
		t.Fatalf("役割語表の解析: %v", err)
	}
	for _, sex := range []string{"Male", "Female"} {
		named := BuildToneTraits(model.LinePersonaInput{
			AttitudeBand: tone.AttitudeNeutral, EmotionBand: tone.EmotionMid,
			Sex: sex, RaceEDID: "NordRace",
		}, roles)
		free := BuildGenericToneTraits("汎用台詞。", tone.EmotionMid, sex, roles, "")
		want := sexTrait(sex)
		if !hasLine(named, want) || !hasLine(free, want) {
			t.Errorf("性別 %q で名指し話者と汎用台詞の性別の行が揃わない:\n  名指し=%v\n  汎用=%v", sex, named, free)
		}
	}
}

// R-4-3: 性別を取れない話者の台詞に、性別を示す行が出ないこと。
// プレイヤーの性別が未設定の選択肢（空文字）も、性別を取れない話者と同じ扱いにする。
func TestSexLineAbsentWhenSexUnknown(t *testing.T) {
	roles, err := rolespeech.ParseRoleSpeech(strings.NewReader("adult\t*\t*\t私\t\n"))
	if err != nil {
		t.Fatalf("役割語表の解析: %v", err)
	}
	for _, sex := range []string{"", "  ", "Unknown"} {
		named := BuildToneTraits(model.LinePersonaInput{
			AttitudeBand: tone.AttitudeNeutral, EmotionBand: tone.EmotionMid,
			Sex: sex, RaceEDID: "NordRace",
		}, roles)
		free := BuildGenericToneTraits("汎用台詞。", tone.EmotionMid, sex, roles, "")
		if hasPrefix(named, "- 性別: ") {
			t.Errorf("性別 %q（取れない）の名指し話者に性別の行が出た: %v", sex, named)
		}
		if hasPrefix(free, "- 性別: ") {
			t.Errorf("性別 %q（取れない）の汎用台詞に性別の行が出た: %v", sex, free)
		}
	}
}

// hasLine は口調指示の行に want と完全一致する行があるかを返す。
func hasLine(traits []string, want string) bool {
	for _, l := range traits {
		if l == want {
			return true
		}
	}
	return false
}

// hasPrefix は口調指示の行に prefix で始まる行があるかを返す。
func hasPrefix(traits []string, prefix string) bool {
	for _, l := range traits {
		if strings.HasPrefix(l, prefix) {
			return true
		}
	}
	return false
}
