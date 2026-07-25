// Package rolespeech は役割語テンプレート（一人称・語尾 register）の解析と照合を行う純粋ルール。
// ファイル読み込みは持たない。呼び出し側が io.Reader を渡す（os 読みは composition root と cmd の責務）。
package rolespeech

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Template は役割語の出力。一人称と、言い回しの傾向（語尾の register）と、口調の例文。
// 一人称が空なら基底口調へ委ね、言い回しが空なら指定しない。例文が空なら例示しない。
type Template struct {
	FirstPerson string  // 一人称（例「わたし」「ぼく」。空なら委ねる）
	Register    string  // 言い回しの傾向（例「年配の女性らしい、落ち着いた言い回しにする。」。空なら指定なし）
	Example     Example // 口調の例文（英語原文と日本語訳文の 1 対）。空なら例示なし
}

// Example は口調の例文 1 対。説明文だけでは揺れる一人称・語尾を、訳例で固定するために持つ。
// 台詞は 1 件 1 リクエストで送るため、同じ話者の複数台詞が独立に推論される。その揺れを実例で抑える。
type Example struct {
	Source string // 英語原文（短い台詞 1 行）
	Dest   string // 日本語訳文（そのセル・性別・年齢区分で期待する訳）
}

// IsZero は例文が空（原文か訳文のどちらかを欠く）かを返す。片側だけの行は例示に使わない。
func (e Example) IsZero() bool { return e.Source == "" || e.Dest == "" }

// roleSpeechRow は role-speech テンプレートの 1 行。キー 3 列（race 区分・性別・基底口調セル）と値。
// キー列の "*" はワイルドカード（任意一致）。
type roleSpeechRow struct {
	race, sex, cell string
	tmpl            Template
}

// exampleRow は例文テンプレートの 1 行。キー 3 列は roleSpeechRow と同じで、値が例文 1 対。
type exampleRow struct {
	race, sex, cell string
	example         Example
}

// roleWildcard はキー列のワイルドカード（任意一致）。
const roleWildcard = "*"

// Table は一人称・語尾テンプレートと口調の例文の表。話者の（race 区分・性別・基底口調セル）で引く。
// 中身は実画面確認で見直すため assets/role-speech.tsv と assets/role-speech-examples.tsv に外部化し、
// composition root が 2 ファイルを読んで 1 つの Table を組む。1 行が長くなるのを避けて表を分けたが、
// キーと照合規則は同じなので Lookup が 2 つの行集合をそれぞれ最具体一致で引き、1 つの Template に束ねる。
// 照合（キー → テンプレート、ワイルドカード優先順位）は純粋関数で、DB・prose に依存しない。
type Table struct {
	rows     []roleSpeechRow
	examples []exampleRow
}

// ParseRoleSpeech はタブ区切りの role-speech テンプレート（race・sex・cell・一人称・言い回し）を読む。
// "#" で始まる行と空行は読み飛ばす。列が 5 未満の行はエラーにする。
func ParseRoleSpeech(r io.Reader) (*Table, error) {
	var rows []roleSpeechRow
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := strings.TrimRight(sc.Text(), "\r")
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		f := strings.Split(line, "\t")
		if len(f) < 5 {
			return nil, fmt.Errorf("役割語テンプレートの列が足りない（5 列必要）: %q", line)
		}
		rows = append(rows, roleSpeechRow{
			race: f[0], sex: f[1], cell: f[2],
			tmpl: Template{FirstPerson: f[3], Register: f[4]},
		})
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("役割語テンプレートの読み取り: %w", err)
	}
	return &Table{rows: rows}, nil
}

// ParseRoleSpeechExamples はタブ区切りの例文テンプレート（race・sex・cell・英語原文・日本語訳文）を読み、
// 既存の Table へ例文行として束ねて返す。"#" で始まる行と空行は読み飛ばす。列が 5 未満の行はエラーにする。
// 役割語表と例文表はキーが同じで値だけが違うため、別ファイルで読んで 1 つの Table に持たせる。
// base が nil なら例文だけを持つ Table を返す（役割語未配線でも例文の解析だけは通す）。
func ParseRoleSpeechExamples(base *Table, r io.Reader) (*Table, error) {
	var rows []exampleRow
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := strings.TrimRight(sc.Text(), "\r")
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		f := strings.Split(line, "\t")
		if len(f) < 5 {
			return nil, fmt.Errorf("口調例文テンプレートの列が足りない（5 列必要）: %q", line)
		}
		rows = append(rows, exampleRow{
			race: f[0], sex: f[1], cell: f[2],
			example: Example{Source: f[3], Dest: f[4]},
		})
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("口調例文テンプレートの読み取り: %w", err)
	}
	if base == nil {
		return &Table{examples: rows}, nil
	}
	return &Table{rows: base.rows, examples: rows}, nil
}

// Lookup は（race 区分・性別・セル）に最も具体的に一致するテンプレートを返す。
// 各キー列は完全一致かワイルドカード "*" で一致し、ワイルドカードが少ない行（具体的な行）を優先する。
// 同点は先に現れた行を採る。役割語行と例文行は同じキーでそれぞれ独立に引き、1 つの Template へ束ねる。
// どちらか一方でも一致すれば ok=true。両方とも一致が無ければ ok=false。
// nil レシーバは ok=false（テンプレート未配線でも安全）。
func (t *Table) Lookup(race, sex, cell string) (Template, bool) {
	if t == nil {
		return Template{}, false
	}
	best := -1
	var found Template
	for _, row := range t.rows {
		if score, ok := matchScore(row.race, row.sex, row.cell, race, sex, cell); ok && score > best {
			best = score
			found = row.tmpl
		}
	}
	bestExample := -1
	for _, row := range t.examples {
		if score, ok := matchScore(row.race, row.sex, row.cell, race, sex, cell); ok && score > bestExample {
			bestExample = score
			found.Example = row.example
		}
	}
	return found, best >= 0 || bestExample >= 0
}

// matchScore は行のキー（rowRace・rowSex・rowCell）が引き当て先（race・sex・cell）に一致するかと、
// 具体度（非ワイルドカード一致の数）を返す。いずれかのキー列が一致しなければ ok=false。
func matchScore(rowRace, rowSex, rowCell, race, sex, cell string) (score int, ok bool) {
	for _, p := range [...][2]string{{rowRace, race}, {rowSex, sex}, {rowCell, cell}} {
		switch p[0] {
		case roleWildcard:
			// ワイルドカードは一致するが具体度は上げない。
		case p[1]:
			score++
		default:
			return 0, false
		}
	}
	return score, true
}

// RoleClassOfRace は race EditorID を役割区分（child / elder / adult）へ畳む。
// *Child は子供、ElderRace は老人、それ以外は成人。役割語テンプレートのキーに使う。
func RoleClassOfRace(raceEDID string) string {
	switch {
	case strings.Contains(raceEDID, "Child"):
		return "child"
	case strings.Contains(raceEDID, "Elder"):
		return "elder"
	default:
		return "adult"
	}
}
