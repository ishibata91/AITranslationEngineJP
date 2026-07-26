package model

// ProperNoun は概念モデルの「固有名」の訳の単位（方針A・plugin スコープの非共有）。同一実行内で AI 訳を留め、
// 横断永続辞書 master_term へは昇格しない。Plugin は非共有スコープ（同綴りの mod 固有名を plugin ごとに別行にする）。
// Source（原語）が本文機械置換の照合キー。Category（種別＝rec）は同綴り異義の区別用（concept-model 弱点1）。
// db タグは db/migrations の proper_noun 列に対応する。
// Origin は 1 行の出どころ（抽出由来か派生か）。OriginDerived 以外は抽出由来＝翻訳対象。
type ProperNoun struct {
	ID       int64  `db:"id"`
	Plugin   string `db:"plugin"`
	Source   string `db:"source"`
	Category string `db:"category"`
	Dest     string `db:"dest"`
	Status   int    `db:"status"`
	Origin   string `db:"origin"`
}

// OriginDerived は機械派生した人名の部分形を表す ProperNoun.Origin の値。
// 部分形は機械置換辞書の材料であり、翻訳対象ではない。原文 record を持たないため、
// 翻訳対象として数える経路（固有名の一覧、対象 plugin の進捗件数）はこの値で外す。
// 原文位置を解決する経路（言及の位置解決、出力位置）は Category を rec と突き合わせるので、
// Category が空の派生行を元から拾わない。書き込み側（engine）と絞り込み側（store）が同じ値を参照する。
const OriginDerived = "derive"

// ConfirmedName は実行内で訳が確定した固有名 1 件を、原文の field つきで表す。
// ProperNoun は氏名（FULL）と短縮名（SHRT）を分ける列を持たないため、原文（extracted_field）と
// (plugin, category = rec, source) で結んで field を取り戻す。人名の部分形を派生する入力になる。
type ConfirmedName struct {
	Field  string `db:"field"`
	Source string `db:"source"`
	Dest   string `db:"dest"`
}
