package service

import (
	"context"
	"fmt"
	"io"
	"strings"
	"testing"
)

var importServiceRecords = []xmlStringRecord{
	{REC: "WEAP:FULL", EDID: "DLC1AurielsBow", Source: "Auriel's Bow", Dest: "アーリエルの弓"},
	{REC: "BOOK:FULL", EDID: "BookSnowElf", Source: "Snow Elf History", Dest: "スノーエルフ史"},
	{REC: "ACTI:FULL", EDID: "DeniedActi", Source: "Crossbow Mount", Dest: "クロスボウ台座"},
	{REC: "NPC_:FULL", EDID: "MissingDest", Source: "Missing Translation", Dest: "   "},
}

var importServiceExpectedProgress = []int{0, 25, 50, 75, 100, 100}

type importServiceSuccessResult struct {
	repo       *repositoryStub
	xmlFiles   *xmlFilePortStub
	xmlRecords *xmlRecordReaderStub
	runtime    *importProgressRecorder
	summary    MasterDictionaryImportSummary
}

func TestMasterDictionaryImportServiceResolvesOriginalXMLPath(t *testing.T) {
	result := runSuccessfulImport(t)

	if result.xmlFiles.resolvedPaths[0] != " input/import.xml " {
		t.Fatalf("expected one resolve call with raw path, got %v", result.xmlFiles.resolvedPaths)
	}
}

func TestMasterDictionaryImportServiceOpensXMLForCountAndRead(t *testing.T) {
	result := runSuccessfulImport(t)

	if len(result.xmlFiles.openedPaths) != 2 {
		t.Fatalf("expected xml file to be opened for count and read, got %v", result.xmlFiles.openedPaths)
	}
}

func TestMasterDictionaryImportServiceCountsAndReadsXMLRecordsOnce(t *testing.T) {
	result := runSuccessfulImport(t)

	if result.xmlRecords.countCalls != 1 || result.xmlRecords.readCalls != 1 {
		t.Fatalf("expected one count/read call, got count=%d read=%d", result.xmlRecords.countCalls, result.xmlRecords.readCalls)
	}
}

func TestMasterDictionaryImportServiceReturnsResolvedFileInfo(t *testing.T) {
	result := runSuccessfulImport(t)

	if result.summary.FilePath != "/resolved/import.xml" || result.summary.FileName != "import.xml" {
		t.Fatalf("expected resolved file info, got %+v", result.summary)
	}
}

func TestMasterDictionaryImportServiceReturnsImportCounters(t *testing.T) {
	result := runSuccessfulImport(t)

	if result.summary.ImportedCount != 1 || result.summary.UpdatedCount != 1 || result.summary.SkippedCount != 2 {
		t.Fatalf("unexpected summary: %+v", result.summary)
	}
}

func TestMasterDictionaryImportServiceReturnsLastEntryID(t *testing.T) {
	result := runSuccessfulImport(t)

	if result.summary.LastEntryID != 2 {
		t.Fatalf("expected last entry id 2, got %d", result.summary.LastEntryID)
	}
}

func TestMasterDictionaryImportServiceUpsertsOnlyImportableRecords(t *testing.T) {
	result := runSuccessfulImport(t)

	if len(result.repo.upsertRecords) != 2 {
		t.Fatalf("expected 2 upsert calls, got %d", len(result.repo.upsertRecords))
	}
}

func TestMasterDictionaryImportServiceMapsWeaponCategory(t *testing.T) {
	result := runSuccessfulImport(t)

	if result.repo.upsertRecords[0].Category != "装備" {
		t.Fatalf("expected WEAP record category 装備, got %q", result.repo.upsertRecords[0].Category)
	}
}

func TestMasterDictionaryImportServiceMapsBookCategory(t *testing.T) {
	result := runSuccessfulImport(t)

	if result.repo.upsertRecords[1].Category != "書籍" {
		t.Fatalf("expected BOOK record category 書籍, got %q", result.repo.upsertRecords[1].Category)
	}
}

func TestMasterDictionaryImportServiceAssignsXMLImportOrigin(t *testing.T) {
	result := runSuccessfulImport(t)

	for _, record := range result.repo.upsertRecords {
		if record.Origin != "XML取込" {
			t.Fatalf("expected import origin XML取込, got %q", record.Origin)
		}
	}
}

func TestMasterDictionaryImportServiceAssignsFixedImportTimestamp(t *testing.T) {
	result := runSuccessfulImport(t)

	for _, record := range result.repo.upsertRecords {
		if !record.UpdatedAt.Equal(fixedMasterDictionaryNow()) {
			t.Fatalf("expected import timestamp %s, got %s", fixedMasterDictionaryNow(), record.UpdatedAt)
		}
	}
}

func TestMasterDictionaryImportServicePublishesExpectedProgressValues(t *testing.T) {
	result := runSuccessfulImport(t)

	if len(result.runtime.values) != len(importServiceExpectedProgress) {
		t.Fatalf("expected %d progress events, got %d: %v", len(importServiceExpectedProgress), len(result.runtime.values), result.runtime.values)
	}
	for index, progress := range importServiceExpectedProgress {
		if result.runtime.values[index] != progress {
			t.Fatalf("expected progress[%d]=%d, got %d", index, progress, result.runtime.values[index])
		}
	}
}

func TestMasterDictionaryImportServicePropagatesReadFailure(t *testing.T) {
	repo := &repositoryStub{}
	xmlFiles := &xmlFilePortStub{
		resolvePathFunc: func(rawPath string) (string, error) {
			return rawPath, nil
		},
		openFunc: func(_ string) (io.ReadCloser, error) {
			return io.NopCloser(strings.NewReader("unused")), nil
		},
	}
	xmlRecords := &xmlRecordReaderStub{
		countStringRecordsFunc: func(_ io.Reader) (int, error) {
			return 1, nil
		},
		readStringRecordsFunc: func(_ io.Reader, _ func(xmlStringRecord) error) error {
			return fmt.Errorf("boom")
		},
	}
	service := NewMasterDictionaryImportService(repo, xmlFiles, xmlRecords, nil, fixedMasterDictionaryNow)

	_, err := service.ImportXML(context.Background(), "broken.xml")
	if err == nil {
		t.Fatal("expected import to fail")
	}
	if !strings.Contains(err.Error(), "read xml records") {
		t.Fatalf("expected read xml records context, got %v", err)
	}
	if len(repo.upsertRecords) != 0 {
		t.Fatalf("expected no upsert on read failure, got %d", len(repo.upsertRecords))
	}
}

func runSuccessfulImport(t *testing.T) importServiceSuccessResult {
	t.Helper()

	repo := newImportRepositoryStub()
	xmlFiles := newImportXMLFileStub()
	xmlRecords := newImportXMLRecordReaderStub(importServiceRecords)
	runtime := &importProgressRecorder{}
	service := NewMasterDictionaryImportService(repo, xmlFiles, xmlRecords, runtime, fixedMasterDictionaryNow)

	summary, err := service.ImportXML(context.Background(), " input/import.xml ")
	if err != nil {
		t.Fatalf("expected import to succeed: %v", err)
	}

	return importServiceSuccessResult{
		repo:       repo,
		xmlFiles:   xmlFiles,
		xmlRecords: xmlRecords,
		runtime:    runtime,
		summary:    summary,
	}
}

func newImportRepositoryStub() *repositoryStub {
	repo := &repositoryStub{}
	repo.upsertBySourceAndRECFunc = func(_ context.Context, record MasterDictionaryImportRecord) (MasterDictionaryEntry, bool, error) {
		entry := MasterDictionaryEntry{
			ID:          int64(len(repo.upsertRecords)),
			Source:      record.Source,
			Translation: record.Translation,
			Category:    record.Category,
			Origin:      record.Origin,
			REC:         record.REC,
			EDID:        record.EDID,
			UpdatedAt:   record.UpdatedAt,
		}
		return entry, record.Source == "Snow Elf History", nil
	}
	return repo
}

func newImportXMLFileStub() *xmlFilePortStub {
	return &xmlFilePortStub{
		resolvePathFunc: func(_ string) (string, error) {
			return "/resolved/import.xml", nil
		},
		openFunc: func(_ string) (io.ReadCloser, error) {
			return io.NopCloser(strings.NewReader("unused")), nil
		},
	}
}

func newImportXMLRecordReaderStub(records []xmlStringRecord) *xmlRecordReaderStub {
	return &xmlRecordReaderStub{
		countStringRecordsFunc: func(_ io.Reader) (int, error) {
			return len(records), nil
		},
		readStringRecordsFunc: func(_ io.Reader, handle func(xmlStringRecord) error) error {
			for _, record := range records {
				if err := handle(record); err != nil {
					return err
				}
			}
			return nil
		},
	}
}
