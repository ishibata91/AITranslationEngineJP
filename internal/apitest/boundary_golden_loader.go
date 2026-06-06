// Package apitest provides shared test infrastructure for Wails controller boundary tests.
// This file provides the boundary golden loader used by both backend boundary tests
// and as the canonical path reference for frontend boundary tests.
package apitest

import (
	"embed"
	"encoding/json"
	"fmt"
)

// goldenFS embeds all golden JSON files under testdata/boundary/master_dictionary/.
//
//go:embed testdata/boundary/master_dictionary/*.golden.json
var goldenFS embed.FS

// BoundaryGoldenInput describes the input side of a golden case.
type BoundaryGoldenInput struct {
	Description string          `json:"description"`
	Method      string          `json:"method"`
	Case        string          `json:"case"`
	Request     json.RawMessage `json:"request"`
}

// BoundaryGoldenFile is the top-level structure of a single golden JSON file.
// Input holds the input-side description; Expected holds the expected DTO field values
// as a raw JSON value so each test can decode it into the appropriate response DTO type.
type BoundaryGoldenFile struct {
	Input    BoundaryGoldenInput `json:"input"`
	Expected json.RawMessage     `json:"expected"`
}

// LoadBoundaryGolden reads the named golden file from the embedded FS and returns the parsed
// BoundaryGoldenFile. The name must be the bare file name (e.g. "list_normal.golden.json").
// The Expected field is left as raw JSON; callers decode it into the appropriate type.
func LoadBoundaryGolden(name string) (BoundaryGoldenFile, error) {
	path := "testdata/boundary/master_dictionary/" + name
	data, err := goldenFS.ReadFile(path)
	if err != nil {
		return BoundaryGoldenFile{}, fmt.Errorf("load boundary golden %q: %w", name, err)
	}

	var gf BoundaryGoldenFile
	if err := json.Unmarshal(data, &gf); err != nil {
		return BoundaryGoldenFile{}, fmt.Errorf("parse boundary golden %q: %w", name, err)
	}
	return gf, nil
}

// MustLoadBoundaryGolden is like LoadBoundaryGolden but panics on error.
// Use only in test helpers where a missing golden file is always a programming error.
func MustLoadBoundaryGolden(name string) BoundaryGoldenFile {
	gf, err := LoadBoundaryGolden(name)
	if err != nil {
		panic(err)
	}
	return gf
}
