package model

// LineAnalysis は 1 本文の言語特徴キャッシュ（line_analysis）。最も重い prose の品詞解析を本文ハッシュで 1 度だけ持つ。
// db タグは db/migrations の line_analysis 列に対応する。
type LineAnalysis struct {
	ID            int64  `db:"id"`
	SourceHash    string `db:"source_hash"`
	SentenceCount int    `db:"sentence_count"`
	PoliteCount   int    `db:"polite_count"`
	InsultCount   int    `db:"insult_count"`
	IsImperative  int    `db:"is_imperative"`
	ExclaimCount  int    `db:"exclaim_count"`
	ElongCount    int    `db:"elong_count"`
	EmotionCount  int    `db:"emotion_count"`
}

// PersonaCharacter は対話由来で生成した話者の基底口調（persona_character）。
// 話者の安定識別（解決した base NPC の plugin・form_id）をキーにする。
type PersonaCharacter struct {
	ID            int64  `db:"id"`
	SpeakerPlugin string `db:"speaker_plugin"`
	SpeakerFormID string `db:"speaker_form_id"`
	AttitudeBand  int    `db:"attitude_band"`
	EmotionBand   int    `db:"emotion_band"`
	Marked        int    `db:"marked"`
	DecisionPath  string `db:"decision_path"`
	HandEdited    int    `db:"hand_edited"`
}

// SpeakerLineSource は生成の入力 1 行。話者の安定識別・声型 EditorID・種族 EditorID と、その話者が話す 1 本文。
// 1 話者は複数行を持つため、呼び出し側が話者ごとに束ねて基底口調を集計する。
type SpeakerLineSource struct {
	SpeakerPlugin string `db:"speaker_plugin"`
	SpeakerFormID string `db:"speaker_form_id"`
	VoiceEDID     string `db:"voice_edid"`
	RaceEDID      string `db:"race_edid"`
	Source        string `db:"source"`
}

// LineSpeaker は台詞 1 行の話者識別と属性。誰の台詞かと、口調指示の根拠（性別・年齢＝種族・声型）を表示するために読む。
// 1 台詞に複数話者が紐づく場合は呼び出し側が id 昇順の先頭話者を採る。
type LineSpeaker struct {
	LineID    int64  `db:"line_id"`
	EDID      string `db:"edid"`       // 話者の EditorID（例 AventusAretino）
	Sex       string `db:"sex"`        // 性別（"Female"/"Male"/""）
	RaceEDID  string `db:"race_edid"`  // 種族 EditorID（年齢区分の判定元。ElderRace/*Child 等）
	VoiceEDID string `db:"voice_edid"` // 声型 EditorID（例 FemaleOldGrumpy）
}

// LineSpeakerSexSet は台詞に結ばれた全話者の性別集合である。
// Sex は全話者が Male または Female で揃う場合だけその値を持つ。
// 性別不明または男女混在では空である。
type LineSpeakerSexSet struct {
	Count int
	Sex   string
}

// LinePersonaInput は注入の入力 1 件。台詞 1 行へ与える話者の生成済み基底口調と、役割語マーカー用の性別・種族 EditorID。
// persona_character を持つ話者の台詞だけが現れる。Sex・RaceEDID は一人称・語尾テンプレートと種族訛りの引きに使う。
type LinePersonaInput struct {
	LineID       int64  `db:"line_id"`
	AttitudeBand int    `db:"attitude_band"`
	EmotionBand  int    `db:"emotion_band"`
	Marked       int    `db:"marked"`
	DecisionPath string `db:"decision_path"`
	Sex          string `db:"sex"` // 話者の性別（"Female"/"Male"/""）。役割語テンプレートの引きに使う。
	RaceEDID     string `db:"race_edid"`
	// EmotionType は台詞の TRDT 感情型（非 Neutral のみ）。engine が model.Line から重ねる。DB からは load しない。
	EmotionType string
}
