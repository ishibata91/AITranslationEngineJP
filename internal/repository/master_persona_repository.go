package repository

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

var (
	// ErrMasterPersonaEntryNotFound reports that the requested master persona entry does not exist.
	ErrMasterPersonaEntryNotFound = errors.New("master persona entry not found")
)

const (
	masterPersonaDefaultPageSize      = 30
	masterPersonaMaxPageSize          = 100
	masterPersonaIdentityKeyErrorFmt  = "%w: identity_key=%s"
	masterPersonaSeedFollowersPlugin  = "FollowersPlus.esp"
	masterPersonaSeedNightCourtPlugin = "NightCourt.esp"
)

// MasterPersonaEntry stores one master persona record.
type MasterPersonaEntry struct {
	IdentityKey          string
	TargetPlugin         string
	FormID               string
	RecordType           string
	EditorID             string
	DisplayName          string
	Race                 *string
	Sex                  *string
	VoiceType            string
	ClassName            string
	SourcePlugin         string
	PersonaSummary       string
	SpeechStyle          string
	PersonaBody          string
	GenerationSourceJSON string
	BaselineApplied      bool
	DialogueCount        int
	Dialogues            []string
	UpdatedAt            time.Time
}

// MasterPersonaListQuery describes list conditions for master persona entries.
type MasterPersonaListQuery struct {
	Keyword      string
	PluginFilter string
	Page         int
	PageSize     int
}

// MasterPersonaPluginGroup stores one plugin group summary for filters.
type MasterPersonaPluginGroup struct {
	TargetPlugin string
	Count        int
}

// MasterPersonaListResult stores list items, paging, and plugin groups.
type MasterPersonaListResult struct {
	Items        []MasterPersonaEntry
	TotalCount   int
	Page         int
	PageSize     int
	PluginGroups []MasterPersonaPluginGroup
}

// MasterPersonaDraft stores the persistence payload for create-like upsert and update.
type MasterPersonaDraft struct {
	IdentityKey          string
	TargetPlugin         string
	FormID               string
	RecordType           string
	EditorID             string
	DisplayName          string
	Race                 *string
	Sex                  *string
	VoiceType            string
	ClassName            string
	SourcePlugin         string
	PersonaSummary       string
	SpeechStyle          string
	PersonaBody          string
	GenerationSourceJSON string
	BaselineApplied      bool
	Dialogues            []string
	UpdatedAt            time.Time
}

// MasterPersonaAISettingsRecord stores persisted page-local AI settings except secrets.
type MasterPersonaAISettingsRecord struct {
	Provider string
	Model    string
}

// MasterPersonaRunStatusRecord stores persisted generation run status.
type MasterPersonaRunStatusRecord struct {
	RunState              string
	TargetPlugin          string
	ProcessedCount        int
	SuccessCount          int
	ExistingSkipCount     int
	ZeroDialogueSkipCount int
	GenericNPCCount       int
	CurrentActorLabel     string
	Message               string
	StartedAt             *time.Time
	FinishedAt            *time.Time
}

// MasterPersonaQueryRepository defines read-only master persona persistence.
type MasterPersonaQueryRepository interface {
	List(ctx context.Context, query MasterPersonaListQuery) (MasterPersonaListResult, error)
	GetByIdentityKey(ctx context.Context, identityKey string) (MasterPersonaEntry, error)
}

// MasterPersonaCommandRepository defines mutating master persona persistence.
type MasterPersonaCommandRepository interface {
	GetByIdentityKey(ctx context.Context, identityKey string) (MasterPersonaEntry, error)
	UpsertIfAbsent(ctx context.Context, draft MasterPersonaDraft) (MasterPersonaEntry, bool, error)
	Update(ctx context.Context, identityKey string, draft MasterPersonaDraft) (MasterPersonaEntry, error)
	Delete(ctx context.Context, identityKey string) error
}

// MasterPersonaAISettingsRepository defines page-local AI settings persistence.
type MasterPersonaAISettingsRepository interface {
	LoadAISettings(ctx context.Context) (MasterPersonaAISettingsRecord, error)
	SaveAISettings(ctx context.Context, record MasterPersonaAISettingsRecord) error
}

// MasterPersonaRunRepository defines generation run status persistence.
type MasterPersonaRunRepository interface {
	LoadRunStatus(ctx context.Context) (MasterPersonaRunStatusRecord, error)
	SaveRunStatus(ctx context.Context, status MasterPersonaRunStatusRecord) error
}

// SecretStore defines secret access for page-local AI settings.
type SecretStore interface {
	Load(ctx context.Context, key string) (string, error)
	Save(ctx context.Context, key string, value string) error
	Delete(ctx context.Context, key string) error
}

// InMemoryMasterPersonaRepository provides an in-memory backend seam for master persona data.
type InMemoryMasterPersonaRepository struct {
	mu         sync.RWMutex
	entries    map[string]MasterPersonaEntry
	aiSettings MasterPersonaAISettingsRecord
	runStatus  MasterPersonaRunStatusRecord
}

// InMemorySecretStore provides an in-memory secret seam for tests and local wiring.
type InMemorySecretStore struct {
	mu      sync.RWMutex
	secrets map[string]string
}

// NewInMemoryMasterPersonaRepository creates an in-memory master persona repository with seed data.
func NewInMemoryMasterPersonaRepository(seed []MasterPersonaEntry) *InMemoryMasterPersonaRepository {
	entries := make(map[string]MasterPersonaEntry, len(seed))
	for _, entry := range seed {
		entries[entry.IdentityKey] = cloneMasterPersonaEntry(entry)
	}
	return &InMemoryMasterPersonaRepository{
		entries: entries,
		runStatus: MasterPersonaRunStatusRecord{
			RunState: "入力待ち",
		},
	}
}

// NewInMemorySecretStore creates an in-memory secret store.
func NewInMemorySecretStore() *InMemorySecretStore {
	return &InMemorySecretStore{secrets: map[string]string{}}
}

// List returns a filtered, paged list of master persona entries.
func (repository *InMemoryMasterPersonaRepository) List(
	_ context.Context,
	query MasterPersonaListQuery,
) (MasterPersonaListResult, error) {
	repository.mu.RLock()
	defer repository.mu.RUnlock()

	keyword := strings.ToLower(strings.TrimSpace(query.Keyword))
	pluginFilter := strings.TrimSpace(query.PluginFilter)

	keywordMatched := make([]MasterPersonaEntry, 0, len(repository.entries))
	pluginCounts := map[string]int{}
	for _, entry := range repository.entries {
		if !masterPersonaMatchesKeyword(entry, keyword) {
			continue
		}
		keywordMatched = append(keywordMatched, cloneMasterPersonaEntry(entry))
		pluginCounts[entry.TargetPlugin]++
	}

	sort.Slice(keywordMatched, func(left, right int) bool {
		if keywordMatched[left].UpdatedAt.Equal(keywordMatched[right].UpdatedAt) {
			return keywordMatched[left].IdentityKey < keywordMatched[right].IdentityKey
		}
		return keywordMatched[left].UpdatedAt.After(keywordMatched[right].UpdatedAt)
	})

	filtered := make([]MasterPersonaEntry, 0, len(keywordMatched))
	for _, entry := range keywordMatched {
		if pluginFilter != "" && entry.TargetPlugin != pluginFilter {
			continue
		}
		filtered = append(filtered, entry)
	}

	page, pageSize := normalizeMasterPersonaPagination(query.Page, query.PageSize, len(filtered))
	start := (page - 1) * pageSize
	end := start + pageSize
	if start > len(filtered) {
		start = len(filtered)
	}
	if end > len(filtered) {
		end = len(filtered)
	}

	pluginGroups := make([]MasterPersonaPluginGroup, 0, len(pluginCounts))
	for targetPlugin, count := range pluginCounts {
		pluginGroups = append(pluginGroups, MasterPersonaPluginGroup{TargetPlugin: targetPlugin, Count: count})
	}
	sort.Slice(pluginGroups, func(left, right int) bool {
		return pluginGroups[left].TargetPlugin < pluginGroups[right].TargetPlugin
	})

	return MasterPersonaListResult{
		Items:        append([]MasterPersonaEntry(nil), filtered[start:end]...),
		TotalCount:   len(filtered),
		Page:         page,
		PageSize:     pageSize,
		PluginGroups: pluginGroups,
	}, nil
}

// GetByIdentityKey loads one master persona entry by identity key.
func (repository *InMemoryMasterPersonaRepository) GetByIdentityKey(
	_ context.Context,
	identityKey string,
) (MasterPersonaEntry, error) {
	repository.mu.RLock()
	defer repository.mu.RUnlock()

	entry, ok := repository.entries[strings.TrimSpace(identityKey)]
	if !ok {
		return MasterPersonaEntry{}, fmt.Errorf(masterPersonaIdentityKeyErrorFmt, ErrMasterPersonaEntryNotFound, identityKey)
	}
	return cloneMasterPersonaEntry(entry), nil
}

// UpsertIfAbsent inserts an entry only when the identity key does not already exist.
func (repository *InMemoryMasterPersonaRepository) UpsertIfAbsent(
	_ context.Context,
	draft MasterPersonaDraft,
) (MasterPersonaEntry, bool, error) {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	identityKey := strings.TrimSpace(draft.IdentityKey)
	if _, exists := repository.entries[identityKey]; exists {
		return cloneMasterPersonaEntry(repository.entries[identityKey]), false, nil
	}
	entry := entryFromDraft(draft)
	repository.entries[identityKey] = entry
	return cloneMasterPersonaEntry(entry), true, nil
}

// Update replaces one existing master persona entry.
func (repository *InMemoryMasterPersonaRepository) Update(
	_ context.Context,
	identityKey string,
	draft MasterPersonaDraft,
) (MasterPersonaEntry, error) {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	trimmedIdentityKey := strings.TrimSpace(identityKey)
	if _, exists := repository.entries[trimmedIdentityKey]; !exists {
		return MasterPersonaEntry{}, fmt.Errorf(masterPersonaIdentityKeyErrorFmt, ErrMasterPersonaEntryNotFound, identityKey)
	}

	nextIdentityKey := strings.TrimSpace(draft.IdentityKey)
	if nextIdentityKey != trimmedIdentityKey {
		if _, exists := repository.entries[nextIdentityKey]; exists {
			return MasterPersonaEntry{}, fmt.Errorf("master persona identity key already exists: %s", nextIdentityKey)
		}
		delete(repository.entries, trimmedIdentityKey)
	}

	entry := entryFromDraft(draft)
	repository.entries[nextIdentityKey] = entry
	return cloneMasterPersonaEntry(entry), nil
}

// Delete removes one master persona entry by identity key.
func (repository *InMemoryMasterPersonaRepository) Delete(_ context.Context, identityKey string) error {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	trimmedIdentityKey := strings.TrimSpace(identityKey)
	if _, exists := repository.entries[trimmedIdentityKey]; !exists {
		return fmt.Errorf(masterPersonaIdentityKeyErrorFmt, ErrMasterPersonaEntryNotFound, identityKey)
	}
	delete(repository.entries, trimmedIdentityKey)
	return nil
}

// LoadAISettings returns the current page-local AI settings record.
func (repository *InMemoryMasterPersonaRepository) LoadAISettings(_ context.Context) (MasterPersonaAISettingsRecord, error) {
	repository.mu.RLock()
	defer repository.mu.RUnlock()
	return repository.aiSettings, nil
}

// SaveAISettings stores the current page-local AI settings record.
func (repository *InMemoryMasterPersonaRepository) SaveAISettings(_ context.Context, record MasterPersonaAISettingsRecord) error {
	repository.mu.Lock()
	defer repository.mu.Unlock()
	repository.aiSettings = record
	return nil
}

// LoadRunStatus returns the persisted generation run status.
func (repository *InMemoryMasterPersonaRepository) LoadRunStatus(_ context.Context) (MasterPersonaRunStatusRecord, error) {
	repository.mu.RLock()
	defer repository.mu.RUnlock()
	return cloneMasterPersonaRunStatus(repository.runStatus), nil
}

// SaveRunStatus stores the persisted generation run status.
func (repository *InMemoryMasterPersonaRepository) SaveRunStatus(_ context.Context, status MasterPersonaRunStatusRecord) error {
	repository.mu.Lock()
	defer repository.mu.Unlock()
	repository.runStatus = cloneMasterPersonaRunStatus(status)
	return nil
}

// Load returns one secret value by key.
func (store *InMemorySecretStore) Load(_ context.Context, key string) (string, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	return store.secrets[strings.TrimSpace(key)], nil
}

// Save stores one secret value by key.
func (store *InMemorySecretStore) Save(_ context.Context, key string, value string) error {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.secrets[strings.TrimSpace(key)] = value
	return nil
}

// Delete removes one secret value by key.
func (store *InMemorySecretStore) Delete(_ context.Context, key string) error {
	store.mu.Lock()
	defer store.mu.Unlock()
	delete(store.secrets, strings.TrimSpace(key))
	return nil
}

// InMemoryMasterPersonaRunStatusRepository provides a session-local, non-persistent run status store.
// It satisfies MasterPersonaRunRepository and resets to idle on construction, so app restart clears state.
type InMemoryMasterPersonaRunStatusRepository struct {
	mu        sync.RWMutex
	runStatus MasterPersonaRunStatusRecord
}

// NewInMemoryMasterPersonaRunStatusRepository creates a session-local run status repository seeded with idle state.
func NewInMemoryMasterPersonaRunStatusRepository() *InMemoryMasterPersonaRunStatusRepository {
	return &InMemoryMasterPersonaRunStatusRepository{
		runStatus: MasterPersonaRunStatusRecord{RunState: "入力待ち"},
	}
}

// LoadRunStatus returns the current session run status.
func (r *InMemoryMasterPersonaRunStatusRepository) LoadRunStatus(_ context.Context) (MasterPersonaRunStatusRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return cloneMasterPersonaRunStatus(r.runStatus), nil
}

// SaveRunStatus stores the current session run status.
func (r *InMemoryMasterPersonaRunStatusRepository) SaveRunStatus(_ context.Context, status MasterPersonaRunStatusRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.runStatus = cloneMasterPersonaRunStatus(status)
	return nil
}

// DefaultMasterPersonaSeed returns deterministic seed entries for bootstrap wiring.
func DefaultMasterPersonaSeed(now time.Time) []MasterPersonaEntry {
	lysmarenRace := "Breton"
	lysmarenSex := "Female"
	kaelRace := "Nord"
	kaelSex := "Male"
	return []MasterPersonaEntry{
		{
			IdentityKey:    BuildMasterPersonaIdentityKey(masterPersonaSeedFollowersPlugin, "FE01A812", "NPC_"),
			TargetPlugin:   masterPersonaSeedFollowersPlugin,
			FormID:         "FE01A812",
			RecordType:     "NPC_",
			EditorID:       "FP_LysMaren",
			DisplayName:    "Lys Maren",
			Race:           &lysmarenRace,
			Sex:            &lysmarenSex,
			VoiceType:      "FemaleYoungEager",
			ClassName:      "FPScoutClass",
			SourcePlugin:   masterPersonaSeedFollowersPlugin,
			PersonaSummary: "乾いた率直さで応じ、必要な場面だけ短く本音を置く。",
			PersonaBody:    "口調は丁寧語へ寄せず、中性的な温度を保つ。会話の主導権は急いで取らず、相手の出方を見てから短く返す。",
			DialogueCount:  3,
			Dialogues: []string{
				"ここで待って。まだ相手の出方が見えていない。",
				"急がなくていい。必要になったら、わたしから声をかける。",
				"本音を聞きたいなら、先にそっちが隠し事をやめて。",
			},
			UpdatedAt: now,
		},
		{
			IdentityKey:    BuildMasterPersonaIdentityKey(masterPersonaSeedFollowersPlugin, "FE01A813", "NPC_"),
			TargetPlugin:   masterPersonaSeedFollowersPlugin,
			FormID:         "FE01A813",
			RecordType:     "NPC_",
			EditorID:       "FP_KaelRuun",
			DisplayName:    "Kael Ruun",
			Race:           &kaelRace,
			Sex:            &kaelSex,
			VoiceType:      "MaleCommander",
			ClassName:      "FPMercenaryClass",
			SourcePlugin:   masterPersonaSeedFollowersPlugin,
			PersonaSummary: "判断を先に示し、無駄なく短く指示を伝える。",
			PersonaBody:    "判断を先に述べ、必要な指示だけを短く渡す。曖昧な慰めより役割と責任を優先する。",
			DialogueCount:  2,
			Dialogues: []string{
				"命令は簡潔でいい。動けるなら動け。",
				"状況確認を先に済ませる。感想は後だ。",
			},
			UpdatedAt: now.Add(-time.Minute),
		},
		{
			IdentityKey:     BuildMasterPersonaIdentityKey(masterPersonaSeedNightCourtPlugin, "FE01A814", "NPC_"),
			TargetPlugin:    masterPersonaSeedNightCourtPlugin,
			FormID:          "FE01A814",
			RecordType:      "NPC_",
			EditorID:        "FP_WatcherHusk",
			DisplayName:     "Watcher Husk",
			VoiceType:       "FemaleCondescending",
			ClassName:       "FPOccultClass",
			SourcePlugin:    masterPersonaSeedNightCourtPlugin,
			PersonaSummary:  "含みのある言い回しで相手を試し、答えを急がせない。",
			PersonaBody:     "含みを残した言い回しで相手の反応を測る。欠落属性は見せず、観察を優先する話し方に寄せる。",
			BaselineApplied: true,
			DialogueCount:   2,
			Dialogues: []string{
				"急いで答えを出す必要はないわ。迷い方にも価値がある。",
				"隠していることがあるなら、声の揺れで十分にわかる。",
			},
			UpdatedAt: now.Add(-2 * time.Minute),
		},
	}
}

// BuildMasterPersonaIdentityKey builds the no-overwrite identity key for master persona entries.
func BuildMasterPersonaIdentityKey(targetPlugin string, formID string, recordType string) string {
	return strings.TrimSpace(targetPlugin) + ":" + strings.TrimSpace(formID) + ":" + strings.TrimSpace(recordType)
}

func cloneMasterPersonaEntry(entry MasterPersonaEntry) MasterPersonaEntry {
	cloned := entry
	if entry.Race != nil {
		race := *entry.Race
		cloned.Race = &race
	}
	if entry.Sex != nil {
		sex := *entry.Sex
		cloned.Sex = &sex
	}
	cloned.Dialogues = append([]string(nil), entry.Dialogues...)
	return cloned
}

func cloneMasterPersonaRunStatus(status MasterPersonaRunStatusRecord) MasterPersonaRunStatusRecord {
	cloned := status
	if status.StartedAt != nil {
		startedAt := *status.StartedAt
		cloned.StartedAt = &startedAt
	}
	if status.FinishedAt != nil {
		finishedAt := *status.FinishedAt
		cloned.FinishedAt = &finishedAt
	}
	return cloned
}

func entryFromDraft(draft MasterPersonaDraft) MasterPersonaEntry {
	entry := MasterPersonaEntry{
		IdentityKey:          strings.TrimSpace(draft.IdentityKey),
		TargetPlugin:         strings.TrimSpace(draft.TargetPlugin),
		FormID:               strings.TrimSpace(draft.FormID),
		RecordType:           strings.TrimSpace(draft.RecordType),
		EditorID:             strings.TrimSpace(draft.EditorID),
		DisplayName:          strings.TrimSpace(draft.DisplayName),
		VoiceType:            strings.TrimSpace(draft.VoiceType),
		ClassName:            strings.TrimSpace(draft.ClassName),
		SourcePlugin:         strings.TrimSpace(draft.SourcePlugin),
		PersonaSummary:       strings.TrimSpace(draft.PersonaSummary),
		SpeechStyle:          strings.TrimSpace(draft.SpeechStyle),
		PersonaBody:          strings.TrimSpace(draft.PersonaBody),
		GenerationSourceJSON: strings.TrimSpace(draft.GenerationSourceJSON),
		BaselineApplied:      draft.BaselineApplied,
		DialogueCount:        len(draft.Dialogues),
		Dialogues:            append([]string(nil), draft.Dialogues...),
		UpdatedAt:            draft.UpdatedAt.UTC(),
	}
	if draft.Race != nil {
		race := strings.TrimSpace(*draft.Race)
		entry.Race = &race
	}
	if draft.Sex != nil {
		sex := strings.TrimSpace(*draft.Sex)
		entry.Sex = &sex
	}
	return entry
}

func masterPersonaMatchesKeyword(entry MasterPersonaEntry, keyword string) bool {
	if keyword == "" {
		return true
	}
	searchTarget := strings.ToLower(strings.Join([]string{
		entry.DisplayName,
		entry.FormID,
		entry.EditorID,
		pointerString(entry.Race),
		entry.VoiceType,
	}, " "))
	return strings.Contains(searchTarget, keyword)
}

func pointerString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func normalizeMasterPersonaPagination(page int, pageSize int, total int) (int, int) {
	normalizedPageSize := pageSize
	if normalizedPageSize <= 0 {
		normalizedPageSize = masterPersonaDefaultPageSize
	}
	if normalizedPageSize > masterPersonaMaxPageSize {
		normalizedPageSize = masterPersonaMaxPageSize
	}

	normalizedPage := page
	if normalizedPage <= 0 {
		normalizedPage = 1
	}

	maxPage := 1
	if total > 0 {
		maxPage = ((total - 1) / normalizedPageSize) + 1
	}
	if normalizedPage > maxPage {
		normalizedPage = maxPage
	}
	return normalizedPage, normalizedPageSize
}
