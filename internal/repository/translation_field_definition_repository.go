package repository

import "context"

// TranslationFieldDefinition は TRANSLATION_FIELD_DEFINITION テーブルの 1 レコードを表す。
type TranslationFieldDefinition struct {
	ID                   int64
	RecordType           string
	SubrecordType        string
	DisplayName          string
	AIDescription        string
	Translatable         bool
	Ordered              bool
	OrderScope           string
	ReferenceRequirement string
}

// TranslationFieldDefinitionDraft は TRANSLATION_FIELD_DEFINITION の作成ペイロードを表す。
type TranslationFieldDefinitionDraft struct {
	RecordType           string
	SubrecordType        string
	DisplayName          string
	AIDescription        string
	Translatable         bool
	Ordered              bool
	OrderScope           string
	ReferenceRequirement string
}

// TranslationFieldDefinitionUpdateDraft は TRANSLATION_FIELD_DEFINITION の更新ペイロードを表す。
type TranslationFieldDefinitionUpdateDraft struct {
	DisplayName          string
	AIDescription        string
	Translatable         bool
	Ordered              bool
	OrderScope           string
	ReferenceRequirement string
}

// TranslationFieldDefinitionRepository はフィールド定義メタデータの永続化操作を定義する。
// 扱うテーブル: TRANSLATION_FIELD_DEFINITION.
type TranslationFieldDefinitionRepository interface {
	Create(ctx context.Context, draft TranslationFieldDefinitionDraft) (TranslationFieldDefinition, error)
	GetByID(ctx context.Context, id int64) (TranslationFieldDefinition, error)
	GetByRecordTypeAndSubrecordType(ctx context.Context, recordType, subrecordType string) (TranslationFieldDefinition, error)
	List(ctx context.Context) ([]TranslationFieldDefinition, error)
	Update(ctx context.Context, id int64, draft TranslationFieldDefinitionUpdateDraft) (TranslationFieldDefinition, error)
}
