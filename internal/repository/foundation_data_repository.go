package repository

import (
	"context"
	"time"
)

// XTranslatorTranslationXML は XTRANSLATOR_TRANSLATION_XML テーブルの 1 レコードを表す。
type XTranslatorTranslationXML struct {
	ID               int64
	FilePath         string
	TargetPluginName string
	TargetPluginType string
	TermCount        int
	ImportedAt       time.Time
}

// XTranslatorTranslationXMLDraft は XTRANSLATOR_TRANSLATION_XML の作成ペイロードを表す。
type XTranslatorTranslationXMLDraft struct {
	FilePath         string
	TargetPluginName string
	TargetPluginType string
	TermCount        int
	ImportedAt       time.Time
}

// Persona は PERSONA テーブルの 1 レコードを表す。
type Persona struct {
	ID                     int64
	NpcProfileID           int64
	TranslationJobID       *int64
	PersonaLifecycle       string
	PersonaScope           string
	PersonaSource          string
	PersonaDescription     string
	SpeechStyle            string
	PersonalitySummary     string
	EvidenceUtteranceCount int
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

// PersonaDraft は PERSONA の作成ペイロードを表す。
type PersonaDraft struct {
	NpcProfileID           int64
	TranslationJobID       *int64
	PersonaLifecycle       string
	PersonaScope           string
	PersonaSource          string
	PersonaDescription     string
	SpeechStyle            string
	PersonalitySummary     string
	EvidenceUtteranceCount int
}

// PersonaUpdateDraft は PERSONA の更新ペイロードを表す。
type PersonaUpdateDraft struct {
	PersonaLifecycle       string
	PersonaScope           string
	PersonaSource          string
	PersonaDescription     string
	SpeechStyle            string
	PersonalitySummary     string
	EvidenceUtteranceCount int
}

// PersonaFieldEvidence は PERSONA_FIELD_EVIDENCE テーブルの 1 レコードを表す。
type PersonaFieldEvidence struct {
	ID                 int64
	PersonaID          int64
	TranslationFieldID int64
	EvidenceRole       string
}

// PersonaFieldEvidenceDraft は PERSONA_FIELD_EVIDENCE の作成ペイロードを表す。
type PersonaFieldEvidenceDraft struct {
	PersonaID          int64
	TranslationFieldID int64
	EvidenceRole       string
}

// DictionaryEntry は DICTIONARY_ENTRY テーブルの 1 レコードを表す。
type DictionaryEntry struct {
	ID                          int64
	XTranslatorTranslationXMLID *int64
	TranslationJobID            *int64
	DictionaryLifecycle         string
	DictionaryScope             string
	DictionarySource            string
	SourceTerm                  string
	TranslatedTerm              string
	TermKind                    string
	Reusable                    bool
	CreatedAt                   time.Time
	UpdatedAt                   time.Time
}

// DictionaryEntryDraft は DICTIONARY_ENTRY の作成ペイロードを表す。
type DictionaryEntryDraft struct {
	XTranslatorTranslationXMLID *int64
	TranslationJobID            *int64
	DictionaryLifecycle         string
	DictionaryScope             string
	DictionarySource            string
	SourceTerm                  string
	TranslatedTerm              string
	TermKind                    string
	Reusable                    bool
}

// DictionaryEntryUpdateDraft は DICTIONARY_ENTRY の更新ペイロードを表す。
type DictionaryEntryUpdateDraft struct {
	DictionaryLifecycle string
	DictionaryScope     string
	DictionarySource    string
	SourceTerm          string
	TranslatedTerm      string
	TermKind            string
	Reusable            bool
}

// FoundationDataRepository はペルソナ・辞書・翻訳 XML の永続化操作を定義する。
// 扱うテーブル: PERSONA, PERSONA_FIELD_EVIDENCE, DICTIONARY_ENTRY, XTRANSLATOR_TRANSLATION_XML.
type FoundationDataRepository interface {
	// XTranslatorTranslationXML
	CreateXTranslatorTranslationXML(ctx context.Context, draft XTranslatorTranslationXMLDraft) (XTranslatorTranslationXML, error)
	GetXTranslatorTranslationXMLByID(ctx context.Context, id int64) (XTranslatorTranslationXML, error)

	// Persona
	CreatePersona(ctx context.Context, draft PersonaDraft) (Persona, error)
	GetPersonaByID(ctx context.Context, id int64) (Persona, error)
	GetPersonaByNpcProfileID(ctx context.Context, npcProfileID int64) (Persona, error)
	UpdatePersona(ctx context.Context, id int64, draft PersonaUpdateDraft) (Persona, error)

	// PersonaFieldEvidence
	CreatePersonaFieldEvidence(ctx context.Context, draft PersonaFieldEvidenceDraft) (PersonaFieldEvidence, error)
	ListPersonaFieldEvidenceByPersonaID(ctx context.Context, personaID int64) ([]PersonaFieldEvidence, error)

	// DictionaryEntry
	CreateDictionaryEntry(ctx context.Context, draft DictionaryEntryDraft) (DictionaryEntry, error)
	GetDictionaryEntryByID(ctx context.Context, id int64) (DictionaryEntry, error)
	UpdateDictionaryEntry(ctx context.Context, id int64, draft DictionaryEntryUpdateDraft) (DictionaryEntry, error)
	DeleteDictionaryEntry(ctx context.Context, id int64) error
}
