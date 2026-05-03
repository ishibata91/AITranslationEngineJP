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

// TranslationArtifact は TRANSLATION_ARTIFACT テーブルの 1 レコードを表す。
type TranslationArtifact struct {
	ID               int64
	TranslationJobID int64
	ArtifactFormat   string
	TargetGame       string
	FilePath         string
	Status           string
	GeneratedAt      *time.Time
}

// TranslationArtifactDraft は TRANSLATION_ARTIFACT の upsert ペイロードを表す。
type TranslationArtifactDraft struct {
	TranslationJobID int64
	ArtifactFormat   string
	TargetGame       string
	FilePath         string
	Status           string
	GeneratedAt      *time.Time
}

// XTranslatorOutputRow は XTRANSLATOR_OUTPUT_ROW テーブルの 1 レコードを表す。
type XTranslatorOutputRow struct {
	ID                    int64
	TranslationArtifactID int64
	JobTranslationFieldID int64
	EDID                  string
	REC                   string
	FIELD                 string
	FORMID                string
	Source                string
	Dest                  string
	Status                int
}

// XTranslatorOutputRowDraft は XTRANSLATOR_OUTPUT_ROW の作成ペイロードを表す。
type XTranslatorOutputRowDraft struct {
	JobTranslationFieldID int64
	EDID                  string
	REC                   string
	FIELD                 string
	FORMID                string
	Source                string
	Dest                  string
	Status                int
}

// JobOutputRepository はジョブ翻訳フィールドの永続化操作を定義する。
// 扱うテーブル: JOB_TRANSLATION_FIELD.
type JobOutputRepository interface {
	CreateJobTranslationField(ctx context.Context, draft JobTranslationFieldDraft) (JobTranslationField, error)
	GetJobTranslationFieldByID(ctx context.Context, id int64) (JobTranslationField, error)
	UpdateJobTranslationField(ctx context.Context, id int64, draft JobTranslationFieldUpdateDraft) (JobTranslationField, error)
	ListJobTranslationFieldsByJobID(ctx context.Context, jobID int64) ([]JobTranslationField, error)
}

// JobOutputArtifactRepository は output artifact と row の永続化境界を定義する。
// 扱うテーブル: TRANSLATION_ARTIFACT, XTRANSLATOR_OUTPUT_ROW.
type JobOutputArtifactRepository interface {
	GetTranslationArtifactByID(ctx context.Context, id int64) (TranslationArtifact, error)
	GetTranslationArtifactByJobID(ctx context.Context, jobID int64) (TranslationArtifact, error)
	ListXTranslatorOutputRowsByArtifactID(ctx context.Context, translationArtifactID int64) ([]XTranslatorOutputRow, error)
	UpsertTranslationArtifact(ctx context.Context, draft TranslationArtifactDraft) (TranslationArtifact, error)
	ReplaceXTranslatorOutputRows(
		ctx context.Context,
		translationArtifactID int64,
		drafts []XTranslatorOutputRowDraft,
	) ([]XTranslatorOutputRow, error)
	CountXTranslatorOutputRowsByArtifactID(ctx context.Context, translationArtifactID int64) (int, error)
}
