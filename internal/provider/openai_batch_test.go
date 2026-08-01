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

// R-1-1: OpenAI batch は共通 Chat Completions JSONL を upload し、24 時間の batch を作ること。
func TestOpenAIBatch送信でLunaのChatCompletionsJSONLを使う(t *testing.T) {
	var uploaded string
	var created map[string]string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer secret" {
			t.Errorf("Authorization = %q", r.Header.Get("Authorization"))
		}
		switch r.URL.Path {
		case "/v1/files":
			if err := r.ParseMultipartForm(1 << 20); err != nil {
				t.Fatalf("multipart: %v", err)
			}
			if r.FormValue("purpose") != "batch" {
				t.Errorf("purpose = %q", r.FormValue("purpose"))
			}
			f, _, err := r.FormFile("file")
			if err != nil {
				t.Fatalf("file: %v", err)
			}
			b, _ := io.ReadAll(f)
			uploaded = string(b)
			_, _ = io.WriteString(w, `{"id":"file-in"}`)
		case "/v1/batches":
			_ = json.NewDecoder(r.Body).Decode(&created)
			_, _ = io.WriteString(w, `{"id":"batch-openai"}`)
		default:
			t.Fatalf("想定外の path: %s", r.URL.Path)
		}
	}))
	defer srv.Close()

	c := NewOpenAIBatch(http.DefaultClient)
	id, err := c.SubmitBatch(context.Background(), Connection{Endpoint: srv.URL, APIKey: "secret"}, "gpt-5.6-luna", []BatchRequest{
		{CustomID: "n:1", Prompt: Prompt{System: "system", User: "user"}},
	})
	if err != nil {
		t.Fatalf("SubmitBatch: %v", err)
	}
	if id != "batch-openai" {
		t.Errorf("id = %q", id)
	}
	if created["input_file_id"] != "file-in" || created["endpoint"] != "/v1/chat/completions" || created["completion_window"] != "24h" {
		t.Errorf("batch 作成 body = %v", created)
	}
	if !strings.Contains(uploaded, `"custom_id":"n:1"`) || !strings.Contains(uploaded, `"model":"gpt-5.6-luna"`) || !strings.Contains(uploaded, `"response_format"`) {
		t.Errorf("JSONL = %s", uploaded)
	}
}

// R-1-1: OpenAI の処理中と完了を request_counts と status から共通進行へ写すこと。
func TestOpenAIBatch状態確認は終端状態だけDoneにする(t *testing.T) {
	response := `{"status":"in_progress","request_counts":{"total":5,"completed":2,"failed":1}}`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { _, _ = io.WriteString(w, response) }))
	defer srv.Close()
	c := NewOpenAIBatch(http.DefaultClient)

	got, err := c.PollBatch(context.Background(), Connection{Endpoint: srv.URL}, "batch-1")
	if err != nil {
		t.Fatalf("PollBatch: %v", err)
	}
	if got.Done || got.Pending != 2 || got.Succeeded != 2 || got.Failed != 1 {
		t.Errorf("処理中 = %+v", got)
	}
	response = `{"status":"expired","request_counts":{"total":5,"completed":2,"failed":3}}`
	got, err = c.PollBatch(context.Background(), Connection{Endpoint: srv.URL}, "batch-1")
	if err != nil || !got.Done || got.Pending != 0 {
		t.Errorf("期限切れ = %+v, err=%v", got, err)
	}
}

// R-2-1: failed は外部 batch ID と OpenAI の code・message を含むエラーにする。
func TestOpenAIBatch状態確認はFailed理由を返す(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, `{"id":"batch-1","status":"failed","errors":{"data":[{"code":"token_limit_exceeded","message":"queued token limit"}]}}`)
	}))
	defer srv.Close()
	_, err := NewOpenAIBatch(http.DefaultClient).PollBatch(context.Background(), Connection{Endpoint: srv.URL}, "batch-1")
	if err == nil || !strings.Contains(err.Error(), "batch-1") || !strings.Contains(err.Error(), "token_limit_exceeded") || !strings.Contains(err.Error(), "queued token limit") {
		t.Fatalf("failed error = %v", err)
	}
}

// R-1-2: 成功行と失敗行が混在しても、成功訳と失敗種別を custom_id ごとに返すこと。
func TestOpenAIBatch成功と失敗の結果を同時に取得する(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/batches/batch-1":
			_, _ = io.WriteString(w, `{"status":"completed","output_file_id":"out","error_file_id":"err"}`)
		case "/v1/files/out/content":
			_, _ = io.WriteString(w, `{"custom_id":"n:1","response":{"status_code":200,"body":{"choices":[{"message":{"content":"{\"translation\":\"成功訳\"}"}}]}},"error":null}`+"\n")
		case "/v1/files/err/content":
			_, _ = io.WriteString(w, `{"custom_id":"n:2","response":null,"error":{"code":"rate_limit","message":"later"}}`+"\n")
		default:
			t.Fatalf("想定外の path: %s", r.URL.Path)
		}
	}))
	defer srv.Close()

	got, err := NewOpenAIBatch(http.DefaultClient).FetchResults(context.Background(), Connection{Endpoint: srv.URL}, "batch-1")
	if err != nil {
		t.Fatalf("FetchResults: %v", err)
	}
	if len(got) != 2 || got[0].Translation != "成功訳" || got[0].Err != nil {
		t.Fatalf("結果 = %+v", got)
	}
	if got[1].CustomID != "n:2" || !errors.Is(got[1].Err, ErrServerTransient) {
		t.Errorf("失敗結果 = %+v", got[1])
	}
}

// R-2-2: failed に理由が無い場合も、外部 batch ID と failed を含むエラーにする。
func TestOpenAIBatch結果取得は理由なしFailedを空結果にしない(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, `{"id":"batch-1","status":"failed","output_file_id":null,"error_file_id":null}`)
	}))
	defer srv.Close()
	got, err := NewOpenAIBatch(http.DefaultClient).FetchResults(context.Background(), Connection{Endpoint: srv.URL}, "batch-1")
	if err == nil || got != nil || !strings.Contains(err.Error(), "batch-1") || !strings.Contains(err.Error(), "failed") {
		t.Errorf("結果 = %+v, err=%v", got, err)
	}
}
