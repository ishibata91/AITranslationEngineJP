package provider

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// SubmitBatch が JSONL を Files API へ upload し、file を参照する batch を作って外部 batch ID を返すこと。
// JSONL 各行に custom_id / method / url / body（同期と同一の chat completions payload、response_format 込み）
// が載り、purpose=batch で upload し、input_file_id で batch を作ることを確かめる。
func TestXAIBatchSubmitUploadsJSONLAndCreatesBatch(t *testing.T) {
	var uploadedJSONL string
	var uploadedPurpose string
	var createBody map[string]string
	var uploadAuth, createAuth string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/v1/files" && r.Method == http.MethodPost:
			uploadAuth = r.Header.Get("Authorization")
			if err := r.ParseMultipartForm(1 << 20); err != nil {
				t.Errorf("ParseMultipartForm: %v", err)
			}
			uploadedPurpose = r.FormValue("purpose")
			f, _, err := r.FormFile("file")
			if err != nil {
				t.Errorf("FormFile: %v", err)
			} else {
				b, _ := io.ReadAll(f)
				uploadedJSONL = string(b)
			}
			_, _ = io.WriteString(w, `{"id":"file-xyz"}`)
		case r.URL.Path == "/v1/batches" && r.Method == http.MethodPost:
			createAuth = r.Header.Get("Authorization")
			b, _ := io.ReadAll(r.Body)
			_ = json.Unmarshal(b, &createBody)
			_, _ = io.WriteString(w, `{"batch_id":"batch-123"}`)
		default:
			t.Errorf("想定外の呼び出し: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	c := NewXAIBatch(http.DefaultClient)
	id, err := c.SubmitBatch(context.Background(), Connection{Endpoint: srv.URL, APIKey: "k"},
		"grok-4", []BatchRequest{
			{CustomID: "n:1", Prompt: Prompt{System: "sys-a", User: "usr-a"}},
			{CustomID: "l:2", Prompt: Prompt{System: "sys-b", User: "usr-b"}},
		})
	if err != nil {
		t.Fatalf("SubmitBatch error: %v", err)
	}
	if id != "batch-123" {
		t.Errorf("外部 batch ID = %q, want batch-123", id)
	}
	if uploadedPurpose != "batch" {
		t.Errorf("purpose = %q, want batch", uploadedPurpose)
	}
	if uploadAuth != "Bearer k" || createAuth != "Bearer k" {
		t.Errorf("Authorization: upload=%q create=%q, want Bearer k", uploadAuth, createAuth)
	}
	if createBody["input_file_id"] != "file-xyz" {
		t.Errorf("input_file_id = %q, want file-xyz", createBody["input_file_id"])
	}

	// JSONL は 2 行。各行の custom_id / url / body（model・messages・response_format）を確かめる。
	lines := strings.Split(strings.TrimSpace(uploadedJSONL), "\n")
	if len(lines) != 2 {
		t.Fatalf("JSONL 行数 = %d, want 2", len(lines))
	}
	var first struct {
		CustomID string `json:"custom_id"`
		Method   string `json:"method"`
		URL      string `json:"url"`
		Body     struct {
			Model    string `json:"model"`
			Messages []struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"messages"`
			ResponseFormat *struct {
				Type       string `json:"type"`
				JSONSchema struct {
					Name   string `json:"name"`
					Strict bool   `json:"strict"`
				} `json:"json_schema"`
			} `json:"response_format"`
		} `json:"body"`
	}
	if err := json.Unmarshal([]byte(lines[0]), &first); err != nil {
		t.Fatalf("1 行目の JSONL 解析: %v", err)
	}
	if first.CustomID != "n:1" || first.Method != http.MethodPost || first.URL != "/v1/chat/completions" {
		t.Errorf("1 行目 custom_id/method/url = %q/%q/%q", first.CustomID, first.Method, first.URL)
	}
	if first.Body.Model != "grok-4" {
		t.Errorf("body.model = %q, want grok-4", first.Body.Model)
	}
	if len(first.Body.Messages) != 2 || first.Body.Messages[0].Content != "sys-a" || first.Body.Messages[1].Content != "usr-a" {
		t.Errorf("body.messages = %+v", first.Body.Messages)
	}
	if first.Body.ResponseFormat == nil || first.Body.ResponseFormat.Type != "json_schema" ||
		!first.Body.ResponseFormat.JSONSchema.Strict {
		t.Errorf("body.response_format = %+v, want json_schema strict", first.Body.ResponseFormat)
	}
}

// PollBatch が state の件数を読み、未処理が 0 のときだけ Done=true にすること。
func TestXAIBatchPollDoneWhenNoPending(t *testing.T) {
	var state string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/batches/batch-123" || r.Method != http.MethodGet {
			t.Errorf("想定外の呼び出し: %s %s", r.Method, r.URL.Path)
		}
		_, _ = io.WriteString(w, state)
	}))
	defer srv.Close()
	c := NewXAIBatch(http.DefaultClient)
	conn := Connection{Endpoint: srv.URL, APIKey: "k"}

	state = `{"state":{"num_requests":10,"num_pending":4,"num_success":5,"num_error":1,"num_cancelled":0}}`
	st, err := c.PollBatch(context.Background(), conn, "batch-123")
	if err != nil {
		t.Fatalf("PollBatch error: %v", err)
	}
	if st.Done {
		t.Errorf("未処理 4 件で Done=true になった: %+v", st)
	}
	if st.Total != 10 || st.Pending != 4 || st.Succeeded != 5 || st.Failed != 1 {
		t.Errorf("件数 = %+v", st)
	}

	state = `{"state":{"num_requests":10,"num_pending":0,"num_success":9,"num_error":1,"num_cancelled":0}}`
	st, err = c.PollBatch(context.Background(), conn, "batch-123")
	if err != nil {
		t.Fatalf("PollBatch error: %v", err)
	}
	if !st.Done {
		t.Errorf("未処理 0 件で Done=false になった: %+v", st)
	}
}

// FetchResults がページングを最後までたどり、custom_id ごとの訳文と失敗種別へ写すこと。
// 成功・error_message あり・choices 空・content 非 JSON の 4 種を、それぞれ正しい結果へ写す。
func TestXAIBatchFetchResultsPaginatesAndMapsFailures(t *testing.T) {
	page1 := `{"results":[
	  {"batch_request_id":"n:1","batch_result":{"response":{"chat_get_completion":{"choices":[{"message":{"content":"{\"translation\":\"訳文A\"}"}}]}}},"error_message":null},
	  {"batch_request_id":"l:2","batch_result":{"response":{"chat_get_completion":{"choices":[]}}},"error_message":"rate limited"}
	],"pagination_token":"tok2"}`
	page2 := `{"results":[
	  {"batch_request_id":"n:3","batch_result":{"response":{"chat_get_completion":{"choices":[]}}},"error_message":null},
	  {"batch_request_id":"p:4","batch_result":{"response":{"chat_get_completion":{"choices":[{"message":{"content":"not json"}}]}}},"error_message":null}
	],"pagination_token":null}`

	var gotTokens []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/batches/batch-123/results" {
			t.Errorf("想定外の path: %s", r.URL.Path)
		}
		tok := r.URL.Query().Get("pagination_token")
		gotTokens = append(gotTokens, tok)
		if tok == "" {
			_, _ = io.WriteString(w, page1)
		} else {
			_, _ = io.WriteString(w, page2)
		}
	}))
	defer srv.Close()

	c := NewXAIBatch(http.DefaultClient)
	results, err := c.FetchResults(context.Background(), Connection{Endpoint: srv.URL, APIKey: "k"}, "batch-123")
	if err != nil {
		t.Fatalf("FetchResults error: %v", err)
	}
	if len(gotTokens) != 2 || gotTokens[0] != "" || gotTokens[1] != "tok2" {
		t.Errorf("ページングの token = %v, want [\"\" tok2]", gotTokens)
	}
	if len(results) != 4 {
		t.Fatalf("結果件数 = %d, want 4", len(results))
	}
	byID := map[string]BatchResult{}
	for _, r := range results {
		byID[r.CustomID] = r
	}
	// 成功: 訳文を取り出す。
	if got := byID["n:1"]; got.Err != nil || got.Translation != "訳文A" {
		t.Errorf("n:1 = %+v, want 訳文A / err nil", got)
	}
	// error_message あり → 一時失敗（ErrServerTransient）。
	if got := byID["l:2"]; !errors.Is(got.Err, ErrServerTransient) {
		t.Errorf("l:2 err = %v, want ErrServerTransient", got.Err)
	}
	// choices 空 → 応答エンベロープの読み取り失敗（ErrResponseUnreadable）。
	if got := byID["n:3"]; !errors.Is(got.Err, ErrResponseUnreadable) {
		t.Errorf("n:3 err = %v, want ErrResponseUnreadable", got.Err)
	}
	// content 非 JSON → 構造化出力の解析失敗（ErrStructuredParse）。
	if got := byID["p:4"]; !errors.Is(got.Err, ErrStructuredParse) {
		t.Errorf("p:4 err = %v, want ErrStructuredParse", got.Err)
	}
}

// 非 200 応答は失敗にすること（batch 状態確認の例）。
func TestXAIBatchPollNon200IsError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()
	c := NewXAIBatch(http.DefaultClient)
	if _, err := c.PollBatch(context.Background(), Connection{Endpoint: srv.URL}, "batch-123"); err == nil {
		t.Errorf("非 200 でも err が nil")
	}
}

// file upload が非 200 のとき SubmitBatch が失敗を伝播すること（batch 作成まで進まない）。
func TestXAIBatchSubmitPropagatesUploadFailure(t *testing.T) {
	var createCalled bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/batches" {
			createCalled = true
		}
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()
	c := NewXAIBatch(http.DefaultClient)
	_, err := c.SubmitBatch(context.Background(), Connection{Endpoint: srv.URL, APIKey: "k"},
		"grok-4", []BatchRequest{{CustomID: "n:1", Prompt: Prompt{System: "s", User: "u"}}})
	if err == nil {
		t.Errorf("upload 失敗でも err が nil")
	}
	if createCalled {
		t.Errorf("upload 失敗後に batch 作成へ進んだ")
	}
}

// batch 作成の応答に外部 batch ID（batch_id / id）が無いとき失敗にすること。
func TestXAIBatchCreateEmptyIDIsError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/files" {
			_, _ = io.WriteString(w, `{"id":"file-1"}`)
			return
		}
		_, _ = io.WriteString(w, `{}`) // batch_id も id も無い
	}))
	defer srv.Close()
	c := NewXAIBatch(http.DefaultClient)
	if _, err := c.SubmitBatch(context.Background(), Connection{Endpoint: srv.URL},
		"grok-4", []BatchRequest{{CustomID: "n:1"}}); err == nil {
		t.Errorf("外部 batch ID が空でも err が nil")
	}
}

// 結果取得が非 200 のとき失敗にすること。
func TestXAIBatchFetchResultsNon200IsError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
	}))
	defer srv.Close()
	c := NewXAIBatch(http.DefaultClient)
	if _, err := c.FetchResults(context.Background(), Connection{Endpoint: srv.URL}, "batch-123"); err == nil {
		t.Errorf("非 200 でも err が nil")
	}
}
