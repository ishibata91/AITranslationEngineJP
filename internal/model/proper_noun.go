package model

// ProperNoun は概念モデルの「固有名」の訳の単位（方針A・plugin スコープの非共有）。同一実行内で AI 訳を留め、
// 横断永続辞書 master_term へは昇格しない。Plugin は非共有スコープ（同綴りの mod 固有名を plugin ごとに別行にする）。
// Source（原語）が本文機械置換の照合キー。Category（種別＝rec）は同綴り異義の区別用（concept-model 弱点1）。
// db タグは db/migrations の proper_noun 列に対応する。
type ProperNoun struct {
	ID       int64  `db:"id"`
	Plugin   string `db:"plugin"`
	Source   string `db:"source"`
	Category string `db:"category"`
	Dest     string `db:"dest"`
	Status   int    `db:"status"`
}
