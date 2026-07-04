package model

// ProperNounPlacement は固有名の原文出現位置（extracted_field 由来）に確定訳（proper_noun 由来）を結んだ行。
// proper_noun 辞書はレコード位置（plugin・form_id 等）を持たないため、xTranslator 書き出しで位置つきの
// <String> 行を作るには extracted_field（box='固有名'）と結ぶ必要がある。その結合結果を運ぶ導出データ。
// db タグは書き出し用 join クエリの SELECT 列に対応する。
type ProperNounPlacement struct {
	Plugin  string `db:"plugin"`
	FormID  string `db:"form_id"`
	EDID    string `db:"edid"`
	Rec     string `db:"rec"`
	Field   string `db:"field"`
	Ordinal int    `db:"ordinal"`
	Source  string `db:"source"`
	Dest    string `db:"dest"`
	Status  int    `db:"status"`
}
