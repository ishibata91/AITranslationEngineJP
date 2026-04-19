package repository

import (
	"context"
	"time"
)

// XEditExtractedData は X_EDIT_EXTRACTED_DATA テーブルの 1 レコードを表す。
type XEditExtractedData struct {
	ID               int64
	SourceFilePath   string
	SourceTool       string
	TargetPluginName string
	TargetPluginType string
	RecordCount      int
	ImportedAt       time.Time
}

// XEditExtractedDataDraft は X_EDIT_EXTRACTED_DATA の作成ペイロードを表す。
type XEditExtractedDataDraft struct {
	SourceFilePath   string
	SourceTool       string
	TargetPluginName string
	TargetPluginType string
	RecordCount      int
	ImportedAt       time.Time
}

// TranslationRecord は TRANSLATION_RECORD テーブルの 1 レコードを表す。
type TranslationRecord struct {
	ID                   int64
	XEditExtractedDataID int64
	FormID               string
	EditorID             string
	RecordType           string
}

// TranslationRecordDraft は TRANSLATION_RECORD の作成ペイロードを表す。
type TranslationRecordDraft struct {
	XEditExtractedDataID int64
	FormID               string
	EditorID             string
	RecordType           string
}

// NpcProfile は NPC_PROFILE テーブルの 1 レコードを表す。
type NpcProfile struct {
	ID               int64
	TargetPluginName string
	FormID           string
	RecordType       string
	EditorID         string
	DisplayName      string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// NpcProfileDraft は NPC_PROFILE の作成・Upsert ペイロードを表す。
type NpcProfileDraft struct {
	TargetPluginName string
	FormID           string
	RecordType       string
	EditorID         string
	DisplayName      string
}

// NpcRecord は NPC_RECORD テーブルの 1 レコードを表す。
type NpcRecord struct {
	TranslationRecordID int64
	NpcProfileID        int64
	Race                *string
	Sex                 *string
	NpcClass            *string
	VoiceType           string
}

// NpcRecordDraft は NPC_RECORD の作成ペイロードを表す。
type NpcRecordDraft struct {
	TranslationRecordID int64
	NpcProfileID        int64
	Race                *string
	Sex                 *string
	NpcClass            *string
	VoiceType           string
}

// TranslationField は TRANSLATION_FIELD テーブルの 1 レコードを表す。
type TranslationField struct {
	ID                           int64
	TranslationRecordID          int64
	TranslationFieldDefinitionID *int64
	SubrecordType                string
	SourceText                   string
	FieldOrder                   int
	PreviousTranslationFieldID   *int64
	NextTranslationFieldID       *int64
}

// TranslationFieldDraft は TRANSLATION_FIELD の作成ペイロードを表す。
type TranslationFieldDraft struct {
	TranslationRecordID          int64
	TranslationFieldDefinitionID *int64
	SubrecordType                string
	SourceText                   string
	FieldOrder                   int
	PreviousTranslationFieldID   *int64
	NextTranslationFieldID       *int64
}

// TranslationFieldRecordReference は TRANSLATION_FIELD_RECORD_REFERENCE テーブルの 1 レコードを表す。
type TranslationFieldRecordReference struct {
	ID                            int64
	TranslationFieldID            int64
	ReferencedTranslationRecordID int64
	ReferenceRole                 string
}

// TranslationFieldRecordReferenceDraft は TRANSLATION_FIELD_RECORD_REFERENCE の作成ペイロードを表す。
type TranslationFieldRecordReferenceDraft struct {
	TranslationFieldID            int64
	ReferencedTranslationRecordID int64
	ReferenceRole                 string
}

// TranslationSourceRepository は入力データ群の永続化操作を定義する。
// 扱うテーブル: X_EDIT_EXTRACTED_DATA, TRANSLATION_RECORD, NPC_PROFILE, NPC_RECORD,
// TRANSLATION_FIELD, TRANSLATION_FIELD_RECORD_REFERENCE.
type TranslationSourceRepository interface {
	// XEditExtractedData
	CreateXEditExtractedData(ctx context.Context, draft XEditExtractedDataDraft) (XEditExtractedData, error)
	GetXEditExtractedDataByID(ctx context.Context, id int64) (XEditExtractedData, error)

	// TranslationRecord
	CreateTranslationRecord(ctx context.Context, draft TranslationRecordDraft) (TranslationRecord, error)
	GetTranslationRecordByID(ctx context.Context, id int64) (TranslationRecord, error)
	ListTranslationRecordsByXEditID(ctx context.Context, xEditID int64) ([]TranslationRecord, error)

	// NpcProfile — target_plugin_name + form_id + record_type を Unique キーとして Upsert する。
	UpsertNpcProfile(ctx context.Context, draft NpcProfileDraft) (NpcProfile, error)
	GetNpcProfileByID(ctx context.Context, id int64) (NpcProfile, error)

	// NpcRecord
	CreateNpcRecord(ctx context.Context, draft NpcRecordDraft) (NpcRecord, error)
	GetNpcRecordByTranslationRecordID(ctx context.Context, translationRecordID int64) (NpcRecord, error)

	// TranslationField
	CreateTranslationField(ctx context.Context, draft TranslationFieldDraft) (TranslationField, error)
	GetTranslationFieldByID(ctx context.Context, id int64) (TranslationField, error)
	ListTranslationFieldsByTranslationRecordID(ctx context.Context, translationRecordID int64) ([]TranslationField, error)

	// TranslationFieldRecordReference
	CreateTranslationFieldRecordReference(ctx context.Context, draft TranslationFieldRecordReferenceDraft) (TranslationFieldRecordReference, error)
	ListTranslationFieldRecordReferencesByFieldID(ctx context.Context, fieldID int64) ([]TranslationFieldRecordReference, error)
}
