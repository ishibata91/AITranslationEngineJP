// Command poc-missing-term は辞書漏れ語抽出（known-issues 1番）の評価ハーネス。
// 既知語 held-out 方式で候補検出（mention.CandidateDetector）を測る:
//  1. C# 抽出器が書いた評価 DB（extracted_field ＋ master_term）を読み、取込段と同じ振り分け
//     （engine.Dispatch）で narration・line の本文と proper_noun を得る。
//  2. 辞書（master_term ∪ proper_noun、stoplist 選別後）の言及が本文に出る行を言及検出
//     （mention.Detector、機械置換と同じ規則）で特定し、決定的に N 件サンプリングする。
//  3. サンプル本文中の言及語を辞書から隠し、候補検出へ「未知語」として見せる。隠した語が正解。
//  4. recall（隠した語のうち検出できた割合）・precision（候補のうち隠した語だった割合）・
//     重複（正規化後同一の候補の複数出力）を測り、取りこぼしと誤検出の一覧を出す。
//
// precision は「候補が隠した語か」の代理指標で、辞書に本当に無い固有名（本 task の本来の獲物）も
// 誤検出側に数える保守的な値になる。誤検出一覧の目視で内訳を確かめる前提の指標。
// 指標の定義と達成基準は docs/exec-plans/active/dictionary-missing-term-detection/goal.md。
//
// 使い方:
//	go run ./cmd/poc-missing-term --db tmp/missing-term-eval/dev-skyrim.sqlite3
//	go run ./cmd/poc-missing-term --db dev.sqlite3,heldout.sqlite3 --n 1000 --seed 1 --ner --runs 3 --dump 40
package main

import (
	"database/sql"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"aitranslationenginejp/internal/core/dictionary"
	"aitranslationenginejp/internal/core/mention"
	"aitranslationenginejp/internal/engine"
	"aitranslationenginejp/internal/model"

	_ "modernc.org/sqlite"
)

const stopwordsPath = "assets/stopwords-en.txt"

// evalRow は評価対象の本文 1 件。key は決定的サンプリングの安定ソートキー。
type evalRow struct {
	plugin string
	key    string
	text   string
	terms  []string // 言及検出が見つけた辞書の原語（隠す対象）
}

func main() {
	dbList := flag.String("db", "", "評価 DB（C# 抽出器の --sqlite 出力）。カンマ区切りで複数可")
	n := flag.Int("n", 1000, "サンプル件数")
	seed := flag.Int64("seed", 1, "サンプリングの乱数シード（決定的）")
	useNER := flag.Bool("ner", false, "prose の固有表現認識で曖昧位置 1 語の判定を補強する")
	runs := flag.Int("runs", 1, "同一サンプルへの実行回数（決定性の確認用）")
	dump := flag.Int("dump", 30, "誤検出・取りこぼし一覧の表示件数")
	trace := flag.String("trace", "", "指定語を含むサンプル本文と候補を表示する（取りこぼしの個別調査用）")
	fpOut := flag.String("fp-out", "", "誤検出の全件を文脈付き TSV で書き出すパス（人手ラベリング用）")
	flag.Parse()
	if *dbList == "" {
		fmt.Fprintln(os.Stderr, "--db は必須。")
		os.Exit(1)
	}
	if err := run(strings.Split(*dbList, ","), *n, *seed, *useNER, *runs, *dump, *trace, *fpOut); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(dbPaths []string, n int, seed int64, useNER bool, runs, dump int, trace, fpOut string) error {
	stopFile, err := os.Open(stopwordsPath)
	if err != nil {
		return fmt.Errorf("stopword リストを開く: %w", err)
	}
	defer stopFile.Close()
	stop, err := dictionary.ParseStoplist(stopFile)
	if err != nil {
		return err
	}

	// 全 DB の本文と辞書原語を合流する（held-out の複数 plugin を 1 つの評価にまとめる）。
	var texts []evalRow
	var dictSources []string
	seenSource := make(map[string]bool)
	for _, path := range dbPaths {
		rows, sources, err := loadDB(strings.TrimSpace(path), stop)
		if err != nil {
			return err
		}
		texts = append(texts, rows...)
		for _, s := range sources {
			if !seenSource[s] {
				seenSource[s] = true
				dictSources = append(dictSources, s)
			}
		}
	}
	fmt.Printf("本文 %d 件・辞書原語 %d 件（stoplist 選別後）\n", len(texts), len(dictSources))

	// 言及付きの本文だけを評価対象にする（言及語を隠せば正解ラベルになる）。
	vocab := make([]mention.Term, len(dictSources))
	for i, s := range dictSources {
		vocab[i] = mention.Term{Kind: "dict", ID: int64(i), Source: s}
	}
	det := mention.NewDetector(vocab)
	var eligible []evalRow
	selfMention := 0
	for _, r := range texts {
		// 言及は動的タグ（<Alias=...> 等）を除いた本文で取る。タグ内の語は地の文の固有名で
		// ないため、隠しても本文からは検出できず、正解ラベルの偽陽性になる。
		found := det.Detect(markupRe.ReplaceAllString(r.text, " "))
		if len(found) == 0 {
			continue
		}
		// 正解ラベルは固有名詞の表記形（大文字始まりの語と接続語 of / the の並び）に限る。
		// xTranslator XML 由来の記号だけの行（"."）や文の形の行（"Let's go"）は固有名詞で
		// ないため、隠す対象にしない（隠さないので既知語として残り、候補からも除かれる）。
		r.terms = nil
		for _, t := range found {
			if titleCaseTerm(t.Source) {
				r.terms = append(r.terms, t.Source)
			}
		}
		if len(r.terms) == 0 {
			continue
		}
		// 本文全体が名前そのもの（レコード名の再掲や内部 ID の 1 行）は地の文・台詞でないため
		// 評価対象から除く。対象語は「本文中に出現する固有名詞」（goal.md）。
		if len(r.terms) == 1 && plainText(r.text) == r.terms[0] {
			selfMention++
			continue
		}
		eligible = append(eligible, r)
	}
	fmt.Printf("言及付き本文 %d 件（本文全体が名前の行 %d 件を除外済み）\n", len(eligible), selfMention)

	// 決定的サンプリング: 安定キーで整列してから固定シードで混ぜ、先頭 N 件を採る。
	sort.Slice(eligible, func(i, j int) bool { return eligible[i].key < eligible[j].key })
	rand.New(rand.NewSource(seed)).Shuffle(len(eligible), func(i, j int) {
		eligible[i], eligible[j] = eligible[j], eligible[i]
	})
	if len(eligible) > n {
		eligible = eligible[:n]
	}
	perPlugin := make(map[string]int)
	for _, r := range eligible {
		perPlugin[r.plugin]++
	}
	fmt.Printf("サンプル %d 件（plugin 内訳: %v）\n", len(eligible), perPlugin)

	// サンプル本文の言及語を正解（隠す対象）に確定し、残りを既知語として検出器へ渡す。
	hidden := make(map[string]bool)
	for _, r := range eligible {
		for _, t := range r.terms {
			hidden[mention.NormalizeCandidate(t)] = true
		}
	}
	var known []string
	for _, s := range dictSources {
		if !hidden[mention.NormalizeCandidate(s)] {
			known = append(known, s)
		}
	}
	sampleTexts := make([]string, len(eligible))
	for i, r := range eligible {
		sampleTexts[i] = r.text
	}
	fmt.Printf("正解（隠した語）%d 語・残り既知語 %d 語\n\n", len(hidden), len(known))

	var ner mention.TextAnalyzer
	if useNER {
		ner = mention.ProseAnalyzer{}
	}
	var prevRecall, prevPrecision float64
	for run := 1; run <= runs; run++ {
		cand := mention.NewCandidateDetector(known, stop, ner).DetectCandidates(sampleTexts)
		if trace != "" && run == 1 {
			traceTerm(trace, eligible, cand)
		}
		if fpOut != "" && run == 1 {
			if err := writeFalsePositives(fpOut, cand, hidden, eligible); err != nil {
				return err
			}
		}
		recall, precision := report(cand, hidden, eligible, dump, run == runs)
		if run > 1 && (recall != prevRecall || precision != prevPrecision) {
			fmt.Printf("!! 決定性違反: run %d で指標が変動（recall %.4f→%.4f、precision %.4f→%.4f）\n",
				run, prevRecall, recall, prevPrecision, precision)
		}
		prevRecall, prevPrecision = recall, precision
		fmt.Printf("run %d/%d: recall %.1f%%・precision %.1f%%\n", run, runs, recall*100, precision*100)
	}
	return nil
}

// report は候補と正解を突き合わせ、指標と誤検出・取りこぼし一覧を出す。verbose の時だけ一覧を出す。
func report(cand []mention.CandidateTerm, hidden map[string]bool, sample []evalRow, dump int, verbose bool) (recall, precision float64) {
	seen := make(map[string]bool, len(cand))
	dups := 0
	tp := 0
	var fps []mention.CandidateTerm
	for _, c := range cand {
		if seen[c.Term] {
			dups++
			continue
		}
		seen[c.Term] = true
		if hidden[c.Term] {
			tp++
		} else {
			fps = append(fps, c)
		}
	}
	var missed []string
	for h := range hidden {
		if !seen[h] {
			missed = append(missed, h)
		}
	}
	sort.Strings(missed)
	if len(hidden) > 0 {
		recall = float64(tp) / float64(len(hidden))
	}
	if len(cand) > 0 {
		precision = float64(tp) / float64(len(cand))
	}
	if !verbose {
		return recall, precision
	}

	fmt.Printf("候補 %d 語（重複 %d）・正解一致 %d・誤検出 %d・取りこぼし %d\n",
		len(cand), dups, tp, len(fps), len(missed))
	sort.Slice(fps, func(i, j int) bool {
		if fps[i].Occurrences != fps[j].Occurrences {
			return fps[i].Occurrences > fps[j].Occurrences
		}
		return fps[i].Term < fps[j].Term
	})
	fmt.Printf("\n## 誤検出（出現数順・上位 %d）\n", dump)
	for i, c := range fps {
		if i >= dump {
			break
		}
		fmt.Printf("  %-40q 出現 %d・文中 %d・NER %d\n", c.Term, c.Occurrences, c.MidSentence, c.NERHits)
	}
	fmt.Printf("\n## 取りこぼし（上位 %d）\n", dump)
	for i, h := range missed {
		if i >= dump {
			break
		}
		fmt.Printf("  %-40q 例: %s\n", h, exampleText(sample, h))
	}
	fmt.Println()
	return recall, precision
}

// traceTerm は指定語の取りこぼし調査用に、語を含むサンプル本文（最大 3 件）と、
// 語を部分文字列に含む候補を表示する。
func traceTerm(term string, sample []evalRow, cand []mention.CandidateTerm) {
	fmt.Printf("## trace %q\n", term)
	re := regexp.MustCompile(`\b` + regexp.QuoteMeta(term) + `\b`)
	shown := 0
	for _, r := range sample {
		if !re.MatchString(r.text) {
			continue
		}
		fmt.Printf("  本文[%s]: %s\n", r.key, strings.ReplaceAll(r.text, "\n", "⏎"))
		shown++
		if shown >= 3 {
			break
		}
	}
	for _, c := range cand {
		if strings.Contains(c.Term, term) || strings.Contains(term, c.Term) {
			fmt.Printf("  候補: %-40q 出現 %d・文中 %d・NER %d\n", c.Term, c.Occurrences, c.MidSentence, c.NERHits)
		}
	}
	fmt.Println()
}

// writeFalsePositives は誤検出（正解に無い候補）の全件を、出現数と例文付きの TSV へ書く。
// 人手ラベリング（候補が固有名詞か否か）で真の精度を測るための材料。代理指標の誤検出には
// 辞書に本当に無い固有名（本 task の本来の獲物）が混ざるため、内訳の確認が要る。
func writeFalsePositives(path string, cand []mention.CandidateTerm, hidden map[string]bool, sample []evalRow) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("誤検出一覧の書き出し: %w", err)
	}
	defer f.Close()
	fmt.Fprintln(f, "term\toccurrences\tmid\tner\texample")
	for _, c := range cand {
		if hidden[c.Term] {
			continue
		}
		fmt.Fprintf(f, "%s\t%d\t%d\t%d\t%s\n", c.Term, c.Occurrences, c.MidSentence, c.NERHits, exampleText(sample, c.Term))
	}
	return nil
}

// markupRe は言及検出・表示・自己言及判定用の動的タグ除去（候補検出側の除去規則と同じ形）。
var markupRe = regexp.MustCompile(`<[^>]*>`)

var (
	termTokenRe     = regexp.MustCompile(`[A-Za-z][A-Za-z'’-]*`)
	termLowerNameRe = regexp.MustCompile(`^[a-z'’]+[-'’][A-Z]`)
)

// titleCaseTerm は辞書の原語が固有名詞の表記形かを返す。語（英字始まり）がすべて大文字始まり
// または小文字接頭の姓（gro-Nolob）で、小文字の接続語（of / the）は内側にだけ許す。
// 英字の語が無い原語（記号だけ）は固有名詞でない。
func titleCaseTerm(source string) bool {
	toks := termTokenRe.FindAllString(source, -1)
	if len(toks) == 0 {
		return false
	}
	for i, t := range toks {
		if t == "of" || t == "the" {
			if i == 0 || i == len(toks)-1 {
				return false
			}
			continue
		}
		r, _ := utf8.DecodeRuneInString(t)
		if !unicode.IsUpper(r) && !termLowerNameRe.MatchString(t) {
			return false
		}
	}
	return true
}

// plainText は本文からタグを除き前後の空白・引用符を落とす（自己言及の判定用）。
func plainText(text string) string {
	return strings.Trim(markupRe.ReplaceAllString(text, " "), " \t\r\n\"'")
}

// exampleText は語を含むサンプル本文の先頭 1 件を返す（取りこぼしの目視確認用）。
func exampleText(sample []evalRow, term string) string {
	re, err := regexp.Compile(`\b` + regexp.QuoteMeta(term) + `\b`)
	if err != nil {
		return "(照合不能)"
	}
	for _, r := range sample {
		if re.MatchString(r.text) {
			t := strings.ReplaceAll(r.text, "\n", " ")
			if len(t) > 90 {
				t = t[:90] + "…"
			}
			return t
		}
	}
	return "(本文に完全一致なし)"
}

// loadDB は評価 DB 1 つから本文（narration・line 相当）と辞書原語（master_term ∪ proper_noun、
// stoplist 選別後）を読む。振り分けは取込段と同じ engine.Dispatch を使い、本番と同じ箱で評価する。
func loadDB(path string, stop *dictionary.Stoplist) ([]evalRow, []string, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, nil, fmt.Errorf("評価 DB を開く（%s）: %w", path, err)
	}
	defer db.Close()

	fields, err := readFields(db)
	if err != nil {
		return nil, nil, fmt.Errorf("extracted_field の読み取り（%s）: %w", path, err)
	}
	master, err := readRecordMaster(db)
	if err != nil {
		return nil, nil, fmt.Errorf("record_type_master の読み取り（%s）: %w", path, err)
	}
	termSources, err := readMasterTermSources(db)
	if err != nil {
		return nil, nil, fmt.Errorf("master_term の読み取り（%s）: %w", path, err)
	}

	d := engine.Dispatch(fields, master)

	// 本文は同一テキストを 1 件に畳む（台詞の重複行がサンプルを占有しないようにする）。
	rows := make(map[string]evalRow)
	keep := func(r evalRow) {
		if old, ok := rows[r.text]; !ok || r.key < old.key {
			rows[r.text] = r
		}
	}
	for _, nr := range d.Narrations {
		keep(evalRow{
			plugin: nr.Plugin, text: nr.Source,
			key: fmt.Sprintf("n|%s|%s|%s|%s|%d", nr.Plugin, nr.FormID, nr.Rec, nr.Field, nr.Ordinal),
		})
	}
	for _, ln := range d.Lines {
		keep(evalRow{
			plugin: ln.Plugin, text: ln.Source,
			key: fmt.Sprintf("l|%s|%s|%s|%s|%d", ln.Plugin, ln.FormID, ln.Rec, ln.Field, ln.Ordinal),
		})
	}
	out := make([]evalRow, 0, len(rows))
	for _, r := range rows {
		out = append(out, r)
	}

	// 辞書原語は本番の供給（translationVocabulary）と同じに master_term → proper_noun の順で
	// stoplist 選別を通す。
	var sources []string
	for _, s := range termSources {
		if !stop.Blocks(s) {
			sources = append(sources, s)
		}
	}
	for _, pn := range d.ProperNouns {
		if !stop.Blocks(pn.Source) {
			sources = append(sources, pn.Source)
		}
	}
	return out, sources, nil
}

func readFields(db *sql.DB) ([]model.ExtractedField, error) {
	rows, err := db.Query(`SELECT plugin, form_id, edid, rec, field, ordinal, source FROM extracted_field`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.ExtractedField
	for rows.Next() {
		var f model.ExtractedField
		if err := rows.Scan(&f.Plugin, &f.FormID, &f.EDID, &f.Rec, &f.Field, &f.Ordinal, &f.Source); err != nil {
			return nil, err
		}
		out = append(out, f)
	}
	return out, rows.Err()
}

func readRecordMaster(db *sql.DB) (map[engine.RecordKey]model.RecordType, error) {
	rows, err := db.Query(`SELECT rec, field, box, directive FROM record_type_master`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make(map[engine.RecordKey]model.RecordType)
	for rows.Next() {
		var rt model.RecordType
		if err := rows.Scan(&rt.Rec, &rt.Field, &rt.Box, &rt.Directive); err != nil {
			return nil, err
		}
		out[engine.RecordKey{Rec: rt.Rec, Field: rt.Field}] = rt
	}
	return out, rows.Err()
}

func readMasterTermSources(db *sql.DB) ([]string, error) {
	rows, err := db.Query(`SELECT source FROM master_term ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}
