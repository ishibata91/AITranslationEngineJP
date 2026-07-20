package provider

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// extractBatchTranslations は返ってきたキーのうち、translation が非空のものだけを載せること（部分成功）。
// 欠けたキー・translation 欠落・空文字・型不正は載せないこと。
func TestExtractBatchTranslationsPartialSuccess(t *testing.T) {
	keys := []string{"1", "2", "3", "4"}
	// 1=正常, 2=空文字(壊れ), 3=translation欠落(壊れ), 4=応答に無い(欠落)
	content := `{"1":{"translation":"訳1"},"2":{"translation":""},"3":{"foo":"bar"}}`
	got := extractBatchTranslations(content, keys)
	if len(got) != 1 || got["1"] != "訳1" {
		t.Fatalf("got = %v, want only {1:訳1}", got)
	}
}

// 構文不正（空応答・途中切れ）は空 map を返し、全キーを欠落扱いにすること。
func TestExtractBatchTranslationsMalformedIsEmpty(t *testing.T) {
	for _, content := range []string{"", "{", "not json", `{"1":{"translation":"x"`} {
		got := extractBatchTranslations(content, []string{"1"})
		if len(got) != 0 {
			t.Errorf("content=%q got=%v, want empty", content, got)
		}
	}
}

// batchResponseFormat は要求キーごとのプロパティを持ち、余分なキーを禁じ、部分成功のため strict=false であること。
func TestBatchResponseFormat(t *testing.T) {
	rf := batchResponseFormat([]string{"10", "20"})
	if rf.Type != "json_schema" {
		t.Errorf("type = %q, want json_schema", rf.Type)
	}
	if rf.JSONSchema.Strict {
		t.Errorf("strict = true, want false（部分成功を許すため）")
	}
	var schema struct {
		Type                 string                     `json:"type"`
		Properties           map[string]json.RawMessage `json:"properties"`
		AdditionalProperties bool                       `json:"additionalProperties"`
	}
	if err := json.Unmarshal(rf.JSONSchema.Schema, &schema); err != nil {
		t.Fatalf("schema unmarshal: %v", err)
	}
	if schema.Type != "object" || schema.AdditionalProperties {
		t.Errorf("schema type=%q additionalProperties=%v, want object/false", schema.Type, schema.AdditionalProperties)
	}
	if _, ok := schema.Properties["10"]; !ok {
		t.Errorf("schema に key 10 が無い: %v", schema.Properties)
	}
	if _, ok := schema.Properties["20"]; !ok {
		t.Errorf("schema に key 20 が無い: %v", schema.Properties)
	}
}

// TranslateBatch は system メッセージと行キー付きの user メッセージを送り、返ったキーの訳文を map で返すこと。
// 欠けたキーは map に含めないこと（部分成功、再送しない）。
func TestTranslateBatchSendsKeyedAndParsesPartial(t *testing.T) {
	var gotSystem, gotUser string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req struct {
			Messages []struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"messages"`
		}
		_ = json.Unmarshal(body, &req)
		for _, m := range req.Messages {
			switch m.Role {
			case "system":
				gotSystem = m.Content
			case "user":
				gotUser = m.Content
			}
		}
		// key 7 だけ返し、key 8 は落とす（部分成功）。
		_, _ = io.WriteString(w, `{"choices":[{"message":{"content":"{\"7\":{\"translation\":\"訳7\"}}"}}]}`)
	}))
	defer srv.Close()

	client := NewOpenAICompatible(http.DefaultClient)
	got, err := client.TranslateBatch(context.Background(), Connection{Endpoint: srv.URL}, "m", "共有 system 指示", []BatchItem{
		{Key: "7", User: "seven"},
		{Key: "8", User: "eight"},
	})
	if err != nil {
		t.Fatalf("TranslateBatch error: %v", err)
	}
	if len(got) != 1 || got["7"] != "訳7" {
		t.Fatalf("got = %v, want {7:訳7}", got)
	}
	if gotSystem != "共有 system 指示" {
		t.Errorf("system = %q, want 共有 system 指示", gotSystem)
	}
	// user メッセージへ各行がキー付きで載ること。
	if !strings.Contains(gotUser, "7: seven") || !strings.Contains(gotUser, "8: eight") {
		t.Errorf("user = %q, want キー付き原文", gotUser)
	}
}

// 通信・HTTP 障害は error で返すこと（部分成功と区別する）。
func TestTranslateBatchHTTPErrorReturnsError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	client := NewOpenAICompatible(http.DefaultClient)
	_, err := client.TranslateBatch(context.Background(), Connection{Endpoint: srv.URL}, "m", "s", []BatchItem{{Key: "1", User: "a"}})
	if err == nil {
		t.Fatalf("want error on HTTP 500, got nil")
	}
}
