package service

import (
	"path/filepath"
	"testing"
)

const xmlAdapterFixtureBaseName = "Dawnguard_english_japanese.xml"

func TestLocalMasterDictionaryXMLFilePortResolvesRepositoryFixtureByBaseName(t *testing.T) {
	port := NewLocalMasterDictionaryXMLFilePort()

	resolvedPath, err := port.ResolvePath(xmlAdapterFixtureBaseName)
	if err != nil {
		t.Fatalf("expected repository fixture to resolve by base name: %v", err)
	}

	if filepath.Base(resolvedPath) != xmlAdapterFixtureBaseName {
		t.Fatalf("expected resolved base name, got %q", filepath.Base(resolvedPath))
	}
}

func TestLocalMasterDictionaryXMLFilePortOpensResolvedFixture(t *testing.T) {
	port := NewLocalMasterDictionaryXMLFilePort()
	resolvedPath, err := port.ResolvePath(xmlAdapterFixtureBaseName)
	if err != nil {
		t.Fatalf("expected repository fixture to resolve by base name: %v", err)
	}

	reader, err := port.Open(resolvedPath)
	if err != nil {
		t.Fatalf("expected resolved fixture to open: %v", err)
	}
	defer func() {
		_ = reader.Close()
	}()

	recordReader := NewXMLDecoderMasterDictionaryRecordReader()
	count, err := recordReader.CountStringRecords(reader)
	if err != nil {
		t.Fatalf("expected string record count to succeed: %v", err)
	}

	if count == 0 {
		t.Fatal("expected repository fixture to contain XML string records")
	}
}
