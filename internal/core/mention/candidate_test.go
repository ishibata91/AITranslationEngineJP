package mention

import (
	"reflect"
	"slices"
	"strings"
	"testing"
)

// fakeAnalyzer は TextAnalyzer のテスト用偽実装。外部学習モデル（prose）を切り離し、
// 固有表現と動詞判定の応答を固定する。
type fakeAnalyzer struct {
	entities []string
	verbs    map[string]bool
}

func (f fakeAnalyzer) Entities(string) []string { return f.entities }
func (f fakeAnalyzer) LeadingVerb(sentence string) bool {
	return f.verbs[strings.TrimSpace(sentence)]
}

// fakeStoplist は Stoplist 境界のテスト用偽実装。1 語だけの原語を小文字化して照合する
// （internal/core/dictionary の Stoplist と同じ規則の最小再現）。
type fakeStoplist map[string]bool

func (f fakeStoplist) Blocks(source string) bool {
	trimmed := strings.TrimSpace(source)
	if trimmed == "" || len(strings.Fields(trimmed)) != 1 {
		return false
	}
	return f[strings.ToLower(trimmed)]
}

// testStoplist は一般語 stoplist の最小セット（stopwords-iso の内容変化から切り離す）。
func testStoplist(t *testing.T) Stoplist {
	t.Helper()
	return fakeStoplist{
		"a": true, "an": true, "and": true, "the": true, "of": true,
		"i": true, "take": true, "well": true, "yes": true, "no": true,
	}
}

// detect はテスト用の短縮呼び出し。既知語・本文・解析器から候補語の表記だけを取り出す。
func detect(t *testing.T, known []string, texts []string, ner TextAnalyzer) []string {
	t.Helper()
	cand := NewCandidateDetector(known, testStoplist(t), ner).DetectCandidates(texts)
	terms := make([]string, len(cand))
	for i, c := range cand {
		terms[i] = c.Term
	}
	return terms
}

// has は候補列に語が含まれるかを返す。
func has(terms []string, want string) bool {
	return slices.Contains(terms, want)
}

// 正規化規則を突く。前後空白・語間空白の畳み込み・所有格（'s・s'・typographic）の除去・
// 大小区別の保持・語中アポストロフィの保持。
func TestNormalizeCandidate(t *testing.T) {
	cases := map[string]string{
		"  Whiterun  ":       "Whiterun",
		"College  of   Wint": "College of Wint",
		"Serana's":           "Serana",
		"Serana’s":           "Serana",
		"Companions'":        "Companions",
		"Mehrunes' Razor":    "Mehrunes' Razor",
		"WHITERUN":           "WHITERUN",
		"Frost-":             "Frost",
	}
	for in, want := range cases {
		if got := NormalizeCandidate(in); got != want {
			t.Errorf("NormalizeCandidate(%q) = %q, want %q", in, got, want)
		}
	}
}

// 文中の大文字始まりの未知語が候補になり、既知語・stoplist 語・1 文字・文頭の一般語は
// 候補にならないこと（基本の受理と除外）。
func TestDetectCandidatesBasic(t *testing.T) {
	texts := []string{
		"I saw Serana near the gate.",   // Serana: 文中の未知語 → 候補
		"We visited Whiterun that day.", // Whiterun: 既知語 → 除外
		"Take the sword now.",           // Take: 文頭の stoplist 語 → 除外
		"Speak to him, he said sadly.",  // Speak: 文頭だけ・小文字用法なし… 後続ケースで確認
		"A dog barked.",                 // A: 1 文字 → 除外
	}
	terms := detect(t, []string{"Whiterun"}, texts, nil)
	if !has(terms, "Serana") {
		t.Errorf("文中の未知語 Serana が候補にない: %v", terms)
	}
	for _, banned := range []string{"Whiterun", "Take", "A", "I"} {
		if has(terms, banned) {
			t.Errorf("%q は候補にしない: %v", banned, terms)
		}
	}
}

// 出力が正規化後表記で一意・辞書順で、同一入力に同一出力（決定性）であること。
func TestDetectCandidatesDeterministicSorted(t *testing.T) {
	texts := []string{
		"They fear Serana and Serana's father.", // 同一語の重複と所有格の畳み込み
		"Beware of Harkon. Harkon rules.",
	}
	d := NewCandidateDetector(nil, testStoplist(t), nil)
	first := d.DetectCandidates(texts)
	second := d.DetectCandidates(texts)
	if !reflect.DeepEqual(first, second) {
		t.Errorf("同一入力で出力が変動した: %v / %v", first, second)
	}
	seen := map[string]bool{}
	for i, c := range first {
		if seen[c.Term] {
			t.Errorf("重複出力: %q", c.Term)
		}
		seen[c.Term] = true
		if i > 0 && first[i-1].Term >= c.Term {
			t.Errorf("辞書順でない: %q の後に %q", first[i-1].Term, c.Term)
		}
	}
	if !seen["Serana"] || !seen["Harkon"] {
		t.Errorf("Serana・Harkon が候補にない: %v", first)
	}
}

// 複数語の句の組み立てを突く。接続語（of・the・&）の取り込み、2 個までの接続語、
// カンマ・2 連空白での句切り、小文字接頭の姓（gro-）、一人称代名詞の除外。
func TestPhraseComposition(t *testing.T) {
	texts := []string{
		"Go to the College of Winterhold today.",     // of 接続
		"They sang of the Bane of the Dead quietly.", // of the の 2 連接続
		"He met Balagog gro-Nolob at noon.",          // 小文字接頭の姓
		"We read Kolb & the Dragon aloud.",           // & 接続
		"I saw Whiterun, Riften and Solitude.",       // カンマで切れる（Whiterun Riften と繋がない）
		"I'm Farkas of the Companions.",              // I'm は句の部材にしない
	}
	terms := detect(t, nil, texts, nil)
	for _, want := range []string{"College of Winterhold", "Bane of the Dead", "Balagog gro-Nolob", "Kolb & the Dragon", "Farkas of the Companions"} {
		if !has(terms, want) {
			t.Errorf("%q が候補にない: %v", want, terms)
		}
	}
	for _, banned := range []string{"Whiterun Riften", "I'm Farkas"} {
		if has(terms, banned) {
			t.Errorf("%q は候補にしない: %v", banned, terms)
		}
	}
	// 2 連空白は整形一覧の項目境界。句を繋がない。
	listTerms := detect(t, nil, []string{"items here: Fire Salts  Bear Pelt"}, nil)
	if has(listTerms, "Fire Salts Bear Pelt") {
		t.Errorf("2 連空白を跨いで句を繋いだ: %v", listTerms)
	}
	if !has(listTerms, "Fire Salts") || !has(listTerms, "Bear Pelt") {
		t.Errorf("一覧の各項目が候補にない: %v", listTerms)
	}
}

// 既知語の彫り出しを突く。句の中の既知語（称号 Jarl）は単独の候補にならず、彫り出し残り
// （Whiterun）が候補になる。全体が未知の句（Jarl of Whiterun・Pillar of Sacrifice）は、
// 辞書が全体と部分の両方を語彙に持つ形に合わせて候補に残る。
func TestKnownCarving(t *testing.T) {
	texts := []string{
		"Ask the Jarl of Whiterun about it. We ask and touch things.",
		"Touch the Pillar of Sacrifice now.",
	}
	terms := detect(t, []string{"Jarl", "Pillar", "Sacrifice"}, texts, nil)
	if !has(terms, "Whiterun") {
		t.Errorf("彫り出し残りの Whiterun が候補にない: %v", terms)
	}
	if has(terms, "Jarl") {
		t.Errorf("既知語 Jarl が単独の候補に残った: %v", terms)
	}
	if !has(terms, "Pillar of Sacrifice") {
		t.Errorf("成分が既知でも全体が未知の句が候補にない: %v", terms)
	}
}

// 文頭（曖昧位置）の扱いを突く。文頭語を含む解釈と外した解釈の両建て（Reduces Health）、
// 閉クラスの機能語先頭の抑止（A Bosmer）、The の例外（The Reach）。
func TestAmbiguousLead(t *testing.T) {
	terms := detect(t, nil, []string{
		"Reduces Health by ten points. It reduces more later.", // reduces に小文字用法
		"A Bosmer arrived late.",
		"The Reach belongs to the Forsworn!",
	}, nil)
	for _, want := range []string{"Reduces Health", "Health", "Bosmer", "The Reach", "Reach"} {
		if !has(terms, want) {
			t.Errorf("%q が候補にない: %v", want, terms)
		}
	}
	if has(terms, "A Bosmer") {
		t.Errorf("機能語 A で始まる句全体は候補にしない: %v", terms)
	}
}

// 引用符で括られた句を突く。全体だけが候補になり（部分は出さない）、文頭でも通る。
// 本文先頭の開き引用符も検出する。
func TestQuotedPhrase(t *testing.T) {
	terms := detect(t, nil, []string{
		`Use the "Clear Skies" Shout to open the path. It must clear.`,
		`"Night of Tears", by Dranor Seleth.`,
	}, nil)
	if !has(terms, "Clear Skies") {
		t.Errorf("引用符付きの Clear Skies が候補にない: %v", terms)
	}
	if !has(terms, "Night of Tears") {
		t.Errorf("本文先頭の引用符付き Night of Tears が候補にない: %v", terms)
	}
}

// 命令形のクエスト目標行を突く。動詞判定（LeadingVerb）が真の 1 行文は、句全体を捨てて
// 先頭抜きだけを候補にする。動詞判定が偽の 1 行文（item 名）は全体を残す。
func TestVerbLeadObjective(t *testing.T) {
	ner := fakeAnalyzer{verbs: map[string]bool{"Kill Vittoria Vici": true}}
	terms := detect(t, nil, []string{
		"Kill Vittoria Vici",
		"Mace Etiquette",
		"People kill for less. A mace helps.", // kill・mace の小文字用法を作る
	}, ner)
	if has(terms, "Kill Vittoria Vici") {
		t.Errorf("動詞先頭の目標行全体は候補にしない: %v", terms)
	}
	if !has(terms, "Vittoria Vici") {
		t.Errorf("目標行の先頭抜き Vittoria Vici が候補にない: %v", terms)
	}
	if !has(terms, "Mace Etiquette") {
		t.Errorf("動詞でない 1 行文 Mace Etiquette が候補にない: %v", terms)
	}
}

// 所有格の分割を突く。所有格の前後（Wylandriah・Ingot）と全体（Wylandriah's Ingot）が
// 候補になる（Azura's Star のような所有格込みの語彙と、所有者名の両対応）。
func TestPossessiveSplit(t *testing.T) {
	terms := detect(t, nil, []string{"Fetch Wylandriah's Ingot quickly. People fetch things."}, nil)
	for _, want := range []string{"Wylandriah's Ingot", "Wylandriah", "Ingot"} {
		if !has(terms, want) {
			t.Errorf("%q が候補にない: %v", want, terms)
		}
	}
}

// 複数語成分の内側の 1 語を突く。小文字用法の無い語（Edda）だけが 1 語候補になり、
// 小文字用法のある語（Door）は出ない。
func TestInnerSingle(t *testing.T) {
	terms := detect(t, nil, []string{
		"He recited the Poetic Edda badly.",
		"Open the Treasury Door. The door creaks.",
	}, nil)
	if !has(terms, "Edda") {
		t.Errorf("小文字用法の無い内側の語 Edda が候補にない: %v", terms)
	}
	if has(terms, "Door") {
		t.Errorf("小文字用法のある内側の語 Door は 1 語候補にしない: %v", terms)
	}
}

// 派生分割を突く。単独でも頻出する称号（Jarl）を外した複数語の派生（Brina Merilis）と、
// 種別語（Shout）を外した派生（Word）が出る。派生の残りの端に接続語（of）が残る形も落とす。
func TestDeriveSplits(t *testing.T) {
	terms := detect(t, nil, []string{
		"He met the Jarl. She thanked the Jarl.",              // Jarl の単独句出現 ×2
		"They wrote to Jarl Brina Merilis kindly.",            // 派生元（先頭が称号）
		"Word came to Jarl of Falkreath. He gave his word.",   // 派生の残りが of 始まり
		"He heard the Shout. She feared the Shout.",           // Shout の単独句出現 ×2
		"It sits by the Word of Shout stone. People sit too.", // 派生元（末尾が種別語）
	}, nil)
	if !has(terms, "Brina Merilis") {
		t.Errorf("称号 Jarl を外した派生 Brina Merilis が候補にない: %v", terms)
	}
	if !has(terms, "Falkreath") {
		t.Errorf("接続語を落とした派生 Falkreath が候補にない: %v", terms)
	}
	if !has(terms, "Word") {
		t.Errorf("種別語 Shout を外した派生 Word が候補にない: %v", terms)
	}
	if has(terms, "of Falkreath") || has(terms, "Word of") {
		t.Errorf("接続語が端に残った派生を出した: %v", terms)
	}
}

// 曖昧位置だけに出た 1 語の最終判定を突く。小文字用法の無い語（Serana）は通り、
// 小文字用法のある語（Calm）は落ち、固有表現認識が entity と判定した語（Ale）は救済する。
func TestAmbiguousSingleAcceptance(t *testing.T) {
	texts := []string{
		"Serana is my daughter.",
		"Calm down. Stay calm.",
		"Ale is cheaper than blood.", // ale の小文字用法を作る
		"Cheap ale flows here.",
	}
	noNER := detect(t, nil, texts, nil)
	if !has(noNER, "Serana") {
		t.Errorf("小文字用法の無い文頭語 Serana が候補にない: %v", noNER)
	}
	for _, banned := range []string{"Calm", "Ale"} {
		if has(noNER, banned) {
			t.Errorf("補強なしで %q は候補にしない: %v", banned, noNER)
		}
	}
	withNER := detect(t, nil, texts, fakeAnalyzer{entities: []string{"Ale"}})
	if !has(withNER, "Ale") {
		t.Errorf("固有表現の救済で Ale が候補にならない: %v", withNER)
	}
}

// 見出し行（title-case の行）を突く。見出しの大文字は文中根拠にならないが、
// 見出しの句自体（書名）は候補になる。機能語先頭の書名（Of Crossed Daggers）も全体が残る。
func TestTitleCaseHeader(t *testing.T) {
	terms := detect(t, nil, []string{
		"Of Crossed Daggers The History of Riften",
		"An Adventure for Nord Boys",
	}, nil)
	if !has(terms, "Of Crossed Daggers The History of Riften") && !has(terms, "Of Crossed Daggers") {
		t.Errorf("見出しの書名が候補にない: %v", terms)
	}
	if !has(terms, "An Adventure for Nord Boys") && !has(terms, "An Adventure") {
		t.Errorf("An 先頭の見出しの句が候補にない: %v", terms)
	}
}

// 動的タグの除去を突く。タグ内の語（Alias・Courier）は候補にならない。
func TestMarkupStripped(t *testing.T) {
	terms := detect(t, nil, []string{"Find the <Alias=Courier> passing through."}, nil)
	for _, banned := range []string{"Alias", "Courier", "Alias=Courier"} {
		if has(terms, banned) {
			t.Errorf("タグ内の %q は候補にしない: %v", banned, terms)
		}
	}
}

// 引用符・コロン等の直後（曖昧位置の再設定）を突く。コロン直後の大文字語は文頭扱いで、
// 小文字用法があれば候補にならない。
func TestResetPositions(t *testing.T) {
	terms := detect(t, nil, []string{
		"Remember: Take nothing. We take risks. Marks--Serana leads.",
		"Wait… Something moved. Any something counts.",
	}, nil)
	if has(terms, "Take") || has(terms, "Something") {
		t.Errorf("区切り直後の一般語を候補にした: %v", terms)
	}
	if !has(terms, "Serana") {
		t.Errorf("ダッシュ直後の未知語 Serana が候補にない: %v", terms)
	}
}

// 空入力・空既知語の境界を突く。本文なしは候補なし。空白だけの既知語は無視する。
func TestEmptyInputs(t *testing.T) {
	if got := detect(t, []string{"  ", ""}, nil, nil); len(got) != 0 {
		t.Errorf("本文なしで候補が出た: %v", got)
	}
	if got := detect(t, nil, []string{"only lowercase words here."}, nil); len(got) != 0 {
		t.Errorf("大文字語なしで候補が出た: %v", got)
	}
}

// 彫り出しの再帰を突く。句の中の複数の既知語（Jarl・Balgruuf）を順に彫り、残り
// （Whiterun）だけが候補になる。彫り出し後に端へ残った接続語（of）は捨てる。
func TestCarveRecursive(t *testing.T) {
	texts := []string{
		"Meet Jarl Balgruuf of Whiterun now. We meet people.",
		"He left the College of Winterhold quietly.",
	}
	terms := detect(t, []string{"Jarl", "Balgruuf", "Winterhold"}, texts, nil)
	if !has(terms, "Whiterun") {
		t.Errorf("複数の既知語を彫った残り Whiterun が候補にない: %v", terms)
	}
	if !has(terms, "College") {
		t.Errorf("接続語の端を落とした College が候補にない: %v", terms)
	}
	for _, banned := range []string{"Jarl", "Balgruuf", "Winterhold", "College of"} {
		if has(terms, banned) {
			t.Errorf("%q は候補にしない: %v", banned, terms)
		}
	}
}

// 派生分割の除外を突く。派生先が既知語なら足さない。派生先が既に候補にあるなら
// 上書きしない（直接の根拠を保つ）。
func TestDeriveSplitSkips(t *testing.T) {
	texts := []string{
		"He met the Jarl. She thanked the Jarl.",
		"They spoke of Jarl Whiterun Kodlak once. People speak.", // 派生 Whiterun Kodlak…でなく検証は下記
		"Runners took news to Jarl Bryling swiftly.",             // 派生先 Bryling は内側 1 語でも出る → 上書きしない
	}
	cand := NewCandidateDetector([]string{"Whiterun Kodlak"}, testStoplist(t), nil).DetectCandidates(texts)
	var bryling *CandidateTerm
	for i := range cand {
		if cand[i].Term == "Whiterun Kodlak" {
			t.Errorf("既知語 Whiterun Kodlak へ派生した: %v", cand)
		}
		if cand[i].Term == "Bryling" {
			bryling = &cand[i]
		}
	}
	if bryling == nil {
		t.Fatalf("Bryling が候補にない: %v", cand)
	}
	// 内側 1 語の直接根拠（文中出現）が派生で上書きされないこと。
	if bryling.MidSentence == 0 {
		t.Errorf("Bryling の直接根拠（文中出現）が失われた: %+v", *bryling)
	}
}

// 彫り出し残りの端の接続語を突く。既知語（College）を彫った残りの先頭に接続語（of）が
// 残る形（of Winterhold）は、接続語を落として Winterhold だけを候補にする。
func TestCarveLeadingConnector(t *testing.T) {
	terms := detect(t, []string{"College"}, []string{"He left the College of Winterhold quietly."}, nil)
	if !has(terms, "Winterhold") {
		t.Errorf("彫り出し残りの Winterhold が候補にない: %v", terms)
	}
	if has(terms, "of Winterhold") {
		t.Errorf("接続語で始まる候補を出した: %v", terms)
	}
}

// prose 実装の決定性と基本動作を突く。同一入力に同一出力を返すこと、明確な固有表現
// （人名 2 語）を拾うこと、文頭が代名詞の文を動詞先頭と判定しないこと。
// 学習モデルの判定品質そのものは仕様にしないため、境界の緩い断定だけを置く。
func TestProseAnalyzer(t *testing.T) {
	a := ProseAnalyzer{}
	text := "George Washington crossed the Delaware."
	first := a.Entities(text)
	second := a.Entities(text)
	if !reflect.DeepEqual(first, second) {
		t.Errorf("Entities が同一入力で変動した: %v / %v", first, second)
	}
	// 文頭が代名詞のため動詞先頭判定は false で固定でき、prose の曖昧さを避ける。
	if a.LeadingVerb("He kills rabbits.") {
		t.Errorf("代名詞先頭の文を動詞先頭と判定した")
	}
	// 記号トークンを読み飛ばして最初の意味トークンで判定すること。
	if a.LeadingVerb("... he ran away.") {
		t.Errorf("記号の後の代名詞先頭の文を動詞先頭と判定した")
	}
	if a.LeadingVerb("") {
		t.Errorf("空文を動詞先頭と判定した")
	}
}
