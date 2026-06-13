// Package model は concept-model.md の箱に対応するデータ構造を持つ。
package model

// Narration は概念モデルの「叙述文」。1 レコード 1 フィールドの叙述テキストと、その訳・状態・配置情報を持つ。
// db タグは db/migrations の narration テーブル列に対応する。
type Narration struct {
	ID      int64  `db:"id"`
	Source  string `db:"source"`
	Dest    string `db:"dest"`
	Status  int    `db:"status"`
	Style   string `db:"style"`
	Plugin  string `db:"plugin"`
	FormID  string `db:"form_id"`
	EDID    string `db:"edid"`
	Rec     string `db:"rec"`
	Field   string `db:"field"`
	Ordinal int    `db:"ordinal"`
}
