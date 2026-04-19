package repository

import (
	"context"
	"time"
)

// TranslationJob は TRANSLATION_JOB テーブルの 1 レコードを表す。
type TranslationJob struct {
	ID                   int64
	XEditExtractedDataID int64
	JobName              string
	State                string
	ProgressPercent      int
	CreatedAt            time.Time
	StartedAt            *time.Time
	FinishedAt           *time.Time
}

// TranslationJobDraft は TRANSLATION_JOB の作成ペイロードを表す。
type TranslationJobDraft struct {
	XEditExtractedDataID int64
	JobName              string
	State                string
	ProgressPercent      int
}

// TranslationJobUpdateDraft は TRANSLATION_JOB の更新ペイロードを表す。
type TranslationJobUpdateDraft struct {
	JobName         string
	State           string
	ProgressPercent int
	StartedAt       *time.Time
	FinishedAt      *time.Time
}

// JobPhaseRun は JOB_PHASE_RUN テーブルの 1 レコードを表す。
type JobPhaseRun struct {
	ID                  int64
	TranslationJobID    int64
	PhaseType           string
	State               string
	ExecutionOrder      int
	ProgressPercent     int
	AIProvider          string
	ModelName           string
	ExecutionMode       string
	CredentialRef       string
	InstructionKind     string
	LatestExternalRunID string
	LatestError         string
	StartedAt           *time.Time
	FinishedAt          *time.Time
}

// JobPhaseRunDraft は JOB_PHASE_RUN の作成ペイロードを表す。
type JobPhaseRunDraft struct {
	TranslationJobID int64
	PhaseType        string
	State            string
	ExecutionOrder   int
	AIProvider       string
	ModelName        string
	ExecutionMode    string
	CredentialRef    string
	InstructionKind  string
}

// JobPhaseRunUpdateDraft は JOB_PHASE_RUN の更新ペイロードを表す。
type JobPhaseRunUpdateDraft struct {
	State               string
	ProgressPercent     int
	LatestExternalRunID string
	LatestError         string
	StartedAt           *time.Time
	FinishedAt          *time.Time
}

// PhaseRunTranslationField は PHASE_RUN_TRANSLATION_FIELD テーブルの 1 レコードを表す。
type PhaseRunTranslationField struct {
	ID                    int64
	PhaseRunID            int64
	JobTranslationFieldID int64
	Role                  string
}

// PhaseRunTranslationFieldDraft は PHASE_RUN_TRANSLATION_FIELD の作成ペイロードを表す。
type PhaseRunTranslationFieldDraft struct {
	PhaseRunID            int64
	JobTranslationFieldID int64
	Role                  string
}

// PhaseRunPersona は PHASE_RUN_PERSONA テーブルの 1 レコードを表す。
type PhaseRunPersona struct {
	ID         int64
	PhaseRunID int64
	PersonaID  int64
	Role       string
}

// PhaseRunPersonaDraft は PHASE_RUN_PERSONA の作成ペイロードを表す。
type PhaseRunPersonaDraft struct {
	PhaseRunID int64
	PersonaID  int64
	Role       string
}

// PhaseRunDictionaryEntry は PHASE_RUN_DICTIONARY_ENTRY テーブルの 1 レコードを表す。
type PhaseRunDictionaryEntry struct {
	ID                int64
	PhaseRunID        int64
	DictionaryEntryID int64
	Role              string
}

// PhaseRunDictionaryEntryDraft は PHASE_RUN_DICTIONARY_ENTRY の作成ペイロードを表す。
type PhaseRunDictionaryEntryDraft struct {
	PhaseRunID        int64
	DictionaryEntryID int64
	Role              string
}

// JobLifecycleRepository は翻訳ジョブとフェーズ実行の永続化操作を定義する。
// 扱うテーブル: TRANSLATION_JOB, JOB_PHASE_RUN, PHASE_RUN_TRANSLATION_FIELD,
// PHASE_RUN_PERSONA, PHASE_RUN_DICTIONARY_ENTRY.
type JobLifecycleRepository interface {
	// TranslationJob
	CreateTranslationJob(ctx context.Context, draft TranslationJobDraft) (TranslationJob, error)
	GetTranslationJobByID(ctx context.Context, id int64) (TranslationJob, error)
	UpdateTranslationJob(ctx context.Context, id int64, draft TranslationJobUpdateDraft) (TranslationJob, error)

	// JobPhaseRun
	CreateJobPhaseRun(ctx context.Context, draft JobPhaseRunDraft) (JobPhaseRun, error)
	GetJobPhaseRunByID(ctx context.Context, id int64) (JobPhaseRun, error)
	UpdateJobPhaseRun(ctx context.Context, id int64, draft JobPhaseRunUpdateDraft) (JobPhaseRun, error)
	ListJobPhaseRunsByJobID(ctx context.Context, jobID int64) ([]JobPhaseRun, error)

	// PhaseRunTranslationField
	CreatePhaseRunTranslationField(ctx context.Context, draft PhaseRunTranslationFieldDraft) (PhaseRunTranslationField, error)

	// PhaseRunPersona
	CreatePhaseRunPersona(ctx context.Context, draft PhaseRunPersonaDraft) (PhaseRunPersona, error)

	// PhaseRunDictionaryEntry
	CreatePhaseRunDictionaryEntry(ctx context.Context, draft PhaseRunDictionaryEntryDraft) (PhaseRunDictionaryEntry, error)
}
