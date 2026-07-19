package model

// ReferenceTranslation は record 単位の既存訳（参照訳）。xTranslator 英日 XML から取り込み、
// 翻訳前の完全一致置換（known-issues 項目7）で (rec, field, source) が一致する叙述文・台詞へ dest を流用する。
// db タグは db/migrations の reference_translation テーブル列に対応する。
type ReferenceTranslation struct {
	Rec    string `db:"rec"`
	Field  string `db:"field"`
	Source string `db:"source"`
	Dest   string `db:"dest"`
}
