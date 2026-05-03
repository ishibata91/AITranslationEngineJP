package usecase

import "context"

// TranslationOutputArtifactUsecase freezes the public output artifact seam.
type TranslationOutputArtifactUsecase interface {
	GetTranslationOutputReview(context.Context, GetTranslationOutputReviewRequest) (TranslationOutputReviewResult, error)
	GetTranslationOutputDiffPreview(context.Context, GetTranslationOutputDiffPreviewRequest) (TranslationOutputDiffPreviewResult, error)
	GenerateXTranslatorOutputArtifact(context.Context, GenerateXTranslatorOutputArtifactRequest) (XTranslatorOutputArtifactCommandResult, error)
	RegenerateXTranslatorOutputArtifact(context.Context, RegenerateXTranslatorOutputArtifactRequest) (XTranslatorOutputArtifactCommandResult, error)
}

// TranslationOutputArtifactErrorKind identifies public rejected outcomes.
type TranslationOutputArtifactErrorKind = string

const (
	// TranslationOutputArtifactErrorKindNotCompleted identifies an incomplete job rejection.
	TranslationOutputArtifactErrorKindNotCompleted TranslationOutputArtifactErrorKind = "not_completed"
	// TranslationOutputArtifactErrorKindCanceled identifies a canceled job rejection.
	TranslationOutputArtifactErrorKindCanceled TranslationOutputArtifactErrorKind = "canceled"
	// TranslationOutputArtifactErrorKindStatusMismatch identifies row or job status mismatch.
	TranslationOutputArtifactErrorKindStatusMismatch TranslationOutputArtifactErrorKind = "status_mismatch"
	// TranslationOutputArtifactErrorKindMissingRequiredRowField identifies missing XML row columns.
	TranslationOutputArtifactErrorKindMissingRequiredRowField TranslationOutputArtifactErrorKind = "missing_required_row_field"
	// TranslationOutputArtifactErrorKindUnknownOutputStatus identifies an unmapped output status.
	TranslationOutputArtifactErrorKindUnknownOutputStatus TranslationOutputArtifactErrorKind = "unknown_output_status"
	// TranslationOutputArtifactErrorKindXMLSerializationFailed identifies XML serialization failure.
	TranslationOutputArtifactErrorKindXMLSerializationFailed TranslationOutputArtifactErrorKind = "xml_serialization_failed"
	// TranslationOutputArtifactErrorKindFileWriteFailed identifies filesystem write failure.
	TranslationOutputArtifactErrorKindFileWriteFailed TranslationOutputArtifactErrorKind = "file_write_failed"
	// TranslationOutputArtifactErrorKindArtifactSaveFailed identifies artifact persistence failure.
	TranslationOutputArtifactErrorKindArtifactSaveFailed TranslationOutputArtifactErrorKind = "artifact_save_failed"
	// TranslationOutputArtifactErrorKindCompatibilityRejected identifies compatibility validation rejection.
	TranslationOutputArtifactErrorKindCompatibilityRejected TranslationOutputArtifactErrorKind = "compatibility_rejected"
	// TranslationOutputArtifactErrorKindSecretRedacted identifies a redacted secret-related failure.
	TranslationOutputArtifactErrorKindSecretRedacted TranslationOutputArtifactErrorKind = "secret_redacted"
)

// TranslationOutputArtifactRedactionFieldObligation freezes the public redaction boundary.
type TranslationOutputArtifactRedactionFieldObligation struct {
	AllowedReferenceFields []string
	ForbiddenOutputFields  []string
}

// TranslationOutputArtifactPublicRedactionFieldObligation returns the frozen public redaction boundary.
func TranslationOutputArtifactPublicRedactionFieldObligation() TranslationOutputArtifactRedactionFieldObligation {
	return TranslationOutputArtifactRedactionFieldObligation{
		AllowedReferenceFields: []string{
			"job_id",
			"artifact_id",
			"field_id",
			"row_digest",
			"count",
			"status",
			"file_path",
			"target_game",
			"error_kind",
			"retryable",
		},
		ForbiddenOutputFields: []string{
			"secret",
			"api_key",
			"token",
			"authorization_header",
			"provider_raw_request",
			"provider_raw_response",
			"decryptable_value",
			"full_source_text",
			"full_dest_text",
		},
	}
}

// GetTranslationOutputReviewRequest identifies the review selection target.
type GetTranslationOutputReviewRequest struct {
	SelectedJobID *int64
}

// GetTranslationOutputDiffPreviewRequest identifies the preview target pair.
type GetTranslationOutputDiffPreviewRequest struct {
	JobID      int64
	ArtifactID int64
}

// GenerateXTranslatorOutputArtifactRequest identifies one output generation request.
type GenerateXTranslatorOutputArtifactRequest struct {
	JobID      int64
	TargetGame string
	OutputPath string
}

// RegenerateXTranslatorOutputArtifactRequest identifies one output regeneration request.
type RegenerateXTranslatorOutputArtifactRequest struct {
	JobID      int64
	ArtifactID int64
	TargetGame string
	OutputPath string
}

// TranslationOutputReviewResult returns the frozen Output Review summary.
type TranslationOutputReviewResult struct {
	CompletedJobs    []TranslationOutputCompletedJobSummary
	SelectedJob      TranslationOutputSelectedJobSummary
	OutputReadiness  TranslationOutputReadinessSummary
	ArtifactStatus   TranslationOutputArtifactStatusSummary
	RejectionReasons []TranslationOutputArtifactErrorSummary
}

// TranslationOutputCompletedJobSummary summarizes one completed job candidate.
type TranslationOutputCompletedJobSummary struct {
	JobID                    int64
	JobStatus                string
	ArtifactStatus           string
	OutputReady              bool
	TranslatedCount          int
	OutputStatusDistribution map[string]int
}

// TranslationOutputSelectedJobSummary summarizes one selected job.
type TranslationOutputSelectedJobSummary struct {
	JobID           int64
	JobStatus       string
	BodyPhaseStatus string
	OutputReady     bool
	ResultSummary   TranslationOutputResultSummary
}

// TranslationOutputResultSummary summarizes output counts and provenance.
type TranslationOutputResultSummary struct {
	TranslatedCount int
	RowCount        int
	InputProvenance TranslationOutputInputProvenanceSummary
}

// TranslationOutputInputProvenanceSummary summarizes frozen provenance digests.
type TranslationOutputInputProvenanceSummary struct {
	InputSnapshotDigest string
	SourceFileDigest    string
}

// TranslationOutputReadinessSummary summarizes output readiness and rejection kind.
type TranslationOutputReadinessSummary struct {
	Ready         bool
	Retryable     bool
	RejectionKind TranslationOutputArtifactErrorKind
}

// TranslationOutputArtifactStatusSummary summarizes current artifact state.
type TranslationOutputArtifactStatusSummary struct {
	ArtifactID     int64
	Status         string
	RowCount       int
	CurrentVersion bool
}

// TranslationOutputDiffPreviewResult returns the frozen diff preview contract.
type TranslationOutputDiffPreviewResult struct {
	JobID                int64
	ArtifactID           int64
	Rows                 []TranslationOutputDiffPreviewRow
	CompatibilitySummary TranslationOutputCompatibilitySummary
}

// TranslationOutputDiffPreviewRow summarizes one xTranslator-compatible row.
type TranslationOutputDiffPreviewRow struct {
	FieldID              int64
	RowDigest            string
	EDID                 string
	REC                  string
	FIELD                string
	FORMID               string
	SourceExcerpt        string
	DestExcerpt          string
	XTranslatorStatus    int
	InternalOutputStatus string
	RowReflectionSummary string
	StaleReason          string
	CanRegenerate        bool
}

// TranslationOutputCompatibilitySummary summarizes compatibility gate counts.
type TranslationOutputCompatibilitySummary struct {
	Passed       bool
	WarningCount int
	RejectCount  int
}

// XTranslatorOutputArtifactCommandResult returns the frozen write command response.
type XTranslatorOutputArtifactCommandResult struct {
	JobID            int64
	ArtifactID       int64
	ArtifactStatus   string
	RowCount         int
	FilePath         string
	TargetGame       string
	ErrorSummary     *TranslationOutputArtifactErrorSummary
	OperationSummary TranslationOutputOperationSummary
}

// TranslationOutputArtifactErrorSummary exposes only redacted error information.
type TranslationOutputArtifactErrorSummary struct {
	ErrorKind  TranslationOutputArtifactErrorKind
	Reason     string
	Retryable  bool
	IsRedacted bool
}

// TranslationOutputOperationSummary summarizes generate or regenerate effects.
type TranslationOutputOperationSummary struct {
	OperationKind       string
	ReplacedArtifactID  int64
	AffectedFieldIDs    []int64
	DuplicateRowCreated bool
}
