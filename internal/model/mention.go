package model

// NarrationMention は叙述文の本文中言及（e4）の 1 件。言及の相手は機械置換辞書の供給源に合わせ、
// master_term か proper_noun のどちらか一方を指す（0 は「指さない」）。
// db タグは db/migrations の narration_mention 列に対応する（0 の列は NULL で書く）。
type NarrationMention struct {
	NarrationID  int64 `db:"narration_id"`
	ProperNounID int64 `db:"proper_noun_id"` // 実行内固有名（proper_noun.id）。0 なら master_term 側。
	MasterTermID int64 `db:"master_term_id"` // 横断マスター辞書（master_term.id）。0 なら proper_noun 側。
}

// LineMention は台詞の本文中言及（e5）の 1 件。NarrationMention と対称。
type LineMention struct {
	LineID       int64 `db:"line_id"`
	ProperNounID int64 `db:"proper_noun_id"`
	MasterTermID int64 `db:"master_term_id"`
}
