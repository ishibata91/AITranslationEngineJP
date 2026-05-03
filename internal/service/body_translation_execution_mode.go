package service

import (
	"fmt"
	"strings"
)

const bodyTranslationExecutionModeSync = "sync"

func normalizeBodyTranslationExecutionMode(executionMode string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(executionMode))
	switch normalized {
	case "", BodyTranslationExecutionModeSingleRequest, bodyTranslationExecutionModeSync:
		return BodyTranslationExecutionModeSingleRequest, nil
	default:
		return "", fmt.Errorf("unsupported body translation execution mode: %s", executionMode)
	}
}
