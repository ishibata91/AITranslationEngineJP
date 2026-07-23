package model

// ExtractedField は C# 抽出器が素朴吸い出しした原文 1 件。箱・directive の判定を持たない受け皿。
// Go 取込段が record_type_master を rec/field で引いて box と directive を得て、
// narration / proper_noun / line へ振り分ける。
// db タグは db/migrations の extracted_field 列に対応する。
type ExtractedField struct {
	ID      int64  `db:"id"`
	Plugin  string `db:"plugin"`
	FormID  string `db:"form_id"`
	EDID    string `db:"edid"`
	Rec     string `db:"rec"`
	Field   string `db:"field"`
	Ordinal int    `db:"ordinal"`
	Source  string `db:"source"`
	// Dest は日本語既訳（英日 Strings が両方あった field のみ非空）。参照訳・固有名の確定訳語の供給源。
	Dest string `db:"dest"`
}
