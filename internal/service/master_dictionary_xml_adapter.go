package service

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const masterDictionaryOpenXMLFileError = "open xml file: %w"

// LocalMasterDictionaryXMLFilePort resolves XML fixtures from the current workspace.
type LocalMasterDictionaryXMLFilePort struct{}

// NewLocalMasterDictionaryXMLFilePort creates the default XML file adapter.
func NewLocalMasterDictionaryXMLFilePort() *LocalMasterDictionaryXMLFilePort {
	return &LocalMasterDictionaryXMLFilePort{}
}

// ResolvePath validates and resolves the requested XML path.
func (*LocalMasterDictionaryXMLFilePort) ResolvePath(rawPath string) (string, error) {
	trimmedPath := strings.TrimSpace(rawPath)
	if trimmedPath == "" {
		return "", fmt.Errorf("%w: xml path is required", ErrMasterDictionaryValidation)
	}

	for _, candidate := range masterDictionaryXMLPathCandidates(trimmedPath) {
		info, statErr := os.Stat(candidate)
		if statErr != nil {
			if errors.Is(statErr, os.ErrNotExist) {
				continue
			}
			return "", fmt.Errorf("stat xml file: %w", statErr)
		}
		if info.IsDir() {
			continue
		}
		if !strings.EqualFold(filepath.Ext(candidate), ".xml") {
			return "", fmt.Errorf("%w: xml file extension is required", ErrMasterDictionaryValidation)
		}
		return candidate, nil
	}

	return "", fmt.Errorf(masterDictionaryOpenXMLFileError, os.ErrNotExist)
}

// Open opens the resolved XML file.
func (*LocalMasterDictionaryXMLFilePort) Open(path string) (io.ReadCloser, error) {
	file, err := os.Open(path) // #nosec G304 -- path is validated by ResolvePath.
	if err != nil {
		return nil, fmt.Errorf(masterDictionaryOpenXMLFileError, err)
	}
	return file, nil
}

// XMLDecoderMasterDictionaryRecordReader decodes XML string records with encoding/xml.
type XMLDecoderMasterDictionaryRecordReader struct{}

// NewXMLDecoderMasterDictionaryRecordReader creates the default XML record adapter.
func NewXMLDecoderMasterDictionaryRecordReader() *XMLDecoderMasterDictionaryRecordReader {
	return &XMLDecoderMasterDictionaryRecordReader{}
}

// CountStringRecords counts <String> records in the XML stream.
func (*XMLDecoderMasterDictionaryRecordReader) CountStringRecords(reader io.Reader) (int, error) {
	decoder := xml.NewDecoder(reader)
	total := 0
	for {
		token, tokenErr := decoder.Token()
		if tokenErr != nil {
			if errors.Is(tokenErr, io.EOF) {
				break
			}
			return 0, fmt.Errorf("read xml token: %w", tokenErr)
		}

		startElement, ok := token.(xml.StartElement)
		if ok && startElement.Name.Local == "String" {
			total++
		}
	}
	return total, nil
}

// ReadStringRecords streams each <String> record to the provided callback.
func (*XMLDecoderMasterDictionaryRecordReader) ReadStringRecords(
	reader io.Reader,
	handle func(xmlStringRecord) error,
) error {
	decoder := xml.NewDecoder(reader)
	for {
		token, err := decoder.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return fmt.Errorf("read xml token: %w", err)
		}

		start, ok := token.(xml.StartElement)
		if !ok || start.Name.Local != "String" {
			continue
		}

		record := xmlStringRecord{}
		if err := decoder.DecodeElement(&record, &start); err != nil {
			return fmt.Errorf("decode xml string record: %w", err)
		}
		if err := handle(record); err != nil {
			return err
		}
	}
}

func masterDictionaryXMLPathCandidates(rawPath string) []string {
	cleanedPath := filepath.Clean(rawPath)
	baseName := filepath.Base(cleanedPath)

	candidates := make([]string, 0, 8)
	seen := make(map[string]struct{}, 8)
	appendCandidate := func(path string) {
		if strings.TrimSpace(path) == "" {
			return
		}
		cleaned := filepath.Clean(path)
		if _, exists := seen[cleaned]; exists {
			return
		}
		seen[cleaned] = struct{}{}
		candidates = append(candidates, cleaned)
	}

	appendCandidate(cleanedPath)
	appendCandidate(baseName)

	if cwd, err := os.Getwd(); err == nil {
		directory := cwd
		for depth := 0; depth < 6; depth++ {
			appendCandidate(filepath.Join(directory, "dictionaries", baseName))
			parent := filepath.Dir(directory)
			if parent == directory {
				break
			}
			directory = parent
		}
	}

	return candidates
}

func baseName(path string) string {
	return filepath.Base(path)
}
