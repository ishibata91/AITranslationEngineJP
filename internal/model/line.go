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
}

// SpeakerIdentity は台詞に紐づく話者の事実上の同定情報。extractor が master 連鎖から解決して書いた
// 種族・声型・所属勢力の識別子（EditorID）をそのまま読む。口調などの解釈は持たない。
// 識別子から口調 traits への変換は engine のルール（解釈）が行う（responsibility 分担）。
type SpeakerIdentity struct {
	RaceEDID     string
	VoiceEDID    string
	FactionEDIDs []string
}

// SpeakerPersona は話者の口調 traits の投影。engine のルールが SpeakerIdentity から組む、
// ペルソナ口調指示文の組み立て入力。種族の気質・声質・所属の気風を持つ。
type SpeakerPersona struct {
	RaceNature     string
	VoiceNature    string
	FactionNatures []string
}
