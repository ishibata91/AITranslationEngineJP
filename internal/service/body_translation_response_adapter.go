package service

import (
	"context"
	"fmt"
	"reflect"
	"strings"
)

type bodyTranslationProviderClientResponse struct {
	Items         []bodyTranslationProviderClientItem
	ExecutionMode string
	PromptDigest  string
}

type bodyTranslationProviderClientItem struct {
	RequestUnitID       string
	FieldCorrelationKey string
	TranslatedText      string
}

type bodyTranslationFailureMetadata interface {
	error
	FailureKind() string
	FailureRetryable() bool
}

func invokeBodyTranslationClientGenerateBodyTranslation(
	ctx context.Context,
	client any,
	providerID string,
	model string,
	executionMode string,
	credentialRef string,
	endpointSummary string,
	prompt string,
) (bodyTranslationProviderClientResponse, error) {
	if client == nil {
		return bodyTranslationProviderClientResponse{}, fmt.Errorf("body translation provider client is required")
	}
	method := reflectValueOfClientMethod(client, "GenerateBodyTranslation")
	if !method.IsValid() {
		return bodyTranslationProviderClientResponse{}, fmt.Errorf("body translation provider client does not implement GenerateBodyTranslation")
	}
	if (method.Type().NumIn() != 6 && method.Type().NumIn() != 7) || method.Type().NumOut() != 2 {
		return bodyTranslationProviderClientResponse{}, fmt.Errorf("body translation provider client has incompatible GenerateBodyTranslation signature")
	}
	args := []reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(providerID),
		reflect.ValueOf(model),
		reflect.ValueOf(executionMode),
		reflect.ValueOf(credentialRef),
	}
	if method.Type().NumIn() == 7 {
		args = append(args, reflect.ValueOf(strings.TrimSpace(endpointSummary)))
	}
	args = append(args, reflect.ValueOf(prompt))
	results := method.Call(args)
	if errValue := results[1]; !errValue.IsNil() {
		err, _ := errValue.Interface().(error)
		return bodyTranslationProviderClientResponse{}, err
	}
	return mapBodyTranslationClientResponse(results[0])
}

func mapBodyTranslationClientResponse(value reflect.Value) (bodyTranslationProviderClientResponse, error) {
	value, err := bodyTranslationStructValue(value, "body translation provider response")
	if err != nil {
		return bodyTranslationProviderClientResponse{}, err
	}
	itemsField := value.FieldByName("Items")
	executionModeField := value.FieldByName("ExecutionMode")
	promptDigestField := value.FieldByName("PromptDigest")
	if !itemsField.IsValid() || itemsField.Kind() != reflect.Slice ||
		!executionModeField.IsValid() || executionModeField.Kind() != reflect.String ||
		!promptDigestField.IsValid() || promptDigestField.Kind() != reflect.String {
		return bodyTranslationProviderClientResponse{}, fmt.Errorf("body translation provider response has incompatible shape")
	}

	items := make([]bodyTranslationProviderClientItem, 0, itemsField.Len())
	for index := 0; index < itemsField.Len(); index++ {
		item, itemErr := mapBodyTranslationClientItem(itemsField.Index(index))
		if itemErr != nil {
			return bodyTranslationProviderClientResponse{}, itemErr
		}
		items = append(items, item)
	}

	return bodyTranslationProviderClientResponse{
		Items:         items,
		ExecutionMode: executionModeField.String(),
		PromptDigest:  promptDigestField.String(),
	}, nil
}

func mapBodyTranslationClientItem(value reflect.Value) (bodyTranslationProviderClientItem, error) {
	value, err := bodyTranslationStructValue(value, "body translation provider item")
	if err != nil {
		return bodyTranslationProviderClientItem{}, err
	}
	requestUnitIDField := value.FieldByName("RequestUnitID")
	fieldCorrelationKeyField := value.FieldByName("FieldCorrelationKey")
	translatedTextField := value.FieldByName("TranslatedText")
	if !requestUnitIDField.IsValid() || requestUnitIDField.Kind() != reflect.String ||
		!fieldCorrelationKeyField.IsValid() || fieldCorrelationKeyField.Kind() != reflect.String ||
		!translatedTextField.IsValid() || translatedTextField.Kind() != reflect.String {
		return bodyTranslationProviderClientItem{}, fmt.Errorf("body translation provider item has incompatible shape")
	}
	return bodyTranslationProviderClientItem{
		RequestUnitID:       requestUnitIDField.String(),
		FieldCorrelationKey: fieldCorrelationKeyField.String(),
		TranslatedText:      translatedTextField.String(),
	}, nil
}

func mapBodyTranslationProviderResponse(
	baseResult BodyTranslationProviderResult,
	clientResponse bodyTranslationProviderClientResponse,
) BodyTranslationProviderResult {
	if len(clientResponse.Items) != 1 {
		return bodyTranslationProviderFailureResult(
			baseResult,
			BodyTranslationProviderErrorKindInvalidProviderResponse,
			bodyTranslationProviderResponseInvalidReason,
			true,
		)
	}

	item := clientResponse.Items[0]
	if strings.TrimSpace(item.RequestUnitID) != baseResult.RequestUnitID ||
		strings.TrimSpace(item.FieldCorrelationKey) != baseResult.FieldCorrelationKey ||
		strings.TrimSpace(item.TranslatedText) == "" {
		return bodyTranslationProviderFailureResult(
			baseResult,
			BodyTranslationProviderErrorKindInvalidProviderResponse,
			bodyTranslationProviderResponseInvalidReason,
			true,
		)
	}

	translatedText := strings.TrimSpace(item.TranslatedText)
	baseResult.TranslatedCandidate = &BodyTranslationTranslatedCandidate{
		RequestUnitID:       baseResult.RequestUnitID,
		FieldCorrelationKey: baseResult.FieldCorrelationKey,
		RecordType:          baseResult.RecordType,
		FieldType:           baseResult.FieldType,
		TranslatedText:      translatedText,
	}
	baseResult.ProtectionValidationTarget = &BodyTranslationProtectionValidationTarget{
		RequestUnitID:          baseResult.RequestUnitID,
		FieldCorrelationKey:    baseResult.FieldCorrelationKey,
		ProtectionSourceDigest: baseResult.AuditSummary.RequestSummary.ProtectionSourceDigest,
		TranslatedText:         translatedText,
	}
	baseResult.Failure = nil
	return baseResult
}

func bodyTranslationStructValue(value reflect.Value, label string) (reflect.Value, error) {
	if !value.IsValid() {
		return reflect.Value{}, fmt.Errorf("%s is invalid", label)
	}
	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return reflect.Value{}, fmt.Errorf("%s is nil", label)
		}
		value = value.Elem()
	}
	return value, nil
}

func reflectValueOfClientMethod(client any, methodName string) reflect.Value {
	return reflect.ValueOf(client).MethodByName(methodName)
}
