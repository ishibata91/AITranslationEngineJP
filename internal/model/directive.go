package model

// Directive は指示文の正本 1 件。プロンプト合成の「Base 指示 ＋ directive」の directive にあたる。
// 口調・文体（説明体/書物体/日記体/世界観断片）・固有名・定型句を 1 つの概念へ揃える。
// Variables は実行時に埋める変数の宣言（JSON 配列 [{token, description}] の生文字列）。口調だけ {traits} を持つ。
// db タグは db/migrations の directive 列に対応する。
type Directive struct {
	Key         string `db:"key"`
	Instruction string `db:"instruction"`
	Variables   string `db:"variables"`
}

// RecordType は record_type_master の 1 行。REC:FIELD の分類（box）と directive 割り当ての正本。
// LogicalName は REC:FIELD だけでは分からない人間向けの種別名。
// db タグは db/migrations の record_type_master 列に対応する。
type RecordType struct {
	Rec         string `db:"rec"`
	Field       string `db:"field"`
	Box         string `db:"box"`
	Directive   string `db:"directive"`
	LogicalName string `db:"logical_name"`
}
