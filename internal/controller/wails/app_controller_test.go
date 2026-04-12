package wails

import (
	"context"
	"path/filepath"
	"testing"
)

const (
	controllerSource                 = "Controller Source"
	frontendContractSource           = "Frontend Contract Source"
	frontendContractCategory         = "固有名詞"
	frontendContractOrigin           = "手動登録"
	testMasterDictionaryImportTarget = "Dawnguard_english_japanese.xml"
)

type testContextKey string

func TestAppControllerLifecycleHooks(_ *testing.T) {
	controller := NewAppController()
	controller.OnStartup(context.Background())
	controller.OnShutdown(context.Background())
}

func TestAppControllerHealthReturnsOkStatus(t *testing.T) {
	controller := NewAppController()

	response := controller.Health()

	if response.Status != "ok" {
		t.Fatalf("expected status ok, got %q", response.Status)
	}
}

func TestAppControllerMasterDictionaryCRUDFlow(t *testing.T) {
	controller := NewAppController()

	createResponse, err := controller.MasterDictionaryCreate(MasterDictionaryCreateRequestDTO{
		Entry: MasterDictionaryMutationInputDTO{
			Source:      controllerSource,
			Translation: "コントローラ訳語",
			Category:    frontendContractCategory,
			Origin:      frontendContractOrigin,
		},
		Refresh: MasterDictionaryRefreshQueryDTO{Page: 1, PageSize: 30},
	})
	if err != nil {
		t.Fatalf("expected create to succeed: %v", err)
	}
	if createResponse.ChangedEntry == nil {
		t.Fatal("expected create to return changed entry")
	}

	detailResponse, err := controller.MasterDictionaryGetDetail(MasterDictionaryDetailRequestDTO{ID: createResponse.ChangedEntry.ID})
	if err != nil {
		t.Fatalf("expected detail fetch to succeed: %v", err)
	}
	if detailResponse.Entry.Source != controllerSource {
		t.Fatalf("expected source to match, got %q", detailResponse.Entry.Source)
	}

	updateResponse, err := controller.MasterDictionaryUpdate(MasterDictionaryUpdateRequestDTO{
		ID: createResponse.ChangedEntry.ID,
		Entry: MasterDictionaryMutationInputDTO{
			Source:      controllerSource,
			Translation: "コントローラ更新訳語",
			Category:    frontendContractCategory,
			Origin:      frontendContractOrigin,
		},
		Refresh: MasterDictionaryRefreshQueryDTO{Page: 1, PageSize: 30},
	})
	if err != nil {
		t.Fatalf("expected update to succeed: %v", err)
	}
	if updateResponse.ChangedEntry == nil || updateResponse.ChangedEntry.Translation != "コントローラ更新訳語" {
		t.Fatal("expected update to return changed entry with updated translation")
	}

	deleteResponse, err := controller.MasterDictionaryDelete(MasterDictionaryDeleteRequestDTO{
		ID:      createResponse.ChangedEntry.ID,
		Refresh: MasterDictionaryRefreshQueryDTO{Page: 1, PageSize: 30},
	})
	if err != nil {
		t.Fatalf("expected delete to succeed: %v", err)
	}
	if deleteResponse.DeletedEntryID == nil || *deleteResponse.DeletedEntryID != createResponse.ChangedEntry.ID {
		t.Fatal("expected delete response to include deleted entry id")
	}
}

func TestAppControllerImportRespectsAllowlist(t *testing.T) {
	controller := NewAppController()

	importResponse, err := controller.MasterDictionaryImportXML(MasterDictionaryImportRequestDTO{
		XMLPath: filepath.Clean("../../../dictionaries/Dawnguard_english_japanese.xml"),
		Refresh: MasterDictionaryRefreshQueryDTO{Page: 1, PageSize: 30},
	})
	if err != nil {
		t.Fatalf("expected import to succeed: %v", err)
	}
	if importResponse.Summary.ImportedCount == 0 {
		t.Fatal("expected imported count to be greater than zero")
	}

	allowedPage, err := controller.MasterDictionaryGetPage(MasterDictionaryPageRequestDTO{
		Refresh: MasterDictionaryRefreshQueryDTO{
			SearchTerm: "Auriel's Bow",
			Page:       1,
			PageSize:   30,
		},
	})
	if err != nil {
		t.Fatalf("expected allowed search to succeed: %v", err)
	}
	if len(allowedPage.Page.Items) == 0 {
		t.Fatal("expected allowed REC entry to be searchable")
	}

	blockedPage, err := controller.MasterDictionaryGetPage(MasterDictionaryPageRequestDTO{
		Refresh: MasterDictionaryRefreshQueryDTO{
			SearchTerm: "Crossbow Mount",
			Page:       1,
			PageSize:   30,
		},
	})
	if err != nil {
		t.Fatalf("expected blocked search to succeed: %v", err)
	}
	if len(blockedPage.Page.Items) != 0 {
		t.Fatal("expected blocked REC entry to remain filtered out")
	}
}

func TestAppControllerFrontendContractCRUDFlow(t *testing.T) {
	controller := NewAppController()

	created := mustCreateFrontendContractEntry(t, controller, frontendContractSource, "契約訳語", 2)
	mustListFrontendContractEntry(t, controller, frontendContractSource)
	mustFrontendContractDetailMatchSource(t, controller, created.RefreshTargetID, frontendContractSource)
	mustUpdateFrontendContractEntry(t, controller, created.RefreshTargetID, frontendContractSource, "契約更新訳語", 3)
	mustDeleteFrontendContractEntry(t, controller, created.RefreshTargetID, 4)
}

func mustCreateFrontendContractEntry(
	t *testing.T,
	controller *AppController,
	source string,
	translation string,
	page int,
) CreateMasterDictionaryEntryResponseDTO {
	t.Helper()

	createRequest := CreateMasterDictionaryEntryRequestDTO{}
	createRequest.Payload.Source = source
	createRequest.Payload.Translation = translation
	createRequest.Payload.Category = frontendContractCategory
	createRequest.Payload.Origin = frontendContractOrigin
	createRequest.Refresh = &MasterDictionaryFrontendRefreshDTO{Page: page, PageSize: 1}

	created, err := controller.CreateMasterDictionaryEntry(createRequest)
	if err != nil {
		t.Fatalf("expected create contract endpoint to succeed: %v", err)
	}
	if created.RefreshTargetID == "" {
		t.Fatal("expected refresh target id to be set")
	}
	mustFrontendRefreshPage(t, created.Page, page, 1, "create")

	return created
}

func mustListFrontendContractEntry(t *testing.T, controller *AppController, source string) {
	t.Helper()

	listRequest := ListMasterDictionaryEntriesRequestDTO{}
	listRequest.Filters.Query = source
	listRequest.Filters.Page = 1
	listRequest.Filters.PageSize = 30

	listed, err := controller.ListMasterDictionaryEntries(listRequest)
	if err != nil {
		t.Fatalf("expected list contract endpoint to succeed: %v", err)
	}
	if len(listed.Entries) == 0 {
		t.Fatal("expected created entry to be listed")
	}
}

func mustFrontendContractDetailMatchSource(
	t *testing.T,
	controller *AppController,
	id string,
	expectedSource string,
) {
	t.Helper()

	detail, err := controller.GetMasterDictionaryEntry(GetMasterDictionaryEntryRequestDTO{ID: id})
	if err != nil {
		t.Fatalf("expected detail contract endpoint to succeed: %v", err)
	}
	if detail.Entry == nil || detail.Entry.Source != expectedSource {
		t.Fatal("expected detail entry to match created source")
	}
}

func mustUpdateFrontendContractEntry(
	t *testing.T,
	controller *AppController,
	id string,
	source string,
	expectedTranslation string,
	page int,
) {
	t.Helper()

	updateRequest := UpdateMasterDictionaryEntryRequestDTO{ID: id}
	updateRequest.Payload.Source = source
	updateRequest.Payload.Translation = expectedTranslation
	updateRequest.Payload.Category = frontendContractCategory
	updateRequest.Payload.Origin = frontendContractOrigin
	updateRequest.Refresh = &MasterDictionaryFrontendRefreshDTO{Page: page, PageSize: 1}

	updated, err := controller.UpdateMasterDictionaryEntry(updateRequest)
	if err != nil {
		t.Fatalf("expected update contract endpoint to succeed: %v", err)
	}
	if updated.Entry.Translation != expectedTranslation {
		t.Fatalf("expected updated translation, got %q", updated.Entry.Translation)
	}
	mustFrontendRefreshPage(t, updated.Page, page, 1, "update")
}

func mustDeleteFrontendContractEntry(t *testing.T, controller *AppController, id string, page int) {
	t.Helper()

	deleted, err := controller.DeleteMasterDictionaryEntry(DeleteMasterDictionaryEntryRequestDTO{
		ID:      id,
		Refresh: &MasterDictionaryFrontendRefreshDTO{Page: page, PageSize: 1},
	})
	if err != nil {
		t.Fatalf("expected delete contract endpoint to succeed: %v", err)
	}
	if deleted.DeletedID != id {
		t.Fatalf("expected deleted id %s, got %s", id, deleted.DeletedID)
	}
	mustFrontendRefreshPage(t, deleted.Page, page, 1, "delete")
}

func mustFrontendRefreshPage(
	t *testing.T,
	pageDTO *MasterDictionaryPageDTO,
	expectedPage int,
	expectedPageSize int,
	operation string,
) {
	t.Helper()

	if pageDTO == nil {
		t.Fatalf("expected %s response to include refreshed page", operation)
	}
	if pageDTO.Page != expectedPage || pageDTO.PageSize != expectedPageSize {
		t.Fatalf(
			"expected %s refresh page/pageSize to be preserved, got page=%d size=%d",
			operation,
			pageDTO.Page,
			pageDTO.PageSize,
		)
	}
}

func TestAppControllerFrontendContractImportAlias(t *testing.T) {
	controller := NewAppController()

	request := ImportMasterDictionaryXMLRequestDTO{
		FilePath: filepath.Clean("../../../dictionaries/" + testMasterDictionaryImportTarget),
		Refresh:  &MasterDictionaryFrontendRefreshDTO{Page: 2, PageSize: 1},
	}

	response, err := controller.ImportMasterDictionaryXML(request)
	if err != nil {
		t.Fatalf("expected import endpoint to succeed: %v", err)
	}
	if !response.Accepted {
		t.Fatal("expected import endpoint to accept request")
	}
	if response.Page == nil || response.Summary == nil {
		t.Fatal("expected import endpoint to return refresh payload")
	}
	if response.Page.Page != 2 || response.Page.PageSize != 1 {
		t.Fatalf("expected import refresh page/pageSize to be preserved, got page=%d size=%d", response.Page.Page, response.Page.PageSize)
	}
	if response.Summary.ImportedCount+response.Summary.UpdatedCount == 0 {
		t.Fatal("expected import endpoint to report imported or updated count")
	}

	aliasResponse, err := controller.ImportMasterDictionaryXml(request)
	if err != nil {
		t.Fatalf("expected alias import endpoint to succeed: %v", err)
	}
	if !aliasResponse.Accepted {
		t.Fatal("expected alias import endpoint to accept request")
	}
	if aliasResponse.Page == nil || aliasResponse.Summary == nil {
		t.Fatal("expected alias import endpoint to return refresh payload")
	}
}

func TestAppControllerFrontendContractImportUsesFileReference(t *testing.T) {
	controller := NewAppController()

	request := ImportMasterDictionaryXMLRequestDTO{
		FilePath:      "ignored-by-reference.xml",
		FileReference: testMasterDictionaryImportTarget,
	}

	canonical, err := controller.ImportMasterDictionaryXML(request)
	if err != nil {
		t.Fatalf("expected import endpoint to resolve file reference: %v", err)
	}
	if !canonical.Accepted {
		t.Fatal("expected file-reference import endpoint to accept request")
	}
	if canonical.Page == nil || canonical.Summary == nil {
		t.Fatal("expected file-reference import endpoint to return refresh payload")
	}
	if canonical.Summary.FileName != testMasterDictionaryImportTarget {
		t.Fatalf("expected canonical file name from fileReference, got %q", canonical.Summary.FileName)
	}

	alias, err := controller.ImportMasterDictionaryXml(request)
	if err != nil {
		t.Fatalf("expected alias endpoint to resolve file reference: %v", err)
	}
	if !alias.Accepted {
		t.Fatal("expected alias file-reference import endpoint to accept request")
	}
	if alias.Page == nil || alias.Summary == nil {
		t.Fatal("expected alias file-reference import endpoint to return refresh payload")
	}
	if alias.Summary.FileName != testMasterDictionaryImportTarget {
		t.Fatalf("expected alias file name from fileReference, got %q", alias.Summary.FileName)
	}
}

func TestResolveMasterDictionaryImportReference(t *testing.T) {
	requestWithBoth := ImportMasterDictionaryXMLRequestDTO{
		FilePath:      " /tmp/path.xml ",
		FileReference: " " + testMasterDictionaryImportTarget + " ",
	}
	resolved := resolveMasterDictionaryImportReference(requestWithBoth)
	if resolved != testMasterDictionaryImportTarget {
		t.Fatalf("expected fileReference to take priority, got %q", resolved)
	}

	requestWithFilePathOnly := ImportMasterDictionaryXMLRequestDTO{FilePath: " /tmp/path.xml "}
	resolved = resolveMasterDictionaryImportReference(requestWithFilePathOnly)
	if resolved != "/tmp/path.xml" {
		t.Fatalf("expected trimmed filePath when fileReference is empty, got %q", resolved)
	}
}

func TestAppControllerGetMasterDictionaryEntryNotFoundReturnsNil(t *testing.T) {
	controller := NewAppController()

	response, err := controller.GetMasterDictionaryEntry(GetMasterDictionaryEntryRequestDTO{ID: "999999999"})
	if err != nil {
		t.Fatalf("expected not found to be mapped to nil entry: %v", err)
	}
	if response.Entry != nil {
		t.Fatal("expected entry to be nil for not found")
	}
}

func TestAppControllerContractRejectsInvalidID(t *testing.T) {
	controller := NewAppController()

	if _, err := controller.GetMasterDictionaryEntry(GetMasterDictionaryEntryRequestDTO{ID: ""}); err == nil {
		t.Fatal("expected empty id error")
	}

	updateRequest := UpdateMasterDictionaryEntryRequestDTO{ID: "abc"}
	if _, err := controller.UpdateMasterDictionaryEntry(updateRequest); err == nil {
		t.Fatal("expected parse error for update id")
	}

	if _, err := controller.DeleteMasterDictionaryEntry(DeleteMasterDictionaryEntryRequestDTO{ID: "-1"}); err == nil {
		t.Fatal("expected validation error for delete id")
	}
}

func TestParseStringID(t *testing.T) {
	validID, err := parseStringID(" 123 ")
	if err != nil {
		t.Fatalf("expected parse success: %v", err)
	}
	if validID != 123 {
		t.Fatalf("expected parsed id 123, got %d", validID)
	}

	if _, err := parseStringID(" "); err == nil {
		t.Fatal("expected required id error")
	}
	if _, err := parseStringID("abc"); err == nil {
		t.Fatal("expected parse error")
	}
	if _, err := parseStringID("0"); err == nil {
		t.Fatal("expected positive id validation error")
	}
}

func TestToEntryDetailDTONoteFormatting(t *testing.T) {
	base := MasterDictionaryEntryDTO{ID: 1, Source: "src", Translation: "dst", Category: "cat", Origin: "org", UpdatedAt: "2026-01-01T00:00:00Z"}

	plain := toEntryDetailDTO(base)
	if plain.Note != "マスター辞書エントリ" {
		t.Fatalf("expected default note, got %q", plain.Note)
	}

	withREC := base
	withREC.REC = "BOOK:FULL"
	withRECDetail := toEntryDetailDTO(withREC)
	if withRECDetail.Note != "REC: BOOK:FULL" {
		t.Fatalf("expected rec note, got %q", withRECDetail.Note)
	}

	withEDID := withREC
	withEDID.EDID = "BookAuriel"
	withEDIDDetail := toEntryDetailDTO(withEDID)
	if withEDIDDetail.Note != "REC: BOOK:FULL / EDID: BookAuriel" {
		t.Fatalf("expected rec+edid note, got %q", withEDIDDetail.Note)
	}
}

func TestAppControllerRequestContextAlwaysReturnsBackground(t *testing.T) {
	controller := NewAppController()

	if controller.requestContext() == nil {
		t.Fatal("expected request context to always be available")
	}

	startup := context.WithValue(context.Background(), testContextKey("key"), "value")
	controller.OnStartup(startup)
	if controller.requestContext() == startup {
		t.Fatal("expected request context to not retain startup context")
	}

	controller.OnShutdown(context.Background())
	if controller.requestContext() == startup {
		t.Fatal("expected request context to remain detached from startup context")
	}
}
