package provider

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// LM Studio など base URL が /v1 を含まない指定でも、/v1 配下のエンドポイントへ届くこと。
// API キーが空のときは Authorization ヘッダを付けないこと。
func TestListModelsNormalizesV1AndOmitsAuthWhenKeyEmpty(t *testing.T) {
	var gotPath string
	var hadAuth bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		_, hadAuth = r.Header["Authorization"]
		_, _ = io.WriteString(w, `{"data":[{"id":"qwen2.5-7b"},{"id":"llama-3.1-8b"}]}`)
	}))
	defer srv.Close()

	client := NewOpenAICompatible(http.DefaultClient)
	models, err := client.ListModels(context.Background(), Connection{Endpoint: srv.URL, APIKey: ""})
	if err != nil {
		t.Fatalf("ListModels error: %v", err)
	}

	if gotPath != "/v1/models" {
		t.Errorf("path = %q, want /v1/models", gotPath)
	}
	if hadAuth {
		t.Errorf("Authorization header sent though API key is empty")
	}
	if len(models) != 2 || models[0] != "qwen2.5-7b" || models[1] != "llama-3.1-8b" {
		t.Errorf("models = %v, want [qwen2.5-7b llama-3.1-8b]", models)
	}
}

// Translate は指定 model と原文を /v1/chat/completions へ送り、応答本文を訳文として返すこと。
// API キーがあるときは Bearer で送ること。
func TestTranslateSendsSourceAndReturnsContent(t *testing.T) {
	var gotPath, gotAuth, gotModel string
	var gotBodyHasSource bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotAuth = r.Header.Get("Authorization")
		body, _ := io.ReadAll(r.Body)
		var req struct {
			Model    string `json:"model"`
			Messages []struct {
				Content string `json:"content"`
			} `json:"messages"`
		}
		_ = json.Unmarshal(body, &req)
		gotModel = req.Model
		for _, m := range req.Messages {
			if contains(m.Content, "Ancient Nord text") {
				gotBodyHasSource = true
			}
		}
		_, _ = io.WriteString(w, `{"choices":[{"message":{"content":"古代ノルドの文章"}}]}`)
	}))
	defer srv.Close()

	client := NewOpenAICompatible(http.DefaultClient)
	dest, err := client.Translate(context.Background(),
		Connection{Endpoint: srv.URL, APIKey: "sk-test"}, "qwen2.5-7b", "Ancient Nord text")
	if err != nil {
		t.Fatalf("Translate error: %v", err)
	}

	if gotPath != "/v1/chat/completions" {
		t.Errorf("path = %q, want /v1/chat/completions", gotPath)
	}
	if gotAuth != "Bearer sk-test" {
		t.Errorf("Authorization = %q, want Bearer sk-test", gotAuth)
	}
	if gotModel != "qwen2.5-7b" {
		t.Errorf("model = %q, want qwen2.5-7b", gotModel)
	}
	if !gotBodyHasSource {
		t.Errorf("request body did not include the source text")
	}
	if dest != "古代ノルドの文章" {
		t.Errorf("dest = %q, want 古代ノルドの文章", dest)
	}
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
