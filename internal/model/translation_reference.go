package model

// PrebuiltDictionaryReference は事前作成済み辞書readerが返す意味単位の候補である。
// Meaning はreaderの検証と重複判定だけに使い、本文翻訳の入力へは渡さない。
type PrebuiltDictionaryReference struct {
	Source         string `db:"source"`
	Dest           string `db:"dest"`
	PartOfSpeech   string `db:"part_of_speech"`
	Meaning        string `db:"meaning"`
	SkyrimCategory string `db:"skyrim_category"`
}

// TranslationReference は本文翻訳と結果表示が共有する参考語である。
// Meaning を持たないため、本文prompt、snapshot、API、UIへ意味欄を流せない。
type TranslationReference struct {
	Source         string `json:"source"`
	Dest           string `json:"dest"`
	PartOfSpeech   string `json:"partOfSpeech"`
	SkyrimCategory string `json:"skyrimCategory"`
	Origin         string `json:"origin"`
}

// TranslationReferenceSnapshot は送信時の本文参考語とprompt hashを保存する単位である。
type TranslationReferenceSnapshot struct {
	Plugin     string
	Kind       string
	RowID      int64
	References []TranslationReference
	PromptHash string
}
