// Package rolespeech は役割語テンプレート（一人称・語尾 register）の解析と照合を行う純粋ルール。
// ファイル読み込みは持たない。呼び出し側が io.Reader を渡す（os 読みは composition root と cmd の責務）。
package rolespeech

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Template は役割語の出力。一人称と、言い回しの傾向（語尾の register）。
// 一人称が空なら基底口調へ委ね、言い回しが空なら指定しない。
type Template struct {
	FirstPerson string // 一人称（例「わたし」「ぼく」。空なら委ねる）
	Register    string // 言い回しの傾向（例「年配の女性らしい、落ち着いた言い回しにする。」。空なら指定なし）
}

// roleSpeechRow は role-speech テンプレートの 1 行。キー 3 列（race 区分・性別・基底口調セル）と値。
// キー列の "*" はワイルドカード（任意一致）。
type roleSpeechRow struct {
	race, sex, cell string
	tmpl            Template
}

// roleWildcard はキー列のワイルドカード（任意一致）。
const roleWildcard = "*"

// Table は一人称・語尾テンプレートの表。話者の（race 区分・性別・基底口調セル）で引く。
// 中身は実画面確認で見直すため assets/role-speech.tsv に外部化し、bootstrap で読む。
// 照合（キー → テンプレート、ワイルドカード優先順位）は純粋関数で、DB・prose に依存しない。
type Table struct{ rows []roleSpeechRow }

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

// Lookup は（race 区分・性別・セル）に最も具体的に一致するテンプレートを返す。
// 各キー列は完全一致かワイルドカード "*" で一致し、ワイルドカードが少ない行（具体的な行）を優先する。
// 同点は先に現れた行を採る。一致が無ければ ok=false。nil レシーバは ok=false（テンプレート未配線でも安全）。
func (t *Table) Lookup(race, sex, cell string) (Template, bool) {
	if t == nil {
		return Template{}, false
	}
	best := -1
	var found Template
	for _, row := range t.rows {
		if score, ok := matchScore(row, race, sex, cell); ok && score > best {
			best = score
			found = row.tmpl
		}
	}
	return found, best >= 0
}

// matchScore は行が（race・sex・cell）に一致するかと、具体度（非ワイルドカード一致の数）を返す。
// いずれかのキー列が一致しなければ ok=false。
func matchScore(row roleSpeechRow, race, sex, cell string) (score int, ok bool) {
	for _, p := range [...][2]string{{row.race, race}, {row.sex, sex}, {row.cell, cell}} {
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
