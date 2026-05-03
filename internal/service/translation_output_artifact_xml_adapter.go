package service

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type translationOutputArtifactXMLSerializer interface {
	Serialize(targetGame string, rows []xTranslatorArtifactRow) ([]byte, error)
}

type translationOutputArtifactFileWriter interface {
	WriteFile(path string, payload []byte) error
}

type translationOutputArtifactTransactionalFileWriter interface {
	WriteTemporaryFile(path string, payload []byte) (string, error)
	PublishTemporaryFile(tempPath string, finalPath string) error
	RemoveFile(path string) error
}

type xTranslatorOutputArtifactXMLSerializer struct{}

type localTranslationOutputArtifactFileWriter struct{}

type xTranslatorOutputArtifactXMLDocument struct {
	XMLName xml.Name
	Strings []xTranslatorOutputArtifactXMLRow `xml:"String"`
}

type xTranslatorOutputArtifactXMLRow struct {
	EDID   string `xml:"EDID"`
	REC    string `xml:"REC"`
	FIELD  string `xml:"FIELD"`
	FORMID string `xml:"FORMID"`
	Source string `xml:"Source"`
	Dest   string `xml:"Dest"`
	Status int    `xml:"Status"`
}

// NewXTranslatorOutputArtifactXMLSerializer returns the default xTranslator XML serializer.
func NewXTranslatorOutputArtifactXMLSerializer() translationOutputArtifactXMLSerializer {
	return xTranslatorOutputArtifactXMLSerializer{}
}

// NewLocalTranslationOutputArtifactFileWriter returns the local filesystem writer.
func NewLocalTranslationOutputArtifactFileWriter() translationOutputArtifactFileWriter {
	return localTranslationOutputArtifactFileWriter{}
}

func (xTranslatorOutputArtifactXMLSerializer) Serialize(
	targetGame string,
	rows []xTranslatorArtifactRow,
) ([]byte, error) {
	rootName, err := xTranslatorRootElementName(targetGame)
	if err != nil {
		return nil, err
	}
	document := xTranslatorOutputArtifactXMLDocument{
		XMLName: xml.Name{Local: rootName},
		Strings: make([]xTranslatorOutputArtifactXMLRow, 0, len(rows)),
	}
	for _, row := range rows {
		document.Strings = append(document.Strings, xTranslatorOutputArtifactXMLRow{
			EDID:   row.EDID,
			REC:    row.REC,
			FIELD:  row.FIELD,
			FORMID: row.FORMID,
			Source: row.Source,
			Dest:   row.Dest,
			Status: row.Status,
		})
	}

	payload, err := xml.MarshalIndent(document, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal xtranslator xml: %w", err)
	}
	return bytes.Join([][]byte{[]byte(xml.Header), payload, []byte("\n")}, nil), nil
}

func validateTranslationOutputArtifactPath(path string) (string, error) {
	trimmed := strings.TrimSpace(path)
	if trimmed == "" {
		return "", fmt.Errorf("output path is required")
	}
	if strings.HasSuffix(trimmed, string(os.PathSeparator)) {
		return "", fmt.Errorf("output path must be a file")
	}
	cleaned := filepath.Clean(trimmed)
	if !filepath.IsAbs(cleaned) {
		return "", fmt.Errorf("output path must be absolute")
	}
	if !strings.EqualFold(filepath.Ext(cleaned), ".xml") {
		return "", fmt.Errorf("output path must end with .xml")
	}
	if isReadonlyTranslationOutputPath(cleaned) {
		return "", fmt.Errorf("output path is not writable")
	}
	if info, err := os.Stat(cleaned); err == nil && info.IsDir() {
		return "", fmt.Errorf("output path points to a directory")
	}
	return cleaned, nil
}

func (localTranslationOutputArtifactFileWriter) WriteTemporaryFile(path string, payload []byte) (string, error) {
	trimmed, pathErr := validateTranslationOutputArtifactPath(path)
	if pathErr != nil {
		return "", pathErr
	}
	directory := filepath.Dir(trimmed)
	if err := os.MkdirAll(directory, 0o750); err != nil {
		return "", fmt.Errorf("create output artifact directory: %w", err)
	}
	tempFile, err := os.CreateTemp(directory, ".translation-output-artifact-*.xml")
	if err != nil {
		return "", fmt.Errorf("create output artifact temp file: %w", err)
	}
	tempPath := tempFile.Name()
	if _, err := tempFile.Write(payload); err != nil {
		_ = tempFile.Close()
		_ = os.Remove(tempPath)
		return "", fmt.Errorf("write output artifact temp file: %w", err)
	}
	if err := tempFile.Close(); err != nil {
		_ = os.Remove(tempPath)
		return "", fmt.Errorf("close output artifact temp file: %w", err)
	}
	return tempPath, nil
}

func (writer localTranslationOutputArtifactFileWriter) WriteFile(path string, payload []byte) error {
	tempPath, err := writer.WriteTemporaryFile(path, payload)
	if err != nil {
		return err
	}
	if err := writer.PublishTemporaryFile(tempPath, path); err != nil {
		_ = writer.RemoveFile(tempPath)
		return err
	}
	return nil
}

func (localTranslationOutputArtifactFileWriter) PublishTemporaryFile(tempPath string, finalPath string) error {
	if err := os.Rename(tempPath, finalPath); err != nil {
		return fmt.Errorf("publish output artifact file: %w", err)
	}
	return nil
}

func (localTranslationOutputArtifactFileWriter) RemoveFile(path string) error {
	if strings.TrimSpace(path) == "" {
		return nil
	}
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("remove output artifact file: %w", err)
	}
	return nil
}

func isReadonlyTranslationOutputPath(path string) bool {
	for _, prefix := range []string{"/System/", "/usr/", "/bin/", "/sbin/", "/etc/", "/private/etc/"} {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}

func xTranslatorRootElementName(targetGame string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(targetGame)) {
	case "skyrim_se", "skyrimse", "sse":
		return "SSETranslator", nil
	case "skyrim_le", "skyrimle", "tesv", "skyrim":
		return "TESVTranslator", nil
	default:
		return "", fmt.Errorf("unsupported target game %q", targetGame)
	}
}
