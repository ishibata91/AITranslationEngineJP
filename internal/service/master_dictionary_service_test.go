package service

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"aitranslationenginejp/internal/repository"
)

func TestMasterDictionaryImportFromXMLFiltersByAllowedREC(t *testing.T) {
	repo := repository.NewInMemoryMasterDictionaryRepository(nil)
	svc := NewMasterDictionaryService(repo, func() time.Time {
		return time.Date(2026, 4, 12, 10, 0, 0, 0, time.UTC)
	})

	xmlPath := filepath.Clean("../../dictionaries/Dawnguard_english_japanese.xml")
	summary, err := svc.ImportFromXML(context.Background(), xmlPath)
	if err != nil {
		t.Fatalf("expected xml import to succeed: %v", err)
	}
	if summary.ImportedCount == 0 {
		t.Fatal("expected at least one imported entry")
	}

	allowedRows, err := svc.ListEntries(context.Background(), MasterDictionaryQuery{
		SearchTerm: "Auriel's Bow",
		Page:       1,
		PageSize:   30,
	})
	if err != nil {
		t.Fatalf("expected list for allowed REC to succeed: %v", err)
	}
	if len(allowedRows.Items) == 0 {
		t.Fatal("expected allowed REC entry to exist after import")
	}

	notAllowedRows, err := svc.ListEntries(context.Background(), MasterDictionaryQuery{
		SearchTerm: "Crossbow Mount",
		Page:       1,
		PageSize:   30,
	})
	if err != nil {
		t.Fatalf("expected list for denied REC to succeed: %v", err)
	}
	if len(notAllowedRows.Items) != 0 {
		t.Fatal("expected denied REC entry to be filtered out")
	}

	spellRows, err := svc.ListEntries(context.Background(), MasterDictionaryQuery{
		SearchTerm: "Transform into the vampire lord.",
		Page:       1,
		PageSize:   30,
	})
	if err != nil {
		t.Fatalf("expected list for denied REC field to succeed: %v", err)
	}
	if len(spellRows.Items) != 0 {
		t.Fatal("expected denied REC field entry to be filtered out")
	}
}

func TestMasterDictionaryImportFromXMLAllowlistWithMinimalFixture(t *testing.T) {
	repo := repository.NewInMemoryMasterDictionaryRepository(nil)
	svc := NewMasterDictionaryService(repo, func() time.Time {
		return time.Date(2026, 4, 12, 10, 0, 0, 0, time.UTC)
	})

	xmlContent := `<?xml version="1.0" encoding="utf-8"?>
<Root>
	<String>
		<REC>BOOK:FULL</REC>
		<EDID>AllowedBook</EDID>
		<Source>Allowed Source</Source>
		<Dest>許可訳</Dest>
	</String>
	<String>
		<REC>ACTI:FULL</REC>
		<EDID>DeniedActi</EDID>
		<Source>Denied Source</Source>
		<Dest>拒否訳</Dest>
	</String>
	<String>
		<REC>NPC_:FULL</REC>
		<EDID>MissingDest</EDID>
		<Source>Missing Dest Source</Source>
		<Dest></Dest>
	</String>
</Root>`

	tmpDir := t.TempDir()
	xmlPath := filepath.Join(tmpDir, "allowlist.xml")
	if err := os.WriteFile(xmlPath, []byte(xmlContent), 0o600); err != nil {
		t.Fatalf("write xml fixture: %v", err)
	}

	summary, err := svc.ImportFromXML(context.Background(), xmlPath)
	if err != nil {
		t.Fatalf("expected xml import to succeed: %v", err)
	}
	if summary.ImportedCount != 1 {
		t.Fatalf("expected exactly one imported record, got %d", summary.ImportedCount)
	}
	if summary.SkippedCount != 2 {
		t.Fatalf("expected exactly two skipped records, got %d", summary.SkippedCount)
	}

	allowedRows, err := svc.ListEntries(context.Background(), MasterDictionaryQuery{
		SearchTerm: "Allowed Source",
		Page:       1,
		PageSize:   30,
	})
	if err != nil {
		t.Fatalf("expected list for allowed REC to succeed: %v", err)
	}
	if len(allowedRows.Items) != 1 {
		t.Fatalf("expected one allowed entry, got %d", len(allowedRows.Items))
	}

	deniedRows, err := svc.ListEntries(context.Background(), MasterDictionaryQuery{
		SearchTerm: "Denied Source",
		Page:       1,
		PageSize:   30,
	})
	if err != nil {
		t.Fatalf("expected list for denied REC to succeed: %v", err)
	}
	if len(deniedRows.Items) != 0 {
		t.Fatalf("expected denied REC entry to be filtered out, got %d", len(deniedRows.Items))
	}
}

func TestMasterDictionaryServiceValidationRejectsInvalidInput(t *testing.T) {
	svc := newTestMasterDictionaryService()

	if _, err := svc.GetEntry(context.Background(), 0); err == nil {
		t.Fatal("expected get with invalid id to fail")
	}
	if _, err := svc.CreateEntry(context.Background(), MasterDictionaryMutationInput{Source: "  ", Translation: "訳語"}); err == nil {
		t.Fatal("expected create with empty source to fail")
	}
	if _, err := svc.CreateEntry(context.Background(), MasterDictionaryMutationInput{Source: "Source", Translation: "  "}); err == nil {
		t.Fatal("expected create with empty translation to fail")
	}
	if _, err := svc.UpdateEntry(context.Background(), 0, MasterDictionaryMutationInput{Source: "S", Translation: "T"}); err == nil {
		t.Fatal("expected update with invalid id to fail")
	}
	if _, err := svc.UpdateEntry(context.Background(), 1, MasterDictionaryMutationInput{Source: "", Translation: "T"}); err == nil {
		t.Fatal("expected update with empty source to fail")
	}
	if err := svc.DeleteEntry(context.Background(), 0); err == nil {
		t.Fatal("expected delete with invalid id to fail")
	}
}

func TestMasterDictionaryServiceCRUDFlow(t *testing.T) {
	svc := newTestMasterDictionaryService()

	created, err := svc.CreateEntry(context.Background(), MasterDictionaryMutationInput{
		Source:      "  Source A  ",
		Translation: "  訳語A  ",
		REC:         " BOOK:FULL ",
		EDID:        " EDID_A ",
	})
	if err != nil {
		t.Fatalf("expected create to succeed: %v", err)
	}
	if created.Category != "固有名詞" {
		t.Fatalf("expected default category, got %q", created.Category)
	}
	if created.Origin != "手動登録" {
		t.Fatalf("expected default origin, got %q", created.Origin)
	}
	if created.Source != "Source A" || created.Translation != "訳語A" {
		t.Fatal("expected source and translation to be trimmed")
	}

	fetched, err := svc.GetEntry(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("expected get to succeed: %v", err)
	}
	if fetched.ID != created.ID {
		t.Fatalf("expected fetched id %d, got %d", created.ID, fetched.ID)
	}

	updated, err := svc.UpdateEntry(context.Background(), created.ID, MasterDictionaryMutationInput{
		Source:      "Source B",
		Translation: "訳語B",
		Category:    "地名",
		Origin:      "更新",
	})
	if err != nil {
		t.Fatalf("expected update to succeed: %v", err)
	}
	if updated.Translation != "訳語B" || updated.Category != "地名" {
		t.Fatal("expected updated values to be reflected")
	}
}

func TestMasterDictionaryServiceDeleteReportsNotFound(t *testing.T) {
	svc := newTestMasterDictionaryService()

	created, err := svc.CreateEntry(context.Background(), MasterDictionaryMutationInput{
		Source:      "Source A",
		Translation: "訳語A",
	})
	if err != nil {
		t.Fatalf("expected create to succeed: %v", err)
	}
	err = svc.DeleteEntry(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("expected delete to succeed: %v", err)
	}

	_, err = svc.GetEntry(context.Background(), created.ID)
	if err == nil {
		t.Fatal("expected not found after delete")
	}
	if !IsNotFoundError(err) {
		t.Fatal("expected IsNotFoundError to detect repository not found error")
	}
	if IsNotFoundError(errors.New("different")) {
		t.Fatal("expected unrelated error to return false")
	}
}

func newTestMasterDictionaryService() *MasterDictionaryService {
	repo := repository.NewInMemoryMasterDictionaryRepository(nil)
	return NewMasterDictionaryService(repo, func() time.Time {
		return time.Date(2026, 4, 12, 10, 0, 0, 0, time.UTC)
	})
}

func TestMasterDictionaryServiceDefaultFactoryAndImportValidation(t *testing.T) {
	svc := NewDefaultMasterDictionaryService(func() time.Time {
		return time.Date(2026, 4, 12, 10, 0, 0, 0, time.UTC)
	})

	listed, err := svc.ListEntries(context.Background(), MasterDictionaryQuery{Page: 1, PageSize: 30})
	if err != nil {
		t.Fatalf("expected default service list to succeed: %v", err)
	}
	if listed.TotalCount == 0 {
		t.Fatal("expected default service to contain seeded entries")
	}

	_, err = svc.ImportFromXML(context.Background(), " ")
	if err == nil {
		t.Fatal("expected empty xml path to fail")
	}

	summary, err := svc.ImportFromXML(context.Background(), "Dawnguard_english_japanese.xml")
	if err != nil {
		t.Fatalf("expected filename-based xml reference to succeed: %v", err)
	}
	if summary.FileName != "Dawnguard_english_japanese.xml" {
		t.Fatalf("expected imported filename to be preserved, got %q", summary.FileName)
	}
}
