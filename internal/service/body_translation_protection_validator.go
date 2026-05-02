package service

import (
	"fmt"
	"regexp"
)

const (
	bodyTranslationProtectionValidationPassed          = "passed"
	bodyTranslationProtectionValidationFailed          = "failed"
	bodyTranslationProtectionFailureMissing            = "missing"
	bodyTranslationProtectionFailureModified           = "modified"
	bodyTranslationProtectionFailureDuplicate          = "duplicate"
	bodyTranslationProtectionFailureReordered          = "reordered"
	bodyTranslationProtectionFailureExtra              = "extra"
	bodyTranslationProtectionValidationReasonMissing   = "protected element is missing"
	bodyTranslationProtectionValidationReasonModified  = "protected element was modified"
	bodyTranslationProtectionValidationReasonDuplicate = "protected element was duplicated"
	bodyTranslationProtectionValidationReasonReordered = "protected element order changed"
	bodyTranslationProtectionValidationReasonExtra     = "unexpected protected element was added"
)

var bodyTranslationProtectedElementPattern = regexp.MustCompile(`(<[^>\n]+>|\{[^}\n]+\})`)

// BodyTranslationProtectionValidationResult summarizes one validation outcome without raw provider payload.
type BodyTranslationProtectionValidationResult struct {
	Status      string
	FailureKind string
	Summary     string
}

// ValidateBodyTranslationProtection checks whether translated text preserves protected elements exactly.
func ValidateBodyTranslationProtection(
	target *BodyTranslationProtectionValidationTarget,
	expectedElements []BodyTranslationProtectedElement,
) (BodyTranslationProtectionValidationResult, error) {
	if target == nil {
		return BodyTranslationProtectionValidationResult{}, fmt.Errorf("body translation protection validation target is required")
	}
	normalizedElements, err := normalizeBodyTranslationProtectedElements(expectedElements)
	if err != nil {
		return BodyTranslationProtectionValidationResult{}, err
	}

	expectedTokens := make([]string, 0, len(normalizedElements))
	for _, element := range normalizedElements {
		expectedTokens = append(expectedTokens, element.SourceText)
	}
	actualTokens := bodyTranslationProtectedElementPattern.FindAllString(target.TranslatedText, -1)

	if len(expectedTokens) == 0 && len(actualTokens) == 0 {
		return BodyTranslationProtectionValidationResult{Status: bodyTranslationProtectionValidationPassed}, nil
	}
	if len(actualTokens) < len(expectedTokens) {
		return BodyTranslationProtectionValidationResult{
			Status:      bodyTranslationProtectionValidationFailed,
			FailureKind: bodyTranslationProtectionFailureMissing,
			Summary:     bodyTranslationProtectionValidationReasonMissing,
		}, nil
	}
	if len(actualTokens) > len(expectedTokens) {
		if bodyTranslationContainsDuplicateToken(actualTokens, expectedTokens) {
			return BodyTranslationProtectionValidationResult{
				Status:      bodyTranslationProtectionValidationFailed,
				FailureKind: bodyTranslationProtectionFailureDuplicate,
				Summary:     bodyTranslationProtectionValidationReasonDuplicate,
			}, nil
		}
		return BodyTranslationProtectionValidationResult{
			Status:      bodyTranslationProtectionValidationFailed,
			FailureKind: bodyTranslationProtectionFailureExtra,
			Summary:     bodyTranslationProtectionValidationReasonExtra,
		}, nil
	}
	if bodyTranslationStringSlicesEqual(expectedTokens, actualTokens) {
		return BodyTranslationProtectionValidationResult{Status: bodyTranslationProtectionValidationPassed}, nil
	}
	if bodyTranslationContainsSameTokenMultiset(expectedTokens, actualTokens) {
		return BodyTranslationProtectionValidationResult{
			Status:      bodyTranslationProtectionValidationFailed,
			FailureKind: bodyTranslationProtectionFailureReordered,
			Summary:     bodyTranslationProtectionValidationReasonReordered,
		}, nil
	}
	return BodyTranslationProtectionValidationResult{
		Status:      bodyTranslationProtectionValidationFailed,
		FailureKind: bodyTranslationProtectionFailureModified,
		Summary:     bodyTranslationProtectionValidationReasonModified,
	}, nil
}

func bodyTranslationContainsDuplicateToken(actual []string, expected []string) bool {
	actualCounts := bodyTranslationTokenCounts(actual)
	expectedCounts := bodyTranslationTokenCounts(expected)
	for token, count := range actualCounts {
		if count > expectedCounts[token] {
			return true
		}
	}
	return false
}

func bodyTranslationContainsSameTokenMultiset(left []string, right []string) bool {
	leftCounts := bodyTranslationTokenCounts(left)
	rightCounts := bodyTranslationTokenCounts(right)
	if len(leftCounts) != len(rightCounts) {
		return false
	}
	for token, leftCount := range leftCounts {
		if rightCounts[token] != leftCount {
			return false
		}
	}
	return true
}

func bodyTranslationTokenCounts(values []string) map[string]int {
	counts := make(map[string]int, len(values))
	for _, value := range values {
		counts[value]++
	}
	return counts
}

func bodyTranslationStringSlicesEqual(left []string, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}
