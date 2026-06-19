package model

// MasterTerm は Mod 横断マスター辞書の 1 件。固有名の「原語 → 確定訳語」の対応。
// 叙述文・台詞の本文へ固有名を機械置換するための辞書語。db タグは db/migrations の master_term 列に対応する。
type MasterTerm struct {
	ID       int64  `db:"id"`
	Source   string `db:"source"`   // 原語（英語 FULL）
	Dest     string `db:"dest"`     // 確定訳語（日本語の公式既訳）
	Category string `db:"category"` // 種別（同綴り異義の区別用。本文置換は source で照合する）
}
