package repository

import (
	"context"
	"time"
)

// JobTranslationField は JOB_TRANSLATION_FIELD テーブルの 1 レコードを表す。
type JobTranslationField struct {
	ID                 int64
	TranslationJobID   int64
	TranslationFieldID int64
	AppliedPersonaID   *int64
	TranslatedText     string
	OutputStatus       string
	RetryCount         int
	UpdatedAt          time.Time
}

// JobTranslationFieldDraft は JOB_TRANSLATION_FIELD の作成ペイロードを表す。
type JobTranslationFieldDraft struct {
	TranslationJobID   int64
	TranslationFieldID int64
	AppliedPersonaID   *int64
	TranslatedText     string
	OutputStatus       string
	RetryCount         int
}

// JobTranslationFieldUpdateDraft は JOB_TRANSLATION_FIELD の更新ペイロードを表す。
type JobTranslationFieldUpdateDraft struct {
	AppliedPersonaID *int64
	TranslatedText   string
	OutputStatus     string
	RetryCount       int
}

// JobOutputRepository はジョブ翻訳フィールドの永続化操作を定義する。
// 扱うテーブル: JOB_TRANSLATION_FIELD.
type JobOutputRepository interface {
	CreateJobTranslationField(ctx context.Context, draft JobTranslationFieldDraft) (JobTranslationField, error)
	GetJobTranslationFieldByID(ctx context.Context, id int64) (JobTranslationField, error)
	UpdateJobTranslationField(ctx context.Context, id int64, draft JobTranslationFieldUpdateDraft) (JobTranslationField, error)
	ListJobTranslationFieldsByJobID(ctx context.Context, jobID int64) ([]JobTranslationField, error)
}
