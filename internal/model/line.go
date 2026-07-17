package model

// Line は概念モデルの「台詞」。話者で口調が変わるため重複排除せず、自前の訳・状態・配置情報を持つ。
// db タグは db/migrations の line テーブル列に対応する。
type Line struct {
	ID            int64  `db:"id"`
	Source        string `db:"source"`
	Dest          string `db:"dest"`
	Status        int    `db:"status"`
	ResponseOrder int    `db:"response_order"`
	Plugin        string `db:"plugin"`
	FormID        string `db:"form_id"`
	EDID          string `db:"edid"`
	Rec           string `db:"rec"`
	Field         string `db:"field"`
	Ordinal       int    `db:"ordinal"`
	// EmotionType は TRDT 由来の台詞感情型（非 Neutral のみ、line.emotion_type 列）。無しは空。
	EmotionType string `db:"emotion_type"`
}
