package usecase

// TranslationJobManagementReasonCategory identifies user-visible reasons and error kinds.
type TranslationJobManagementReasonCategory string

// TranslationJobManagementReasonCategory values define the frozen reason strings.
const (
	TranslationJobManagementReasonCategoryCacheMissing                   TranslationJobManagementReasonCategory = "cache_missing"
	TranslationJobManagementReasonCategoryTerminalState                  TranslationJobManagementReasonCategory = "terminal_state"
	TranslationJobManagementReasonCategoryStateProjectionInconsistent    TranslationJobManagementReasonCategory = "state_projection_inconsistent"
	TranslationJobManagementReasonCategoryPhaseProgressAggregationFailed TranslationJobManagementReasonCategory = "phase_progress_aggregation_failed"
	TranslationJobManagementReasonCategoryStaleSelection                 TranslationJobManagementReasonCategory = "stale_selection"
	TranslationJobManagementReasonCategoryListLoadFailure                TranslationJobManagementReasonCategory = "list_load_failure"
	TranslationJobManagementReasonCategoryRunningDeleteBlocked           TranslationJobManagementReasonCategory = "running_delete_blocked"
	TranslationJobManagementReasonCategoryStopRequested                  TranslationJobManagementReasonCategory = "stop_requested"
	TranslationJobManagementReasonCategoryStopFailed                     TranslationJobManagementReasonCategory = "stop_failed"
	TranslationJobManagementReasonCategoryDeleteFailed                   TranslationJobManagementReasonCategory = "delete_failed"
	TranslationJobManagementReasonCategoryResumeFailed                   TranslationJobManagementReasonCategory = "resume_failed"
)

// TranslationJobManagementInputSourceSummary summarizes one persisted input source.
type TranslationJobManagementInputSourceSummary struct {
	InputSourceID        int64
	InputSourceLabel     string
	InputSourceKindLabel string
	SourcePath           string
	PluginName           string
	ExtractedJSONLabel   string
}

// TranslationJobManagementProgressSummary summarizes the current job progress.
type TranslationJobManagementProgressSummary struct {
	CurrentPhase      string
	CurrentPhaseLabel string
	Percent           int
	ProgressLabel     string
	LastUpdatedLabel  string
}

// TranslationJobManagementProtectedSettingSummary exposes redacted runtime settings.
type TranslationJobManagementProtectedSettingSummary struct {
	ProviderLabel        string
	ModelLabel           string
	ExecutionModeLabel   string
	CredentialState      string
	CredentialStateLabel string
}

// TranslationJobManagementOperationAvailability describes one operation guard.
type TranslationJobManagementOperationAvailability struct {
	Kind           string
	Enabled        bool
	Label          string
	HelperText     string
	ReasonCategory string
	ReasonText     string
}

// TranslationJobManagementBlockedReason explains one blocked or warning state.
type TranslationJobManagementBlockedReason struct {
	Category string
	Title    string
	Detail   string
}

// TranslationJobManagementJobSummary describes one list item.
type TranslationJobManagementJobSummary struct {
	JobID              int64
	JobState           string
	JobStateLabel      string
	StateTone          string
	CanOpenPhase       bool
	OpenBlockedReason  *TranslationJobManagementBlockedReason
	InputSource        TranslationJobManagementInputSourceSummary
	Progress           TranslationJobManagementProgressSummary
	StopAvailability   TranslationJobManagementOperationAvailability
	ResumeAvailability TranslationJobManagementOperationAvailability
	DeleteAvailability TranslationJobManagementOperationAvailability
}

// TranslationJobManagementJobDetail expands one selected job.
type TranslationJobManagementJobDetail struct {
	TranslationJobManagementJobSummary
	CacheState           string
	CacheStateLabel      string
	RuntimeSummary       TranslationJobManagementProtectedSettingSummary
	ResumeBlockedReasons []TranslationJobManagementBlockedReason
	Warnings             []TranslationJobManagementBlockedReason
	DeleteImpactLines    []string
}

// TranslationJobManagementListResult returns the incomplete-job list.
type TranslationJobManagementListResult struct {
	Jobs []TranslationJobManagementJobSummary
}

// TranslationJobManagementGetDetailRequest identifies one detail target.
type TranslationJobManagementGetDetailRequest struct {
	JobID int64
}

// TranslationJobManagementDeleteRequest identifies one delete target.
type TranslationJobManagementDeleteRequest struct {
	JobID int64
}

// TranslationJobManagementActionRequest identifies one stop or resume target.
type TranslationJobManagementActionRequest struct {
	JobID int64
}

// TranslationJobManagementActionResult returns the public action outcome.
type TranslationJobManagementActionResult struct {
	Message        string
	Tone           string
	Detail         *TranslationJobManagementJobDetail
	DeletedJobID   *int64
	ReasonCategory string
}
