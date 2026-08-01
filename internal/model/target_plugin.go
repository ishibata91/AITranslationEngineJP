package model

// TargetPlugin は翻訳した対象 plugin の登録単位（target_plugin テーブル）。plugin と 1 対 1。
// Plugin は plugin ファイル名（例 Dawnguard.esm）で識別子。SourcePath は選んだ plugin のフルパス。
// CreatedAt は初回登録時刻の文字列。Total / Translated は一覧表示用の計算値で、
// この plugin が束ねる翻訳対象（叙述文＋台詞＋固有名）の総数と、うち訳済み（status != 0）の数。
// sync_retry_ready は同期実行の内部判断だけに使い、一覧 DTO へは出さない。
// Total / Translated は List クエリの計算列から埋める。
type TargetPlugin struct {
	Plugin     string `db:"plugin"`
	SourcePath string `db:"source_path"`
	CreatedAt  string `db:"created_at"`
	Total      int    `db:"total"`
	Translated int    `db:"translated"`
}
