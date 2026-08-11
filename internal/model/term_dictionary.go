package model

// TermDictionaryFilter は用語辞書の一覧を絞り込む入力を表す。
type TermDictionaryFilter struct {
	Source       string
	Destination  string
	PartOfSpeech string
	Category     string
}

// TermDictionaryEntry は用語辞書の一行を表す。
type TermDictionaryEntry struct {
	ID           int64
	Source       string
	Destination  string
	PartOfSpeech string
	Revision     int64
	Categories   []string
}

// TermDictionaryPage は固定件数の用語辞書一覧を表す。
type TermDictionaryPage struct {
	Entries    []TermDictionaryEntry
	TotalCount int
	PageNumber int
}

// TermDictionaryCreate は用語辞書を作成する入力を表す。
type TermDictionaryCreate struct {
	Source       string
	Destination  string
	PartOfSpeech string
}

// TermDictionaryPatch は用語辞書の変更項目を表す。
type TermDictionaryPatch struct {
	ID           int64
	Revision     int64
	Source       *string
	Destination  *string
	PartOfSpeech *string
}
