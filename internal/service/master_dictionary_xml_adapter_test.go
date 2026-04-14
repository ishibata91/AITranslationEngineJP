package service

import (
	"path/filepath"
	"testing"
)

func TestLocalMasterDictionaryXMLFilePortResolvesRepositoryFixtureByBaseName(t *testing.T) {
	port := NewLocalMasterDictionaryXMLFilePort()
	resolvedPath, err := port.ResolvePath("Dawnguard_english_japanese.xml")
	if err != nil {
		t.Fatalf("expected repository fixture to resolve by base name: %v", err)
	}
	if filepath.Base(resolvedPath) != "Dawnguard_english_japanese.xml" {
		t.Fatalf("expected resolved base name, got %q", filepath.Base(resolvedPath))
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
